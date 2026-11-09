package cache

import (
	"context"
)

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	r := &keys
	_, err := c.handle(ctx, "Del", r, func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*[]string)
		err := c.del(ctx, *r...)
		return nil, err
	})
	return err
}

func (c *Cache) del(ctx context.Context, keys ...string) error {
	err := c.cacheDB.Del(ctx, keys...)
	return err
}
