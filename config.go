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

	defCacheDB_Type = "bigcache"

	defCacheDB_BigCache_Shards             = 1024
	defCacheDB_BigCache_CleanTimeMs        = 1
	defCacheDB_BigCache_MaxEntriesInWindow = 1000 * 10 * 60
	defCacheDB_BigCache_MaxEntrySize       = 500

	defCacheDB_FreeCache_SizeMB = 1

	defCacheDB_Redis_Address         = "127.0.0.1:6379"
	defCacheDB_Redis_DB              = 0
	defCacheDB_Redis_IsCluster       = false
	defCacheDB_Redis_MinIdleConns    = 2
	defCacheDB_Redis_PoolSize        = 5
	defCacheDB_Redis_ReadTimeoutSec  = 5
	defCacheDB_Redis_WriteTimeoutSec = 5
	defCacheDB_Redis_DialTimeoutSec  = 5
)

type Config struct {
	Compactor        string // 默认压缩器名, 可选 raw, zstd, gzip
	Serializer       string // 默认序列化器名, 可选 msgpack, jsoniter_standard, jsoniter, json, yaml
	SingleFlight     string // 默认单跑模块, 可选 no, single
	ExpireSec        int    // 默认过期时间, 秒, <= 0 表示永久
	IgnoreCacheFault bool   // 是否忽略缓存数据库故障, 如果设为true, 在缓存数据库故障时从加载器获取数据, 这会导致缓存击穿. 如果设为false, 在缓存数据库故障时直接返回错误
	CacheDB          struct {
		Type     string // 缓存数据库类型, 支持 no, bigcache, freecache, redis
		BigCache struct {
			Shards             int // 分片数, 必须是2的幂
			CleanTimeMs        int // 清理周期秒数, 为 0 时不自动清理
			MaxEntriesInWindow int // 初始化时申请允许储存的条目数的内存, 当实际使用量超过当前最大量时会触发内存重分配
			MaxEntrySize       int // 初始化时申请的每个条目的占用内存, 单位字节, 当实际使用量超过当前最大量时会触发内存重分配
		}
		FreeCache struct {
			SizeMB int // 分配内存大小, 单位mb, 单条数据大小不能超过该值的 1/1024
		}
		Redis struct {
			Address         string // 地址: host1:port1,host2:port2
			UserName        string // 用户名
			Password        string // 密码
			DB              int    // db, 只有非集群有效
			IsCluster       bool   // 是否为集群
			MinIdleConns    int    // 最小空闲连接数
			PoolSize        int    // 客户端池大小
			ReadTimeoutSec  int    // 读取超时, 单位秒
			WriteTimeoutSec int    // 写入超时, 单位秒
			DialTimeoutSec  int    // 连接超时, 单位秒
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

	conf.CacheDB.BigCache.CleanTimeMs = defCacheDB_BigCache_CleanTimeMs

	conf.CacheDB.FreeCache.SizeMB = defCacheDB_FreeCache_SizeMB

	conf.CacheDB.Redis.Address = defCacheDB_Redis_Address
	conf.CacheDB.Redis.DB = defCacheDB_Redis_DB
	conf.CacheDB.Redis.IsCluster = defCacheDB_Redis_IsCluster
	conf.CacheDB.Redis.MinIdleConns = defCacheDB_Redis_MinIdleConns
	conf.CacheDB.Redis.PoolSize = defCacheDB_Redis_PoolSize
	conf.CacheDB.Redis.ReadTimeoutSec = defCacheDB_Redis_ReadTimeoutSec
	conf.CacheDB.Redis.WriteTimeoutSec = defCacheDB_Redis_WriteTimeoutSec
	conf.CacheDB.Redis.DialTimeoutSec = defCacheDB_Redis_DialTimeoutSec
	return conf
}

func (conf *Config) Check() error {
	switch v := strings.ToLower(conf.CacheDB.Type); v {
	case "":
		conf.CacheDB.Type = defCacheDB_Type
	case "bigcache","freecache", "redis":
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

	if conf.CacheDB.BigCache.Shards < 1 {
		conf.CacheDB.BigCache.Shards = defCacheDB_BigCache_Shards
	}
	if conf.CacheDB.BigCache.MaxEntriesInWindow < 1 {
		conf.CacheDB.BigCache.MaxEntriesInWindow = defCacheDB_BigCache_MaxEntriesInWindow
	}
	if conf.CacheDB.BigCache.MaxEntrySize < 1 {
		conf.CacheDB.BigCache.MaxEntrySize = defCacheDB_BigCache_MaxEntrySize
	}

	if conf.CacheDB.FreeCache.SizeMB < 1 {
		conf.CacheDB.FreeCache.SizeMB = defCacheDB_FreeCache_SizeMB
	}

	if conf.CacheDB.Redis.Address == "" {
		conf.CacheDB.Redis.Address = defCacheDB_Redis_Address
	}
	if conf.CacheDB.Redis.ReadTimeoutSec < 1 {
		conf.CacheDB.Redis.ReadTimeoutSec = defCacheDB_Redis_ReadTimeoutSec
	}
	if conf.CacheDB.Redis.WriteTimeoutSec < 1 {
		conf.CacheDB.Redis.WriteTimeoutSec = defCacheDB_Redis_WriteTimeoutSec
	}
	if conf.CacheDB.Redis.DialTimeoutSec < 1 {
		conf.CacheDB.Redis.DialTimeoutSec = defCacheDB_Redis_DialTimeoutSec
	}
	return nil
}
