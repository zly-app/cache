package bigcache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

type bigCache struct {
	cache *bigcache.BigCache
}

func (m *bigCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, info, err := m.cache.GetWithInfo(key)
	if err == bigcache.ErrEntryNotFound || info.EntryStatus != 0 {
		return nil, errs.CacheMiss
	}
	return data, err
}

func (m *bigCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	return m.cache.Set(key, data)
}

func (m *bigCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		_ = m.cache.Delete(key)
	}
	return nil
}

func (m *bigCache) Close() error {
	return m.cache.Close()
}

func NewCache(shards, expireSec, cleanTimeMs, maxEntriesInWindow, maxEntrySize, hardMaxCacheSize int) (core.ICacheDB, error) {
	if expireSec <= 0 {
		expireSec = 31536000000 // 1000年
	}
	conf := bigcache.Config{
		Shards:             shards,
		LifeWindow:         time.Duration(expireSec) * time.Second,
		CleanWindow:        time.Duration(cleanTimeMs) * time.Second,
		MaxEntriesInWindow: maxEntriesInWindow,
		MaxEntrySize:       maxEntrySize,
		HardMaxCacheSize:   hardMaxCacheSize,
	}
	cache, err := bigcache.New(context.Background(), conf)
	return &bigCache{
		cache: cache,
	}, err
}
