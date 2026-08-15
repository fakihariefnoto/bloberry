package repository

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repo struct {
	folders *mongo.Collection
	objects *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{folders: db.Collection("folders"), objects: db.Collection("objects")}
}

func (r *repo) GetByID(ctx context.Context, tenantID, id string) (*domain.Folder, error) {
	var f domain.Folder
	err := r.folders.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}).Decode(&f)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &f, nil
}

func (r *repo) GetByPath(ctx context.Context, tenantID, path string) (*domain.Folder, error) {
	var f domain.Folder
	err := r.folders.FindOne(ctx, bson.M{"tenant_id": tenantID, "path": path}).Decode(&f)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &f, nil
}

func (r *repo) Create(ctx context.Context, f *domain.Folder) error {
	f.ID = crypto.NewID()
	f.CreatedAt = time.Now().UTC()
	f.UpdatedAt = f.CreatedAt
	_, err := r.folders.InsertOne(ctx, f)
	return err
}

func (r *repo) Update(ctx context.Context, f *domain.Folder) error {
	f.UpdatedAt = time.Now().UTC()
	_, err := r.folders.ReplaceOne(ctx, bson.M{"_id": f.ID}, f)
	return err
}

func (r *repo) Delete(ctx context.Context, tenantID, id string) error {
	_, err := r.folders.DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
	return err
}

func (r *repo) ListChildren(ctx context.Context, tenantID string, parentID *string) ([]domain.Folder, error) {
	q := bson.M{"tenant_id": tenantID}
	if parentID == nil {
		q["parent_id"] = bson.M{"$exists": false}
	} else {
		q["parent_id"] = *parentID
	}
	cur, err := r.folders.Find(ctx, q, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []domain.Folder
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) Descendants(ctx context.Context, tenantID, id string) ([]domain.Folder, error) {
	cur, err := r.folders.Find(ctx, bson.M{"tenant_id": tenantID, "ancestors": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []domain.Folder
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) DescendantObjects(ctx context.Context, tenantID string, folderIDs []string) ([]domain.Object, error) {
	cur, err := r.objects.Find(ctx, bson.M{"tenant_id": tenantID, "folder_id": bson.M{"$in": folderIDs}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []domain.Object
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
