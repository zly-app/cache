package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

func (c *Cache) Get(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	bs, cacheErr := c.cacheDB.Get(ctx, key)
	if cacheErr == nil {
		return c.unmarshalQuery(bs, aPtr, opt.Serializer, opt.Compactor)
	}
	if cacheErr != ErrCacheMiss { // 缓存故障
		if c.ignoreCacheFault {
			logger.Log.Error("从缓存数据库加载数据故障", zap.String("key", key), zap.Error(cacheErr))
		}
		cacheErr = fmt.Errorf("从缓存数据库加载数据故障: err: %v", cacheErr)
		if !c.ignoreCacheFault { // 如果不忽略缓存故障则直接报告错误
			return cacheErr
		}
	}

	if opt.LoadFn == nil { // 加载器不存在, 直接报告错误
		return cacheErr
	}

	// 加载数据
	bs, err := c.sf.Do(ctx, key, c.load(opt))
	if err != nil {
		return err
	}
	return c.unmarshalQuery(bs, aPtr, opt.Serializer, opt.Compactor)
}

func (c *Cache) MGet(ctx context.Context, aPtrMap map[string]interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	keys := make([]string, len(aPtrMap))
	aPtrList := make([]interface{}, len(aPtrMap))
	index := 0
	for k, v := range aPtrMap {
		keys[index] = k
		aPtrList[index] = v
		index++
	}

	// 从缓存获取结果
	cacheResults := c.cacheDB.MGet(ctx, keys...)

	// 整理已得到的结果
	result := make(map[string]error, len(aPtrMap))  // 结果数据
	needLoadKeys := make([]string, 0, len(aPtrMap)) // 需要额外加载的key
	for index, key := range keys {
		cacheResult, ok := cacheResults[key]
		if !ok || cacheResult.Err == ErrCacheMiss {
			result[key] = ErrCacheMiss
			needLoadKeys = append(needLoadKeys, key)
			continue
		}

		if cacheResult.Err == nil {
			err := c.unmarshalQuery(cacheResult.Data, aPtrList[index], opt.Serializer, opt.Compactor)
			if err != nil {
				result[key] = err
			}
			continue
		}

		// 缓存故障
		result[key] = fmt.Errorf("从缓存数据库加载数据故障: err: %v", cacheResult.Err)
		if c.ignoreCacheFault { // 如果忽略缓存故障, 则这些key也需要加载数据
			logger.Log.Error("从缓存数据库加载数据故障", zap.String("key", key), zap.Error(cacheResult.Err))
			needLoadKeys = append(needLoadKeys, key)
		}
	}

	// 如果不需要额外加载或者加载数据函数为空则直接返回
	if len(needLoadKeys) == 0 || opt.LoadFn == nil {
		return errs.NewQueryErr(result)
	}

	// 加载数据
	for _, key := range needLoadKeys {
		bs, err := c.sf.Do(ctx, key, c.load(opt))
		if err != nil {
			result[key] = err // 重新设置err
			continue
		}
		result[key] = c.unmarshalQuery(bs, aPtrMap[key], opt.Serializer, opt.Compactor)
	}
	return errs.NewQueryErr(result)
}

func (c *Cache) MGetSlice(ctx context.Context, keys []string, slicePtr interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	// 从缓存获取结果
	cacheResults := c.cacheDB.MGet(ctx, keys...)

	// 整理已得到的结果
	result := make(map[string]error, len(keys))  // 结果数据
	needLoadKeys := make([]string, 0, len(keys)) // 需要额外加载的key
	for _, key := range keys {
		cacheResult, ok := cacheResults[key]
		if !ok || cacheResult.Err == ErrCacheMiss {
			result[key] = ErrCacheMiss
			needLoadKeys = append(needLoadKeys, key)
			continue
		}

		if cacheResult.Err == nil {
			continue
		}

		// 缓存故障
		result[key] = fmt.Errorf("从缓存数据库加载数据故障: err: %v", cacheResult.Err)
		if c.ignoreCacheFault { // 如果忽略缓存故障, 则这些key也需要加载数据
			logger.Log.Error("从缓存数据库加载数据故障", zap.String("key", key), zap.Error(cacheResult.Err))
			needLoadKeys = append(needLoadKeys, key)
		}
	}

	// 加载数据
	if len(needLoadKeys) > 0 && opt.LoadFn != nil {
		for _, key := range needLoadKeys {
			bs, err := c.sf.Do(ctx, key, c.load(opt))
			if err != nil {
				result[key] = err // 重新设置err
				continue
			}
			cacheResults[key] = core.CacheResult{Data: bs} // 替换数据
		}
	}

	// 反序列化
	unResults := c.unmarshalMQuerySlice(keys, cacheResults, slicePtr, opt.Serializer, opt.Compactor)
	for key, err := range unResults {
		result[key] = err // 重新设置err
	}
	return errs.NewQueryErr(result)
}

func (c *Cache) load(opt *options) core.LoadInvoke {
	return func(ctx context.Context, key string) (bs []byte, err error) {
		err = utils.Recover.WrapCall(func() error {
			// 加载数据
			data, err := opt.LoadFn(ctx, key)
			if err != nil {
				return fmt.Errorf("从加载函数加载数据失败: %v", err)
			}

			// 编码数据
			bs, err = c.marshalQuery(data, opt.Serializer, opt.Compactor)
			if err != nil {
				return fmt.Errorf("编码数据失败: %v", err)
			}

			// 写入缓存
			cacheErr := c.cacheDB.Set(ctx, key, bs, opt.ExpireSec)
			if cacheErr != nil {
				if !c.ignoreCacheFault {
					return fmt.Errorf("写入缓存失败: %v", cacheErr)
				}
				logger.Log.Error("写入缓存失败", zap.String("key", key), zap.Error(err))
			}
			return nil
		})
		return bs, err
	}
}
