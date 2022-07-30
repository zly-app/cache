package cache

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"

	"github.com/zly-app/cache/errs"
)

func makeMemoryCache() ICache {
	cache, err := NewCache(NewConfig())
	if err != nil {
		panic(fmt.Errorf("创建Cache失败: %v", err))
	}
	return cache
}

func makeRedisCache() ICache {
	conf := NewConfig()
	conf.CacheDB.Type = "redis"
	conf.CacheDB.Redis.Address = "127.0.0.1:6379"
	cache, err := NewCache(conf)
	if err != nil {
		panic(fmt.Errorf("创建Cache失败: %v", err))
	}
	return cache
}

func TestMemoryCache(t *testing.T) {
	t.Run("testSetGet", func(t *testing.T) { testSetGet(t, makeMemoryCache()) })
	t.Run("testSetGetSlice", func(t *testing.T) { testSetGetSlice(t, makeMemoryCache()) })
	t.Run("testDel", func(t *testing.T) { testDel(t, makeMemoryCache()) })
	t.Run("testExpire", func(t *testing.T) { testExpire(t, makeMemoryCache()) })
	t.Run("testDefaultExpire", func(t *testing.T) {
		conf := NewConfig()
		conf.ExpireSec = 1
		cache, err := NewCache(conf)
		if err != nil {
			panic(fmt.Errorf("创建Cache失败: %v", err))
		}
		testDefaultExpire(t, cache)
	})
	t.Run("testLoadFn", func(t *testing.T) { testLoadFn(t, makeMemoryCache()) })
	t.Run("testMSet", func(t *testing.T) { testMSet(t, makeMemoryCache()) })
	t.Run("testMGet", func(t *testing.T) { testMGet(t, makeMemoryCache()) })
	t.Run("testMGetSlice", func(t *testing.T) { testMGetSlice(t, makeMemoryCache()) })
	t.Run("testClose", func(t *testing.T) { testClose(t, makeMemoryCache()) })
}

func TestRedisCache(t *testing.T) {
	t.Run("testSetGet", func(t *testing.T) { testSetGet(t, makeRedisCache()) })
	t.Run("testSetGetSlice", func(t *testing.T) { testSetGetSlice(t, makeRedisCache()) })
	t.Run("testDel", func(t *testing.T) { testDel(t, makeRedisCache()) })
	t.Run("testExpire", func(t *testing.T) { testExpire(t, makeRedisCache()) })
	t.Run("testLoadFn", func(t *testing.T) { testLoadFn(t, makeRedisCache()) })
	t.Run("testMSet", func(t *testing.T) { testMSet(t, makeRedisCache()) })
	t.Run("testMGet", func(t *testing.T) { testMGet(t, makeRedisCache()) })
	t.Run("testMGetSlice", func(t *testing.T) { testMGetSlice(t, makeRedisCache()) })
	t.Run("testClose", func(t *testing.T) { testClose(t, makeRedisCache()) })
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

	time.Sleep(time.Second)

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

	time.Sleep(time.Second)

	var c int
	err = cache.Get(context.Background(), key, &c)
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
	}))
	require.Nil(t, err)
	require.Equal(t, true, load)
	require.Equal(t, a, b)
}
func testMSet(t *testing.T, cache ICache) {
	const key1 = "testMSet1"
	const key2 = "testMSet2"

	var a = map[string]interface{}{
		key1: 1,
		key2: 2,
	}
	err := cache.MSet(context.Background(), a)
	require.Nil(t, err)

	var b int
	err = cache.Get(context.Background(), key1, &b)
	require.Nil(t, err)
	require.Equal(t, a[key1], b)

	var c int
	err = cache.Get(context.Background(), key2, &c)
	require.Nil(t, err)
	require.Equal(t, a[key2], c)
}
func testMGet(t *testing.T, cache ICache) {
	const key1 = "testMGet1"
	const key2 = "testMGet2"
	const key3 = "testMGet3"

	var a = 1
	err := cache.Set(context.Background(), key1, a)
	require.Nil(t, err)

	var b1, b2, b3 int
	b := map[string]interface{}{
		key1: &b1,
		key2: &b2,
		key3: &b3,
	}
	err = cache.MGet(context.Background(), b)
	require.NotNil(t, err)
	require.Equal(t, 1, b1)
	require.Equal(t, nil, GetKeyError(err, key1))
	require.Equal(t, 0, b2)
	require.Equal(t, ErrCacheMiss, GetKeyError(err, key2))
	require.Equal(t, 0, b3)
	require.Equal(t, ErrCacheMiss, GetKeyError(err, key3))

	var c1, c2, c3 int
	var loadC2, loadC3, loadOther bool
	c := map[string]interface{}{
		key1: &c1,
		key2: &c2,
		key3: &c3,
	}
	var loadFn = func(ctx context.Context, key string) (interface{}, error) {
		switch key {
		case key2:
			loadC2 = true
			return 2, nil
		case key3:
			loadC3 = true
			return nil, nil
		}
		loadOther = true
		return nil, fmt.Errorf("意外的加载key: %v", err)
	}
	err = cache.MGet(context.Background(), c, WithLoadFn(loadFn))
	require.NotNil(t, err)
	require.Equal(t, false, loadOther)
	require.Equal(t, true, loadC2)
	require.Equal(t, true, loadC3)
	require.Equal(t, 1, c1)
	require.Equal(t, nil, GetKeyError(err, key1))
	require.Equal(t, 2, c2)
	require.Equal(t, nil, GetKeyError(err, key2))
	require.Equal(t, 0, c3)
	require.Equal(t, ErrDataIsNil, GetKeyError(err, key3))
}
func testMGetSlice(t *testing.T, cache ICache) {
	const key1 = "testMGetSlice1"
	const key2 = "testMGetSlice2"
	const key3 = "testMGetSlice3"

	var a = map[string]interface{}{
		key1: 1,
		key2: 2,
	}
	err := cache.MSet(context.Background(), a)
	require.Nil(t, err)

	var b []int
	err = cache.MGetSlice(context.Background(), []string{key1, key2}, &b)
	require.Nil(t, err)
	require.Equal(t, []int{1, 2}, b)

	var c []int
	err = cache.MGetSlice(context.Background(), []string{key1, key2, key3}, &c)
	require.NotNil(t, err)
	require.Equal(t, []int{1, 2}, c)
	require.Equal(t, nil, GetKeyError(err, key1))
	require.Equal(t, nil, GetKeyError(err, key2))
	require.Equal(t, ErrCacheMiss, GetKeyError(err, key3))
}
func testClose(t *testing.T, cache ICache) {
	const key = "testClose"

	err := cache.Close()
	require.Nil(t, err)

	var b []byte
	err = cache.Get(context.Background(), key, &b)
	require.NotNil(t, err)
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
	conf.CacheDB.Memory.SizeMB = sizeMB
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
