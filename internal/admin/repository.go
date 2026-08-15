package admin

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/storage"
	"github.com/fakihariefnoto/bloberry/internal/storage/registry"
)

// Repository persists storage backends (install-level + BYO).
type Repository interface {
	InsertBackend(ctx context.Context, b *domain.StorageBackend) error
	GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error)
	ListBackends(ctx context.Context) ([]domain.StorageBackend, error)
	// ListForTenant returns install-level backends plus those owned by the tenant.
	ListForTenant(ctx context.Context, tenantID string) ([]domain.StorageBackend, error)
	UpdateBackend(ctx context.Context, b *domain.StorageBackend) error
	DeleteBackend(ctx context.Context, id string) error
	CountTenantsOnBackend(ctx context.Context, backendID string) (int64, error)
}

// Registry is the narrow storage-registry interface admin depends on.
type Registry interface {
	Register(record registry.BackendRecord) (storage.Driver, error)
	Remove(id string)
}

type Usecase interface {
	CreateBackend(ctx context.Context, name, driverType string, config, credentials, rateCard map[string]interface{}) (*domain.StorageBackend, error)
	GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error)
	ListBackends(ctx context.Context) ([]domain.StorageBackend, error)
	ListAvailable(ctx context.Context, tenantID string) ([]domain.StorageBackend, error)
	DeleteBackend(ctx context.Context, id string) error
	CheckHealth(ctx context.Context, id string) (*domain.StorageBackend, error)
	InstallStats(ctx context.Context) (*InstallStats, error)
	ListAllTenants(ctx context.Context) ([]domain.Tenant, error)
}

type InstallStats struct {
	Tenants       int64 `json:"tenants"`
	Users         int64 `json:"users"`
	Objects       int64 `json:"objects"`
	Backends      int   `json:"backends"`
	ActiveJobs    int64 `json:"active_jobs"`
	StorageBytes  int64 `json:"storage_bytes"`
}

// Counters for install stats.
type Counters interface {
	CountTenants(ctx context.Context) (int64, error)
	CountUsers(ctx context.Context) (int64, error)
	CountObjects(ctx context.Context) (int64, error)
	CountActiveJobs(ctx context.Context) (int64, error)
	SumObjectBytes(ctx context.Context) (int64, error)
}

// AllTenantsReader lists every tenant (platform admin).
type AllTenantsReader interface {
	ListAll(ctx context.Context) ([]domain.Tenant, error)
}
