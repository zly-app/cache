package no_sf

import (
	"context"

	"github.com/zly-app/cache/core"
)

type NoSingleFlight struct{}

func (n NoSingleFlight) Do(ctx context.Context, key string, invoke core.LoadInvoke) ([]byte, error) {
	return invoke(ctx, key)
}

// 一个关闭并发查询控制的ISingleFlight
func NewNoSingleFlight() core.ISingleFlight {
	return NoSingleFlight{}
}
