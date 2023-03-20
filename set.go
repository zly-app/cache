package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/zapp/pkg/utils"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/pkg"
)

func (c *Cache) Set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	attr := []utils.OtelSpanKV{
		pkg.Trace.AttrKey(key),
		pkg.Trace.AttrData(data),
	}
	attr = append(attr, opt.MakeTraceAttr()...)
	ctx = pkg.Trace.TraceStart(ctx, c.cacheName, "Set", attr...)
	defer pkg.Trace.TraceEnd(ctx)

	bs, err := c.marshalQuery(data, opt.Serializer, opt.Compactor)
	if err == nil {
		err = c.set(ctx, key, bs, opt)
	}

	pkg.Trace.TraceReply(ctx, "ok", err)
	return err
}

func (c *Cache) set(ctx context.Context, key string, bs []byte, opt *options) error {
	err := c.cacheDB.Set(ctx, key, bs, opt.ExpireSec)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %v", err)
	}
	return nil
}
