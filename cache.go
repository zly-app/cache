package cache

import (
	"fmt"
	"strings"

	"github.com/zly-app/component/redis"

	"github.com/zly-app/cache/v2/cachedb/bigcache"
	"github.com/zly-app/cache/v2/cachedb/freecache"
	"github.com/zly-app/cache/v2/cachedb/redis_cache"
	"github.com/zly-app/cache/v2/core"
	"github.com/zly-app/cache/v2/single_flight"
)

type Cache struct {
	cacheName        string
	cacheDB          core.ICacheDB
	compactor        core.ICompactor
	serializer       core.ISerializer
	sf               core.ISingleFlight // 单跑模块
	expireSec        int                // 默认过期时间
	ignoreCacheFault bool               // 是否忽略缓存数据库故障
}

func (c *Cache) Close() error {
	return c.cacheDB.Close()
}

func (c *Cache) marshalQuery(data interface{}, serializer core.ISerializer, compactor core.ICompactor) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	rawData, err := serializer.MarshalBytes(data)
	if err != nil {
		return nil, fmt.Errorf("序列化失败: %v", err)
	}

	comData, err := compactor.CompressBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf("压缩失败: %v", err)
	}
	return comData, nil
}

func (c *Cache) unmarshalQuery(comData []byte, aPtr interface{}, serializer core.ISerializer, compactor core.ICompactor) error {
	if len(comData) == 0 {
		return ErrDataIsNil
	}

	rawData, err := compactor.UnCompressBytes(comData)
	if err != nil {
		return fmt.Errorf("解压缩失败: %v", err)
	}

	err = serializer.UnmarshalBytes(rawData, aPtr)
	if err != nil {
		return fmt.Errorf("反序列化失败: %v", err)
	}
	return nil
}

func NewCache(name string, conf *Config) (ICache, error) {
	err := conf.Check()
	if err != nil {
		return nil, fmt.Errorf("cache配置检查失败: %v", err)
	}

	cache := &Cache{
		cacheName:        name,
		expireSec:        conf.ExpireSec,
		ignoreCacheFault: conf.IgnoreCacheFault,
	}

	switch v := strings.ToLower(conf.CacheDB.Type); v {
	case "bigcache":
		cache.cacheDB, err = bigcache.NewCache(
			conf.CacheDB.BigCache.Shards,
			conf.ExpireSec,
			conf.CacheDB.BigCache.CleanTimeSec,
			conf.CacheDB.BigCache.MaxEntriesInWindow,
			conf.CacheDB.BigCache.MaxEntrySize,
			conf.CacheDB.BigCache.HardMaxCacheSize,
			conf.CacheDB.BigCache.ExactExpire,
		)
		if err != nil {
			return nil, fmt.Errorf("创建bigcache失败: %v", err)
		}
	case "freecache":
		cache.cacheDB = freecache.NewCache(conf.CacheDB.FreeCache.SizeMB)
	case "redis":
		var redisClient redis.UniversalClient
		if conf.CacheDB.RedisName != "" {
			redisClient = redis.GetClient(conf.CacheDB.RedisName)
		} else {
			redisClient, err = redis.NewClient(&conf.CacheDB.Redis, "cache")
		}
		if err != nil {
			return nil, fmt.Errorf("创建redis客户端失败: %v", err)
		}
		cache.cacheDB = redis_cache.NewRedisCache(redisClient)
	}

	cache.compactor = GetCompactor(strings.ToLower(conf.Compactor))
	cache.serializer = GetSerializer(strings.ToLower(conf.Serializer))
	cache.sf = single_flight.GetSingleFlight(strings.ToLower(conf.SingleFlight))

	return cache, nil
}
