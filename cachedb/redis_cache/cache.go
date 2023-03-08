package redis_cache

import (
	"context"
	"time"

	"github.com/zly-app/component/redis"

	"github.com/zly-app/cache/core"
	"github.com/zly-app/cache/errs"
)

type redisCache struct {
	client redis.UniversalClient
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	s, err := r.client.Get(ctx, key).Result()
	if err == nil {
		return []byte(s), nil
	}
	if err == redis.Nil {
		return nil, errs.CacheMiss
	}
	return nil, err
}

func (r *redisCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	var ex time.Duration
	if expireSec > 0 {
		ex = time.Duration(expireSec) * time.Second
	}

	return r.client.Set(ctx, key, data, ex).Err()
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err == redis.Nil { // 虽然不会出现 redis.Nil
		return nil
	}
	return err
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

func NewRedisCache(client redis.UniversalClient) core.ICacheDB {
	return &redisCache{client: client}
}
