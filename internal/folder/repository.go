package folder

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists folders and their subtree relations.
type Repository interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Folder, error)
	GetByPath(ctx context.Context, tenantID, path string) (*domain.Folder, error)
	Create(ctx context.Context, f *domain.Folder) error
	Update(ctx context.Context, f *domain.Folder) error
	Delete(ctx context.Context, tenantID, id string) error
	ListChildren(ctx context.Context, tenantID string, parentID *string) ([]domain.Folder, error)
	// Descendants returns the full subtree (excluding self).
	Descendants(ctx context.Context, tenantID, id string) ([]domain.Folder, error)
	// DescendantObjects returns object IDs + their folder for a subtree.
	DescendantObjects(ctx context.Context, tenantID string, folderIDs []string) ([]domain.Object, error)
}

// Reader is the narrow read interface object depends on.
type Reader interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Folder, error)
}

// Writer is the narrow interface job depends on.
type Writer interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Folder, error)
}

type Usecase interface {
	Create(ctx context.Context, tenantID, parentID, name string) (*domain.Folder, error)
	Get(ctx context.Context, tenantID, id string) (*domain.Folder, error)
	GetRoot(ctx context.Context, tenantID string) (*domain.Folder, error)
	Rename(ctx context.Context, tenantID, id, name string) (*domain.Folder, error)
	Move(ctx context.Context, tenantID, id, targetParentID string) (*domain.Folder, error)
	Delete(ctx context.Context, tenantID, id string) error
	SetPolicy(ctx context.Context, tenantID, id string, p *domain.UploadPolicy) (*domain.Folder, error)
	ListChildren(ctx context.Context, tenantID string, parentID *string) ([]domain.Folder, error)
}
