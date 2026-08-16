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
	// GetTenant returns a tenant by id (for key name enrichment).
	GetTenant(ctx context.Context, id string) (*domain.Tenant, error)
	ListApplications(ctx context.Context, tenantID string) ([]domain.Application, error)
	DeleteApplication(ctx context.Context, tenantID, id string) error

	InsertKey(ctx context.Context, k *domain.AccessKey) error
	GetKeyByHash(ctx context.Context, hash string) (*domain.AccessKey, error)
	ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error)
	// ListAllKeys returns every key in the tenant (all applications), for the
	// SDK-facing API keys page.
	ListAllKeys(ctx context.Context, tenantID string) ([]domain.AccessKey, error)
	// ListKeysForAdmin returns keys across ALL tenants (platform admin view).
	ListKeysForAdmin(ctx context.Context) ([]domain.AccessKey, error)
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
	CreateKey(ctx context.Context, tenantID, applicationID, name string, scope []string, perms []string, expiresAt *time.Time) (*CreatedKey, error)
	// CreateTenantKey creates a tenant-scoped key not tied to an application.
	CreateTenantKey(ctx context.Context, tenantID, name string, scope []string, perms []string, expiresAt *time.Time) (*CreatedKey, error)
	ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error)
	// ListKeysPage returns keys for the API keys page (all tenants for admin,
	// own tenant otherwise).
	ListKeysPage(ctx context.Context, tenantID string, isPlatformAdmin bool) ([]KeyWithApp, error)
	RevokeKey(ctx context.Context, tenantID, keyID string) error
	// RevokeKeyAny revokes a key by id regardless of application.
	RevokeKeyAny(ctx context.Context, tenantID, keyID string) error
}

// KeyWithApp is an access key joined with its application/tenant names for the
// API keys page.
type KeyWithApp struct {
	domain.AccessKey
	ApplicationName string `json:"application_name"`
	TenantName      string `json:"tenant_name"`
}
