package repository

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repo struct {
	backends *mongo.Collection
	tenants  *mongo.Collection
	users    *mongo.Collection
	objects  *mongo.Collection
	jobs     *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{
		backends: db.Collection("storage_backends"),
		tenants:  db.Collection("tenants"),
		users:    db.Collection("users"),
		objects:  db.Collection("objects"),
		jobs:     db.Collection("jobs"),
	}
}

func (r *repo) InsertBackend(ctx context.Context, b *domain.StorageBackend) error {
	b.ID = crypto.NewID()
	b.CreatedAt = time.Now().UTC()
	_, err := r.backends.InsertOne(ctx, b)
	return err
}

func (r *repo) GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error) {
	var b domain.StorageBackend
	err := r.backends.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &b, nil
}

func (r *repo) ListBackends(ctx context.Context) ([]domain.StorageBackend, error) {
	cur, err := r.backends.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []domain.StorageBackend
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) UpdateBackend(ctx context.Context, b *domain.StorageBackend) error {
	_, err := r.backends.ReplaceOne(ctx, bson.M{"_id": b.ID}, b)
	return err
}

func (r *repo) DeleteBackend(ctx context.Context, id string) error {
	_, err := r.backends.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *repo) CountTenantsOnBackend(ctx context.Context, backendID string) (int64, error) {
	return r.tenants.CountDocuments(ctx, bson.M{"default_backend_id": backendID})
}

func (r *repo) CountTenants(ctx context.Context) (int64, error) {
	return r.tenants.CountDocuments(ctx, bson.M{})
}

func (r *repo) CountUsers(ctx context.Context) (int64, error) {
	return r.users.CountDocuments(ctx, bson.M{})
}

func (r *repo) CountObjects(ctx context.Context) (int64, error) {
	return r.objects.CountDocuments(ctx, bson.M{"state": "active", "deleted_at": bson.M{"$exists": false}})
}

func (r *repo) CountActiveJobs(ctx context.Context) (int64, error) {
	return r.jobs.CountDocuments(ctx, bson.M{"state": bson.M{"$in": []string{"queued", "running"}}})
}

func (r *repo) SumObjectBytes(ctx context.Context) (int64, error) {
	cur, err := r.objects.Find(ctx, bson.M{"state": "active", "deleted_at": bson.M{"$exists": false}})
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)
	var total int64
	for cur.Next(ctx) {
		var o struct {
			Size int64 `bson:"size_bytes"`
		}
		if cur.Decode(&o) == nil {
			total += o.Size
		}
	}
	return total, cur.Err()
}
