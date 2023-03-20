package cache

import (
	"context"
	"errors"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/pkg"
)

func (c *Cache) SingleFlightDo(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	opts = append([]core.Option{WithForceLoad(true)}, opts...)
	opt := c.newOptions(opts)
	defer putOptions(opt)
	opt.ForceLoad = true // 强行从加载函数加载

	ctx = pkg.Trace.TraceStart(ctx, "SingleFlightDo", pkg.Trace.AttrKey(key), opt.MakeTraceAttr()...)
	defer pkg.Trace.TraceEnd(ctx)

	comData, err := c.singleFlightDo(ctx, key, opt)
	if err == nil {
		err = c.unmarshalQuery(comData, aPtr, opt.Serializer, opt.Compactor)
	}

	pkg.Trace.TraceReply(ctx, aPtr, err)
	return err
}

func (c *Cache) singleFlightDo(ctx context.Context, key string, opt *options) ([]byte, error) {
	if opt.LoadFn == nil {
		return nil, errors.New("LoadFn is nil")
	}

	bs, err := c.sf.Do(ctx, key, c.load(opt))
	return bs, err
}
