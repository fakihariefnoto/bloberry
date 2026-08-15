package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	goredis "github.com/redis/go-redis/v9"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type SessionData struct {
	UserID   string `json:"user_id"`
	Platform string `json:"platform"`
	IssuedAt int64  `json:"issued_at"`
}

// Store owns refresh tokens in Redis (domains.md §7).
type Store struct {
	rdb *goredis.Client
}

func NewStore(rdb *goredis.Client) *Store { return &Store{rdb: rdb} }

func (s *Store) Create(ctx context.Context, userID, platform string, ttlSeconds int64) (string, error) {
	return s.create(ctx, userID, platform, time.Duration(ttlSeconds)*time.Second)
}

func (s *Store) create(ctx context.Context, userID, platform string, ttl time.Duration) (string, error) {
	sid := crypto.NewToken(24)
	data, _ := json.Marshal(SessionData{UserID: userID, Platform: platform, IssuedAt: time.Now().Unix()})
	if err := s.rdb.Set(ctx, "refresh:"+sid, data, ttl).Err(); err != nil {
		return "", err
	}
	return sid, nil
}

func (s *Store) Get(ctx context.Context, sid string) (string, string, error) {
	val, err := s.rdb.Get(ctx, "refresh:"+sid).Result()
	if errors.Is(err, goredis.Nil) {
		return "", "", errors.New("session not found")
	}
	if err != nil {
		return "", "", err
	}
	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return "", "", err
	}
	return data.UserID, data.Platform, nil
}

// RevokeAll deletes every refresh session for a user (password reset,
// logout-all). Iterates the keyspace; acceptable at v1 scale.
func (s *Store) RevokeAll(ctx context.Context, userID string) error {
	iter := s.rdb.Scan(ctx, 0, "refresh:*", 1000).Iterator()
	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		var data SessionData
		if json.Unmarshal([]byte(val), &data) == nil && data.UserID == userID {
			_ = s.rdb.Del(ctx, iter.Val()).Err()
		}
	}
	return iter.Err()
}

// Rotate deletes the old session and issues a new one (rotation, not renewal —
// domains.md §4.3).
func (s *Store) Rotate(ctx context.Context, sid, userID, platform string, ttlSeconds int64) (string, error) {
	if err := s.rdb.Del(ctx, "refresh:"+sid).Err(); err != nil {
		return "", err
	}
	return s.create(ctx, userID, platform, time.Duration(ttlSeconds)*time.Second)
}

func (s *Store) Revoke(ctx context.Context, sid string) error {
	return s.rdb.Del(ctx, "refresh:"+sid).Err()
}
