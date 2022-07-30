package no_cache

import (
	"context"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

var _ core.ICacheDB = (*noCache)(nil)

type noCache struct{}

// 创建一个不会缓存数据的ICacheDB
func NoCache() core.ICacheDB { return noCache{} }

func (n noCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, errs.CacheMiss
}

func (n noCache) MGet(ctx context.Context, keys ...string) map[string]*core.CacheResult {
	result := make(map[string]*core.CacheResult, len(keys))
	for _, key := range keys {
		result[key] = &core.CacheResult{Err: errs.CacheMiss}
	}
	return result
}

func (n noCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	return nil
}

func (n noCache) MSet(ctx context.Context, dataMap map[string][]byte, expireSec int) map[string]error {
	result := make(map[string]error, len(dataMap))
	for key := range dataMap {
		result[key] = nil
	}
	return result
}

func (n noCache) Del(ctx context.Context, keys ...string) map[string]error {
	result := make(map[string]error, len(keys))
	for _, key := range keys {
		result[key] = nil
	}
	return result
}

func (n noCache) Close() error {
	return nil
}
