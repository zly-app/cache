package cache

import (
	"context"
	"fmt"

	open_log "github.com/opentracing/opentracing-go/log"
	"github.com/zly-app/zapp/pkg/utils"

	"github.com/zly-app/cache/core"
)

func (c *Cache) Set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	span := utils.Trace.GetChildSpan(ctx, "cache.Set")
	defer span.Finish()
	ctx = utils.Trace.SaveSpan(ctx, span)

	span.LogFields(open_log.String("key", key))

	opt := c.newOptions(opts)
	defer putOptions(opt)

	bs, err := c.marshalQuery(data, opt.Serializer, opt.Compactor)
	if err == nil {
		err = c.set(ctx, key, bs, opt)
	}
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	}
	return err
}

func (c *Cache) set(ctx context.Context, key string, bs []byte, opt *options) error {
	err := c.cacheDB.Set(ctx, key, bs, opt.ExpireSec)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %v", err)
	}
	return nil
}
