package cache

import (
	"context"
	"fmt"

	open_log "github.com/opentracing/opentracing-go/log"
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/cache/v2/core"
)

func (c *Cache) Get(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	span := utils.Trace.GetChildSpan(ctx, "cache.Get")
	defer span.Finish()
	ctx = utils.Trace.SaveSpan(ctx, span)

	span.LogFields(open_log.String("key", key))

	err := c.get(ctx, key, aPtr, opts...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	}
	return err
}
func (c *Cache) get(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	var bs []byte
	cacheErr := ErrCacheMiss
	if !opt.ForceLoad {
		bs, cacheErr = c.cacheDB.Get(ctx, key)
	}

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
	if opt.LoadFn == nil {
		return cacheErr
	}

	// 加载数据
	bs, err := c.sf.Do(ctx, key, c.load(opt))
	if err != nil {
		return err
	}
	return c.unmarshalQuery(bs, aPtr, opt.Serializer, opt.Compactor)
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
