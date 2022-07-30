package core

import (
	"context"
)

type CacheResult struct {
	Data []byte
	Err  error
}

// 缓存数据库接口
type ICacheDB interface {
	// 获取一个值
	Get(ctx context.Context, key string) ([]byte, error)
	// 批量获取
	MGet(ctx context.Context, keys ...string) map[string]CacheResult

	// 设置一个值, expireSec <= 0 时表示永不过期
	Set(ctx context.Context, key string, data []byte, expireSec int) error
	// 批量设置, expireSec <= 0 时表示永不过期
	MSet(ctx context.Context, data map[string][]byte, expireSec int) map[string]error

	// 删除数据
	Del(ctx context.Context, keys ...string) error

	// 关闭
	Close() error
}
