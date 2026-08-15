package httpx

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/authz"
	"github.com/fakihariefnoto/bloberry/internal/platform/redis"
	"github.com/golang-jwt/jwt/v5"
	goredis "github.com/redis/go-redis/v9"
)

type ctxKey int

const principalKey ctxKey = 0

// Claims is the JWT payload for human sessions.
type Claims struct {
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type PrincipalResolver interface {
	// ResolveAccessKey returns the Principal for a hashed access-key secret.
	// Second return is a terminal flag: true means the key is revoked/expired
	// (do not retry), false means invalid_credentials semantics.
	ResolveAccessKey(ctx context.Context, secretHash string) (*authz.Principal, bool, error)
	// ResolveUser returns the Principal for a human user id.
	ResolveUser(ctx context.Context, userID string) (*authz.Principal, error)
}

// Middleware authenticates requests. Human JWT and application access keys
// both resolve to one authz.Principal (§4.7). Nothing downstream branches on
// the scheme. Per-access-key rate limiting (PRD-Q5).
type Middleware struct {
	jwtSecret []byte
	resolver  PrincipalResolver
	rdb       *goredis.Client
	rateLimit int
	rateWindow time.Duration
}

func NewMiddleware(jwtSecret []byte, resolver PrincipalResolver, rdb *goredis.Client, rateLimit int, rateWindow time.Duration) *Middleware {
	return &Middleware{
		jwtSecret:  jwtSecret,
		resolver:   resolver,
		rdb:        rdb,
		rateLimit:  rateLimit,
		rateWindow: rateWindow,
	}
}

// Authenticate wraps a handler requiring a valid principal.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := m.resolve(r)
		if !ok {
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), principalKey, principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) resolve(r *http.Request) (*authz.Principal, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return nil, false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return nil, false
	}

	if strings.HasPrefix(token, "blob_") {
		return m.resolveAccessKey(r, token)
	}
	return m.resolveJWT(r, token)
}

func (m *Middleware) resolveJWT(r *http.Request, token string) (*authz.Principal, bool) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))
	claims := &Claims{}
	parsed, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, false
	}
	p, err := m.resolver.ResolveUser(r.Context(), claims.UserID)
	if err != nil {
		return nil, false
	}
	return p, true
}

func (m *Middleware) resolveAccessKey(r *http.Request, token string) (*authz.Principal, bool) {
	hash := token // argon2id hash of the secret; hashed at rest per PRD D5
	p, terminal, err := m.resolver.ResolveAccessKey(r.Context(), hash)
	if err != nil {
		return nil, false
	}
	if p == nil {
		return nil, terminal
	}
	// Rate-limit per access key (PRD-Q5): one misbehaving key must not
	// throttle its tenant's dashboard.
	tb := redis.NewTokenBucket(m.rdb)
	if ok, _ := tb.Allow(r.Context(), token, m.rateLimit, m.rateWindow); !ok {
		return nil, false
	}
	return p, true
}

func PrincipalFrom(ctx context.Context) *authz.Principal {
	p, _ := ctx.Value(principalKey).(*authz.Principal)
	return p
}
