package cache

import (
	"fmt"

	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/component/conn"
	"github.com/zly-app/zapp/consts"
	"github.com/zly-app/zapp/core"
)

// 默认组件类型
const defComponentType core.ComponentType = "cache"

type ICacheCreator interface {
	// 获取cache, 每次请求应该尽量重新调用这个方法
	GetCache(name string) ICache
	// 获取cache, 每次请求应该尽量重新调用这个方法
	GetDefCache() ICache
}

type instance struct {
	cache ICache
}

func (i *instance) Close() {
	_ = i.cache.Close()
}

type cacheCreatorAdapter struct {
	conn *conn.Conn
}

func (c *cacheCreatorAdapter) GetCache(name string) ICache {
	ins, err := c.conn.GetConn(c.makeCache, name)
	if err != nil {
		return newErrCache(err)
	}
	return ins.(*instance).cache
}

func (c *cacheCreatorAdapter) GetDefCache() ICache {
	return c.GetCache(consts.DefaultComponentName)
}

func (c *cacheCreatorAdapter) Close() {
	c.conn.CloseAll()
}

func (c *cacheCreatorAdapter) makeCache(name string) (conn.IInstance, error) {
	conf := NewConfig()
	err := zapp.App().GetConfig().ParseComponentConfig(defComponentType, name, conf, true)
	if err != nil {
		return nil, fmt.Errorf("cache配置错误: %v", err)
	}

	cache, err := NewCache(name, conf)
	if err != nil {
		return nil, fmt.Errorf("cache创建失败: %v", err)
	}
	return &instance{cache: cache}, nil
}

// deprecated: use GetCacheCreator
func NewCacheCreator(_ core.IApp) ICacheCreator {
	return defCreator
}

func GetCacheCreator() ICacheCreator {
	return defCreator
}
