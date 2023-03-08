package cache

import (
	"fmt"
	"strings"

	"github.com/zly-app/component/redis"
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/v2/cachedb/freecache"
	"github.com/zly-app/cache/v2/cachedb/redis_cache"
	"github.com/zly-app/cache/v2/core"
	"github.com/zly-app/cache/v2/single_flight"
)

type Cache struct {
	cacheDB          core.ICacheDB
	compactor        compactor.ICompactor
	serializer       serializer.ISerializer
	sf               core.ISingleFlight // 单跑模块
	expireSec        int                // 默认过期时间
	ignoreCacheFault bool               // 是否忽略缓存数据库故障
}

func (c *Cache) Close() error {
	return c.cacheDB.Close()
}

func (c *Cache) marshalQuery(data interface{}, serializer serializer.ISerializer, compactor compactor.ICompactor) ([]byte, error) {
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

func (c *Cache) unmarshalQuery(comData []byte, aPtr interface{}, serializer serializer.ISerializer, compactor compactor.ICompactor) error {
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

func NewCache(conf *Config) (ICache, error) {
	err := conf.Check()
	if err != nil {
		return nil, fmt.Errorf("cache配置检查失败: %v", err)
	}

	cache := &Cache{
		expireSec:        conf.ExpireSec,
		ignoreCacheFault: conf.IgnoreCacheFault,
	}

	switch v := strings.ToLower(conf.CacheDB.Type); v {
	case "freecache":
		cache.cacheDB = freecache.NewMemoryCache(conf.CacheDB.FreeCache.SizeMB)
	case "redis":
		redisClient, err := redis.NewClient(&redis.RedisConfig{
			Address:         conf.CacheDB.Redis.Address,
			UserName:        conf.CacheDB.Redis.UserName,
			Password:        conf.CacheDB.Redis.Password,
			DB:              conf.CacheDB.Redis.DB,
			IsCluster:       conf.CacheDB.Redis.IsCluster,
			MinIdleConns:    conf.CacheDB.Redis.MinIdleConns,
			PoolSize:        conf.CacheDB.Redis.PoolSize,
			ReadTimeoutSec:  conf.CacheDB.Redis.ReadTimeoutSec,
			WriteTimeoutSec: conf.CacheDB.Redis.WriteTimeoutSec,
			DialTimeoutSec:  conf.CacheDB.Redis.DialTimeoutSec,
		})
		if err != nil {
			return nil, fmt.Errorf("创建redis客户端失败: %v", err)
		}
		cache.cacheDB = redis_cache.NewRedisCache(redisClient)
	}

	cache.compactor = compactor.GetCompactor(strings.ToLower(conf.Compactor))
	cache.serializer = serializer.GetSerializer(strings.ToLower(conf.Serializer))
	cache.sf = single_flight.GetSingleFlight(strings.ToLower(conf.SingleFlight))

	return cache, nil
}
