package cache

import (
	"context"
	"errors"

	"github.com/zly-app/zapp/filter"

	"github.com/zly-app/cache/core"
)

func (c *Cache) SingleFlightDo(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	opts = append([]core.Option{WithForceLoad(true)}, opts...)
	opt := c.newOptions(opts)
	opt.ForceLoad = true // 强行从加载函数加载
	defer putOptions(opt)

	ctx, chain := filter.GetClientFilter(ctx, string(defComponentType), c.cacheName, "SingleFlightDo")
	r := &getReq{
		Key:            key,
		opt:            opt,
		ExpireSec:      opt.ExpireSec,
		ForceLoad:      opt.ForceLoad,
		DontWriteCache: opt.DontWriteCache,
	}
	sp := aPtr
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		r := req.(*getReq)
		sp := rsp

		comData, err := c.singleFlightDo(ctx, r.Key, r.opt)
		if err == nil {
			err = c.unmarshalQuery(comData, sp, r.opt.Serializer, r.opt.Compactor)
		}
		return err
	})
	return err
}

func (c *Cache) singleFlightDo(ctx context.Context, key string, opt *options) ([]byte, error) {
	if opt.LoadFn == nil {
		return nil, errors.New("LoadFn is nil")
	}

	bs, err := c.sf.Do(ctx, key, c.load(opt))
	return bs, err
}
