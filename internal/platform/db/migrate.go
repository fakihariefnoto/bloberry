package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Migrate runs the idempotent index migration set in dependency order.
// MongoDB needs no table-creation step; collections are created lazily on
// first write, so migrations only ensure indexes + validators.
func Migrate(ctx context.Context, db *mongo.Database) error {
	indexes := []Index{
		// users
		{Collection: "users", Keys: bsonD("email", 1), Unique: true},
		// tenants
		{Collection: "tenants", Keys: bsonD("slug", 1), Unique: true},
		// memberships
		{Collection: "memberships", Keys: bsonD("user_id", 1, "tenant_id", 1), Unique: true},
		// invitations
		{Collection: "invitations", Keys: bsonD("token_hash", 1), Unique: true},
		{Collection: "invitations", Keys: bsonD("expires_at", 1), TTL: 0},
		// access_keys
		{Collection: "access_keys", Keys: bsonD("secret_hash", 1), Unique: true},
		{Collection: "access_keys", Keys: bsonD("tenant_id", 1, "application_id", 1)},
		// folders
		{Collection: "folders", Keys: bsonD("tenant_id", 1, "parent_id", 1, "name", 1), Unique: true},
		{Collection: "folders", Keys: bsonD("tenant_id", 1, "ancestors", 1)},
		{Collection: "folders", Keys: bsonD("tenant_id", 1, "path", 1), Unique: true},
		// objects
		{Collection: "objects", Keys: bsonD("tenant_id", 1, "folder_id", 1, "name", 1), Unique: true, Partial: `{"deleted_at": null}`},
		{Collection: "objects", Keys: bsonD("tenant_id", 1, "ancestors", 1)},
		{Collection: "objects", Keys: bsonD("tenant_id", 1, "content_hash", 1)},
		{Collection: "objects", Keys: bsonD("state", 1, "created_at", 1), Partial: `{"state": "pending"}`},
		{Collection: "objects", Keys: bsonD("backend_id", 1)},
		// multipart_uploads
		{Collection: "multipart_uploads", Keys: bsonD("expires_at", 1), TTL: 0},
		// grants
		{Collection: "grants", Keys: bsonD("tenant_id", 1, "principal_type", 1, "principal_id", 1), Partial: `{"revoked_at": null}`},
		{Collection: "grants", Keys: bsonD("tenant_id", 1, "folder_id", 1)},
		// share_links
		{Collection: "share_links", Keys: bsonD("slug", 1), Unique: true},
		{Collection: "share_links", Keys: bsonD("tenant_id", 1, "object_id", 1)},
		// jobs
		{Collection: "jobs", Keys: bsonD("tenant_id", 1, "state", 1, "created_at", 1)},
		// usage_snapshots
		{Collection: "usage_snapshots", Keys: bsonD("tenant_id", 1, "period", 1), Unique: true},
		// audit_events
		{Collection: "audit_events", Keys: bsonD("tenant_id", 1, "created_at", -1)},
		{Collection: "audit_events", Keys: bsonD("tenant_id", 1, "target_type", 1, "target_id", 1)},
	}
	return EnsureIndexes(ctx, db, indexes)
}

func bsonD(kv ...interface{}) bson.D {
	var d bson.D
	for i := 0; i+1 < len(kv); i += 2 {
		d = append(d, bson.E{Key: kv[i].(string), Value: kv[i+1]})
	}
	return d
}
