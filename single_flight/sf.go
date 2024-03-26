package single_flight

import (
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/cache/v2/core"
	no_sf "github.com/zly-app/cache/v2/single_flight/no-sf"
	single_sf "github.com/zly-app/cache/v2/single_flight/single-sf"
)

type SingleFlightCreator = func() core.ISingleFlight

var sfs = map[string]SingleFlightCreator{
	"no": func() core.ISingleFlight {
		return no_sf.NewNoSingleFlight()
	},
	"single": func() core.ISingleFlight {
		return single_sf.NewSingleFlight()
	},
}

// 注册, 重复注册会panic
func RegistrySingleFlightCreator(name string, creator SingleFlightCreator) {
	if _, ok := sfs[name]; ok {
		logger.Log.Panic("SingleFlight建造者重复注册", zap.String("name", name))
	}
	sfs[name] = creator
}

// 获取, 不存在会panic
func GetSingleFlight(name string) core.ISingleFlight {
	creator, ok := sfs[name]
	if !ok {
		logger.Log.Panic("未定义的SingleFlightName", zap.String("name", name))
	}
	return creator()
}

// 获取
func TryGetSingleFlight(name string) (core.ISingleFlight, bool) {
	creator, ok := sfs[name]
	return creator(), ok
}
