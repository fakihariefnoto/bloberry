package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func errNotFound() error { return httpx.ErrResourceNotFound }

type repo struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) *repo {
	return &repo{coll: db.Collection("users")}
}

func (r *repo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errNotFound()
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errNotFound()
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) Insert(ctx context.Context, u *domain.User) error {
	u.ID = crypto.NewID()
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = u.CreatedAt
	if u.Settings == (domain.UserSettings{}) {
		u.Settings = domain.UserSettings{NotificationsEnabled: true, Locale: "en"}
	}
	_, err := r.coll.InsertOne(ctx, u)
	return err
}

func (r *repo) Update(ctx context.Context, u *domain.User) error {
	u.UpdatedAt = time.Now().UTC()
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": u.ID}, u)
	return err
}

func (r *repo) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"last_login_at": now}})
	return err
}
