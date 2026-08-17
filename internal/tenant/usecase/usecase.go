package usecase

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/tenant"
)

type usecase struct {
	repo tenant.Repository
}

func NewUsecase(repo tenant.Repository) tenant.Usecase {
	return &usecase{repo: repo}
}

var _ tenant.Usecase = (*usecase)(nil)
var _ tenant.QuotaChecker = (*usecase)(nil)

func (u *usecase) Create(ctx context.Context, name, slug string, quotaBytes, quotaObjects int64, defaultBackendID, ownerUserID string) (*domain.Tenant, error) {
	if name == "" || slug == "" {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	t := &domain.Tenant{
		Name: name, Slug: slug,
		QuotaBytes: quotaBytes, QuotaObjects: quotaObjects,
		DefaultBackendID: defaultBackendID,
	}
	if err := u.repo.Insert(ctx, t); err != nil {
		return nil, err
	}
	// root folder for the tenant
	if err := u.repo.InsertRootFolder(ctx, &domain.Folder{TenantID: t.ID, Name: "", ParentID: nil}); err != nil {
		return nil, err
	}
	if ownerUserID != "" {
		if err := u.repo.InsertMember(ctx, &domain.Membership{UserID: ownerUserID, TenantID: t.ID, Role: "tenant_owner"}); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (u *usecase) Get(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	return u.repo.GetByID(ctx, tenantID)
}

func (u *usecase) ListForUser(ctx context.Context, userID string) ([]domain.Tenant, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *usecase) Update(ctx context.Context, tenantID string, upd tenant.Update) (*domain.Tenant, error) {
	t, err := u.repo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if upd.Name != nil {
		t.Name = *upd.Name
	}
	if upd.QuotaBytes != nil {
		t.QuotaBytes = *upd.QuotaBytes
	}
	if upd.QuotaObjects != nil {
		t.QuotaObjects = *upd.QuotaObjects
	}
	if upd.DefaultBackendID != nil {
		t.DefaultBackendID = *upd.DefaultBackendID
	}
	if upd.StorageEngines != nil {
		t.StorageEngines = *upd.StorageEngines
	}
	if upd.Status != nil {
		t.Status = *upd.Status
	}
	if err := u.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (u *usecase) CheckQuota(ctx context.Context, tenantID string, addBytes, addObjects int64) error {
	t, err := u.repo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	if t.QuotaBytes > 0 && t.UsedBytes+addBytes > t.QuotaBytes {
		return httpx.NewError(httpx.ErrQuotaExceeded, 422)
	}
	if t.QuotaObjects > 0 && t.UsedObjects+addObjects > t.QuotaObjects {
		return httpx.NewError(httpx.ErrQuotaExceeded, 422)
	}
	return nil
}

func (u *usecase) IncrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error {
	return u.repo.IncrementUsed(ctx, tenantID, bytes, objects)
}

func (u *usecase) DecrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error {
	return u.repo.DecrementUsed(ctx, tenantID, bytes, objects)
}

func (u *usecase) ListMembers(ctx context.Context, tenantID string) ([]domain.Membership, error) {
	return u.repo.ListMembers(ctx, tenantID)
}

func (u *usecase) AddMember(ctx context.Context, tenantID, userID, role string) error {
	return u.repo.InsertMember(ctx, &domain.Membership{UserID: userID, TenantID: tenantID, Role: role})
}

func (u *usecase) UpdateMemberRole(ctx context.Context, tenantID, membershipID, role string) error {
	return u.repo.UpdateMemberRole(ctx, membershipID, role)
}

func (u *usecase) RemoveMember(ctx context.Context, tenantID, membershipID string) error {
	return u.repo.RemoveMember(ctx, membershipID)
}

func (u *usecase) CreateInvitation(ctx context.Context, tenantID, email, role, invitedBy string) (*domain.Invitation, error) {
	token := crypto.NewToken(24)
	inv := &domain.Invitation{
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		TokenHash: crypto.HashToken(token),
		InvitedBy: invitedBy,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := u.repo.InsertInvitation(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (u *usecase) ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error) {
	return u.repo.ListInvitations(ctx, tenantID)
}

func (u *usecase) AssignBackend(ctx context.Context, tenantID, backendID string) error {
	t, err := u.repo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	t.DefaultBackendID = backendID
	return u.repo.Update(ctx, t)
}

func (u *usecase) GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error) {
	return u.repo.GetBackend(ctx, id)
}
