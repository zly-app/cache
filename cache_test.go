package cache

import (
	"testing"

	"github.com/zly-app/cache/core"
)

func makeMemoryCache(t *testing.T) core.ICache {
	cache, err := NewCache(NewConfig())
	if err != nil {
		t.Fatalf("创建Cache失败: %v", err)
	}
	return cache
}
