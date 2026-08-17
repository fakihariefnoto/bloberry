package usecase

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/grant"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
)

type usecase struct {
	repo        grant.Repository
	invalidator principalInvalidator
	folders     folderReader
}

type principalInvalidator interface {
	InvalidatePrincipal(ctx context.Context, principalType, principalID, tenantID string) error
}

type folderReader interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Folder, error)
}

type Deps struct {
	Repo        grant.Repository
	Invalidator principalInvalidator
	Folders     folderReader
}

func NewUsecase(d Deps) grant.Usecase {
	return &usecase{repo: d.Repo, invalidator: d.Invalidator, folders: d.Folders}
}

var _ grant.Usecase = (*usecase)(nil)

func (u *usecase) Create(ctx context.Context, tenantID, folderID, principalType, principalID string, perms []string, expiresAt *time.Time, grantedBy string) (*domain.Grant, error) {
	if _, err := u.folders.GetByID(ctx, tenantID, folderID); err != nil {
		return nil, err
	}
	if principalType != "user" && principalType != "application" {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	g := &domain.Grant{
		TenantID: tenantID, FolderID: folderID,
		PrincipalType: principalType, PrincipalID: principalID,
		Permissions: perms, ExpiresAt: expiresAt, GrantedBy: grantedBy,
	}
	if err := u.repo.Insert(ctx, g); err != nil {
		return nil, err
	}
	if u.invalidator != nil {
		_ = u.invalidator.InvalidatePrincipal(ctx, principalType, principalID, tenantID)
	}
	return g, nil
}

func (u *usecase) Revoke(ctx context.Context, tenantID, id string) error {
	g, err := u.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if err := u.repo.Revoke(ctx, tenantID, id); err != nil {
		return err
	}
	if u.invalidator != nil {
		_ = u.invalidator.InvalidatePrincipal(ctx, g.PrincipalType, g.PrincipalID, tenantID)
	}
	return nil
}

func (u *usecase) ListByFolder(ctx context.Context, tenantID, folderID string) ([]domain.Grant, error) {
	return u.repo.ListByFolder(ctx, tenantID, folderID)
}
