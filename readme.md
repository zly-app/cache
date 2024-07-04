
# 透明读缓存工具

# 示例

从缓存加载, 如果加载失败自动从db加载并自动写入缓存, db加载时自动使用SingleFlight

```go
func main() {
	c, _ := cache.NewCache(cache.NewConfig())

	// 加载函数
	load := func(ctx context.Context, key string) (interface{}, error) { // db加载函数
		return "hello", nil
	}

	var a string
	_ = c.Get(context.Background(), "key", &a, // 获取数据
		cache.WithLoadFn(load),
	)

	print(a) // hello
}
```

# 多级缓存

首先从本地缓存加载, 如果加载失败从redis缓存加载并自动写入本地缓存, 如果仍然失败从db加载并自动写入redis缓存, 默认开启SingleFlight

```go
package main

import (
	"context"

	"github.com/zly-app/cache/v2"
)

func main() {
	localCache, _ := cache.NewCache(cache.NewConfig()) // 模拟本地缓存
	redisCache, _ := cache.NewCache(cache.NewConfig()) // 模拟redis缓存

	// redis缓存失效时从db加载
	loadFromDB := func(ctx context.Context, key string) (interface{}, error) {
		// 从db加载
		return "hello", nil
	}
	// 本地缓存失效时从redis加载
	loadFromRedisCache := func(ctx context.Context, key string) (interface{}, error) {
		var result string
		// 从redis加载
		err := redisCache.Get(ctx, key, &result,
			cache.WithLoadFn(loadFromDB), // 如果redis缓存未命中会从db加载
		)
		return result, err
	}

	var a string
	_ = localCache.Get(context.Background(), "key", &a,
		// 设置本地缓存缓存加载数据的方式
		cache.WithLoadFn(loadFromRedisCache),
	)

	print(a) // hello
}
```

# zapp 接入

```go
func main() {
	app := zapp.NewApp("test")
	defer app.Exit()

	creator := cache.NewCacheCreator(app) // 创建cache建造者

	cacheDef := creator.GetCache("default") // 通过cache建造者获取cache, 每次请求应该尽量重新获取cache

	var a string
	err := cacheDef.Get(context.Background(), "key", &a,
		cache.WithLoadFn(func(ctx context.Context, key string) (interface{}, error) {
			return "hello", nil
		}))
	app.Info(a, zap.Error(err))
}
```

## 添加配置文件 `configs/default.yml`. 更多配置说明参考[这里](./config.go)

```yaml
components:
  cache:
    default:
      Compactor: raw # 默认压缩器名, 可选 raw, zstd, gzip
      Serializer: sonic_std # 默认序列化器名, 可选 sonic, sonic_std, msgpack, jsoniter, jsoniter_standard, json, yaml
      SingleFlight: single # 默认单跑模块, 可选 no, single
      ExpireSec: 300 # 默认过期时间, 秒, < 1 表示永久
      IgnoreCacheFault: false # 是否忽略缓存数据库故障, 如果设为true, 在缓存数据库故障时从加载器获取数据, 这会导致缓存击穿. 如果设为false, 在缓存数据库故障时直接返回错误
      CacheDB:
        Type: bigcache # 缓存数据库类型, 支持 no, bigcache, freecache, redis
        BigCache: # 注意: bigcache 仅支持整体的过期时间, 不支持对单个key设置过期时间.
          Shards: 1024 # 分片数, 必须是2的幂
          CleanTimeSec: 60 # 清理周期秒数, 为 0 时不自动清理.
          MaxEntriesInWindow: 600000 # 初始化时申请允许储存的条目数的内存, 当实际使用量超过当前最大量时会触发内存重分配
          MaxEntrySize: 500 # 初始化时申请的每个条目的占用内存, 单位字节, 当实际使用量超过当前最大量时会触发内存重分配
          HardMaxCacheSize: 0 # 最大占用内存大小, 单位 mb, 0 表示不限制
          ExactExpire: false # 精准过期时间, 官方库的过期时间在 [Expire, Expire+CleanTimeSec] 区间. 如果设为true, 则过期时间精确为 Expire
        FreeCache: # memory 内存配置
          SizeMB: 1 # 分配内存大小, 单位mb, 单条数据大小不能超过该值的 1/1024
        Redis: # redis 内存配置
          Address: localhost:6379 # 地址: host1:port1,host2:port2
          UserName: "" # 用户名                     
          Password: "" # 密码
          DB: 0 # db, 只有非集群有效
          MinIdle: 2 # 最小闲置连接数
          MaxIdle: 4 # 最大闲置连接数
          PoolSize: 10 # 客户端池大小
          IdleTimeout: 3600 # 空闲链接超时时间, 单位秒, 如果一个连接长时间未使用将被视为连接无效, 小于1表示永不超时
          WaitTimeout: 5 # 等待获取连接的超时时间, 单位秒
          ConnectTimeout: 5 # 连接超时, 单位秒
          MaxConnLifetime: 3600 # 一个连接最大存活时间, 单位秒, 小于1表示不限制
          MaxRetries: 0 # 操作尝试次数, <1 表示不重试
          ReadTimeoutSec: 5 # 超时, 秒
          WriteTimeoutSec: 5 # 超时, 秒
```

# 支持的数据库

+ 支持任何数据库, 不关心用户如何加载数据

# 支持的缓存数据库

+ [no](./cachedb/no_cache/cache.go)
+ [bigcache](./cachedb/bigcache/cache.go)
+ [freecache](./cachedb/freecache/cache.go)
+ [redis](./cachedb/redis_cache/cache.go)

# 支持的序列化器

+ msgpack (推荐) . msgPack 序列化器
+ jsoniter_standard . jsonIter 实现的模拟内置 json 序列化器
+ jsoniter . jsonIter 序列化器
+ json
+ yaml (不推荐, 最慢)

# 支持的压缩器

+ raw . 不进行任何压缩
+ zstd
+ gzip

# 如何解决缓存击穿

+ 可以启用SingleFlight(默认开启), 当有多个进程同时获取一个相同的数据时, 只有一个进程会真的去加载函数读取数据, 其他的进程会等待该进程结束直接收到结果.

# 如何解决缓存雪崩

+ 为数据设置不同的过期时间甚至永不过期, 可以有效减小缓存雪崩的风险.
+ 预热数据

# 如何解决缓存穿透

+ 我们提供了一个占位符, 如果在loader结果中返回 `nil`, 我们会将它存入缓存并返回 `cache.ErrDataIsNil`, 当你再次获取它的时候会仍然会收到错误 `cache.ErrDataIsNil`
+ 在用户请求key的时候预判断它是否可能不存在, 比如判断id长度不等于16(不符合业务逻辑)的请求直接返回数据不存在错误
