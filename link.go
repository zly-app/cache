package cache

import (
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

var (
	// msgPack 序列化器
	MsgPackSerializer = GetSerializer(serializer.MsgPackSerializerName)
	// jsonIter 实现的模拟内置 json 序列化器
	JsonIterStandardSerializer = GetSerializer(serializer.JsonIterStandardSerializerName)
	// jsonIter序列化器
	JsonIterSerializer = GetSerializer(serializer.JsonIterSerializerName)
	// json序列化器
	JsonSerializer = GetSerializer(serializer.JsonSerializerName)
	// yaml 序列化器
	YamlSerializer = GetSerializer(serializer.YamlSerializerName)

	GetSerializer = func(name string) core.ISerializer {
		return serializer.GetSerializer(name)
	}
)

var (
	// 不压缩
	NoCompactor = GetCompactor(compactor.RawCompactorName)
	// zStd 压缩器
	ZStdCompactor = GetCompactor(compactor.ZStdCompactorName)
	// Gzip 压缩器
	GzipCompactor = GetCompactor(compactor.GzipCompactorName)
	// 获取压缩器, 压缩器不存在会panic
	GetCompactor = func(name string) core.ICompactor {
		return compactor.GetCompactor(name)
	}
)

var (
	// 缓存不存在
	ErrCacheMiss = errs.CacheMiss
	// 数据为nil
	ErrDataIsNil = errs.DataIsNil
)

type (
	ICache = core.ICache
	LoadFn = core.LoadFn
)
