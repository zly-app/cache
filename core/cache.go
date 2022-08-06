package core

import (
	"context"
)

type Option func(opts interface{})

type LoadFn func(ctx context.Context, key string) (interface{}, error)

type ICache interface {
	// 获取数据并放入 aPtr 中
	Get(ctx context.Context, key string, aPtr interface{}, opts ...Option) error
	// 批量获取数据
	MGet(ctx context.Context, aPtrMap map[string]interface{}, opts ...Option) error
	// 批量获取数据并将所有数据都放入 slicePtr 中, slicePtr 是一个带指针的切片, 结果中只存有error的key
	MGetSlice(ctx context.Context, keys []string, slicePtr interface{}, opts ...Option) error

	// 设置数据
	Set(ctx context.Context, key string, data interface{}, opts ...Option) error
	// 批量设置数据
	MSet(ctx context.Context, dataMap map[string]interface{}, opts ...Option) error

	// 单跑执行
	SingleFlightDo(ctx context.Context, key string, opts ...Option) error

	// 删除
	Del(ctx context.Context, keys ...string) error
	// 关闭
	Close() error
}
