package cache

import (
	"context"

	"github.com/zly-app/cache/v2/core"
)

type errCache struct {
	err error
}

func (e errCache) Get(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	return e.err
}

func (e errCache) Set(ctx context.Context, key string, data interface{}, opts ...core.Option) error {
	return e.err
}

func (e errCache) SingleFlightDo(ctx context.Context, key string, aPtr interface{}, opts ...core.Option) error {
	return e.err
}

func (e errCache) Del(ctx context.Context, keys ...string) error {
	return e.err
}

func (e errCache) Close() error {
	return e.err
}

func newErrCache(err error) ICache {
	return errCache{err}
}
