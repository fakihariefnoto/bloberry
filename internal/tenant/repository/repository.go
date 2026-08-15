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
	tenants    *mongo.Collection
	members    *mongo.Collection
	invites    *mongo.Collection
	folders    *mongo.Collection
	backends   *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{
		tenants:  db.Collection("tenants"),
		members:  db.Collection("memberships"),
		invites:  db.Collection("invitations"),
		folders:  db.Collection("folders"),
		backends: db.Collection("storage_backends"),
	}
}

func (r *repo) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.tenants.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &t, nil
}

func (r *repo) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.tenants.FindOne(ctx, bson.M{"slug": slug}).Decode(&t)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &t, nil
}

func (r *repo) ListByUser(ctx context.Context, userID string) ([]domain.Tenant, error) {
	cur, err := r.members.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var memberships []domain.Membership
	if err := cur.All(ctx, &memberships); err != nil {
		return nil, err
	}
	out := make([]domain.Tenant, 0)
	for _, m := range memberships {
		if t, err := r.GetByID(ctx, m.TenantID); err == nil {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *repo) Insert(ctx context.Context, t *domain.Tenant) error {
	t.ID = crypto.NewID()
	t.CreatedAt = time.Now().UTC()
	t.Status = "active"
	_, err := r.tenants.InsertOne(ctx, t)
	return err
}

func (r *repo) Update(ctx context.Context, t *domain.Tenant) error {
	_, err := r.tenants.ReplaceOne(ctx, bson.M{"_id": t.ID}, t)
	return err
}

func (r *repo) ListAll(ctx context.Context) ([]domain.Tenant, error) {
	cur, err := r.tenants.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Tenant, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) InsertRootFolder(ctx context.Context, f *domain.Folder) error {
	f.ID = crypto.NewID()
	f.Ancestors = []string{}
	f.Depth = 0
	f.Path = "/"
	f.CreatedAt = time.Now().UTC()
	f.UpdatedAt = f.CreatedAt
	_, err := r.folders.InsertOne(ctx, f)
	return err
}

func (r *repo) IncrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error {
	_, err := r.tenants.UpdateOne(ctx, bson.M{"_id": tenantID}, bson.M{"$inc": bson.M{"used_bytes": bytes, "used_objects": objects}})
	return err
}

func (r *repo) DecrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error {
	_, err := r.tenants.UpdateOne(ctx, bson.M{"_id": tenantID}, bson.M{"$inc": bson.M{"used_bytes": -bytes, "used_objects": -objects}})
	return err
}

func (r *repo) ListMembers(ctx context.Context, tenantID string) ([]domain.Membership, error) {
	cur, err := r.members.Find(ctx, bson.M{"tenant_id": tenantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Membership, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) ListMembershipsByUser(ctx context.Context, userID string) ([]domain.Membership, error) {
	cur, err := r.members.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Membership, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) InsertMember(ctx context.Context, m *domain.Membership) error {
	m.ID = crypto.NewID()
	m.CreatedAt = time.Now().UTC()
	_, err := r.members.InsertOne(ctx, m)
	return err
}

func (r *repo) UpdateMemberRole(ctx context.Context, membershipID, role string) error {
	_, err := r.members.UpdateOne(ctx, bson.M{"_id": membershipID}, bson.M{"$set": bson.M{"role": role}})
	return err
}

func (r *repo) RemoveMember(ctx context.Context, membershipID string) error {
	_, err := r.members.DeleteOne(ctx, bson.M{"_id": membershipID})
	return err
}

func (r *repo) CountOwners(ctx context.Context, tenantID string) (int64, error) {
	return r.members.CountDocuments(ctx, bson.M{"tenant_id": tenantID, "role": "tenant_owner"})
}

func (r *repo) ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error) {
	cur, err := r.invites.Find(ctx, bson.M{"tenant_id": tenantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	out := make([]domain.Invitation, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repo) InsertInvitation(ctx context.Context, inv *domain.Invitation) error {
	inv.ID = crypto.NewID()
	inv.CreatedAt = time.Now().UTC()
	_, err := r.invites.InsertOne(ctx, inv)
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

func (r *repo) CountTenantsOnBackend(ctx context.Context, backendID string) (int64, error) {
	return r.tenants.CountDocuments(ctx, bson.M{"default_backend_id": backendID})
}
