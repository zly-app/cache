package core

import (
	"context"
)

type Option func(opts interface{})

type ICache interface {
	// 获取数据并放入 aPtr 中
	Get(ctx context.Context, key string, aPtr interface{}, opts ...Option) error
	// 批量获取数据
	MGet(ctx context.Context, aPtrMap map[string]interface{}, opts ...Option) map[string]error
	// 批量获取数据并将所有数据都放入 slicePtr 中, slicePtr 是一个带指针的切片
	MGetSlice(ctx context.Context, keys []string, slicePtr interface{}, opts ...Option) map[string]error

	// 设置数据
	Set(ctx context.Context, key string, data interface{}, opts ...Option) error
	// 批量设置数据
	MSet(ctx context.Context, dataMap map[string]interface{}, opts ...Option) map[string]error

	// 删除
	Del(ctx context.Context, keys ...string) map[string]error
	// 关闭
	Close() error
}
