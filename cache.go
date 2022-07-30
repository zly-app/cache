package cache

import (
	"bytes"
	"context"
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

func (c *Cache) Del(ctx context.Context, keys ...string) map[string]error {
	//TODO implement me
	panic("implement me")
}

func (c *Cache) Close() error {
	//TODO implement me
	panic("implement me")
}

func (c *Cache) marshalQuery(aPtr interface{}, serializer serializer.ISerializer, compactor compactor.ICompactor) ([]byte, error) {
	if aPtr == nil {
		return nil, nil
	}

	var rawData bytes.Buffer
	err := serializer.Marshal(aPtr, &rawData)
	if err != nil {
		return nil, fmt.Errorf("序列化失败: %v", err)
	}

	var comData bytes.Buffer
	err = compactor.Compress(&rawData, &comData)
	if err != nil {
		return nil, fmt.Errorf("压缩失败: %v", err)
	}
	return comData.Bytes(), nil
}

func (c *Cache) unmarshalQuery(data []byte, aPtr interface{}, serializer serializer.ISerializer, compactor compactor.ICompactor) error {
	if len(data) == 0 {
		return errs.DataIsNil
	}

	var rawData bytes.Buffer
	err := compactor.UnCompress(bytes.NewReader(data), &rawData)
	if err != nil {
		return fmt.Errorf("解压缩失败: %v", err)
	}

	err = serializer.Unmarshal(&rawData, aPtr)
	if err != nil {
		return fmt.Errorf("反序列化失败: %v", err)
	}
	return nil
}

func NewCache(conf *Config) (core.ICache, error) {
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
		memory_cache.NewMemoryCache(conf.CacheDB.Memory.SizeMB)
	case "redis":
		return nil, fmt.Errorf("暂未实现redis缓存数据库")
	}

	cache.compactor = compactor.GetCompactor(strings.ToLower(conf.Compactor))
	cache.serializer = serializer.GetSerializer(strings.ToLower(conf.Serializer))
	cache.sf = single_flight.GetSingleFlight(strings.ToLower(conf.SingleFlight))

	return cache, nil
}
