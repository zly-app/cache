package cache

import (
	"context"

	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/core"
)

type LoadFn func(ctx context.Context, key string) (interface{}, error)

type options struct {
	Serializer serializer.ISerializer
	Compactor  compactor.ICompactor
	ExpireSec  int
	LoadFn     LoadFn
}

func (c *Cache) newOptions(opts []core.Option) options {
	opt := options{}
	for _, o := range opts {
		o(opt)
	}
	if opt.Serializer == nil {
		opt.Serializer = c.serializer
	}
	if opt.Compactor == nil {
		opt.Serializer = c.serializer
	}
	if opt.ExpireSec == 0 {
		opt.ExpireSec = c.expireSec
	}
	return opt
}

// 设置序列化器
func WithSerializer(serializer serializer.ISerializer) core.Option {
	return func(opts interface{}) {
		opts.(*options).Serializer = serializer
	}
}

// 设置压缩器
func WithCompactor(compactor compactor.ICompactor) core.Option {
	return func(opts interface{}) {
		opts.(*options).Compactor = compactor
	}
}

// 设置有效期
func WithExpire(expireSec int) core.Option {
	return func(opts interface{}) {
		opts.(*options).ExpireSec = expireSec
	}
}

// 设置加载数据函数
func WithLoadFn(fn LoadFn) core.Option {
	return func(opts interface{}) {
		opts.(*options).LoadFn = fn
	}
}
