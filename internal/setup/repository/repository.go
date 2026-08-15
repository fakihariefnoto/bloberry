package repository

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repo struct {
	users    *mongo.Collection
	tenants  *mongo.Collection
	members  *mongo.Collection
	folders  *mongo.Collection
	backends *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{
		users:    db.Collection("users"),
		tenants:  db.Collection("tenants"),
		members:  db.Collection("memberships"),
		folders:  db.Collection("folders"),
		backends: db.Collection("storage_backends"),
	}
}

func (r *repo) CountUsers(ctx context.Context) (int64, error) {
	return r.users.CountDocuments(ctx, bson.M{})
}

func (r *repo) InsertUser(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		u.ID = crypto.NewID()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	u.UpdatedAt = u.CreatedAt
	_, err := r.users.InsertOne(ctx, u)
	return err
}

func (r *repo) InsertTenant(ctx context.Context, t *domain.Tenant) error {
	if t.ID == "" {
		t.ID = crypto.NewID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	_, err := r.tenants.InsertOne(ctx, t)
	return err
}

func (r *repo) InsertMembership(ctx context.Context, m *domain.Membership) error {
	if m.ID == "" {
		m.ID = crypto.NewID()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	_, err := r.members.InsertOne(ctx, m)
	return err
}

func (r *repo) InsertRootFolder(ctx context.Context, f *domain.Folder) error {
	if f.ID == "" {
		f.ID = crypto.NewID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	f.UpdatedAt = f.CreatedAt
	_, err := r.folders.InsertOne(ctx, f)
	return err
}

func (r *repo) InsertBackend(ctx context.Context, b *domain.StorageBackend) error {
	if b.ID == "" {
		b.ID = crypto.NewStorageID(b.Driver)
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	_, err := r.backends.InsertOne(ctx, b)
	return err
}
