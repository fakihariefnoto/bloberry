package job

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// TransferDeps is everything the transfer runner needs to copy active objects
// from one storage engine to another.
type TransferDeps struct {
	Objects  objectRepo
	Registry registry
}

type objectRepo interface {
	ListActiveByBackend(ctx context.Context, tenantID, backendID string) ([]domain.Object, error)
	Update(ctx context.Context, o *domain.Object) error
}

type registry interface {
	Get(id string) (storage.Driver, error)
}

// TransferRunner copies active objects from a source storage engine to a
// target one, then re-points each object's metadata at the target. Progress is
// objects-done/total on the job record.
type TransferRunner struct {
	deps TransferDeps
}

func NewTransferRunner(deps TransferDeps) *TransferRunner {
	return &TransferRunner{deps: deps}
}

func (r *TransferRunner) Run(ctx context.Context, j *domain.Job) error {
	sourceID := str(j.Payload, "source_storage_id")
	targetID := str(j.Payload, "target_storage_id")
	if sourceID == "" || targetID == "" {
		return httpx.NewError(httpx.ErrBadRequest, 400)
	}
	if sourceID == targetID {
		return httpx.NewError(httpx.ErrBadRequest, 400)
	}
	src, err := r.deps.Registry.Get(sourceID)
	if err != nil {
		return err
	}
	dst, err := r.deps.Registry.Get(targetID)
	if err != nil {
		return err
	}

	objs, err := r.deps.Objects.ListActiveByBackend(ctx, j.TenantID, sourceID)
	if err != nil {
		return err
	}
	j.ProgressTotal = len(objs)
	for i, o := range objs {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		rc, _, err := src.Get(ctx, o.StorageKey, nil)
		if err != nil {
			// skip objects we can't read (already deleted upstream) but keep going
			j.ProgressDone = i + 1
			continue
		}
		err = dst.Put(ctx, o.StorageKey, rc, o.SizeBytes, o.ContentType)
		rc.Close()
		if err != nil {
			return err
		}
		// re-point metadata at the target engine
		o.BackendID = targetID
		o.UpdatedAt = time.Now().UTC()
		if err := r.deps.Objects.Update(ctx, &o); err != nil {
			return err
		}
		j.ProgressDone = i + 1
	}
	return nil
}

func str(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

var _ Runner = (*TransferRunner)(nil)
