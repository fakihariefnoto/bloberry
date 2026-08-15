package grant

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists folder grants.
type Repository interface {
	Insert(ctx context.Context, g *domain.Grant) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.Grant, error)
	ListByFolder(ctx context.Context, tenantID, folderID string) ([]domain.Grant, error)
	ListByPrincipal(ctx context.Context, tenantID, principalType, principalID string) ([]domain.Grant, error)
	Revoke(ctx context.Context, tenantID, id string) error
}

type Usecase interface {
	Create(ctx context.Context, tenantID, folderID, principalType, principalID string, perms []string, expiresAt *time.Time, grantedBy string) (*domain.Grant, error)
	Revoke(ctx context.Context, tenantID, id string) error
	ListByFolder(ctx context.Context, tenantID, folderID string) ([]domain.Grant, error)
}
