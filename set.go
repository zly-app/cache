package cache

import (
	"context"
	"fmt"

	"github.com/zly-app/zapp/filter"

	"github.com/zly-app/cache/core"
)

type setReq struct {
	Key       string
	Data      interface{}
	opt       *options
	ExpireSec int
}

func (c *Cache) Set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	opt := c.newOptions(opts)
	defer putOptions(opt)

	ctx, chain := filter.GetClientFilter(ctx, string(defComponentType), c.cacheName, "Set")
	r := &setReq{
		Key:       key,
		Data:      data,
		opt:       opt,
		ExpireSec: opt.ExpireSec,
	}
	_, err := chain.Handle(ctx, r, func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*setReq)
		bs, err := c.marshalQuery(r.Data, r.opt.Serializer, r.opt.Compactor)
		if err == nil {
			err = c.set(ctx, key, bs, opt)
		}
		return nil, err
	})
	return err
}

func (c *Cache) set(ctx context.Context, key string, bs []byte, opt *options) error {
	err := c.cacheDB.Set(ctx, key, bs, opt.ExpireSec)
	if err != nil {
		return fmt.Errorf("写入缓存失败: %v", err)
	}
	return nil
}
