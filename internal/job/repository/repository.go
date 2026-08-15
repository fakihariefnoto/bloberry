package repository

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repo struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{coll: db.Collection("jobs")}
}

func (r *repo) Insert(ctx context.Context, j *domain.Job) error {
	if j.CreatedAt.IsZero() {
		j.CreatedAt = time.Now().UTC()
	}
	_, err := r.coll.InsertOne(ctx, j)
	return err
}

func (r *repo) GetByID(ctx context.Context, tenantID, id string) (*domain.Job, error) {
	q := bson.M{"_id": id}
	if tenantID != "" {
		q["tenant_id"] = tenantID
	}
	var j domain.Job
	err := r.coll.FindOne(ctx, q).Decode(&j)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &j, nil
}

func (r *repo) Update(ctx context.Context, j *domain.Job) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": j.ID}, j)
	return err
}

func (r *repo) ListByTenant(ctx context.Context, tenantID string) ([]domain.Job, error) {
	cur, err := r.coll.Find(ctx, bson.M{"tenant_id": tenantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Job, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
