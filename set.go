package cache

import (
	"context"
	"fmt"

	open_log "github.com/opentracing/opentracing-go/log"
	"github.com/zly-app/zapp/pkg/utils"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

func (c *Cache) Set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	if c.disableOpenTrace {
		return c.set(ctx, key, data, opts...)
	}

	span := utils.Trace.GetChildSpan(ctx, "cache.Set")
	defer span.Finish()
	ctx = utils.Trace.SaveSpan(ctx, span)

	span.LogFields(open_log.String("key", key))

	err := c.set(ctx, key, data, opts...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	}
	return err
}

func (c *Cache) set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	bs, err := c.marshalQuery(data, opt.Serializer, opt.Compactor)
	if err != nil {
		return fmt.Errorf("编码数据失败: %v", err)
	}

	err = c.cacheDB.Set(ctx, key, bs, opt.ExpireSec)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %v", err)
	}
	return nil
}

func (c *Cache) MSet(ctx context.Context, dataMap map[string]interface{}, opts ...core.Option) error {
	if c.disableOpenTrace {
		return c.mSet(ctx, dataMap, opts...)
	}

	span := utils.Trace.GetChildSpan(ctx, "cache.MSet")
	defer span.Finish()
	ctx = utils.Trace.SaveSpan(ctx, span)

	keys := make([]string, 0, len(dataMap))
	for k := range dataMap {
		keys = append(keys, k)
	}
	span.LogFields(open_log.Object("keys", keys))

	err := c.mSet(ctx, dataMap, opts...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	}
	return err
}
func (c *Cache) mSet(ctx context.Context, dataMap map[string]interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	data := make(map[string][]byte, len(dataMap))
	result := make(map[string]error, len(dataMap))
	for k, v := range dataMap {
		bs, err := c.marshalQuery(v, opt.Serializer, opt.Compactor)
		if err != nil {
			result[k] = fmt.Errorf("编码数据失败: %v", err)
			continue
		}
		data[k] = bs
	}

	if len(data) > 0 {
		cacheResult := c.cacheDB.MSet(ctx, data, opt.ExpireSec)
		for k := range data {
			err := cacheResult[k]
			if err != nil {
				result[k] = err
			}
		}
	}

	return errs.NewQueryErr(result)
}
