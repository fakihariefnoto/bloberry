package repository

import (
	"context"
	"time"

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
	return &repo{coll: db.Collection("usage_snapshots")}
}

func (r *repo) UpsertSnapshot(ctx context.Context, s *domain.UsageSnapshot) error {
	s.ID = crypto.NewID()
	s.CreatedAt = time.Now().UTC()
	_, err := r.coll.ReplaceOne(ctx,
		bson.M{"tenant_id": s.TenantID, "period": s.Period},
		s, options.Replace().SetUpsert(true))
	return err
}

func (r *repo) Latest(ctx context.Context, tenantID string) (*domain.UsageSnapshot, error) {
	var s domain.UsageSnapshot
	err := r.coll.FindOne(ctx, bson.M{"tenant_id": tenantID},
		options.FindOne().SetSort(bson.M{"created_at": -1})).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repo) History(ctx context.Context, tenantID string, limit int) ([]domain.UsageSnapshot, error) {
	if limit <= 0 || limit > 200 {
		limit = 24
	}
	cur, err := r.coll.Find(ctx, bson.M{"tenant_id": tenantID},
		options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.UsageSnapshot, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) AllTenantsLatest(ctx context.Context) ([]domain.UsageSnapshot, error) {
	// For admin usage: latest snapshot per tenant (approximate v1: all recent).
	cur, err := r.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.UsageSnapshot, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
