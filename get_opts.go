package cache

import (
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/core"
)

type getOptions struct {
	serializer serializer.ISerializer
	compactor  compactor.ICompactor
	expireSec  int
}

func (c *Cache) newGetOptions(opts []core.GetOption) getOptions {
	opt := getOptions{}
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
func WithGetSerializer(serializer serializer.ISerializer) core.GetOption {
	return func(opts interface{}) {
		opts.(*setOptions).serializer = serializer
	}
}

// 设置压缩器
func WithGetCompactor(compactor compactor.ICompactor) core.GetOption {
	return func(opts interface{}) {
		opts.(*setOptions).compactor = compactor
	}
}

// 设置有效期
func WithGetExpire(expireSec int) core.GetOption {
	return func(opts interface{}) {
		opts.(*setOptions).expireSec = expireSec
	}
}
