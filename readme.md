
# 透明读缓存工具

# 示例

```go
func main() {
	c, _ := cache.NewCache(cache.NewConfig())

	var a string
	_ = c.Get(context.Background(), "key", &a, // 获取数据
		cache.WithLoadFn(func(ctx context.Context, key string) (interface{}, error) { // db加载函数
			return "hello", nil
		}),
	)

	print(a) // hello
}
```

# 支持的数据库

+ 支持任何数据库, 不关心用户如何加载数据

# 支持的缓存数据库

+ [no-cache](./cachedb/no_cache/cache.go)
+ [memory-cache](./cachedb/memory_cache/cache.go)

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

+ 可以在启用SingleFlight, 当有多个进程同时获取一个相同的数据时, 只有一个进程会真的去加载函数读取数据, 其他的进程会等待该进程结束直接收到结果.

# 如何解决缓存雪崩

+ 为数据设置不同的过期时间甚至永不过期, 可以有效减小缓存雪崩的风险.
+ 预热数据

# 如何解决缓存穿透

+ 我们提供了一个占位符, 如果在loader结果中返回 `nil`, 我们会将它存入缓存, 当你在获取它的时候会收到错误 `cache.ErrDataIsNil`
+ 在用户请求key的时候预判断它是否可能不存在, 比如判断id长度不等于32(uuid去掉横杠的长度)的请求直接返回数据不存在错误
