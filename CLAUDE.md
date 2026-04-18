# Cache 项目 AI 开发指南

## 项目概述

这是一个 Go 语言的透明读缓存工具库，提供缓存读取、写入、删除功能，支持多种缓存后端、序列化器和压缩器。核心特性是"透明读"——当缓存未命中时自动从加载数据源获取数据并写入缓存，同时内置 SingleFlight 防止缓存击穿。

## 模块结构

```
cache/
├── cache.go          # Cache 结构体定义、创建(Cache构造)、序列化/反序列化内部方法
├── config.go         # Config 配置结构体及默认值、校验逻辑
├── options.go        # 单次操作的可选参数(Option 模式)，含 sync.Pool 复用
├── get.go            # Get 方法实现，含透明读逻辑：缓存未命中→SingleFlight→LoadFn→写缓存
├── set.go            # Set 方法实现
├── del.go            # Del 方法实现
├── sf_do.go          # SingleFlightDo 方法实现，强制从加载函数读取
├── err_cache.go      # 错误缓存实现(errCache)，创建失败时返回此实例，所有方法直接返回错误
├── link.go           # 导出的序列化器/压缩器/错误变量，类型别名 ICache/LoadFn
├── zapp_adapter.go   # zapp 框架适配器(ICacheCreator)，通过 zapp 配置自动创建缓存实例
├── zapp_default.go   # 全局默认缓存创建器，提供 GetCache/GetDefCache 快捷方法
├── core/             # 核心接口定义
│   ├── cache.go      # ICache 接口 (Get/Set/SingleFlightDo/Del/Close)、Option 类型、LoadFn 类型
│   ├── cachedb.go    # ICacheDB 接口 (Get/Set/Del/Close)
│   ├── compactor.go  # ICompactor 接口 (CompressBytes/UnCompressBytes)
│   ├── serializer.go # ISerializer 接口 (MarshalBytes/UnmarshalBytes)
│   └── single_flight.go # ISingleFlight 接口、LoadInvoke 类型
├── cachedb/          # 缓存后端实现
│   ├── bigcache/     # BigCache 本地内存缓存
│   ├── freecache/    # FreeCache 本地内存缓存
│   ├── redis_cache/  # Redis 远程缓存
│   └── no_cache/     # 空缓存(禁用缓存)
├── errs/             # 错误定义 (CacheMiss, DataIsNil)
└── single_flight/    # SingleFlight 实现
```

## 核心接口

### ICache (core/cache.go)

```go
type ICache interface {
    Get(ctx context.Context, key string, aPtr interface{}, opts ...Option) error
    Set(ctx context.Context, key string, data interface{}, opts ...Option) error
    SingleFlightDo(ctx context.Context, key string, aPtr interface{}, opts ...Option) error
    Del(ctx context.Context, keys ...string) error
    Close() error
}
```

### ICacheDB (core/cachedb.go)

所有缓存后端必须实现的接口，操作原始 `[]byte` 数据。

### ICompactor / ISerializer / ISingleFlight

分别定义压缩、序列化、单跑行为的接口。

## 核心流程

### Get 透明读流程 (get.go)

1. 若非 ForceLoad，先从 CacheDB 读取
2. 命中 → 解压缩 + 反序列化 → 返回
3. 未命中(CacheMiss) → 记录 TraceEvent
4. 缓存故障 → 检查 IgnoreCacheFault 决定是否继续
5. 无 LoadFn → 返回错误
6. 通过 SingleFlight 执行 LoadFn 加载数据
7. 加载成功 → 编码数据(序列化+压缩) → 写入缓存(DontWriteCache 除外) → 返回

### Set 流程 (set.go)

1. 序列化 + 压缩数据
2. 写入 CacheDB

### SingleFlightDo 流程 (sf_do.go)

1. 强制设置 ForceLoad=true（跳过缓存读取）
2. 必须提供 LoadFn，否则报错
3. 通过 SingleFlight 执行加载
4. 默认不写缓存（除非通过 WithForceLoad(false) 设置 DontWriteCache=false）

## Option 参数体系 (options.go)

所有 `Get`/`Set`/`SingleFlightDo` 方法都支持 Option 模式：

| Option 函数 | 说明 |
|---|---|
| `WithSerializer(s)` | 覆盖默认序列化器 |
| `WithCompactor(c)` | 覆盖默认压缩器 |
| `WithExpire(sec)` | 覆盖过期时间，<0 永不过期，0 使用默认值 |
| `WithLoadFn(fn)` | 设置数据加载函数 |
| `WithForceLoad(dontWriteCache)` | 忽略缓存强制从加载函数读取，dontWriteCache 控制是否回写缓存 |

options 使用 `sync.Pool` 复用以减少 GC 压力。

## 配置 (config.go)

```go
type Config struct {
    Compactor        string // 压缩器: raw, zstd, gzip (默认 raw)
    Serializer       string // 序列化器: sonic, sonic_std, msgpack, jsoniter, jsoniter_standard, json, yaml (默认 sonic_std)
    SingleFlight     string // 单跑模块: no, single (默认 single)
    ExpireSec        int    // 默认过期秒数, <1 表示永久 (默认 300)
    IgnoreCacheFault bool   // 忽略缓存数据库故障 (默认 false)
    CacheDB          struct {
        Type     string // 缓存类型: no, bigcache, freecache, redis (默认 bigcache)
        BigCache struct { ... }
        FreeCache struct { ... }
        RedisName string // 引用已有的 redis 组件名
        Redis     redis.RedisConfig
    }
}
```

## 错误处理

- `ErrCacheMiss` — 缓存未命中
- `ErrDataIsNil` — 加载函数返回 nil 数据（存入占位符，再次读取仍返回此错误）
- `errCache` — 创建失败时的兜底实现，所有方法直接返回创建错误

## zapp 框架集成

- `zapp_adapter.go`: `ICacheCreator` 接口，`NewCacheCreator()` 创建适配器
- `zapp_default.go`: 全局单例 `defCreator`，提供 `GetCache(name)` / `GetDefCache()` 快捷方法
- zapp 应用通过配置文件 `components.cache.<name>` 自动注入配置

## 缓存穿透/击穿/雪崩防护

- **击穿**: SingleFlight 合并并发请求
- **穿透**: LoadFn 返回 nil 时存入占位符，返回 `ErrDataIsNil`
- **雪崩**: 设置不同过期时间、数据预热

## 开发注意事项

1. 新增缓存后端需实现 `core.ICacheDB` 接口
2. 新增序列化器/压缩器需分别在 zapp 的 serializer/compactor 包注册
3. `options` 对象通过 `sync.Pool` 管理，使用后必须调用 `putOptions` 归还
4. `Get`/`Set`/`Del`/`SingleFlightDo` 均通过 zapp filter 链处理，支持链路追踪
5. `NewCache` 的 name 参数在直接使用时仅用于标识，在 zapp 模式下用于匹配配置
6. BigCache 仅支持全局过期时间，不支持单个 key 独立过期（除非开启 ExactExpire）
