package authz

import (
	"context"
	"os"
	"testing"

	goredis "github.com/redis/go-redis/v9"
)

// TestTenantScopedPrincipalCache guards the multi-tenant switcher: a user in
// two tenants must get a distinct cached principal per tenant. Regression for
// the bug where Cache.Set wrote under the bare user id while Get read
// "user:tenant", so the last-resolved tenant's principal leaked into requests
// for any other tenant (all tenant roots pointed at the first tenant).
func TestTenantScopedPrincipalCache(t *testing.T) {
	dsn := os.Getenv("REDIS_TEST_URI")
	if dsn == "" {
		dsn = "redis://127.0.0.1:6379"
	}
	opts, err := goredis.ParseURL(dsn)
	if err != nil {
		t.Skipf("bad REDIS_TEST_URI: %v", err)
	}
	rdb := goredis.NewClient(opts)
	defer rdb.Close()
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skipf("redis unavailable: %v", err)
	}

	cache := NewCache(rdb)
	_ = rdb.Del(context.Background(), "principal:user:u:t1", "principal:user:u:t2").Err()
	defer rdb.Del(context.Background(), "principal:user:u:t1", "principal:user:u:t2").Err()

	// Seeded principal for tenant t1 (first membership).
	if err := cache.Set(context.Background(), PrincipalUser, "u:t1", &Principal{
		Type: PrincipalUser, ID: "u", TenantID: "t1", Role: RoleTenantOwner,
	}); err != nil {
		t.Fatalf("seed t1: %v", err)
	}

	// The buggy Set wrote "principal:user:u"; Get("u:t2") would then miss and
	// rebuild, but Get("u") (the no-header path) would return the wrong tenant.
	if _, ok := cache.Get(context.Background(), PrincipalUser, "u"); ok {
		t.Fatal("bare user key must not exist — Set must key by user:tenant")
	}

	p1, ok := cache.Get(context.Background(), PrincipalUser, "u:t1")
	if !ok || p1.TenantID != "t1" {
		t.Fatalf("expected cached t1 principal, got ok=%v p=%+v", ok, p1)
	}

	// Resolve tenant t2 into the same cache.
	if err := cache.Set(context.Background(), PrincipalUser, "u:t2", &Principal{
		Type: PrincipalUser, ID: "u", TenantID: "t2", Role: RoleMember,
	}); err != nil {
		t.Fatalf("seed t2: %v", err)
	}

	// t1 principal must be untouched.
	p1, ok = cache.Get(context.Background(), PrincipalUser, "u:t1")
	if !ok || p1.TenantID != "t1" {
		t.Fatalf("t1 principal clobbered by t2: ok=%v p=%+v", ok, p1)
	}
	p2, ok := cache.Get(context.Background(), PrincipalUser, "u:t2")
	if !ok || p2.TenantID != "t2" {
		t.Fatalf("expected cached t2 principal, got ok=%v p=%+v", ok, p2)
	}
}
