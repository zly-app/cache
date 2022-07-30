package cache

import (
	"fmt"
	"strings"
)

const (
	defCompactor        = "raw"
	defSerializer       = "msgpack"
	defSingleFlight     = "single"
	defExpireSec        = 0
	defIgnoreCacheFault = false

	defCacheDB_Type          = "memory"
	defCacheDB_Memory_SizeMD = 1
)

type Config struct {
	Compactor        string // 默认压缩器名, 可选 raw, zstd, gzip
	Serializer       string // 默认序列化器名, 可选 msgpack, jsoniter_standard, jsoniter, json, yaml
	SingleFlight     string // 默认单跑模块, 可选 no, single
	ExpireSec        int    // 默认有效时间, 秒, <= 0 表示永久
	IgnoreCacheFault bool   // 是否忽略缓存数据库故障, 如果设为true, 在缓存数据库故障时从加载器获取数据, 这会导致缓存击穿. 如果设为false, 在缓存数据库故障时直接返回错误
	CacheDB          struct {
		Type   string // 缓存数据库类型, 支持 no, memory, redis
		Memory struct {
			SizeMB int // 分配内存大小, 单位mb, 单条数据大小不能超过该值的 1/1024
		}
	}
}

func NewConfig() *Config {
	conf := &Config{
		Compactor:        defCompactor,
		Serializer:       defSerializer,
		SingleFlight:     defSingleFlight,
		ExpireSec:        defExpireSec,
		IgnoreCacheFault: defIgnoreCacheFault,
	}
	conf.CacheDB.Type = defCacheDB_Type
	conf.CacheDB.Memory.SizeMB = defCacheDB_Memory_SizeMD
	return conf
}

func (conf *Config) Check() error {
	switch v := strings.ToLower(conf.CacheDB.Type); v {
	case "":
		conf.CacheDB.Type = defCacheDB_Type
	case "memory", "redis":
	default:
		return fmt.Errorf("不支持的CacheDB: %v", v)
	}

	switch v := strings.ToLower(conf.Compactor); v {
	case "":
		conf.Compactor = defCompactor
	case "raw", "gzip", "zstd":
	default:
		return fmt.Errorf("不支持的Compactor: %v", v)
	}

	switch v := strings.ToLower(conf.Serializer); v {
	case "":
		conf.Serializer = defSerializer
	case "json", "jsoniter", "jsoniter_standard", "msgpack", "yaml":
	default:
		return fmt.Errorf("不支持的Serializer: %v", v)
	}

	switch v := strings.ToLower(conf.SingleFlight); v {
	case "":
		conf.SingleFlight = defSingleFlight
	case "no", "single":
	default:
		return fmt.Errorf("不支持的Serializer: %v", v)
	}

	if conf.CacheDB.Memory.SizeMB < 1 {
		conf.CacheDB.Memory.SizeMB = defCacheDB_Memory_SizeMD
	}
	return nil
}
