package pkg

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/zly-app/zapp/pkg/utils"
)

var Trace = new(traceCli)

type traceCli struct{}

func (c *traceCli) TraceStart(ctx context.Context, method string, attribute utils.OtelSpanKV,
	attributes ...utils.OtelSpanKV) context.Context {
	// 生成新的 span
	ctx, span := utils.Otel.StartSpan(ctx, "cache.method."+method,
		utils.OtelSpanKey("method").String(method))

	attr := []utils.OtelSpanKV{
		attribute,
		c.getOtelSpanKVWithDeadline(ctx),
	}
	attr = append(attr, attributes...)
	utils.Otel.AddSpanEvent(span, "req", attr...)
	return ctx
}

func (*traceCli) TraceEnd(ctx context.Context) {
	span := utils.Otel.GetSpan(ctx)
	utils.Otel.EndSpan(span)
}

func (*traceCli) getOtelSpanKVWithDeadline(ctx context.Context) utils.OtelSpanKV {
	deadline, deadlineOK := ctx.Deadline()
	if !deadlineOK {
		return utils.OtelSpanKey("ctx.deadline").Bool(false)
	}
	d := deadline.Sub(time.Now()) // 剩余时间
	return utils.OtelSpanKey("ctx.deadline").String(d.String())
}

func (c *traceCli) TraceReply(ctx context.Context, reply interface{}, err error) {
	span := utils.Otel.GetSpan(ctx)
	if err == nil {
		text, _ := jsoniter.MarshalToString(reply)
		utils.Otel.AddSpanEvent(span, "reply",
			utils.OtelSpanKey("reply").String(text),
			c.getOtelSpanKVWithDeadline(ctx),
		)
		return
	}

	utils.Otel.MarkSpanAnError(span, true)

	utils.Otel.AddSpanEvent(span, "reply",
		utils.OtelSpanKey("err.detail").String(err.Error()),
		c.getOtelSpanKVWithDeadline(ctx),
	)
}

func (c *traceCli) AttrKey(key string) utils.OtelSpanKV {
	return utils.OtelSpanKey("key").String(key)
}
func (c *traceCli) AttrKeys(keys []string) utils.OtelSpanKV {
	return utils.OtelSpanKey("keys").StringSlice(keys)
}
