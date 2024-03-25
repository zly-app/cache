package cache

import (
	"context"

	"github.com/zly-app/zapp/filter"
)

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	ctx, chain := filter.GetClientFilter(ctx, string(defComponentType), c.cacheName, "Del")
	r := &keys
	_, err := chain.Handle(ctx, r, func(ctx context.Context, req interface{}) (interface{}, error) {
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
