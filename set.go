package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
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

func (c *Cache) MSet(ctx context.Context, dataMap map[string]interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	data := make(map[string][]byte, len(dataMap))
	result := make(map[string]error, len(dataMap))
	for k, v := range dataMap {
		bs, err := c.marshalQuery(v, opt.Serializer, opt.Compactor)
		if err != nil {
			result[k] = fmt.Errorf("编码数据失败: %v", err)
			continue
		}
		data[k] = bs
	}

	if len(data) > 0 {
		cacheResult := c.cacheDB.MSet(ctx, data, opt.ExpireSec)
		for k := range data {
			err := cacheResult[k]
			if err != nil {
				result[k] = err
			}
		}
	}

	return errs.NewQueryErr(result)
}
