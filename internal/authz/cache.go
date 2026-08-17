package authz

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	goredis "github.com/redis/go-redis/v9"
)

type ctxKey int

const tenantCtxKey ctxKey = 1

// TenantIDFrom returns the tenant requested via the X-Tenant-ID header, if any.
func TenantIDFrom(ctx context.Context) string {
	tid, _ := ctx.Value(tenantCtxKey).(string)
	return tid
}

// WithTenantID returns a context carrying the requested tenant.
func WithTenantID(ctx context.Context, tid string) context.Context {
	return context.WithValue(ctx, tenantCtxKey, tid)
}

// Cache stores assembled Principals in Redis with explicit invalidation
// (ADR-6): revocation takes effect on the next request, not after a TTL.
type Cache struct {
	rdb *goredis.Client
}

func NewCache(rdb *goredis.Client) *Cache { return &Cache{rdb: rdb} }

func (c *Cache) Get(ctx context.Context, ptype PrincipalType, id string) (*Principal, bool) {
	key := "principal:" + string(ptype) + ":" + id
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	var p Principal
	if json.Unmarshal([]byte(val), &p) != nil {
		return nil, false
	}
	return &p, true
}

func (c *Cache) Set(ctx context.Context, ptype PrincipalType, id string, p *Principal) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, "principal:"+string(ptype)+":"+id, b, 0).Err()
}

func (c *Cache) Invalidate(ctx context.Context, ptype PrincipalType, id string) error {
	return c.rdb.Del(ctx, "principal:"+string(ptype)+":"+id).Err()
}

// InvalidateKey clears a key's cached lookup after revoke.
func (c *Cache) InvalidateKey(ctx context.Context, secretHash string) error {
	return c.rdb.Del(ctx, "apikey:"+secretHash).Err()
}

// Loader assembles Principals from the domain repositories. It implements
// httpx.PrincipalResolver.
type Loader struct {
	Cache *Cache
	Users userRepo
	Keys  keyRepo
	Grant grantRepo
	Member memberRepo
}

type userRepo interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
type keyRepo interface {
	GetKeyByHash(ctx context.Context, hash string) (*domain.AccessKey, error)
}
type grantRepo interface {
	ListByPrincipal(ctx context.Context, tenantID, principalType, principalID string) ([]domain.Grant, error)
}
type memberRepo interface {
	ListMembershipsByUser(ctx context.Context, userID string) ([]domain.Membership, error)
}

func (l *Loader) ResolveUser(ctx context.Context, userID string) (*Principal, error) {
	// The platform-admin flag is not cached-stable: a user promoted after their
	// principal was cached must see the promotion immediately. Always read the
	// user doc (single indexed lookup) and refresh the cached principal when
	// the flag differs.
	u, err := l.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	isAdmin := u.PlatformRole != nil && *u.PlatformRole == "platform_admin"
	requested := TenantIDFrom(ctx)
	cacheID := userID
	if requested != "" {
		cacheID = userID + ":" + requested
	}
	if p, ok := l.Cache.Get(ctx, PrincipalUser, cacheID); ok {
		if p.IsPlatformAdmin == isAdmin {
			return p, nil
		}
		p.IsPlatformAdmin = isAdmin
		_ = l.Cache.Set(ctx, PrincipalUser, cacheID, p)
		return p, nil
	}
	members, err := l.Member.ListMembershipsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var m *domain.Membership
	if requested != "" {
		for i := range members {
			if members[i].TenantID == requested {
				m = &members[i]
				break
			}
		}
		if m == nil {
			return nil, errors.New("no membership for requested tenant")
		}
	} else if len(members) > 0 {
		m = &members[0]
	}
	p := &Principal{Type: PrincipalUser, ID: userID, IsPlatformAdmin: isAdmin}
	if m != nil {
		p.TenantID = m.TenantID
		p.Role = Role(m.Role)
		grants, err := l.Grant.ListByPrincipal(ctx, m.TenantID, "user", userID)
		if err == nil {
			p.Grants = grantsToAuthz(grants)
		}
	}
	_ = l.Cache.Set(ctx, PrincipalUser, cacheID, p)
	return p, nil
}

// ResolveAccessKey returns the Principal for a hashed key secret.
// terminal=true means revoked/expired — do not retry.
func (l *Loader) ResolveAccessKey(ctx context.Context, secretHash string) (*Principal, bool, error) {
	k, err := l.Keys.GetKeyByHash(ctx, secretHash)
	if err != nil {
		return nil, false, err
	}
	if k.RevokedAt != nil {
		return nil, true, errors.New("key revoked")
	}
	if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
		return nil, true, errors.New("key expired")
	}
	if p, ok := l.Cache.Get(ctx, PrincipalApplication, k.ID); ok {
		return p, false, nil
	}
	grants, err := l.Grant.ListByPrincipal(ctx, k.TenantID, "application", k.ID)
	if err != nil {
		grants = nil
	}
	p := &Principal{
		Type:           PrincipalApplication,
		ID:             k.ID,
		TenantID:       k.TenantID,
		Role:           RoleMember,
		ScopeFolderIDs: k.ScopeFolderIDs,
		KeyPermissions: toPerms(k.Permissions),
		Grants:         grantsToAuthz(grants),
	}
	_ = l.Cache.Set(ctx, PrincipalApplication, k.ID, p)
	return p, false, nil
}

func (l *Loader) InvalidatePrincipal(ctx context.Context, ptype, id, tenantID string) error {
	key := id
	if tenantID != "" {
		key = id + ":" + tenantID
	}
	return l.Cache.Invalidate(ctx, PrincipalType(ptype), key)
}

// InvalidateKey clears a key's cached lookup after revoke.
func (l *Loader) InvalidateKey(ctx context.Context, secretHash string) error {
	return l.Cache.InvalidateKey(ctx, secretHash)
}

func grantsToAuthz(gs []domain.Grant) []Grant {
	out := make([]Grant, 0, len(gs))
	for _, g := range gs {
		out = append(out, Grant{
			FolderID: g.FolderID,
			Expired:  g.ExpiresAt != nil && time.Now().After(*g.ExpiresAt),
			Revoked:  g.RevokedAt != nil,
			Perms:    toPerms(g.Permissions),
		})
	}
	return out
}

func toPerms(ps []string) []Permission {
	out := make([]Permission, 0, len(ps))
	for _, p := range ps {
		out = append(out, Permission(p))
	}
	return out
}
