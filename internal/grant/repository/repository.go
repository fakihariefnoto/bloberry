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
	return &repo{coll: db.Collection("grants")}
}

func (r *repo) Insert(ctx context.Context, g *domain.Grant) error {
	g.ID = crypto.NewID()
	g.CreatedAt = time.Now().UTC()
	_, err := r.coll.InsertOne(ctx, g)
	return err
}

func (r *repo) GetByID(ctx context.Context, tenantID, id string) (*domain.Grant, error) {
	var g domain.Grant
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}).Decode(&g)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &g, nil
}

func (r *repo) ListByFolder(ctx context.Context, tenantID, folderID string) ([]domain.Grant, error) {
	cur, err := r.coll.Find(ctx, bson.M{"tenant_id": tenantID, "folder_id": folderID, "revoked_at": bson.M{"$exists": false}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Grant, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) ListByPrincipal(ctx context.Context, tenantID, principalType, principalID string) ([]domain.Grant, error) {
	cur, err := r.coll.Find(ctx, bson.M{
		"tenant_id": tenantID, "principal_type": principalType, "principal_id": principalID,
		"revoked_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Grant, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) Revoke(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, bson.M{"$set": bson.M{"revoked_at": now}})
	return err
}
