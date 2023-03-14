package cache

import (
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

var (
	// msgPack 序列化器
	MsgPackSerializer = serializer.GetSerializer(serializer.MsgPackSerializerName)
	// jsonIter 实现的模拟内置 json 序列化器
	JsonIterStandardSerializer = serializer.GetSerializer(serializer.JsonIterStandardSerializerName)
	// jsonIter序列化器
	JsonIterSerializer = serializer.GetSerializer(serializer.JsonIterSerializerName)
	// json序列化器
	JsonSerializer = serializer.GetSerializer(serializer.JsonSerializerName)
	// yaml 序列化器
	YamlSerializer = serializer.GetSerializer(serializer.YamlSerializerName)
)

var (
	// 不压缩
	NoCompactor = compactor.GetCompactor(compactor.RawCompactorName)
	// zStd 压缩器
	ZStdCompactor = compactor.GetCompactor(compactor.ZStdCompactorName)
	// Gzip 压缩器
	GzipCompactor = compactor.GetCompactor(compactor.GzipCompactorName)
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
