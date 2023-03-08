package cache

import (
	"sync"

	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/core"
)

var optionsPool = sync.Pool{New: func() interface{} { return &options{} }}

type options struct {
	Serializer     serializer.ISerializer
	Compactor      compactor.ICompactor
	ExpireSec      int
	LoadFn         LoadFn
	ForceLoad      bool // 忽略缓存从加载函数加载数据
	DontWriteCache bool // 不要刷新到缓存
}

func getOptions() *options {
	return optionsPool.Get().(*options)
}
func putOptions(opt *options) {
	opt.Serializer = nil
	opt.Compactor = nil
	opt.ExpireSec = 0
	opt.LoadFn = nil
	opt.ForceLoad = false
	opt.DontWriteCache = false
	optionsPool.Put(opt)
}

func (c *Cache) newOptions(opts []core.Option) *options {
	opt := getOptions()
	for _, o := range opts {
		o(opt)
	}
	if opt.Serializer == nil {
		opt.Serializer = c.serializer
	}
	if opt.Compactor == nil {
		opt.Compactor = c.compactor
	}
	if opt.ExpireSec == 0 {
		opt.ExpireSec = c.expireSec
	}
	return opt
}

// 设置序列化器, 如果设为nil则使用默认序列化器
func WithSerializer(serializer serializer.ISerializer) core.Option {
	return func(opts interface{}) {
		opts.(*options).Serializer = serializer
	}
}

// 设置压缩器, 如果设为nil则使用默认压缩器
func WithCompactor(compactor compactor.ICompactor) core.Option {
	return func(opts interface{}) {
		opts.(*options).Compactor = compactor
	}
}

// 设置有效期, expireSec < 0 表示永不过期, expireSec = 0 表示使用默认值
func WithExpire(expireSec int) core.Option {
	return func(opts interface{}) {
		opts.(*options).ExpireSec = expireSec
	}
}

// 设置加载数据函数, 当缓存未命中或缓存故障时, 会调用它获取数据, 如果设置了 SingleFlight 会在之前前经过 SingleFlight.
func WithLoadFn(fn LoadFn) core.Option {
	return func(opts interface{}) {
		opts.(*options).LoadFn = fn
	}
}

// 忽略缓存从加载函数加载数据
func WithForceLoad(dontWriteCache bool) core.Option {
	return func(opts interface{}) {
		opt := opts.(*options)
		opt.ForceLoad = true
		opt.DontWriteCache = dontWriteCache
	}
}
