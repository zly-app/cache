package single_sf

import (
	"context"
	"errors"
	"hash/fnv"
	"sync"

	"github.com/zly-app/cache/core"
)

const (
	ShardCount uint32 = 1 << 8 // 分片数
)

type waitResult struct {
	wg sync.WaitGroup
	v  []byte
	e  error
}

type SingleFlight struct {
	mxs           []*sync.RWMutex
	waits         []map[string]*waitResult
	shardAndOPVal uint32 // 按位与操作值
}

// 创建一个单跑, 分片数必须大于0且为2的幂
func NewSingleFlight(shardCount ...uint32) core.ISingleFlight {
	count := ShardCount
	if len(shardCount) > 0 && shardCount[0] > 0 {
		count = shardCount[0]
		if count&(count-1) != 0 {
			panic(errors.New("shardCount must power of 2"))
		}
	}

	mxs := make([]*sync.RWMutex, count)
	mms := make([]map[string]*waitResult, count)
	for i := uint32(0); i < count; i++ {
		mxs[i] = new(sync.RWMutex)
		mms[i] = make(map[string]*waitResult)
	}
	return &SingleFlight{
		mxs:           mxs,
		waits:         mms,
		shardAndOPVal: count - 1,
	}
}

func (m *SingleFlight) getShard(key string) (*sync.RWMutex, map[string]*waitResult) {
	f := fnv.New32a()
	_, _ = f.Write([]byte(key))
	n := f.Sum32()
	shard := n & m.shardAndOPVal
	return m.mxs[shard], m.waits[shard]
}

func (m *SingleFlight) Do(ctx context.Context, key string, invoke core.LoadInvoke) ([]byte, error) {
	mx, wait := m.getShard(key)

	mx.RLock()
	result, ok := wait[key]
	mx.RUnlock()

	// 已经有线程在查询
	if ok {
		result.wg.Wait()
		return result.v, result.e
	}

	mx.Lock()

	// 再检查一下, 因为在拿到锁之前可能被别的线程占了位置
	result, ok = wait[key]
	if ok {
		mx.Unlock()
		result.wg.Wait()
		return result.v, result.e
	}

	// 占位置
	result = new(waitResult)
	result.wg.Add(1)
	wait[key] = result
	mx.Unlock()

	// 执行db加载
	result.v, result.e = invoke(ctx, key)
	result.wg.Done()

	// 删除位置
	mx.Lock()
	delete(wait, key)
	mx.Unlock()

	return result.v, result.e
}
