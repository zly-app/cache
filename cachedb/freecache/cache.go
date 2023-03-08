package freecache

import (
	"context"

	"github.com/coocood/freecache"

	"github.com/zly-app/cache/v2/core"
	"github.com/zly-app/cache/v2/errs"
)

// 最小内存大小
const minMemoryMB = 1

type freeCache struct {
	cache *freecache.Cache
}

func (m *freeCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := m.cache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		return nil, errs.CacheMiss
	}
	return data, err
}

func (m *freeCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	return m.cache.Set([]byte(key), data, expireSec)
}

func (m *freeCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		_ = m.cache.Del([]byte(key))
	}
	return nil
}

func (m *freeCache) Close() error {
	m.cache.Clear()
	return nil
}

// memoryMB 分配内存大小, 单位mb, 单条数据大小不能超过该值的 1/1024
func NewMemoryCache(memoryMB int) core.ICacheDB {
	if memoryMB < minMemoryMB {
		memoryMB = minMemoryMB
	}
	cache := freecache.NewCache(memoryMB << 20)

	return &freeCache{
		cache: cache,
	}
}
