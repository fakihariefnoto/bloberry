package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/admin"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/storage/registry"
)

type usecase struct {
	repo     admin.Repository
	reg      admin.Registry
	counters admin.Counters
	envelope envelope
	allTenants admin.AllTenantsReader
}

type envelope interface {
	Encrypt(plaintext []byte) ([]byte, error)
}

type Deps struct {
	Repo       admin.Repository
	Registry   admin.Registry
	Counters   admin.Counters
	Envelope   envelope
	AllTenants admin.AllTenantsReader
}

func NewUsecase(d Deps) admin.Usecase {
	return &usecase{repo: d.Repo, reg: d.Registry, counters: d.Counters, envelope: d.Envelope, allTenants: d.AllTenants}
}

var _ admin.Usecase = (*usecase)(nil)

func (u *usecase) CreateBackend(ctx context.Context, name, driverType string, config, credentials, rateCard map[string]interface{}) (*domain.StorageBackend, error) {
	if name == "" || driverType == "" {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	credBytes, err := marshalJSON(credentials)
	if err != nil {
		return nil, err
	}
	enc, err := u.envelope.Encrypt(credBytes)
	if err != nil {
		return nil, err
	}
	b := &domain.StorageBackend{
		Driver: driverType, Name: name, Config: config,
		CredentialsEncrypted: enc, RateCard: rateCard,
		HealthStatus: "unchecked",
	}
	if err := u.repo.InsertBackend(ctx, b); err != nil {
		return nil, err
	}
	u.reg.Register(registry.BackendRecord{
		ID: b.ID, DriverType: b.Driver, Config: config, Credentials: credentials,
	})
	return b, nil
}

func (u *usecase) GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error) {
	return u.repo.GetBackend(ctx, id)
}

func (u *usecase) ListBackends(ctx context.Context) ([]domain.StorageBackend, error) {
	return u.repo.ListBackends(ctx)
}

// ListAvailable returns the backends a tenant can write to: install-level
// (tenant_id null) plus the tenant's own BYO backends.
func (u *usecase) ListAvailable(ctx context.Context, tenantID string) ([]domain.StorageBackend, error) {
	return u.repo.ListForTenant(ctx, tenantID)
}

func (u *usecase) DeleteBackend(ctx context.Context, id string) error {
	n, err := u.repo.CountTenantsOnBackend(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return httpx.NewErrorContent(httpx.ErrConflict, 409, "backend has live tenants assigned")
	}
	if err := u.repo.DeleteBackend(ctx, id); err != nil {
		return err
	}
	u.reg.Remove(id)
	return nil
}

func (u *usecase) CheckHealth(ctx context.Context, id string) (*domain.StorageBackend, error) {
	b, err := u.repo.GetBackend(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	b.HealthCheckedAt = &now
	if drv, derr := u.reg.Get(id); derr == nil {
		if herr := drv.HealthCheck(ctx); herr == nil {
			b.HealthStatus = "healthy"
			b.HealthError = ""
		} else {
			b.HealthStatus = "unreachable"
			b.HealthError = herr.Error()
		}
	} else {
		// driver not constructed — report as unreachable with the reason
		b.HealthStatus = "unreachable"
		b.HealthError = derr.Error()
	}
	_ = u.repo.UpdateBackend(ctx, b)
	return b, nil
}

func (u *usecase) InstallStats(ctx context.Context) (*admin.InstallStats, error) {
	if u.counters == nil {
		return &admin.InstallStats{}, nil
	}
	tenants, err := u.counters.CountTenants(ctx)
	if err != nil {
		return nil, err
	}
	users, err := u.counters.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	objects, err := u.counters.CountObjects(ctx)
	if err != nil {
		return nil, err
	}
	jobs, err := u.counters.CountActiveJobs(ctx)
	if err != nil {
		return nil, err
	}
	bytes, err := u.counters.SumObjectBytes(ctx)
	if err != nil {
		return nil, err
	}
	backends, err := u.repo.ListBackends(ctx)
	if err != nil {
		return nil, err
	}
	return &admin.InstallStats{
		Tenants: tenants, Users: users, Objects: objects,
		Backends: len(backends), ActiveJobs: jobs, StorageBytes: bytes,
	}, nil
}

func marshalJSON(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(v)
}

func (u *usecase) ListAllTenants(ctx context.Context) ([]domain.Tenant, error) {
	if u.allTenants == nil {
		return []domain.Tenant{}, nil
	}
	return u.allTenants.ListAll(ctx)
}

// TenantUsage computes install-wide usage live from the tenant records, so the
// admin usage page works even before the hourly metering ticker has run.
func (u *usecase) TenantUsage(ctx context.Context) ([]admin.TenantUsageRow, error) {
	tenants, err := u.ListAllTenants(ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]admin.TenantUsageRow, 0, len(tenants))
	for _, t := range tenants {
		row := admin.TenantUsageRow{
			TenantID: t.ID, Name: t.Name, Slug: t.Slug,
			BytesStored: t.UsedBytes, ObjectCount: t.UsedObjects,
		}
		if t.DefaultBackendID != "" {
			if be, err := u.repo.GetBackend(ctx, t.DefaultBackendID); err == nil {
				if rc, ok := be.RateCard["storage_per_gb_month"].(float64); ok && rc > 0 {
					row.HasRateCard = true
					row.StorageCost = float64(t.UsedBytes) / (1024 * 1024 * 1024) * rc
				}
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
