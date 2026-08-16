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
	apps    *mongo.Collection
	keys    *mongo.Collection
	tenants *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{
		apps:    db.Collection("applications"),
		keys:    db.Collection("access_keys"),
		tenants: db.Collection("tenants"),
	}
}

func (r *repo) InsertApplication(ctx context.Context, a *domain.Application) error {
	a.ID = crypto.NewID()
	a.CreatedAt = time.Now().UTC()
	_, err := r.apps.InsertOne(ctx, a)
	return err
}

func (r *repo) GetApplication(ctx context.Context, tenantID, id string) (*domain.Application, error) {
	q := bson.M{"_id": id}
	if tenantID != "" {
		q["tenant_id"] = tenantID
	}
	var a domain.Application
	err := r.apps.FindOne(ctx, q).Decode(&a)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &a, nil
}

func (r *repo) GetTenant(ctx context.Context, id string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.tenants.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &t, nil
}

func (r *repo) ListApplications(ctx context.Context, tenantID string) ([]domain.Application, error) {
	cur, err := r.apps.Find(ctx, bson.M{"tenant_id": tenantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Application, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) DeleteApplication(ctx context.Context, tenantID, id string) error {
	_, err := r.apps.DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
	return err
}

func (r *repo) InsertKey(ctx context.Context, k *domain.AccessKey) error {
	k.ID = crypto.NewID()
	k.CreatedAt = time.Now().UTC()
	_, err := r.keys.InsertOne(ctx, k)
	return err
}

func (r *repo) GetKeyByHash(ctx context.Context, hash string) (*domain.AccessKey, error) {
	var k domain.AccessKey
	err := r.keys.FindOne(ctx, bson.M{"secret_hash": hash}).Decode(&k)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &k, nil
}

func (r *repo) ListKeys(ctx context.Context, tenantID, applicationID string) ([]domain.AccessKey, error) {
	q := bson.M{"tenant_id": tenantID}
	if applicationID != "" {
		q["application_id"] = applicationID
	}
	cur, err := r.keys.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.AccessKey, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) ListAllKeys(ctx context.Context, tenantID string) ([]domain.AccessKey, error) {
	cur, err := r.keys.Find(ctx, bson.M{"tenant_id": tenantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.AccessKey, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) ListKeysForAdmin(ctx context.Context) ([]domain.AccessKey, error) {
	cur, err := r.keys.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.AccessKey, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) RevokeKey(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	_, err := r.keys.UpdateOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, bson.M{"$set": bson.M{"revoked_at": now}})
	return err
}

func (r *repo) RevokeKeyAny(ctx context.Context, tenantID, id string) error {
	return r.RevokeKey(ctx, tenantID, id)
}

func (r *repo) TouchKey(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := r.keys.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"last_used_at": now}})
	return err
}
