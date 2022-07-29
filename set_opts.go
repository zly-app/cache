package cache

import (
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/core"
)

type setOptions struct {
	serializer serializer.ISerializer
	compactor  compactor.ICompactor
	expireSec  int
}

func (c *Cache) newSetOptions(opts []core.SetOption) setOptions {
	opt := setOptions{}
	for _, o := range opts {
		o(opt)
	}
	if opt.serializer == nil {
		opt.serializer = c.serializer
	}
	if opt.compactor == nil {
		opt.serializer = c.serializer
	}
	if opt.expireSec == 0 {
		opt.expireSec = c.expireSec
	}
	return opt
}

// 设置序列化器
func WithSetSerializer(serializer serializer.ISerializer) core.SetOption {
	return func(opts interface{}) {
		opts.(*setOptions).serializer = serializer
	}
}

// 设置压缩器
func WithSetCompactor(compactor compactor.ICompactor) core.SetOption {
	return func(opts interface{}) {
		opts.(*setOptions).compactor = compactor
	}
}

// 设置有效期
func WithSetExpire(expireSec int) core.SetOption {
	return func(opts interface{}) {
		opts.(*setOptions).expireSec = expireSec
	}
}
