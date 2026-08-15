package share

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists share links.
type Repository interface {
	Insert(ctx context.Context, l *domain.ShareLink) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.ShareLink, error)
	GetBySlug(ctx context.Context, slug string) (*domain.ShareLink, error)
	GetByObject(ctx context.Context, tenantID, objectID string) (*domain.ShareLink, error)
	ListByObject(ctx context.Context, tenantID, objectID string) ([]domain.ShareLink, error)
	Update(ctx context.Context, l *domain.ShareLink) error
	IncrementHit(ctx context.Context, tenantID, id string) error
	Revoke(ctx context.Context, tenantID, id string) error
}

type ShareLink struct {
	ID        string     `json:"id"`
	Kind      string     `json:"kind"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	HitCount  int64      `json:"hit_count"`
}

type Usecase interface {
	CreateSigned(ctx context.Context, tenantID, objectID, createdBy string, ttlSeconds int) (*ShareLink, error)
	CreateShort(ctx context.Context, tenantID, objectID, createdBy string) (*ShareLink, error)
	Revoke(ctx context.Context, tenantID, id string) error
	ListByObject(ctx context.Context, tenantID, objectID string) ([]ShareLink, error)
	Resolve(ctx context.Context, slug string) (*domain.Object, *domain.ShareLink, error)
}
