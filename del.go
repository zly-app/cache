package cache

import (
	"context"
)

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	err := c.cacheDB.Del(ctx, keys...)
	return err
}
