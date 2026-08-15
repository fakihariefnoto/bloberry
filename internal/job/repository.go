package job

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository persists jobs.
type Repository interface {
	Insert(ctx context.Context, j *domain.Job) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.Job, error)
	Update(ctx context.Context, j *domain.Job) error
	ListByTenant(ctx context.Context, tenantID string) ([]domain.Job, error)
}

// Enqueuer is the narrow interface folder depends on (subtree delete).
type Enqueuer interface {
	Enqueue(ctx context.Context, tenantID, kind string, payload map[string]interface{}) (string, error)
}

// Queue is the Redis list backing the worker.
type Queue interface {
	Enqueue(ctx context.Context, jobID string) error
	Dequeue(ctx context.Context) (string, error)
}

type Usecase interface {
	Enqueue(ctx context.Context, tenantID, kind string, payload map[string]interface{}) (string, error)
	Get(ctx context.Context, tenantID, id string) (*domain.Job, error)
	List(ctx context.Context, tenantID string) ([]domain.Job, error)
	StartWorker(ctx context.Context)
}

type Runner interface {
	Run(ctx context.Context, j *domain.Job) error
}
