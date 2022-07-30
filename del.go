package cache

import (
	"context"
)

func (c *Cache) Del(ctx context.Context, keys ...string) map[string]error {
	return c.cacheDB.Del(ctx, keys...)
}
