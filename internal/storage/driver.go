package storage

import (
	"context"
	"io"
	"time"
)

// Capabilities declares what a driver can do (ADR-2). Stored on the
// storage_backends document; the conformance suite asserts against it.
type Capabilities struct {
	Presign          bool
	Multipart        bool
	MinPartSize      int64
	MaxPartCount     int
	StorageClasses   bool
	ServerSideCopy   bool
	RangeRequests    bool
	ObjectAttributes bool
}

type Range struct {
	Start int64
	End   int64 // inclusive; -1 means to end
}

type ObjectInfo struct {
	Size         int64
	ContentType  string
	ETag         string
	LastModified time.Time
}

type Part struct {
	Number int
	ETag   string
}

type PresignedURL struct {
	URL     string
	Headers map[string]string
	Method  string
}

// Driver is the interface the product exists to provide (PRD G1/G2). Designed
// against the hardest drivers: local disk (no external signer) and GCS
// (service-account signer) — not the easiest (S3).
type Driver interface {
	Capabilities() Capabilities

	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string, rng *Range) (io.ReadCloser, *ObjectInfo, error)
	Delete(ctx context.Context, keys []string) error
	Stat(ctx context.Context, key string) (*ObjectInfo, error)

	PresignGet(ctx context.Context, key string, ttl time.Duration) (*PresignedURL, error)
	PresignPut(ctx context.Context, key string, ttl time.Duration, size int64) (*PresignedURL, error)

	MultipartInit(ctx context.Context, key string, contentType string) (uploadID string, err error)
	MultipartPresignPart(ctx context.Context, key, uploadID string, part int, ttl time.Duration) (*PresignedURL, error)
	MultipartComplete(ctx context.Context, key, uploadID string, parts []Part) (*ObjectInfo, error)
	MultipartAbort(ctx context.Context, key, uploadID string) error

	HealthCheck(ctx context.Context) error
}

// Config is the per-driver credential/config payload stored encrypted at rest.
type Config map[string]interface{}
