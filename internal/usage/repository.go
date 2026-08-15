package usage

import (
	"context"
	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists usage snapshots.
type Repository interface {
	UpsertSnapshot(ctx context.Context, s *domain.UsageSnapshot) error
	Latest(ctx context.Context, tenantID string) (*domain.UsageSnapshot, error)
	History(ctx context.Context, tenantID string, limit int) ([]domain.UsageSnapshot, error)
	AllTenantsLatest(ctx context.Context) ([]domain.UsageSnapshot, error)
}

type Meter interface {
	// Snapshot aggregates a tenant's usage into a period bucket (idempotent).
	Snapshot(ctx context.Context, tenantID, period string) error
}

type Usecase interface {
	Snapshot(ctx context.Context, tenantID, period string) error
	Latest(ctx context.Context, tenantID string) (*domain.UsageSnapshot, error)
	History(ctx context.Context, tenantID string, limit int) ([]domain.UsageSnapshot, error)
	EstimatedCost(ctx context.Context, tenantID string) (*CostEstimate, error)
	AllTenants(ctx context.Context) ([]domain.UsageSnapshot, error)
}

type CostEstimate struct {
	TenantID      string  `json:"tenant_id"`
	Bytes         int64   `json:"bytes"`
	Objects       int64   `json:"objects"`
	Egress        int64   `json:"egress_bytes"`
	StorageCost   float64 `json:"storage_cost"`
	EgressCost    float64 `json:"egress_cost"`
	RequestCost   float64 `json:"request_cost"`
	Total         float64 `json:"total"`
	HasRateCard   bool    `json:"has_rate_card"`
}


