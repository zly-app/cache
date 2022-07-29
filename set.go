package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/cache/core"
)

func (c *Cache) Set(ctx context.Context, key string, aPtr interface{}, opts ...core.SetOption) error {
	opt := c.newSetOptions(opts)
	bs, err := c.marshalQuery(aPtr, opt.serializer, opt.compactor)
	if err != nil {
		return fmt.Errorf("编码数据失败: %s", err)
	}

	err = c.cacheDB.Set(ctx, key, bs, opt.expireSec)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %s", err)
	}
	return nil
}

func (c *Cache) MSet(ctx context.Context, aPtrMap map[string]interface{}, opts ...core.SetOption) map[string]error {
	//TODO implement me
	panic("implement me")
}
