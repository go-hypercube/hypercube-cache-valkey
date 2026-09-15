package valkeycache

import (
	"context"
	"time"

	"github.com/go-hypercube/go-hypercube/cache"
	"github.com/valkey-io/valkey-go"
)

// valkeyCache implements cache.Cache on top of a valkey.Client,
// mirroring hypercube-cache-redis so a driver repository only needs to
// depend on the cache package, not the rest of the framework.
type valkeyCache struct {
	client valkey.Client
}

// New wraps an already-constructed valkey.Client as a cache.Cache
// implementation. The caller owns the client's life-cycle (including
// calling client.Close() on shutdown).
func New(client valkey.Client) cache.Cache {
	return &valkeyCache{client: client}
}

func (c *valkeyCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Do(ctx, c.client.B().Get().Key(key).Build()).ToString()
	if valkey.IsValkeyNil(err) {
		return "", cache.ErrNotFound
	}
	return value, err
}

func (c *valkeyCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if ttl < 0 {
		return cache.ErrInvalidTTL
	}

	cmd := c.client.B().Set().Key(key).Value(value)
	if ttl > 0 {
		return c.client.Do(ctx, cmd.Ex(ttl).Build()).Error()
	}
	return c.client.Do(ctx, cmd.Build()).Error()
}

func (c *valkeyCache) Delete(ctx context.Context, key string) error {
	return c.client.Do(ctx, c.client.B().Del().Key(key).Build()).Error()
}

func (c *valkeyCache) Has(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Do(ctx, c.client.B().Exists().Key(key).Build()).ToInt64()
	return n > 0, err
}

func (c *valkeyCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return c.client.Do(ctx, c.client.B().Incrby().Key(key).Increment(delta).Build()).ToInt64()
}

func (c *valkeyCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if ttl < 0 {
		return cache.ErrInvalidTTL
	}
	if ttl == 0 {
		return c.client.Do(ctx, c.client.B().Persist().Key(key).Build()).Error()
	}
	seconds := int64(ttl / time.Second)
	return c.client.Do(ctx, c.client.B().Expire().Key(key).Seconds(seconds).Build()).Error()
}
