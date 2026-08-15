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
	invites    *mongo.Collection
	members    *mongo.Collection
	tenants    *mongo.Collection
	folders    *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{
		invites: db.Collection("invitations"),
		members: db.Collection("memberships"),
		tenants: db.Collection("tenants"),
		folders: db.Collection("folders"),
	}
}

func (r *repo) GetInvitationByTokenHash(ctx context.Context, hash string) (*domain.Invitation, error) {
	var inv domain.Invitation
	err := r.invites.FindOne(ctx, bson.M{"token_hash": hash}).Decode(&inv)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &inv, nil
}

func (r *repo) MarkInvitationAccepted(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := r.invites.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"accepted_at": now}})
	return err
}

func (r *repo) InsertMembership(ctx context.Context, m *domain.Membership) error {
	m.ID = crypto.NewID()
	m.CreatedAt = time.Now().UTC()
	_, err := r.members.InsertOne(ctx, m)
	return err
}

func (r *repo) GetMembership(ctx context.Context, userID, tenantID string) (*domain.Membership, error) {
	var m domain.Membership
	err := r.members.FindOne(ctx, bson.M{"user_id": userID, "tenant_id": tenantID}).Decode(&m)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &m, nil
}

func (r *repo) InsertTenant(ctx context.Context, t *domain.Tenant) error {
	t.ID = crypto.NewID()
	t.CreatedAt = time.Now().UTC()
	t.Status = "active"
	_, err := r.tenants.InsertOne(ctx, t)
	return err
}

func (r *repo) GetTenantBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.tenants.FindOne(ctx, bson.M{"slug": slug}).Decode(&t)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	return &t, nil
}

func (r *repo) InsertRootFolder(ctx context.Context, f *domain.Folder) error {
	f.ID = crypto.NewID()
	f.Ancestors = []string{}
	f.Depth = 0
	f.Path = ""
	f.CreatedAt = time.Now().UTC()
	f.UpdatedAt = f.CreatedAt
	_, err := r.folders.InsertOne(ctx, f)
	return err
}
