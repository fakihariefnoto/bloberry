package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Connect(ctx context.Context, uri, database string) (*DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}
	return &DB{Client: client, DB: client.Database(database)}, nil
}

func (d *DB) Disconnect(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}

type Index struct {
	Collection string
	Keys       bson.D
	Unique     bool
	Partial    string // optional partial filter expression
	TTL        time.Duration
}

// EnsureIndexes is idempotent; safe to run at every boot.
func EnsureIndexes(ctx context.Context, db *mongo.Database, indexes []Index) error {
	collections := map[string]*mongo.Collection{}
	for _, idx := range indexes {
		col, ok := collections[idx.Collection]
		if !ok {
			col = db.Collection(idx.Collection)
			collections[idx.Collection] = col
		}
		mi := mongo.IndexModel{
			Keys:    idx.Keys,
			Options: options.Index().SetUnique(idx.Unique),
		}
		if idx.Partial != "" {
			var f bson.M
			if err := bson.UnmarshalExtJSON([]byte(idx.Partial), false, &f); err != nil {
				return fmt.Errorf("bad partial filter for %s.%s: %w", idx.Collection, idx.Keys, err)
			}
			mi.Options = mi.Options.SetPartialFilterExpression(f)
		}
		if idx.TTL > 0 {
			mi.Options = mi.Options.SetExpireAfterSeconds(int32(idx.TTL.Seconds()))
		}
		if _, err := col.Indexes().CreateOne(ctx, mi); err != nil {
			return fmt.Errorf("create index %s.%v: %w", idx.Collection, idx.Keys, err)
		}
	}
	return nil
}
