package tenant

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists tenants, memberships and invitations.
type Repository interface {
	GetByID(ctx context.Context, id string) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	ListByUser(ctx context.Context, userID string) ([]domain.Tenant, error)
	Insert(ctx context.Context, t *domain.Tenant) error
	Update(ctx context.Context, t *domain.Tenant) error
	ListAll(ctx context.Context) ([]domain.Tenant, error)
	InsertRootFolder(ctx context.Context, f *domain.Folder) error
	IncrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
	DecrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error

	ListMembers(ctx context.Context, tenantID string) ([]domain.Membership, error)
	ListMembershipsByUser(ctx context.Context, userID string) ([]domain.Membership, error)
	InsertMember(ctx context.Context, m *domain.Membership) error
	UpdateMemberRole(ctx context.Context, membershipID, role string) error
	RemoveMember(ctx context.Context, membershipID string) error
	CountOwners(ctx context.Context, tenantID string) (int64, error)

	ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error)
	InsertInvitation(ctx context.Context, inv *domain.Invitation) error

	GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error)
	CountTenantsOnBackend(ctx context.Context, backendID string) (int64, error)
}

// QuotaChecker is the narrow interface object depends on.
type QuotaChecker interface {
	CheckQuota(ctx context.Context, tenantID string, addBytes int64, addObjects int64) error
	IncrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
	DecrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
}

type Update struct {
	Name            *string
	QuotaBytes      *int64
	QuotaObjects    *int64
	DefaultBackendID *string
	StorageEngines  *[]string
	Status          *string
	UploadPolicy    *domain.UploadPolicy
}

// Usecase is the tenant domain service.
type Usecase interface {
	Create(ctx context.Context, name, slug string, quotaBytes, quotaObjects int64, defaultBackendID string, ownerUserID string) (*domain.Tenant, error)
	Get(ctx context.Context, tenantID string) (*domain.Tenant, error)
	ListForUser(ctx context.Context, userID string) ([]domain.Tenant, error)
	Update(ctx context.Context, tenantID string, u Update) (*domain.Tenant, error)
	CheckQuota(ctx context.Context, tenantID string, addBytes, addObjects int64) error
	IncrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
	DecrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
	ListMembers(ctx context.Context, tenantID string) ([]domain.Membership, error)
	AddMember(ctx context.Context, tenantID, userID, role string) error
	UpdateMemberRole(ctx context.Context, tenantID, membershipID, role string) error
	RemoveMember(ctx context.Context, tenantID, membershipID string) error
	CreateInvitation(ctx context.Context, tenantID, email, role, invitedBy string) (*domain.Invitation, error)
	ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error)
	AssignBackend(ctx context.Context, tenantID, backendID string) error
	GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error)
}

var _ = time.Time{}
