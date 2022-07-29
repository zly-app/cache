package cache

import (
	"context"

	"github.com/zly-app/cache/core"
)

func (c *Cache) Get(ctx context.Context, key string, aPtr interface{}, opts ...core.GetOption) error {
	//TODO implement me
	panic("implement me")
}

func (c *Cache) MGet(ctx context.Context, aPtrMap map[string]interface{}, opts ...core.GetOption) map[string]error {
	//TODO implement me
	panic("implement me")
}

func (c *Cache) MGetSlice(ctx context.Context, keys []string, slicePtr interface{}, opts ...core.GetOption) map[string]error {
	//TODO implement me
	panic("implement me")
}
