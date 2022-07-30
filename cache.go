package cache

import (
	"fmt"
	"strings"

	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/cachedb/memory_cache"
	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
	"github.com/zly-app/cache/single_flight"
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
	//TODO implement me
	panic("implement me")
}

func (c *Cache) marshalQuery(aPtr interface{}, serializer serializer.ISerializer, compactor compactor.ICompactor) ([]byte, error) {
	if aPtr == nil {
		return nil, nil
	}

	rawData, err := serializer.MarshalBytes(aPtr)
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
		return errs.DataIsNil
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
	case "memory":
		cache.cacheDB = memory_cache.NewMemoryCache(conf.CacheDB.Memory.SizeMB)
	case "redis":
		return nil, fmt.Errorf("暂未实现redis缓存数据库")
	}

	cache.compactor = compactor.GetCompactor(strings.ToLower(conf.Compactor))
	cache.serializer = serializer.GetSerializer(strings.ToLower(conf.Serializer))
	cache.sf = single_flight.GetSingleFlight(strings.ToLower(conf.SingleFlight))

	return cache, nil
}
