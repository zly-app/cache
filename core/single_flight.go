package core

import (
	"context"
)

type LoadInvoke func(ctx context.Context, key string) ([]byte, error)

// 单跑接口
type ISingleFlight interface {
	// 执行, 当缓存数据库不存在时, 在执行loader加载数据前, 会调用此方法
	Do(ctx context.Context, key string, invoke LoadInvoke) ([]byte, error)
}
