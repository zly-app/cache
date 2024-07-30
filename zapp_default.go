package cache

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/component/conn"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/handler"
)

var defCreator = &cacheCreatorAdapter{
	conn: conn.NewConn(),
}

func init() {
	zapp.AddHandler(zapp.AfterCloseComponent, func(_ core.IApp, _ handler.HandlerType) {
		defCreator.Close()
	})
}

func GetCache(name string) ICache {
	return defCreator.GetCache(name)
}

func GetDefCache() ICache {
	return defCreator.GetDefCache()
}
