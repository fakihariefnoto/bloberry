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
	coll *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{coll: db.Collection("share_links")}
}

func (r *repo) Insert(ctx context.Context, l *domain.ShareLink) error {
	l.ID = crypto.NewID()
	l.CreatedAt = time.Now().UTC()
	if l.Slug == "" {
		l.Slug = crypto.RandomID(4)
	}
	_, err := r.coll.InsertOne(ctx, l)
	return err
}

func (r *repo) GetByID(ctx context.Context, tenantID, id string) (*domain.ShareLink, error) {
	var l domain.ShareLink
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}).Decode(&l)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &l, nil
}

func (r *repo) GetBySlug(ctx context.Context, slug string) (*domain.ShareLink, error) {
	var l domain.ShareLink
	err := r.coll.FindOne(ctx, bson.M{"slug": slug}).Decode(&l)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &l, nil
}

func (r *repo) GetByObject(ctx context.Context, tenantID, objectID string) (*domain.ShareLink, error) {
	var l domain.ShareLink
	err := r.coll.FindOne(ctx, bson.M{"tenant_id": tenantID, "object_id": objectID}).Decode(&l)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &l, nil
}

func (r *repo) ListByObject(ctx context.Context, tenantID, objectID string) ([]domain.ShareLink, error) {
	cur, err := r.coll.Find(ctx, bson.M{"tenant_id": tenantID, "object_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []domain.ShareLink
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) Update(ctx context.Context, l *domain.ShareLink) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": l.ID}, l)
	return err
}

func (r *repo) IncrementHit(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id, "tenant_id": tenantID},
		bson.M{"$inc": bson.M{"hit_count": 1}, "$set": bson.M{"last_accessed_at": now}})
	return err
}

func (r *repo) Revoke(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, bson.M{"$set": bson.M{"revoked_at": now}})
	return err
}
