package usecase

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/apikey"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
)

type usecase struct {
	repo        apikey.Repository
	invalidator apikey.Invalidator
}

type Deps struct {
	Repo        apikey.Repository
	Invalidator apikey.Invalidator
}

func NewUsecase(d Deps) apikey.Usecase {
	return &usecase{repo: d.Repo, invalidator: d.Invalidator}
}

var _ apikey.Usecase = (*usecase)(nil)

func (u *usecase) Register(ctx context.Context, tenantID, name, description string) (*domain.Application, error) {
	a := &domain.Application{TenantID: tenantID, Name: name, Description: description}
	if err := u.repo.InsertApplication(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (u *usecase) Get(ctx context.Context, tenantID, id string) (*domain.Application, error) {
	return u.repo.GetApplication(ctx, tenantID, id)
}

func (u *usecase) List(ctx context.Context, tenantID string) ([]domain.Application, error) {
	return u.repo.ListApplications(ctx, tenantID)
}

func (u *usecase) Delete(ctx context.Context, tenantID, id string) error {
	return u.repo.DeleteApplication(ctx, tenantID, id)
}

func (u *usecase) CreateKey(ctx context.Context, tenantID, applicationID, name string, scope, perms []string, expiresAt *time.Time) (*apikey.CreatedKey, error) {
	secret := "blob_live_" + crypto.NewToken(24)
	hash, err := crypto.HashPassword(secret)
	if err != nil {
		return nil, err
	}
	k := &domain.AccessKey{
		TenantID: tenantID, Name: name, ApplicationID: &applicationID,
		Prefix: "blob_live_", SecretHash: hash,
		LastFour: lastFour(secret),
		ScopeFolderIDs: scope, Permissions: perms,
		ExpiresAt: expiresAt,
	}
	if err := u.repo.InsertKey(ctx, k); err != nil {
		return nil, err
	}
	return &apikey.CreatedKey{KeyID: k.ID, Secret: secret, Prefix: k.Prefix, LastFour: k.LastFour, ExpiresAt: expiresAt}, nil
}

// CreateTenantKey creates a key directly on the tenant (not tied to an
// application). It authenticates as a tenant-scoped principal and can only act
// within the tenant's folder scope + permissions.
func (u *usecase) CreateTenantKey(ctx context.Context, tenantID, name string, scope, perms []string, expiresAt *time.Time) (*apikey.CreatedKey, error) {
	secret := "blob_live_" + crypto.NewToken(24)
	hash, err := crypto.HashPassword(secret)
	if err != nil {
		return nil, err
	}
	k := &domain.AccessKey{
		TenantID: tenantID, Name: name, Prefix: "blob_live_", SecretHash: hash,
		LastFour: lastFour(secret),
		ScopeFolderIDs: scope, Permissions: perms,
		ExpiresAt: expiresAt,
	}
	if err := u.repo.InsertKey(ctx, k); err != nil {
		return nil, err
	}
	return &apikey.CreatedKey{KeyID: k.ID, Secret: secret, Prefix: k.Prefix, LastFour: k.LastFour, ExpiresAt: expiresAt}, nil
}

func (u *usecase) ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error) {
	return u.repo.ListKeys(ctx, tenantID, applicationID)
}

// ListAllKeys enriches each key with its application's + tenant's name.
func (u *usecase) ListAllKeys(ctx context.Context, tenantID string) ([]apikey.KeyWithApp, error) {
	keys, err := u.repo.ListAllKeys(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return u.enrich(ctx, keys)
}

// ListKeysPage returns keys for the API keys page: all tenants for a platform
// admin, or the caller's own tenant otherwise.
func (u *usecase) ListKeysPage(ctx context.Context, tenantID string, isPlatformAdmin bool) ([]apikey.KeyWithApp, error) {
	var keys []domain.AccessKey
	var err error
	if isPlatformAdmin {
		keys, err = u.repo.ListKeysForAdmin(ctx)
	} else {
		keys, err = u.repo.ListAllKeys(ctx, tenantID)
	}
	if err != nil {
		return nil, err
	}
	return u.enrich(ctx, keys)
}

func (u *usecase) enrich(ctx context.Context, keys []domain.AccessKey) ([]apikey.KeyWithApp, error) {
	appNames := map[string]string{}
	tenantNames := map[string]string{}
	appIDs := map[string]struct{}{}
	tenantIDs := map[string]struct{}{}
	for _, k := range keys {
		if k.ApplicationID != nil {
			appIDs[*k.ApplicationID] = struct{}{}
		}
		tenantIDs[k.TenantID] = struct{}{}
	}
	for id := range appIDs {
		if a, err := u.repo.GetApplication(ctx, "", id); err == nil {
			appNames[id] = a.Name
		}
	}
	for id := range tenantIDs {
		if t, err := u.repo.GetTenant(ctx, id); err == nil {
			tenantNames[id] = t.Name
		}
	}
	out := make([]apikey.KeyWithApp, 0, len(keys))
	for _, k := range keys {
		wk := apikey.KeyWithApp{AccessKey: k}
		if k.ApplicationID != nil {
			wk.ApplicationName = appNames[*k.ApplicationID]
		}
		wk.TenantName = tenantNames[k.TenantID]
		out = append(out, wk)
	}
	return out, nil
}

func (u *usecase) RevokeKey(ctx context.Context, tenantID, keyID string) error {
	k, err := u.repo.ListKeys(ctx, tenantID, "")
	if err != nil {
		return err
	}
	for _, key := range k {
		if key.ID == keyID {
			if u.invalidator != nil {
				_ = u.invalidator.InvalidateKey(ctx, key.SecretHash)
			}
			return u.repo.RevokeKey(ctx, tenantID, keyID)
		}
	}
	return u.repo.RevokeKey(ctx, tenantID, keyID)
}

// RevokeKeyAny revokes a key by id across all applications (SDK-facing page).
func (u *usecase) RevokeKeyAny(ctx context.Context, tenantID, keyID string) error {
	k, err := u.repo.ListAllKeys(ctx, tenantID)
	if err != nil {
		return err
	}
	for _, key := range k {
		if key.ID == keyID {
			if u.invalidator != nil {
				_ = u.invalidator.InvalidateKey(ctx, key.SecretHash)
			}
			return u.repo.RevokeKeyAny(ctx, tenantID, keyID)
		}
	}
	return u.repo.RevokeKeyAny(ctx, tenantID, keyID)
}

func lastFour(secret string) string {
	if len(secret) < 4 {
		return secret
	}
	return secret[len(secret)-4:]
}


