package cache

import (
	"context"

	open_log "github.com/opentracing/opentracing-go/log"
	"github.com/zly-app/zapp/pkg/utils"
)

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	span := utils.Trace.GetChildSpan(ctx, "cache.Del")
	defer span.Finish()
	ctx = utils.Trace.SaveSpan(ctx, span)

	span.LogFields(open_log.Object("keys", keys))

	err := c.del(ctx, keys...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	}
	return err
}
func (c *Cache) del(ctx context.Context, keys ...string) error {
	err := c.cacheDB.Del(ctx, keys...)
	return err
}
