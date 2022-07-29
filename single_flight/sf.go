package single_flight

import (
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/cache/core"
	no_sf "github.com/zly-app/cache/single_flight/no-sf"
	single_sf "github.com/zly-app/cache/single_flight/single-sf"
)

var sfs = map[string]core.ISingleFlight{
	"no":     no_sf.NewNoSingleFlight(),
	"single": single_sf.NewSingleFlight(),
}

// 注册, 重复注册会panic
func RegistrySingleFlight(name string, sf core.ISingleFlight) {
	if _, ok := sfs[name]; ok {
		logger.Log.Panic("SingleFlight重复注册", zap.String("name", name))
	}
	sfs[name] = sf
}

// 获取, 不存在会panic
func GetSingleFlight(name string) core.ISingleFlight {
	c, ok := sfs[name]
	if !ok {
		logger.Log.Panic("未定义的SingleFlightName", zap.String("name", name))
	}
	return c
}

// 尝试获取
func TryGetSingleFlight(name string) (core.ISingleFlight, bool) {
	c, ok := sfs[name]
	return c, ok
}
