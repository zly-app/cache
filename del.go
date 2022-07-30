package cache

import (
	"context"

	"github.com/zly-app/cache/errs"
)

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	result := c.cacheDB.Del(ctx, keys...)
	return errs.NewQueryErr(result)
}
