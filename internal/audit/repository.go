package audit

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists audit events (append-only) and queries them.
type Repository interface {
	Insert(ctx context.Context, e *domain.AuditEvent) error
	List(ctx context.Context, tenantID string, filter ListFilter) ([]domain.AuditEvent, error)
}

type ListFilter struct {
	TargetType string
	TargetID   string
	Action     string
	Limit      int
	Cursor     *time.Time
}

// Usecase is the audit domain service.
type Usecase interface {
	Write(ctx context.Context, e *domain.AuditEvent) error
	Query(ctx context.Context, tenantID string, filter ListFilter) ([]domain.AuditEvent, error)
}

// Writer is the narrow write interface other domains depend on.
type Writer interface {
	Write(ctx context.Context, e *domain.AuditEvent) error
}
