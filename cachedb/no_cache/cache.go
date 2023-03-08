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

func (n noCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	return nil
}

func (n noCache) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (n noCache) Close() error {
	return nil
}
