package cache

import (
	"fmt"
	"strings"

	"github.com/zly-app/component/redis"
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/v2/single_flight"
)

const (
	defCompactor        = "raw"
	defSerializer       = "sonic_std"
	defSingleFlight     = "single"
	defExpireSec        = 300
	defIgnoreCacheFault = false

	defCacheDB_Type = "bigcache"

	defCacheDB_BigCache_Shards             = 1024
	defCacheDB_BigCache_CleanTimeSec       = 60
	defCacheDB_BigCache_MaxEntriesInWindow = 1000 * 10 * 60
	defCacheDB_BigCache_MaxEntrySize       = 500

	defCacheDB_FreeCache_SizeMB = 1
)

type Config struct {
	Compactor        string // 默认压缩器名, 可选 raw, zstd, gzip
	Serializer       string // 默认序列化器名, 可选 sonic, sonic_std, msgpack, jsoniter, jsoniter_standard, json, yaml
	SingleFlight     string // 默认单跑模块, 可选 no, single
	ExpireSec        int    // 默认过期时间, 秒, < 1 表示永久
	IgnoreCacheFault bool   // 是否忽略缓存数据库故障, 如果设为true, 在缓存数据库故障时从加载器获取数据, 这会导致缓存击穿. 如果设为false, 在缓存数据库故障时直接返回错误
	CacheDB          struct {
		Type     string // 缓存数据库类型, 支持 no, bigcache, freecache, redis
		BigCache struct {
			Shards             int  // 分片数, 必须是2的幂
			CleanTimeSec       int  // 清理周期秒数, 为 0 时不自动清理.
			MaxEntriesInWindow int  // 初始化时申请允许储存的条目数的内存, 当实际使用量超过当前最大量时会触发内存重分配
			MaxEntrySize       int  // 初始化时申请的每个条目的占用内存, 单位字节, 当实际使用量超过当前最大量时会触发内存重分配
			HardMaxCacheSize   int  // 最大占用内存大小, 单位 mb, 0 表示不限制
			ExactExpire        bool // 精确过期时间
		}
		FreeCache struct {
			SizeMB int // 分配内存大小, 单位mb, 单条数据大小不能超过该值的 1/1024
		}
		Redis redis.RedisConfig
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

	conf.CacheDB.BigCache.CleanTimeSec = defCacheDB_BigCache_CleanTimeSec

	conf.CacheDB.FreeCache.SizeMB = defCacheDB_FreeCache_SizeMB
	return conf
}

func (conf *Config) Check() error {
	if conf.ExpireSec < 1 {
		conf.ExpireSec = 0
	}

	switch v := strings.ToLower(conf.CacheDB.Type); v {
	case "":
		conf.CacheDB.Type = defCacheDB_Type
	case "no", "bigcache", "freecache", "redis":
	default:
		return fmt.Errorf("不支持的CacheDB: %v", v)
	}

	if conf.Compactor == "" {
		conf.Compactor = defCompactor
	}
	_, ok := compactor.TryGetCompactor(conf.Compactor)
	if !ok {
		return fmt.Errorf("不支持的Compactor: %v", conf.Compactor)
	}

	if conf.Serializer == "" {
		conf.Serializer = defSerializer
	}
	_, ok = serializer.TryGetSerializer(conf.Serializer)
	if !ok {
		return fmt.Errorf("不支持的Serializer: %v", conf.Serializer)
	}

	if conf.SingleFlight == "" {
		conf.SingleFlight = defSingleFlight
	}
	_, ok = single_flight.TryGetSingleFlight(conf.SingleFlight)
	if !ok {
		return fmt.Errorf("不支持的Serializer: %v", conf.SingleFlight)
	}

	if conf.CacheDB.BigCache.Shards < 1 {
		conf.CacheDB.BigCache.Shards = defCacheDB_BigCache_Shards
	}
	if conf.CacheDB.BigCache.MaxEntriesInWindow < 1 {
		conf.CacheDB.BigCache.MaxEntriesInWindow = defCacheDB_BigCache_MaxEntriesInWindow
	}
	if conf.CacheDB.BigCache.MaxEntrySize < 1 {
		conf.CacheDB.BigCache.MaxEntrySize = defCacheDB_BigCache_MaxEntrySize
	}
	if conf.CacheDB.BigCache.HardMaxCacheSize < 0 {
		conf.CacheDB.BigCache.HardMaxCacheSize = 0
	}

	if conf.CacheDB.FreeCache.SizeMB < 1 {
		conf.CacheDB.FreeCache.SizeMB = defCacheDB_FreeCache_SizeMB
	}
	return nil
}
