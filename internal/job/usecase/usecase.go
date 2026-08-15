package usecase

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/job"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
)

type usecase struct {
	repo  job.Repository
	queue job.Queue
	run   job.Runner
}

type Deps struct {
	Repo  job.Repository
	Queue job.Queue
	Run   job.Runner
}

func NewUsecase(d Deps) job.Usecase {
	return &usecase{repo: d.Repo, queue: d.Queue, run: d.Run}
}

var _ job.Usecase = (*usecase)(nil)
var _ job.Enqueuer = (*usecase)(nil)

func (u *usecase) Enqueue(ctx context.Context, tenantID, kind string, payload map[string]interface{}) (string, error) {
	j := &domain.Job{
		TenantID: tenantID, Kind: kind, State: "queued",
		Payload: payload, CreatedAt: time.Now().UTC(),
	}
	j.ID = crypto.NewID()
	if err := u.repo.Insert(ctx, j); err != nil {
		return "", err
	}
	if err := u.queue.Enqueue(ctx, j.ID); err != nil {
		return "", err
	}
	return j.ID, nil
}

func (u *usecase) Get(ctx context.Context, tenantID, id string) (*domain.Job, error) {
	return u.repo.GetByID(ctx, tenantID, id)
}

func (u *usecase) List(ctx context.Context, tenantID string) ([]domain.Job, error) {
	return u.repo.ListByTenant(ctx, tenantID)
}

// workerLoop is wired by main: it pulls jobs off the queue and runs them.
func (u *usecase) workerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		id, err := u.queue.Dequeue(ctx)
		if err != nil || id == "" {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}
		j, err := u.repo.GetByID(ctx, "", id)
		if err != nil {
			continue
		}
		now := time.Now().UTC()
		j.State = "running"
		j.StartedAt = &now
		_ = u.repo.Update(ctx, j)
		if err := u.run.Run(ctx, j); err != nil {
			j.State = "failed"
			j.FailureCode = "job_failed"
			j.FailureMessage = err.Error()
			j.FinishedAt = nowPtr()
		} else {
			j.State = "succeeded"
			j.FinishedAt = nowPtr()
		}
		_ = u.repo.Update(ctx, j)
	}
}

func (u *usecase) StartWorker(ctx context.Context) {
	go u.workerLoop(ctx)
}

func nowPtr() *time.Time {
	t := time.Now().UTC()
	return &t
}

var _ = httpx.ErrInternal
