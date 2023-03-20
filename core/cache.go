package core

import (
	"context"
)

type Option func(opts interface{})

type LoadFn func(ctx context.Context, key string) (interface{}, error)

type ICache interface {
	// 获取数据并放入 aPtr 中
	Get(ctx context.Context, key string, aPtr interface{}, opts ...Option) error

	// 设置数据
	Set(ctx context.Context, key string, data interface{}, opts ...Option) error

	// 单跑执行, 忽略缓存直接从db加载数据, 默认不会自动写入缓存, 必须设置 LoadFn
	SingleFlightDo(ctx context.Context, key string, aPtr interface{}, opts ...Option) error

	// 删除
	Del(ctx context.Context, keys ...string) error

	// 关闭
	Close() error
}
