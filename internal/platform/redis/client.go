package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, uri string) (*redis.Client, error) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		opts = &redis.Options{Addr: uri}
	}
	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

// Key helpers — the schema the rest of the plan assumes (domains.md §7).

func RefreshKey(sessionID string) string    { return "refresh:" + sessionID }
func PairKey(token string) string           { return "pair:" + token }
func TotpPendingKey(userID string) string   { return "totp:pending:" + userID }
func OtpLoginKey(email string) string       { return "otp:login:" + email }
func OtpAttemptsKey(email string) string    { return "otp:attempts:" + email }
func ResetKey(tokenHash string) string      { return "reset:" + tokenHash }
func PrincipalKey(ptype, id string) string  { return "principal:" + ptype + ":" + id }
func APIKeyKey(secretHash string) string    { return "apikey:" + secretHash }
func SlugKey(slug string) string            { return "slug:" + slug }
func RateKey(key string) string             { return "rate:" + key }

// TokenBucket is a simple Redis token-bucket rate limiter keyed per access key.
type TokenBucket struct {
	client *redis.Client
}

func NewTokenBucket(client *redis.Client) *TokenBucket {
	return &TokenBucket{client: client}
}

// Allow returns true if the request passes; otherwise the Retry-After seconds.
func (tb *TokenBucket) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, int64) {
	rk := RateKey(key)
	// INCR + PEXPIRE after: if a TTL was set previously, keep it (the bucket
	// replenishes on the window boundary, which is the standard fixed-window
	// approximation the plan's PRD-Q5 asks for).
	script := redis.NewScript(`
local c = redis.call('INCR', KEYS[1])
if c == 1 then
  redis.call('PEXPIRE', KEYS[1], ARGV[1])
end
return c
`)
	n, err := script.Run(ctx, tb.client, []string{rk}, window.Milliseconds()).Int64()
	if err != nil {
		return true, 0
	}
	if n > int64(limit) {
		ttl := tb.client.PTTL(ctx, rk).Val()
		if ttl < 0 {
			ttl = time.Duration(window.Milliseconds()) * time.Millisecond
		}
		return false, int64(ttl) / int64(time.Second)
	}
	return true, 0
}
