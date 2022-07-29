package core

import (
	"context"
)

type GetOption func(opts interface{})
type SetOption func(opts interface{})

type ICache interface {
	// 获取数据并放入 aPtr 中
	Get(ctx context.Context, key string, aPtr interface{}, opts ...GetOption) error
	// 批量获取数据
	MGet(ctx context.Context, aPtrMap map[string]interface{}, opts ...GetOption) map[string]error
	// 批量获取数据并将所有数据都放入 slicePtr 中, slicePtr 是一个带指针的切片
	MGetSlice(ctx context.Context, keys []string, slicePtr interface{}, opts ...GetOption) map[string]error

	// 设置数据
	Set(ctx context.Context, key string, aPtr interface{}, opts ...SetOption) error
	// 批量设置数据
	MSet(ctx context.Context, aPtrMap map[string]interface{}, opts ...SetOption) map[string]error

	// 删除
	Del(ctx context.Context, keys ...string) map[string]error
	// 关闭
	Close() error
}
