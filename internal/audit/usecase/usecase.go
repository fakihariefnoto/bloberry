package usecase

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/audit"
	"github.com/fakihariefnoto/bloberry/internal/domain"
)

type usecase struct {
	repo audit.Repository
}

func NewUsecase(repo audit.Repository) audit.Usecase {
	return &usecase{repo: repo}
}

var _ audit.Usecase = (*usecase)(nil)

func (u *usecase) Write(ctx context.Context, e *domain.AuditEvent) error {
	return u.repo.Insert(ctx, e)
}

func (u *usecase) Query(ctx context.Context, tenantID string, filter audit.ListFilter) ([]domain.AuditEvent, error) {
	return u.repo.List(ctx, tenantID, filter)
}
