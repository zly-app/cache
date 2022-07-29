package core

// 单跑接口
type ISingleFlight interface {
	// 执行, 当缓存数据库不存在时, 在执行loader加载数据前, 会调用此方法
	Do(cacheDB ICacheDB, key string, invoke func(cacheDB ICacheDB, key string) ([]byte, error)) ([]byte, error)
}
