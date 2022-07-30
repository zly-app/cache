package redis_cache

import (
	"context"
	"fmt"
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

func (r *redisCache) MGet(ctx context.Context, keys ...string) map[string]core.CacheResult {
	result := make(map[string]core.CacheResult)

	cacheResult, err := r.client.MGet(ctx, keys...).Result()
	if err != nil && err != redis.Nil {
		for _, key := range keys {
			result[key] = core.CacheResult{Err: err}
		}
		return result
	}

	for i, key := range keys {
		switch v := cacheResult[i].(type) {
		case nil:
			result[key] = core.CacheResult{Err: errs.CacheMiss}
		case string:
			result[key] = core.CacheResult{Data: []byte(v)}
		case []byte:
			result[key] = core.CacheResult{Data: v}
		default: // 虽然不会出现
			result[key] = core.CacheResult{Err: fmt.Errorf("不能识别的redis结果类型 <%T>", v)}
		}
	}
	return result
}

func (r *redisCache) Set(ctx context.Context, key string, data []byte, expireSec int) error {
	if expireSec < 1 {
		expireSec = -1
	}

	return r.client.SetEX(ctx, key, data, time.Duration(expireSec)*time.Second).Err()
}

func (r *redisCache) MSet(ctx context.Context, data map[string][]byte, expireSec int) map[string]error {
	result := make(map[string]error)

	if expireSec < 1 {
		args := make([]interface{}, 0, len(data)*2)
		for k, v := range data {
			args = append(args, k, v)
		}

		err := r.client.MSet(ctx, args...).Err()
		for key := range data {
			result[key] = err
		}
		return result
	}

	for k, v := range data {
		result[k] = r.client.SetEX(ctx, k, v, time.Duration(expireSec)*time.Second).Err()
	}
	return result
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
