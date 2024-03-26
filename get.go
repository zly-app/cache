package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/zapp/filter"
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/cache/core"
)

type getReq struct {
	Key            string
	opt            *options
	ExpireSec      int
	ForceLoad      bool // 忽略缓存从加载函数加载数据
	DontWriteCache bool // 不要刷新到缓存
}

func (c *Cache) Get(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	ctx, chain := filter.GetClientFilter(ctx, string(defComponentType), c.cacheName, "Get")
	r := &getReq{
		Key:            key,
		opt:            opt,
		ExpireSec:      opt.ExpireSec,
		ForceLoad:      opt.ForceLoad,
		DontWriteCache: opt.DontWriteCache,
	}
	sp := aPtr
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		r := req.(*getReq)
		sp := rsp

		comData, err := c.getRaw(ctx, r.Key, r.opt)
		if err == nil {
			err = c.unmarshalQuery(comData, sp, r.opt.Serializer, r.opt.Compactor)
		}
		return err
	})
	return err
}
func (c *Cache) getRaw(ctx context.Context, key string, opt *options) ([]byte, error) {
	var bs []byte
	cacheErr := ErrCacheMiss
	if !opt.ForceLoad {
		bs, cacheErr = c.cacheDB.Get(ctx, key)
	}

	if cacheErr == nil {
		return bs, nil
	}

	if cacheErr == ErrCacheMiss {
		utils.Otel.CtxEvent(ctx, "CacheMiss")
	} else {
		utils.Otel.CtxErrEvent(ctx, "GetCacheErr", cacheErr)
	}

	if cacheErr != ErrCacheMiss { // 缓存故障
		if c.ignoreCacheFault {
			logger.Log.Error("从缓存数据库加载数据故障", zap.String("key", key), zap.Error(cacheErr))
		}
		cacheErr = fmt.Errorf("从缓存数据库加载数据故障: err: %v", cacheErr)
		if !c.ignoreCacheFault { // 如果不忽略缓存故障则直接报告错误
			return nil, cacheErr
		}
	}
	if opt.LoadFn == nil {
		return nil, cacheErr
	}

	// 加载数据
	bs, err := c.sf.Do(ctx, key, c.load(opt))
	return bs, err
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
			if opt.DontWriteCache {
				return nil
			}
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
