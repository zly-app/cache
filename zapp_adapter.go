package cache

import (
	"fmt"

	"github.com/zly-app/zapp/component/conn"
	"github.com/zly-app/zapp/core"
)

// 默认组件类型
const defComponentType core.ComponentType = "cache"

type ICacheCreator interface {
	GetCache(name string) ICache
	Close()
}

type instance struct {
	cache ICache
}

func (i *instance) Close() {
	_ = i.cache.Close()
}

type cacheCreatorAdapter struct {
	app  core.IApp
	conn *conn.Conn
}

func (c *cacheCreatorAdapter) GetCache(name string) ICache {
	return c.conn.GetInstance(c.makeCache, name).(*instance).cache
}

func (c *cacheCreatorAdapter) Close() {
	c.conn.CloseAll()
}

func (c *cacheCreatorAdapter) makeCache(name string) (conn.IInstance, error) {
	conf := NewConfig()
	err := c.app.GetConfig().ParseComponentConfig(defComponentType, name, conf, true)
	if err != nil {
		return nil, fmt.Errorf("cache配置错误: %v", err)
	}

	cache, err := NewCache(conf)
	if err != nil {
		return nil, fmt.Errorf("cache创建失败: %v", err)
	}
	return &instance{cache: cache}, nil
}

func NewCacheCreator(app core.IApp) ICacheCreator {
	creator := &cacheCreatorAdapter{
		app:  app,
		conn: conn.NewConn(),
	}
	return creator
}
