package repository

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/audit"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repo struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{coll: db.Collection("audit_events")}
}

func (r *repo) Insert(ctx context.Context, e *domain.AuditEvent) error {
	e.ID = crypto.NewID()
	e.CreatedAt = time.Now().UTC()
	if e.TenantID == "" {
		return nil // platform-scoped events without a tenant are dropped from audit
	}
	_, err := r.coll.InsertOne(ctx, e)
	return err
}

func (r *repo) List(ctx context.Context, tenantID string, f audit.ListFilter) ([]domain.AuditEvent, error) {
	q := bson.M{"tenant_id": tenantID}
	if f.TargetType != "" {
		q["target_type"] = f.TargetType
	}
	if f.TargetID != "" {
		q["target_id"] = f.TargetID
	}
	if f.Action != "" {
		q["action"] = f.Action
	}
	if f.Cursor != nil {
		q["created_at"] = bson.M{"$lt": f.Cursor}
	}
	limit := f.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	cur, err := r.coll.Find(ctx, q, options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []domain.AuditEvent
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
