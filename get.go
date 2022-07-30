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
	if cacheErr != errs.CacheMiss { // 非缓存未命中错误
		cacheErr = fmt.Errorf("从缓存数据库加载数据失败: err: %v", cacheErr)
		if !c.ignoreCacheFault { // 直接报告错误
			return cacheErr
		}
	}

	if opt.LoadFn == nil { // 加载器不存在, 直接报告错误
		return cacheErr
	}

	// 加载数据
	bs, err := c.sf.Do(ctx, key, c.load(ctx, key, opt))
	if err != nil {
		return err
	}
	return c.unmarshalQuery(bs, aPtr, opt.Serializer, opt.Compactor)
}

func (c *Cache) MGet(ctx context.Context, aPtrMap map[string]interface{}, opts ...core.Option) map[string]error {
	//TODO implement me
	panic("implement me")
}

func (c *Cache) MGetSlice(ctx context.Context, keys []string, slicePtr interface{}, opts ...core.Option) map[string]error {
	//TODO implement me
	panic("implement me")
}

func (c *Cache) load(ctx context.Context, key string, opt *options) core.LoadInvoke {
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
