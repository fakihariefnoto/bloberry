package setup

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository — the persistence boundaries the setup wizard needs.
type Repository interface {
	CountUsers(ctx context.Context) (int64, error)
	InsertUser(ctx context.Context, u *domain.User) error
	InsertTenant(ctx context.Context, t *domain.Tenant) error
	InsertMembership(ctx context.Context, m *domain.Membership) error
	InsertRootFolder(ctx context.Context, f *domain.Folder) error
	InsertBackend(ctx context.Context, b *domain.StorageBackend) error
}

type Status struct {
	NeedsSetup bool `json:"needs_setup"`
}

// Usecase — first-run initialization. Runs once on a fresh install:
// creates the platform admin, the first tenant (owner), its root folder,
// and an install-level disk storage backend.
type Usecase interface {
	Status(ctx context.Context) (*Status, error)
	Run(ctx context.Context, email, password, displayName, tenantName, tenantSlug string) error
}
