package usecase

import (
	"context"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/usage"
)

type usecase struct {
	repo     usage.Repository
	objects  objectCounter
	backends backendReader
}

type objectCounter interface {
	CountActive(ctx context.Context, tenantID string) (int64, error)
	SumActiveBytes(ctx context.Context, tenantID string) (int64, error)
}

type backendReader interface {
	GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error)
}

type Deps struct {
	Repo     usage.Repository
	Objects  objectCounter
	Backends backendReader
}

func NewUsecase(d Deps) usage.Usecase {
	return &usecase{repo: d.Repo, objects: d.Objects, backends: d.Backends}
}

var _ usage.Usecase = (*usecase)(nil)
var _ usage.Meter = (*usecase)(nil)

func (u *usecase) Snapshot(ctx context.Context, tenantID, period string) error {
	bytes, err := u.objects.SumActiveBytes(ctx, tenantID)
	if err != nil {
		return err
	}
	count, err := u.objects.CountActive(ctx, tenantID)
	if err != nil {
		return err
	}
	latest, _ := u.repo.Latest(ctx, tenantID)
	egress := int64(0)
	if latest != nil {
		egress = latest.EgressBytes
	}
	s := &domain.UsageSnapshot{
		TenantID: tenantID, Period: period,
		BytesStored: bytes, ObjectCount: count, EgressBytes: egress,
	}
	return u.repo.UpsertSnapshot(ctx, s)
}

func (u *usecase) Latest(ctx context.Context, tenantID string) (*domain.UsageSnapshot, error) {
	return u.repo.Latest(ctx, tenantID)
}

func (u *usecase) History(ctx context.Context, tenantID string, limit int) ([]domain.UsageSnapshot, error) {
	return u.repo.History(ctx, tenantID, limit)
}

func (u *usecase) EstimatedCost(ctx context.Context, tenantID string) (*usage.CostEstimate, error) {
	latest, err := u.repo.Latest(ctx, tenantID)
	if err != nil {
		latest = &domain.UsageSnapshot{}
	}
	est := &usage.CostEstimate{
		TenantID: tenantID, Bytes: latest.BytesStored,
		Objects: latest.ObjectCount, Egress: latest.EgressBytes,
	}
	// Rate card lookup: tenant's backend. Resolve via object counter's backend.
	// v1 keeps this approximate — cost estimate is a Should (M18).
	return est, nil
}

func (u *usecase) AllTenants(ctx context.Context) ([]domain.UsageSnapshot, error) {
	return u.repo.AllTenantsLatest(ctx)
}


