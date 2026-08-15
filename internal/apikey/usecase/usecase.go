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

func (u *usecase) CreateKey(ctx context.Context, tenantID, applicationID string, scope, perms []string, expiresAt *time.Time) (*apikey.CreatedKey, error) {
	secret := "blob_live_" + crypto.NewToken(24)
	hash, err := crypto.HashPassword(secret)
	if err != nil {
		return nil, err
	}
	k := &domain.AccessKey{
		TenantID: tenantID, ApplicationID: &applicationID,
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

func (u *usecase) ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error) {
	return u.repo.ListKeys(ctx, tenantID, applicationID)
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

func lastFour(secret string) string {
	if len(secret) < 4 {
		return secret
	}
	return secret[len(secret)-4:]
}


