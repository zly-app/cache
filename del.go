package cache

import (
	"context"

	"github.com/zly-app/cache/pkg"
)

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	ctx = pkg.Trace.TraceStart(ctx, "Del", pkg.Trace.AttrKeys(keys))
	defer pkg.Trace.TraceEnd(ctx)

	err := c.del(ctx, keys...)

	pkg.Trace.TraceReply(ctx, "ok", err)
	return err
}
func (c *Cache) del(ctx context.Context, keys ...string) error {
	err := c.cacheDB.Del(ctx, keys...)
	return err
}
