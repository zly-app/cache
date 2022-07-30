package cache

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/zly-app/component/redis"
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/cachedb/memory_cache"
	"github.com/zly-app/cache/cachedb/redis_cache"
	"github.com/zly-app/cache/core"
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

func (c *Cache) unmarshalMQuerySlice(keys []string, cacheResults map[string]core.CacheResult, slicePtr interface{},
	serializer serializer.ISerializer, compactor compactor.ICompactor) map[string]error {
	// 类型检查
	sliceType := reflect.TypeOf(slicePtr)
	if sliceType.Kind() != reflect.Ptr {
		panic(errors.New("slicePtr必须是带指针的切片"))
	}
	sliceType = sliceType.Elem()
	if sliceType.Kind() != reflect.Slice {
		panic(errors.New("slicePtr必须是带指针的切片"))
	}

	// 值检查
	sliceValue := reflect.ValueOf(slicePtr).Elem()
	if sliceValue.Len() != 0 {
		panic(errors.New("slicePtr的长度必须为0"))
	}
	if sliceValue.Kind() == reflect.Invalid {
		panic(errors.New("slicePtr中子类型无法访问"))
	}

	// 获取子类型
	itemType := sliceType.Elem()                // 获取子类型
	itemIsPtr := itemType.Kind() == reflect.Ptr // 检查子类型是否为指针
	if itemIsPtr {
		itemType = itemType.Elem() // 获取指针指向的真正的子类型
	}

	// 数据处理
	result := make(map[string]error, len(cacheResults))
	items := make([]reflect.Value, 0, len(cacheResults))
	for _, key := range keys {
		cacheResult, ok := cacheResults[key]
		if !ok || cacheResult.Err != nil { // 只处理正常的数据
			continue
		}

		child := reflect.New(itemType) // 创建一个相同类型的指针
		err := c.unmarshalQuery(cacheResult.Data, child.Interface(), serializer, compactor)
		if err == nil {
			if !itemIsPtr {
				child = child.Elem() // 如果想要的不是指针那么获取它的内容
			}
			items = append(items, child)
		}
		result[key] = err
	}

	values := reflect.Append(sliceValue, items...) // 构建内容切片
	sliceValue.Set(values)                         // 将内容切片写入原始切片中
	return result
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
		redisClient, err := redis.MakeRedisClient(&redis.RedisConfig{
			Address:         conf.CacheDB.Redis.Address,
			UserName:        conf.CacheDB.Redis.UserName,
			Password:        conf.CacheDB.Redis.Password,
			DB:              conf.CacheDB.Redis.DB,
			IsCluster:       conf.CacheDB.Redis.IsCluster,
			MinIdleConns:    conf.CacheDB.Redis.MinIdleConns,
			PoolSize:        conf.CacheDB.Redis.PoolSize,
			ReadTimeout:     conf.CacheDB.Redis.ReadTimeoutSec * 1000,
			WriteTimeout:    conf.CacheDB.Redis.WriteTimeoutSec * 1000,
			DialTimeout:     conf.CacheDB.Redis.DialTimeoutSec * 1000,
			EnableOpenTrace: false,
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
