package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/cache/core"
)

func (c *Cache) Set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	bs, err := c.marshalQuery(data, opt.Serializer, opt.Compactor)
	if err != nil {
		return fmt.Errorf("编码数据失败: %v", err)
	}

	err = c.cacheDB.Set(ctx, key, bs, opt.ExpireSec)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %v", err)
	}
	return nil
}

func (c *Cache) MSet(ctx context.Context, dataMap map[string]interface{}, opts ...core.Option) map[string]error {
	//TODO implement me
	panic("implement me")
}
