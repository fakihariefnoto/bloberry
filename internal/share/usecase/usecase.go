package usecase

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/share"
)

type usecase struct {
	repo       share.Repository
	objects    objectReader
	baseURL    string
}

type objectReader interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Object, error)
}

type Deps struct {
	Repo    share.Repository
	Objects objectReader
	BaseURL string
}

func NewUsecase(d Deps) share.Usecase {
	return &usecase{repo: d.Repo, objects: d.Objects, baseURL: d.BaseURL}
}

var _ share.Usecase = (*usecase)(nil)

func (u *usecase) CreateSigned(ctx context.Context, tenantID, objectID, createdBy string, ttlSeconds int) (*share.ShareLink, error) {
	if _, err := u.objects.GetByID(ctx, tenantID, objectID); err != nil {
		return nil, err
	}
	slug := slug()
	var exp *time.Time
	if ttlSeconds > 0 {
		t := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		exp = &t
	}
	l := &domain.ShareLink{
		TenantID: tenantID, ObjectID: objectID, Kind: "signed",
		Slug: slug, ExpiresAt: exp, CreatedBy: createdBy,
	}
	if err := u.repo.Insert(ctx, l); err != nil {
		return nil, err
	}
	return &share.ShareLink{
		ID: l.ID, Kind: "signed", ExpiresAt: exp,
		URL: u.baseURL + "/s/" + slug,
	}, nil
}

func (u *usecase) CreateShort(ctx context.Context, tenantID, objectID, createdBy string) (*share.ShareLink, error) {
	if _, err := u.objects.GetByID(ctx, tenantID, objectID); err != nil {
		return nil, err
	}
	slug := slug()
	l := &domain.ShareLink{
		TenantID: tenantID, ObjectID: objectID, Kind: "short",
		Slug: slug, CreatedBy: createdBy,
	}
	if err := u.repo.Insert(ctx, l); err != nil {
		return nil, err
	}
	return &share.ShareLink{ID: l.ID, Kind: "short", URL: u.baseURL + "/s/" + slug}, nil
}

func (u *usecase) Revoke(ctx context.Context, tenantID, id string) error {
	return u.repo.Revoke(ctx, tenantID, id)
}

func (u *usecase) ListByObject(ctx context.Context, tenantID, objectID string) ([]share.ShareLink, error) {
	links, err := u.repo.ListByObject(ctx, tenantID, objectID)
	if err != nil {
		return nil, err
	}
	out := make([]share.ShareLink, 0, len(links))
	for _, l := range links {
		if l.RevokedAt != nil {
			continue
		}
		out = append(out, share.ShareLink{
			ID: l.ID, Kind: l.Kind, ExpiresAt: l.ExpiresAt, HitCount: l.HitCount,
			URL: u.baseURL + "/s/" + l.Slug,
		})
	}
	return out, nil
}

func (u *usecase) Resolve(ctx context.Context, slug string) (*domain.Object, *domain.ShareLink, error) {
	l, err := u.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, nil, httpx.NewErrorContent(httpx.ErrLinkExpired, 410, "This link has expired")
	}
	if l.RevokedAt != nil || (l.ExpiresAt != nil && time.Now().After(*l.ExpiresAt)) {
		return nil, nil, httpx.NewErrorContent(httpx.ErrLinkExpired, 410, "This link has expired")
	}
	obj, err := u.objects.GetByID(ctx, l.TenantID, l.ObjectID)
	if err != nil {
		return nil, nil, httpx.NewErrorContent(httpx.ErrLinkExpired, 410, "This link has expired")
	}
	_ = u.repo.IncrementHit(ctx, l.TenantID, l.ID)
	return obj, l, nil
}

func slug() string {
	return crypto.RandomID(4)
}
