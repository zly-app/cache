package cache

import (
	"fmt"

	"github.com/zly-app/component/redis"
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
	// Close 关闭全部 cache
	Close()
}

type cacheCreatorAdapter struct {
	app          core.IApp
	redisCreator redis.Creator
	conn         *conn.AnyConn[ICache]
}

func (c *cacheCreatorAdapter) GetCache(name string) ICache {
	ins, err := c.conn.GetConn(c.makeCache, name)
	if err != nil {
		return newErrCache(err)
	}
	return ins
}

func (c *cacheCreatorAdapter) GetDefCache() ICache {
	return c.GetCache(consts.DefaultComponentName)
}

func (c *cacheCreatorAdapter) Close() {
	c.conn.CloseAll()
}

func (c *cacheCreatorAdapter) makeCache(name string) (ICache, error) {
	conf := NewConfig()
	if err := c.app.GetConfig().ParseComponentConfig(defComponentType, name, conf, true); err != nil {
		return nil, fmt.Errorf("cache配置错误: %v", err)
	}

	cache, err := NewCache(name, conf, c.redisCreator, c.app)
	if err != nil {
		return nil, fmt.Errorf("cache创建失败: %v", err)
	}
	return cache, nil
}

// NewCacheCreator 创建 cache 集合。app 与 redis creator 由组合根注入，不再有包级单例。
func NewCacheCreator(app core.IApp, redisCreator redis.Creator) ICacheCreator {
	return &cacheCreatorAdapter{
		app:          app,
		redisCreator: redisCreator,
		conn: conn.NewAnyConn[ICache](func(name string, conn ICache) {
			_ = conn.Close()
		}),
	}
}
