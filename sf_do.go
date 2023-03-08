package cache

import (
	"context"
	"errors"

	open_log "github.com/opentracing/opentracing-go/log"
	"github.com/zly-app/zapp/pkg/utils"

	"github.com/zly-app/cache/core"
)

func (c *Cache) SingleFlightDo(ctx context.Context, key string, opts ...core.Option) error {
	span := utils.Trace.GetChildSpan(ctx, "cache.SingleFlightDo")
	defer span.Finish()
	ctx = utils.Trace.SaveSpan(ctx, span)

	span.LogFields(open_log.String("key", key))

	err := c.singleFlightDo(ctx, key, opts...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	}
	return err
}

func (c *Cache) singleFlightDo(ctx context.Context, key string, opts ...core.Option) error {
	opts = append([]core.Option{WithForceLoad(true)}, opts...)
	opt := c.newOptions(opts)
	defer putOptions(opt)
	opt.ForceLoad = true // 强行从加载函数加载

	if opt.LoadFn == nil {
		return errors.New("LoadFn is nil")
	}

	_, err := c.sf.Do(ctx, key, c.load(opt))
	return err
}
