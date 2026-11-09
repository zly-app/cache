package cache

import (
	"context"

	"github.com/zly-app/zapp/filter"
)

// handle 在接入 zapp 时走 filter 链；独立使用（app == nil）时直接执行。
func (c *Cache) handle(ctx context.Context, action string, req interface{}, fn func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
	if c.app == nil {
		return fn(ctx, req)
	}
	ctx, chain := filter.GetClientFilter(ctx, string(defComponentType), c.cacheName, action)
	return chain.Handle(ctx, req, fn)
}

// handleInject 同 handle，用于把结果注入 rsp 的处理器。
func (c *Cache) handleInject(ctx context.Context, action string, req, rsp interface{}, fn func(ctx context.Context, req, rsp interface{}) error) error {
	if c.app == nil {
		return fn(ctx, req, rsp)
	}
	ctx, chain := filter.GetClientFilter(ctx, string(defComponentType), c.cacheName, action)
	return chain.HandleInject(ctx, req, rsp, fn)
}
