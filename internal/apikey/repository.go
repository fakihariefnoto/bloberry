package apikey

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists applications and access keys.
type Repository interface {
	InsertApplication(ctx context.Context, a *domain.Application) error
	GetApplication(ctx context.Context, tenantID, id string) (*domain.Application, error)
	ListApplications(ctx context.Context, tenantID string) ([]domain.Application, error)
	DeleteApplication(ctx context.Context, tenantID, id string) error

	InsertKey(ctx context.Context, k *domain.AccessKey) error
	GetKeyByHash(ctx context.Context, hash string) (*domain.AccessKey, error)
	ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error)
	// ListAllKeys returns every key in the tenant (all applications), for the
	// SDK-facing API keys page.
	ListAllKeys(ctx context.Context, tenantID string) ([]domain.AccessKey, error)
	RevokeKey(ctx context.Context, tenantID, id string) error
	// RevokeKeyAny revokes a key by id regardless of application.
	RevokeKeyAny(ctx context.Context, tenantID, id string) error
	TouchKey(ctx context.Context, id string) error
}

// Invalidator clears the Redis principal cache on key revoke (ADR-6).
type Invalidator interface {
	InvalidateKey(ctx context.Context, hash string) error
}

type CreatedKey struct {
	KeyID      string   `json:"key_id"`
	Secret     string   `json:"secret"` // shown exactly once (PRD D5)
	Prefix     string   `json:"prefix"`
	LastFour   string   `json:"last_four"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type Usecase interface {
	Register(ctx context.Context, tenantID, name, description string) (*domain.Application, error)
	Get(ctx context.Context, tenantID, id string) (*domain.Application, error)
	List(ctx context.Context, tenantID string) ([]domain.Application, error)
	Delete(ctx context.Context, tenantID, id string) error
	CreateKey(ctx context.Context, tenantID, applicationID string, scope []string, perms []string, expiresAt *time.Time) (*CreatedKey, error)
	ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error)
	// ListAllKeys returns keys across all applications, enriched with the app name.
	ListAllKeys(ctx context.Context, tenantID string) ([]KeyWithApp, error)
	RevokeKey(ctx context.Context, tenantID, keyID string) error
	// RevokeKeyAny revokes a key by id regardless of application.
	RevokeKeyAny(ctx context.Context, tenantID, keyID string) error
}

// KeyWithApp is an access key joined with its application name for the
// cross-app API keys page.
type KeyWithApp struct {
	domain.AccessKey
	ApplicationName string `json:"application_name"`
}
