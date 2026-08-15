package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerAddr string
	MongoURI   string
	RedisURI   string
	Database   string

	JWTSecret               string
	JWTWebAccessTTL         time.Duration
	JWTWebRefreshTTL        time.Duration
	JWTMobileAccessTTL      time.Duration
	JWTMobileRefreshTTL     time.Duration
	CredentialEncryptionKey string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	GoogleClientID     string
	GoogleClientSecret string

	BaseURL string

	RateLimitRequests int
	RateLimitWindow   time.Duration

	MaxObjectSize  int64
	MultipartPartSize int64

	AuditRetentionDays int
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddr:          getenv("SERVER_ADDR", "127.0.0.1:8080"),
		MongoURI:            getenv("MONGODB_URI", "mongodb://127.0.0.1:27017"),
		RedisURI:            getenv("REDIS_URI", "redis://127.0.0.1:6379"),
		Database:            getenv("MONGODB_DATABASE", "bloberry"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTWebAccessTTL:     getdur("JWT_WEB_ACCESS_TTL", 48*time.Hour),
		JWTWebRefreshTTL:    getdur("JWT_WEB_REFRESH_TTL", 144*time.Hour),
		JWTMobileAccessTTL:  getdur("JWT_MOB_ACCESS_TTL", 720*time.Hour),
		JWTMobileRefreshTTL: getdur("JWT_MOB_REFRESH_TTL", 2160*time.Hour),
		CredentialEncryptionKey: os.Getenv("CREDENTIAL_ENCRYPTION_KEY"),
		SMTPHost:            os.Getenv("SMTP_HOST"),
		SMTPPort:            getint("SMTP_PORT", 587),
		SMTPUser:            os.Getenv("SMTP_USER"),
		SMTPPassword:        os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:            getenv("SMTP_FROM", "Bloberry <no-reply@bloberry.app>"),
		GoogleClientID:      os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		BaseURL:             getenv("BLOBERRY_BASE_URL", "http://localhost:8080"),
		RateLimitRequests:   getint("RATE_LIMIT_REQUESTS", 1000),
		RateLimitWindow:     getdur("RATE_LIMIT_WINDOW", time.Hour),
		MaxObjectSize:       getint64("MAX_OBJECT_SIZE", 5*1024*1024*1024),
		MultipartPartSize:   getint64("MULTIPART_PART_SIZE", 16*1024*1024),
		AuditRetentionDays:  getint("AUDIT_RETENTION_DAYS", 365),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.CredentialEncryptionKey == "" {
		return nil, fmt.Errorf("CREDENTIAL_ENCRYPTION_KEY is required")
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getint(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getint64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

func getdur(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(strings.TrimSpace(v)); err == nil {
			return d
		}
	}
	return def
}
