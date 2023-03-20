package cache

import (
	"context"
	"errors"

	"github.com/zly-app/zapp/pkg/utils"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/pkg"
)

func (c *Cache) SingleFlightDo(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	opts = append([]core.Option{WithForceLoad(true)}, opts...)
	opt := c.newOptions(opts)
	defer putOptions(opt)
	opt.ForceLoad = true // 强行从加载函数加载

	attr := []utils.OtelSpanKV{
		pkg.Trace.AttrKey(key),
	}
	attr = append(attr, opt.MakeTraceAttr()...)
	ctx = pkg.Trace.TraceStart(ctx, "SingleFlightDo", attr...)
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
