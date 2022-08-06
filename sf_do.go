package cache

import (
	"context"

	"github.com/zly-app/cache/core"
)

func (c *Cache) SingleFlightDo(ctx context.Context, key string, invoke core.LoadInvoke) ([]byte, error) {
	return c.sf.Do(ctx, key, invoke)
}
