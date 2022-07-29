package no_sf

import (
	"github.com/zly-app/cache/core"
)

type NoSingleFlight struct{}

func (n NoSingleFlight) Do(cacheDB core.ICacheDB, key string, invoke func(cacheDB core.ICacheDB, key string) ([]byte, error)) ([]byte, error) {
	return invoke(cacheDB, key)
}

// 一个关闭并发查询控制的ISingleFlight
func NewNoSingleFlight() core.ISingleFlight {
	return NoSingleFlight{}
}
