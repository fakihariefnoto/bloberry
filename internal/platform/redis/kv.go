package redis

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var errKeyNotFound = errors.New("redis: key not found")

// KV is a small Redis key-value adapter implementing the auth domain's
// RedisStore interface.
type KV struct {
	c *goredis.Client
}

func NewKV(c *goredis.Client) *KV { return &KV{c: c} }

func (k *KV) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return k.c.Set(ctx, key, value, ttl).Err()
}

func (k *KV) Get(ctx context.Context, key string) (string, error) {
	val, err := k.c.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", errKeyNotFound
	}
	return val, err
}

func (k *KV) Del(ctx context.Context, key string) error {
	return k.c.Del(ctx, key).Err()
}

func (k *KV) Incr(ctx context.Context, key string) (int64, error) {
	return k.c.Incr(ctx, key).Result()
}

func (k *KV) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return k.c.Expire(ctx, key, ttl).Err()
}
