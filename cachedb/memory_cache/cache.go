package memory_cache

import (
	"context"

	"github.com/coocood/freecache"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

// 最小内存大小
const minMemoryMB = 1

type memoryCache struct {
	cache *freecache.Cache
}

func (m *memoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := m.cache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		return nil, errs.CacheMiss
	}
	return data, err
}

func (m *memoryCache) MGet(ctx context.Context, keys ...string) map[string]core.CacheResult {
	result := make(map[string]core.CacheResult, len(keys))
	for _, key := range keys {
		data, err := m.Get(ctx, key)
		result[key] = core.CacheResult{Data: data, Err: err}
	}
	return result
}

func (m *memoryCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	return m.cache.Set([]byte(key), data, expireSec)
}

func (m *memoryCache) MSet(ctx context.Context, dataMap map[string][]byte, expireSec int) map[string]error {
	result := make(map[string]error, len(dataMap))
	for q, v := range dataMap {
		result[q] = m.Set(ctx, q, v, expireSec)
	}
	return result
}

func (m *memoryCache) Del(ctx context.Context, keys ...string) map[string]error {
	result := make(map[string]error, len(keys))
	for _, key := range keys {
		_ = m.cache.Del([]byte(key))
		result[key] = nil
	}
	return result
}

func (m *memoryCache) Close() error {
	m.cache.Clear()
	return nil
}

func NewMemoryCache(memoryMB int) core.ICacheDB {
	if memoryMB < minMemoryMB {
		memoryMB = minMemoryMB
	}
	cache := freecache.NewCache(memoryMB << 20)

	return &memoryCache{
		cache: cache,
	}
}
