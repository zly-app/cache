package core

import (
	"context"
)

// 缓存数据库接口
type ICacheDB interface {
	// 获取一个值
	Get(ctx context.Context, key string) ([]byte, error)

	// 设置一个值, expireSec <= 0 时表示永不过期
	Set(ctx context.Context, key string, data []byte, expireSec int) error

	// 删除数据
	Del(ctx context.Context, keys ...string) error

	// 关闭
	Close() error
}
