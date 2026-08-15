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
	objects  *mongo.Collection
	mp       *mongo.Collection
	tenants  *mongo.Collection
	backends *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{
		objects:  db.Collection("objects"),
		mp:       db.Collection("multipart_uploads"),
		tenants:  db.Collection("tenants"),
		backends: db.Collection("storage_backends"),
	}
}

func (r *repo) Insert(ctx context.Context, o *domain.Object) error {
	if o.ID == "" {
		o.ID = crypto.NewID()
	}
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	if o.State == "" {
		o.State = "pending"
	}
	if o.Visibility == "" {
		o.Visibility = "private"
	}
	_, err := r.objects.InsertOne(ctx, o)
	return err
}

func (r *repo) GetByID(ctx context.Context, tenantID, id string) (*domain.Object, error) {
	var o domain.Object
	err := r.objects.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID, "deleted_at": bson.M{"$exists": false}}).Decode(&o)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &o, nil
}

func (r *repo) GetByName(ctx context.Context, tenantID, folderID, name string) (*domain.Object, error) {
	var o domain.Object
	err := r.objects.FindOne(ctx, bson.M{"tenant_id": tenantID, "folder_id": folderID, "name": name, "deleted_at": bson.M{"$exists": false}}).Decode(&o)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &o, nil
}

func (r *repo) Update(ctx context.Context, o *domain.Object) error {
	o.UpdatedAt = time.Now().UTC()
	_, err := r.objects.ReplaceOne(ctx, bson.M{"_id": o.ID}, o)
	return err
}

func (r *repo) Delete(ctx context.Context, tenantID, id string) error {
	_, err := r.objects.DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
	return err
}

func (r *repo) SoftDelete(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	_, err := r.objects.UpdateOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, bson.M{"$set": bson.M{"deleted_at": now}})
	return err
}

func (r *repo) ListByFolder(ctx context.Context, tenantID, folderID string) ([]domain.Object, error) {
	cur, err := r.objects.Find(ctx, bson.M{"tenant_id": tenantID, "folder_id": folderID, "deleted_at": bson.M{"$exists": false}}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Object, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) CountActive(ctx context.Context, tenantID string) (int64, error) {	return r.objects.CountDocuments(ctx, bson.M{"tenant_id": tenantID, "state": "active", "deleted_at": bson.M{"$exists": false}})
}

func (r *repo) ListActiveByBackend(ctx context.Context, tenantID, backendID string) ([]domain.Object, error) {
	q := bson.M{"tenant_id": tenantID, "state": "active", "deleted_at": bson.M{"$exists": false}}
	if backendID != "" {
		q["backend_id"] = backendID
	}
	cur, err := r.objects.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Object, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) SumActiveBytes(ctx context.Context, tenantID string) (int64, error) {
	cur, err := r.objects.Find(ctx, bson.M{"tenant_id": tenantID, "state": "active", "deleted_at": bson.M{"$exists": false}}, options.Find().SetProjection(bson.M{"size_bytes": 1}))
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

func (r *repo) GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error) {
	var b domain.StorageBackend
	err := r.backends.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &b, nil
}

func (r *repo) GetTenantBackend(ctx context.Context, tenantID string) (*domain.StorageBackend, error) {
	var t domain.Tenant
	if err := r.tenants.FindOne(ctx, bson.M{"_id": tenantID}).Decode(&t); err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	if t.DefaultBackendID == "" {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	return r.GetBackend(ctx, t.DefaultBackendID)
}

func (r *repo) InsertMultipart(ctx context.Context, m *domain.MultipartUpload) error {
	if m.ID == "" {
		m.ID = crypto.NewID()
	}
	m.CreatedAt = time.Now().UTC()
	_, err := r.mp.InsertOne(ctx, m)
	return err
}

func (r *repo) GetMultipart(ctx context.Context, tenantID, objectID string) (*domain.MultipartUpload, error) {
	var m domain.MultipartUpload
	err := r.mp.FindOne(ctx, bson.M{"object_id": objectID, "tenant_id": tenantID}).Decode(&m)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &m, nil
}

func (r *repo) UpdateMultipartParts(ctx context.Context, tenantID, objectID string, parts []domain.PartRec) error {
	_, err := r.mp.UpdateOne(ctx, bson.M{"object_id": objectID, "tenant_id": tenantID}, bson.M{"$set": bson.M{"parts_received": parts}})
	return err
}

func (r *repo) DeleteMultipart(ctx context.Context, tenantID, objectID string) error {
	_, err := r.mp.DeleteOne(ctx, bson.M{"object_id": objectID, "tenant_id": tenantID})
	return err
}
