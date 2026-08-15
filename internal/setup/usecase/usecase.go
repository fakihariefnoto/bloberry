package usecase

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/setup"
)

type usecase struct {
	repo          setup.Repository
	diskRoot      string
}

type Deps struct {
	Repo     setup.Repository
	DiskRoot string
}

func NewUsecase(d Deps) setup.Usecase {
	return &usecase{repo: d.Repo, diskRoot: d.DiskRoot}
}

var _ setup.Usecase = (*usecase)(nil)

func (u *usecase) Status(ctx context.Context) (*setup.Status, error) {
	n, err := u.repo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	return &setup.Status{NeedsSetup: n == 0}, nil
}

func (u *usecase) Run(ctx context.Context, email, password, displayName, tenantName, tenantSlug string) error {
	if email == "" || password == "" || tenantName == "" || tenantSlug == "" {
		return httpx.NewError(httpx.ErrBadRequest, 400)
	}
	n, err := u.repo.CountUsers(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return httpx.NewError(httpx.ErrConflict, 409)
	}

	// 1. platform admin
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}
	role := "platform_admin"
	now := time.Now().UTC()
	admin := &domain.User{
		ID:            crypto.NewID(),
		Email:         email,
		PasswordHash:  &hash,
		DisplayName:   displayName,
		PlatformRole:  &role,
		EmailVerified: true,
		Settings: domain.UserSettings{
			NotificationsEnabled: true,
			Locale:               "en",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.repo.InsertUser(ctx, admin); err != nil {
		return err
	}

	// 2. install-level disk backend (created before the tenant so the tenant
	// can reference it as its default)
	be := &domain.StorageBackend{
		ID:           crypto.NewID(),
		Driver:       "disk",
		Name:         "local-disk",
		Config:       map[string]interface{}{"root": u.diskRoot},
		RateCard:     map[string]interface{}{"storage_per_gb_month": 0.0, "egress_per_gb": 0.0, "per_1k_requests": 0.0},
		HealthStatus: "unchecked",
		CreatedAt:    now,
	}
	if err := u.repo.InsertBackend(ctx, be); err != nil {
		return err
	}

	// 3. first tenant
	tenant := &domain.Tenant{
		ID:               crypto.NewID(),
		Name:             tenantName,
		Slug:             tenantSlug,
		QuotaBytes:       0,
		Status:           "active",
		DefaultBackendID: be.ID,
		CreatedAt:        now,
	}
	if err := u.repo.InsertTenant(ctx, tenant); err != nil {
		return err
	}

	// 4. owner membership
	if err := u.repo.InsertMembership(ctx, &domain.Membership{
		ID:        crypto.NewID(),
		UserID:    admin.ID,
		TenantID:  tenant.ID,
		Role:      "tenant_owner",
		CreatedAt: now,
	}); err != nil {
		return err
	}

	// 5. tenant root folder
	if err := u.repo.InsertRootFolder(ctx, &domain.Folder{
		ID:        crypto.NewID(),
		TenantID:  tenant.ID,
		Name:      "",
		Path:      "/",
		Ancestors: []string{},
		Depth:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}); err != nil {
		return err
	}

	return nil
}
