package cache

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/errs"
)

func makeBigCache() ICache {
	conf := NewConfig()
	conf.CacheDB.Type = "bigcache"
	cache, err := NewCache(conf)
	if err != nil {
		panic(fmt.Errorf("创建Cache失败: %v", err))
	}
	return cache
}

func makeFreeCache() ICache {
	conf := NewConfig()
	conf.CacheDB.Type = "freecache"
	cache, err := NewCache(conf)
	if err != nil {
		panic(fmt.Errorf("创建Cache失败: %v", err))
	}
	return cache
}

func makeRedisCache() ICache {
	conf := NewConfig()
	conf.CacheDB.Type = "redis"
	conf.CacheDB.Redis.Address = "localhost:6379"
	cache, err := NewCache(conf)
	if err != nil {
		panic(fmt.Errorf("创建Cache失败: %v", err))
	}
	return cache
}

func TestBigCache(t *testing.T) {
	t.Run("testSetGet", func(t *testing.T) { testSetGet(t, makeBigCache()) })
	t.Run("testSetGetSlice", func(t *testing.T) { testSetGetSlice(t, makeBigCache()) })
	t.Run("testDel", func(t *testing.T) { testDel(t, makeBigCache()) })
	t.Run("testExpire", func(t *testing.T) { testExpire(t, makeBigCache()) })
	t.Run("testDefaultExpire", func(t *testing.T) {
		conf := NewConfig()
		conf.ExpireSec = 1
		conf.CacheDB.Type = "bigcache"
		cache, err := NewCache(conf)
		if err != nil {
			panic(fmt.Errorf("创建Cache失败: %v", err))
		}
		testDefaultExpire(t, cache)
	})
	t.Run("testLoadFn", func(t *testing.T) { testLoadFn(t, makeBigCache()) })
	t.Run("testClose", func(t *testing.T) { testClose(t, makeBigCache()) })
	t.Run("testForceLoad", func(t *testing.T) { testForceLoad(t, makeBigCache()) })
	t.Run("testSF", func(t *testing.T) { testSF(t, makeBigCache()) })
}

func TestFreeCache(t *testing.T) {
	t.Run("testSetGet", func(t *testing.T) { testSetGet(t, makeFreeCache()) })
	t.Run("testSetGetSlice", func(t *testing.T) { testSetGetSlice(t, makeFreeCache()) })
	t.Run("testDel", func(t *testing.T) { testDel(t, makeFreeCache()) })
	t.Run("testExpire", func(t *testing.T) { testExpire(t, makeFreeCache()) })
	t.Run("testDefaultExpire", func(t *testing.T) {
		conf := NewConfig()
		conf.ExpireSec = 1
		conf.CacheDB.Type = "freecache"
		cache, err := NewCache(conf)
		if err != nil {
			panic(fmt.Errorf("创建Cache失败: %v", err))
		}
		testDefaultExpire(t, cache)
	})
	t.Run("testLoadFn", func(t *testing.T) { testLoadFn(t, makeFreeCache()) })
	t.Run("testClose", func(t *testing.T) { testClose(t, makeFreeCache()) })
	t.Run("testForceLoad", func(t *testing.T) { testForceLoad(t, makeFreeCache()) })
	t.Run("testSF", func(t *testing.T) { testSF(t, makeFreeCache()) })
}

func TestRedisCache(t *testing.T) {
	t.Run("testSetGet", func(t *testing.T) { testSetGet(t, makeRedisCache()) })
	t.Run("testSetGetSlice", func(t *testing.T) { testSetGetSlice(t, makeRedisCache()) })
	t.Run("testDel", func(t *testing.T) { testDel(t, makeRedisCache()) })
	t.Run("testExpire", func(t *testing.T) { testExpire(t, makeRedisCache()) })
	t.Run("testLoadFn", func(t *testing.T) { testLoadFn(t, makeRedisCache()) })
	t.Run("testClose", func(t *testing.T) { testClose(t, makeRedisCache()) })
	t.Run("testForceLoad", func(t *testing.T) { testForceLoad(t, makeRedisCache()) })
	t.Run("testSF", func(t *testing.T) { testSF(t, makeRedisCache()) })
}

func testSetGet(t *testing.T, cache ICache) {
	const key = "testSetGet"

	var a = 3
	err := cache.Set(context.Background(), key, a)
	require.Nil(t, err)

	var b int
	err = cache.Get(context.Background(), key, &b)
	require.Nil(t, err)
	require.Equal(t, a, b)
}
func testSetGetSlice(t *testing.T, cache ICache) {
	const key = "testSetGetSlice"

	type A struct {
		A int
	}

	var a = []A{
		{1},
		{2},
		{3},
	}
	err := cache.Set(context.Background(), key, a)
	require.Nil(t, err)

	var b []A
	err = cache.Get(context.Background(), key, &b)
	require.Nil(t, err)
	require.Equal(t, a, b)
}
func testDel(t *testing.T, cache ICache) {
	const key = "testDel"

	var a = 3
	err := cache.Set(context.Background(), key, a)
	require.Nil(t, err)

	err = cache.Del(context.Background(), key)
	require.Nil(t, err)

	var b int
	err = cache.Get(context.Background(), key, &b)
	require.Equal(t, errs.CacheMiss, err)
}
func testExpire(t *testing.T, cache ICache) {
	const key = "testExpire"

	var a = 3
	err := cache.Set(context.Background(), key, a, WithExpire(1))
	require.Nil(t, err)

	var b int
	err = cache.Get(context.Background(), key, &b)
	require.Nil(t, err)
	require.Equal(t, a, b)

	time.Sleep(time.Second * 2)

	var c int
	err = cache.Get(context.Background(), key, &c)
	require.Equal(t, errs.CacheMiss, err)
}
func testDefaultExpire(t *testing.T, cache ICache) {
	const key = "testDefaultExpire"

	var a = 3
	err := cache.Set(context.Background(), key, a)
	require.Nil(t, err)

	var b int
	err = cache.Get(context.Background(), key, &b)
	require.Nil(t, err)
	require.Equal(t, a, b)

	time.Sleep(time.Second * 3)

	var c int
	err = cache.Get(context.Background(), key, &c)
	t.Log(c)
	require.Equal(t, errs.CacheMiss, err)
}
func testLoadFn(t *testing.T, cache ICache) {
	const key = "testLoadFn"

	var a = 3
	var load bool

	var b int
	err := cache.Get(context.Background(), key, &b, WithLoadFn(func(ctx context.Context, key string) (interface{}, error) {
		load = true
		return a, nil
	}), WithExpire(3))
	require.Nil(t, err)
	require.Equal(t, true, load)
	require.Equal(t, a, b)
}
func testClose(t *testing.T, cache ICache) {
	const key = "testClose"

	err := cache.Close()
	require.Nil(t, err)

	var b []byte
	err = cache.Get(context.Background(), key, &b)
	require.NotNil(t, err)
}
func testForceLoad(t *testing.T, cache ICache) {
	const key = "testForceLoad"

	var a = 3
	var a2 = 4
	err := cache.Set(context.Background(), key, a)
	require.Nil(t, err)

	var b int
	var load bool
	err = cache.Get(context.Background(), key, &b,
		WithForceLoad(true),
		WithLoadFn(func(ctx context.Context, key string) (interface{}, error) {
			load = true
			return a2, nil
		}))
	require.Nil(t, err)
	require.Equal(t, true, load)
	require.Equal(t, a2, b)

	var c int
	err = cache.Get(context.Background(), key, &c)
	require.Nil(t, err)
	require.Equal(t, a, c)
}
func testSF(t *testing.T, cache ICache) {
	const key = "testSF"

	var a = 3
	var loadB, loadC bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cache.SingleFlightDo(context.Background(), key, WithLoadFn(func(ctx context.Context, key string) (interface{}, error) {
			loadB = true
			return a, nil
		}))
		require.Nil(t, err)
	}()

	wg.Wait()
	var b int
	err := cache.Get(context.Background(), key, &b, WithLoadFn(func(ctx context.Context, key string) (interface{}, error) {
		loadC = true
		return a, nil
	}), WithExpire(3))
	require.Nil(t, err)
	require.Equal(t, a, b)

	require.Equal(t, true, loadB)
	require.Equal(t, true, loadC)
}

func BenchmarkGet(b *testing.B) {
	keyCount := []struct {
		name   string
		count  int
		sizeMB int
	}{
		{"1k", 1000, 10},
		//{"10k", 10000, 100},
		//{"30k", 100000, 300},
	}
	compactors := []struct {
		name string
		compactor.ICompactor
	}{
		{"NoCompactor", NoCompactor},
		//{"ZStdCompactor", ZStdCompactor},
		//{"GzipCompactor", GzipCompactor},
	}
	serializers := []struct {
		name string
		serializer.ISerializer
	}{
		{"MsgPackSerializer", MsgPackSerializer},
		//{"JsonIterStandardSerializer", JsonIterStandardSerializer},
		//{"JsonSerializer", JsonSerializer},
		//{"JsonIterSerializer", JsonIterSerializer},
		//{"YamlSerializer", YamlSerializer},
	}
	for _, k := range keyCount {
		for _, c := range compactors {
			for _, s := range serializers {
				name := fmt.Sprintf("%v_%v_%v", k.name, s.name, c.name)
				b.Run(name, func(b *testing.B) {
					benchGet(b, k.count, k.sizeMB, s, c)
				})
			}
		}
	}
}

func benchGet(b *testing.B, maxKeyCount, sizeMB int, serializer serializer.ISerializer, compactor compactor.ICompactor) {
	rand.Seed(time.Now().UnixNano())
	const dataLen = 512

	conf := NewConfig()
	conf.CacheDB.FreeCache.SizeMB = sizeMB
	cache, err := NewCache(conf)
	require.Nil(b, err)

	expects := make([][]byte, maxKeyCount)
	for i := 0; i < maxKeyCount; i++ {
		bs := make([]byte, dataLen)
		for j := 0; j < dataLen; j++ {
			bs[j] = byte(rand.Int() & 255)
		}
		expects[i] = bs

		key := strconv.Itoa(i)
		err := cache.Set(context.Background(), key, bs, WithSerializer(serializer), WithCompactor(compactor))
		require.NoError(b, err, "数据设置失败")
	}

	// 缓存随机key
	randKeys := make([]int, 1<<20)
	for i := 0; i < len(randKeys); i++ {
		randKeys[i] = rand.Int() % maxKeyCount
	}

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			i++
			key := randKeys[i&(len(randKeys)-1)]
			expect := expects[key]

			var bs []byte
			err := cache.Get(context.Background(), strconv.Itoa(key), &bs, WithSerializer(serializer), WithCompactor(compactor))
			if err != nil {
				b.Fatalf("数据加载失败: key: %v, err %v", key, err)
			}
			if len(bs) != dataLen {
				b.Fatalf("数据长度不一致, key: %v, need %v, got %v", key, dataLen, len(bs))
			}
			if !bytes.Equal(bs, expect) {
				b.Fatalf("数据不一致: key: %v, need %v, got %v", key, expect, bs)
			}
		}
	})
}
