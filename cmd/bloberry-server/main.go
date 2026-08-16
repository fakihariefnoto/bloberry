package main

import (
	"context"
	"errors"
	"fmt"
	"bytes"
	"io"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adminrepo "github.com/fakihariefnoto/bloberry/internal/admin/repository"
	adminuc "github.com/fakihariefnoto/bloberry/internal/admin/usecase"
	"github.com/fakihariefnoto/bloberry/internal/api"
	apikeyrepo "github.com/fakihariefnoto/bloberry/internal/apikey/repository"
	apikeyuc "github.com/fakihariefnoto/bloberry/internal/apikey/usecase"
	auditrepo "github.com/fakihariefnoto/bloberry/internal/audit/repository"
	audituc "github.com/fakihariefnoto/bloberry/internal/audit/usecase"
	"github.com/fakihariefnoto/bloberry/internal/auth"
	authrepo "github.com/fakihariefnoto/bloberry/internal/auth/repository"
	authuc "github.com/fakihariefnoto/bloberry/internal/auth/usecase"
	"github.com/fakihariefnoto/bloberry/internal/authz"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	folderrepo "github.com/fakihariefnoto/bloberry/internal/folder/repository"
	folderuc "github.com/fakihariefnoto/bloberry/internal/folder/usecase"
	grantrepo "github.com/fakihariefnoto/bloberry/internal/grant/repository"
	grantuc "github.com/fakihariefnoto/bloberry/internal/grant/usecase"
	"github.com/fakihariefnoto/bloberry/internal/job"
	jobrepo "github.com/fakihariefnoto/bloberry/internal/job/repository"
	jobuc "github.com/fakihariefnoto/bloberry/internal/job/usecase"
	objectrepo "github.com/fakihariefnoto/bloberry/internal/object/repository"
	objectuc "github.com/fakihariefnoto/bloberry/internal/object/usecase"
	server "github.com/fakihariefnoto/bloberry/internal/platform/api"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/db"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/platform/jwtutil"
	"github.com/fakihariefnoto/bloberry/internal/platform/mailer"
	"github.com/fakihariefnoto/bloberry/internal/platform/session"
	"github.com/fakihariefnoto/bloberry/internal/platform/web"
	sharerepo "github.com/fakihariefnoto/bloberry/internal/share/repository"
	shareuc "github.com/fakihariefnoto/bloberry/internal/share/usecase"
	setuprepo "github.com/fakihariefnoto/bloberry/internal/setup/repository"
	setupuc "github.com/fakihariefnoto/bloberry/internal/setup/usecase"
	"github.com/fakihariefnoto/bloberry/internal/storage/disk"
	"github.com/fakihariefnoto/bloberry/internal/storage/registry"
	tenantrepo "github.com/fakihariefnoto/bloberry/internal/tenant/repository"
	tenantuc "github.com/fakihariefnoto/bloberry/internal/tenant/usecase"
	usagerepo "github.com/fakihariefnoto/bloberry/internal/usage/repository"
	usageuc "github.com/fakihariefnoto/bloberry/internal/usage/usecase"
	userrepo "github.com/fakihariefnoto/bloberry/internal/user/repository"
	useruc "github.com/fakihariefnoto/bloberry/internal/user/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5/middleware"
	"path"
	"strings"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	cfg, err := configLoad()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// datastores
	mdb, err := db.Connect(ctx, cfg.MongoURI, cfg.Database)
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	defer mdb.Disconnect(context.Background())

	rdb, err := redisConnect(ctx, cfg.RedisURI)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	if err := db.Migrate(ctx, mdb.DB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// wire domains
	usrs := userrepo.New(mdb.DB)
	userUC := useruc.NewUsecase(usrs)

	authCache := authz.NewCache(rdb)
	grantRepo := grantrepo.New(mdb.DB)
	authzLoader := &authz.Loader{
		Cache:  authCache,
		Users:  usrs,
		Keys:   apikeyrepo.New(mdb.DB),
		Grant:  grantRepo,
		Member: tenantrepo.New(mdb.DB),
	}

	issuer := jwtutil.NewIssuer([]byte(cfg.JWTSecret))
	sess := session.NewStore(rdb)
	mail := mailer.Noop{Log: func(msg string) { log.Println(msg) }}

	authRepo := authrepo.New(mdb.DB)
	authUC := authuc.NewUsecase(auth.Deps{
		Repo: authRepo, Users: usrs, Sessions: sess, Tokens: issuer,
		Mailer: mail, Google: &googleVerifier{clientID: cfg.GoogleClientID},
		Envelope: crypto.NewEnvelopeOrPanic(cfg.CredentialEncryptionKey),
		Redis: redisKV{rdb: rdb},
	})

	tenantRepo := tenantrepo.New(mdb.DB)
	tenantUC := tenantuc.NewUsecase(tenantRepo)

	folderRepo := folderrepo.New(mdb.DB)
	folderUC := folderuc.NewUsecase(folderRepo)

	reg := registry.NewRegistry(registry.DefaultFactory([]byte(cfg.CredentialEncryptionKey)))
	regInstance = reg
	objectRepo := objectrepo.New(mdb.DB)
	objectUC := objectuc.NewUsecase(objectuc.Deps{
		Repo: objectRepo, Registry: reg, Quota: tenantUC, Folders: folderUC,
		MaxSize: cfg.MaxObjectSize, PartSize: cfg.MultipartPartSize,
		BaseURL: cfg.BaseURL, RawSecret: []byte(cfg.CredentialEncryptionKey),
	})

	shareUC := shareuc.NewUsecase(shareuc.Deps{Repo: sharerepo.New(mdb.DB), Objects: objectRepo, BaseURL: cfg.BaseURL})

	apikeyUC := apikeyuc.NewUsecase(apikeyuc.Deps{Repo: apikeyrepo.New(mdb.DB), Invalidator: authzLoader})

	grantUC := grantuc.NewUsecase(grantuc.Deps{Repo: grantRepo, Invalidator: authzLoader, Folders: folderRepo})

	auditRepo := auditrepo.New(mdb.DB)
	auditUC := audituc.NewUsecase(auditRepo)

	jobQueue := &redisJobQueue{rdb: rdb}
	jobRunner := job.NewTransferRunner(job.TransferDeps{Objects: objectRepo, Registry: reg})
	jobUC := jobuc.NewUsecase(jobuc.Deps{Repo: jobrepo.New(mdb.DB), Queue: jobQueue, Run: jobRunner})

	usageUC := usageuc.NewUsecase(usageuc.Deps{Repo: usagerepo.New(mdb.DB), Objects: objectRepo})

	adminRepo := adminrepo.New(mdb.DB)
	adminUC := adminuc.NewUsecase(adminuc.Deps{Repo: adminRepo, Registry: reg, Envelope: crypto.NewEnvelopeOrPanic(cfg.CredentialEncryptionKey), Counters: adminRepo, AllTenants: tenantRepo})

	setupRepo := setuprepo.New(mdb.DB)
	setupUC := setupuc.NewUsecase(setupuc.Deps{Repo: setupRepo, DiskRoot: envOr("DISK_STORAGE_PATH", defaultDiskRoot()), Registry: reg})

	// Bootstrap: register all stored backends into the in-memory driver
	// registry at boot so the registry survives restarts (ADR-2).
	if err := bootstrapBackends(ctx, mdb.DB, reg, cfg); err != nil {
		log.Printf("bootstrap backends: %v", err)
	}

	handler := &api.Handler{
		Auth: authUC, Users: userUC, Tenants: tenantUC, Folders: folderUC,
		Objects: objectUC, Shares: shareUC, APIKeys: apikeyUC, Grants: grantUC,
		Jobs: jobUC, Usage: usageUC, Audit: auditUC, Admin: adminUC,
		Setup: setupUC,
		Storage: reg,
	}

	// start worker
	jobUC.StartWorker(ctx)

	// router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Minute))

	apiMw := httpx.NewMiddleware([]byte(cfg.JWTSecret), authzLoader, rdb, cfg.RateLimitRequests, cfg.RateLimitWindow)

	// disk driver's raw HMAC endpoint (domains.md §6.3) — token-gated, not JWT.
	r.Get("/v1/objects/raw", diskRawGet(objectRepo, reg, []byte(cfg.CredentialEncryptionKey)))
	r.Put("/v1/objects/raw", diskRawPut(objectRepo, reg, []byte(cfg.CredentialEncryptionKey), cfg.MaxObjectSize))

	// The generated mux owns all spec routes (public + authed). An auth-gate
	// middleware wraps it: public paths pass through, everything else must
	// authenticate. HandlerFromMux mounts at the mux root (spec paths have no
	// /v1 prefix because the spec's `servers.url` is /v1).
	apiRouter := chi.NewRouter()
	server.HandlerFromMux(handler, apiRouter)
	// Embedded dashboard is the SPA catch-all — anything the API doesn't own.
	if webHandler, err := web.Handler(); err == nil {
		apiRouter.NotFound(webHandler.ServeHTTP)
	} else {
		log.Printf("web embed unavailable: %v", err)
	}
	r.Mount("/", authGate(apiMw, apiRouter))

	log.Printf("bloberry-server listening on %s", cfg.ServerAddr)
	srv := &http.Server{Addr: cfg.ServerAddr, Handler: r}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
}

// publicPaths are the spec endpoints that do not require authentication.
var publicPaths = map[string]bool{
	"GET /health":                    true,
	"GET /setup/status":              true,
	"POST /setup":                    true,
	"GET /s/":                        true,
	"POST /auth/signup":              true,
	"POST /auth/login":               true,
	"POST /auth/refresh":             true,
	"POST /auth/forgot-password":     true,
	"POST /auth/reset-password":      true,
	"POST /auth/otp/request":         true,
	"POST /auth/otp/verify":          true,
	"POST /auth/google":              true,
	"POST /auth/pair/verify":         true,
	"POST /auth/login/verify-totp":   true,
}

// apiPrefixes are the top-level path prefixes the JSON API owns. Entries with a
// trailing slash are matched as prefixes; entries without one are matched
// exactly at the root (e.g. /tenants, /folders, /objects — collection roots).
var apiPrefixes = []string{
	"/v1/", "/auth/", "/users/", "/tenants/", "/tenants", "/folders/", "/folders",
	"/objects/", "/objects", "/shares/", "/shares", "/keys", "/applications",
	"/applications/", "/grants/", "/archives", "/jobs", "/transfers", "/audit",
	"/usage/", "/admin/", "/backends", "/backends/",
}

func isAPIPath(path string) bool {
	if path == "/keys" || path == "/transfers" || path == "/jobs" {
		return true
	}
	for _, p := range apiPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func authGate(mw *httpx.Middleware, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		key := r.Method + " " + path
		// Public endpoints pass through (auth + short links + health/setup).
		if publicPaths[key] || (r.Method == "GET" && strings.HasPrefix(path, "/s/")) {
			next.ServeHTTP(w, r)
			return
		}
		// Only API paths require authentication. Everything else — the SPA
		// shell, /assets, and deep links like /app/files — is served as-is so
		// direct navigation never hits a 401 JSON; the dashboard router handles
		// session gating client-side.
		if !isAPIPath(path) {
			next.ServeHTTP(w, r)
			return
		}
		mw.Authenticate(next).ServeHTTP(w, r)
	})
}

// --- wiring helpers ---

func configLoad() (configT, error) {
	// Load .env if present (never overrides real env vars — godotenv skips set ones).
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return configT{}, fmt.Errorf("load .env: %w", err)
		}
	}
	var c configT
	c.ServerAddr = envOr("SERVER_ADDR", "127.0.0.1:8080")
	c.MongoURI = envOr("MONGODB_URI", "mongodb://127.0.0.1:27017")
	c.RedisURI = envOr("REDIS_URI", "redis://127.0.0.1:6379")
	c.Database = envOr("MONGODB_DATABASE", "bloberry")
	c.JWTSecret = os.Getenv("JWT_SECRET")
	c.CredentialEncryptionKey = os.Getenv("CREDENTIAL_ENCRYPTION_KEY")
	c.GoogleClientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	c.BaseURL = envOr("BLOBERRY_BASE_URL", "http://localhost:8080")
	c.MaxObjectSize = int64Env("MAX_OBJECT_SIZE", 5*1024*1024*1024)
	c.MultipartPartSize = int64Env("MULTIPART_PART_SIZE", 16*1024*1024)
	c.RateLimitRequests = intEnv("RATE_LIMIT_REQUESTS", 1000)
	c.RateLimitWindow = durEnv("RATE_LIMIT_WINDOW", time.Hour)
	if c.JWTSecret == "" {
		return c, errors.New("JWT_SECRET required")
	}
	if c.CredentialEncryptionKey == "" {
		return c, errors.New("CREDENTIAL_ENCRYPTION_KEY required")
	}
	return c, nil
}

type configT struct {
	ServerAddr              string
	MongoURI                string
	RedisURI                string
	Database                string
	JWTSecret               string
	CredentialEncryptionKey string
	GoogleClientID          string
	BaseURL                 string
	MaxObjectSize           int64
	MultipartPartSize       int64
	RateLimitRequests       int
	RateLimitWindow         time.Duration
}

func intEnv(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		n := 0
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return def
}

func int64Env(k string, def int64) int64 {
	if v := os.Getenv(k); v != "" {
		n := int64(0)
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return def
}

func durEnv(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func redisConnect(ctx context.Context, uri string) (*redis.Client, error) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		opts = &redis.Options{Addr: uri}
	}
	c := redis.NewClient(opts)
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return c, nil
}

type redisKV struct{ rdb *redis.Client }

func (k redisKV) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return k.rdb.Set(ctx, key, value, ttl).Err()
}
func (k redisKV) Get(ctx context.Context, key string) (string, error) {
	v, err := k.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("not found")
	}
	return v, err
}
func (k redisKV) Del(ctx context.Context, key string) error { return k.rdb.Del(ctx, key).Err() }
func (k redisKV) Incr(ctx context.Context, key string) (int64, error) {
	return k.rdb.Incr(ctx, key).Result()
}
func (k redisKV) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return k.rdb.Expire(ctx, key, ttl).Err()
}

type redisJobQueue struct{ rdb *redis.Client }

func (q *redisJobQueue) Enqueue(ctx context.Context, jobID string) error {
	return q.rdb.LPush(ctx, "job:queue", jobID).Err()
}
func (q *redisJobQueue) Dequeue(ctx context.Context) (string, error) {
	return q.rdb.RPop(ctx, "job:queue").Result()
}

type googleVerifier struct{ clientID string }

func (g *googleVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.GoogleIdentity, error) {
	if g.clientID == "" {
		return nil, errors.New("google oauth not configured")
	}
	return nil, errors.New("google oauth not configured")
}

// --- disk raw endpoints ---

type rawObjectRepo interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Object, error)
}

// diskRawVerify finds the registered disk driver that verifies the HMAC token
// for this key+method, and returns it. The token was signed by whichever disk
// backend the object belongs to.
func diskRawVerify(r *http.Request, key, method string) (*disk.Driver, bool) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return nil, false
	}
	for _, d := range regInstance.All() {
		dd, ok := d.(*disk.Driver)
		if !ok {
			continue
		}
		if dd.VerifyToken(key, method, token) == nil {
			return dd, true
		}
	}
	return nil, false
}

func diskRawGet(repo rawObjectRepo, reg *registry.Registry, secret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request")
			return
		}
		d, ok := diskRawVerify(r, key, "GET")
		if !ok {
			httpx.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		rc, info, err := d.Get(r.Context(), key, nil)
		if err != nil {
			httpx.Error(w, http.StatusNotFound, "not_found")
			return
		}
		defer rc.Close()
		// The browser saves by Content-Disposition; without it it falls back to
		// the URL's last segment (the storage key). Derive the real filename
		// from the key's basename so share links download with the true name.
		filename := path.Base(key)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+sanitizeFilename(filename)+"\"")
		http.ServeContent(w, r, "", info.LastModified, ioReadSeeker(rc))
	}
}

func diskRawPut(repo rawObjectRepo, reg *registry.Registry, secret []byte, maxSize int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request")
			return
		}
		d, ok := diskRawVerify(r, key, "PUT")
		if !ok {
			httpx.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		httpx.MaxBytesReader(w, r, maxSize)
		if err := d.Put(r.Context(), key, r.Body, r.ContentLength, r.Header.Get("Content-Type")); err != nil {
			log.Printf("raw-put failed: key=%s err=%v", key, err)
			httpx.ErrorWithContent(w, http.StatusBadGateway, "backend_unreachable",
				"failed to write bytes to storage: "+err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

var regInstance *registry.Registry

// sanitizeFilename strips characters that would break a Content-Disposition
// header (quotes, CR/LF injection).
func sanitizeFilename(name string) string {
	name = strings.Map(func(r rune) rune {
		if r == '"' || r == '\\' || r == '\r' || r == '\n' {
			return '_'
		}
		return r
	}, name)
	return name
}

// ioReadSeeker adapts an io.ReadCloser into an io.ReadSeeker for ServeContent
// by buffering. Bounded by the disk driver's proxy path.
func ioReadSeeker(rc io.ReadCloser) io.ReadSeeker {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(rc)
	_ = rc.Close()
	return bytes.NewReader(buf.Bytes())
}

// --- helpers used by wiring ---

type countersRepo interface {
	CountTenants(ctx context.Context) (int64, error)
}

// defaultDiskRoot returns a writable default for the setup-created disk
// backend: /var/lib/bloberry/objects when we can create it, otherwise a
// user-writable path (dev-friendly; a real install sets DISK_STORAGE_PATH).
func defaultDiskRoot() string {
	root := "/var/lib/bloberry/objects"
	if os.MkdirAll(root, 0o750) == nil {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "./bloberry-objects"
	}
	return home + "/.bloberry/objects"
}

// bootstrapBackends registers every stored storage_backend into the driver
// registry at boot, decrypting credentials in memory (TRD R7).
func bootstrapBackends(ctx context.Context, db *mongo.Database, reg *registry.Registry, cfg configT) error {
	env, err := crypto.NewEnvelope(cfg.CredentialEncryptionKey)
	if err != nil {
		return err
	}
	cur, err := db.Collection("storage_backends").Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var b domain.StorageBackend
		if err := cur.Decode(&b); err != nil {
			continue
		}
		credBytes, err := env.Decrypt(b.CredentialsEncrypted)
		if err != nil {
			// Empty encrypted credentials (e.g. a setup-created disk backend)
			// decrypt to nothing — that's not a failure, just no credentials.
			if len(b.CredentialsEncrypted) > 0 {
				log.Printf("bootstrap: SKIP backend %s (%s): decrypt failed", b.ID, b.Driver)
				continue
			}
			credBytes = []byte("{}")
		}
		var creds map[string]interface{}
		_ = json.Unmarshal(credBytes, &creds)
		if _, err := reg.Register(registry.BackendRecord{
			ID: b.ID, DriverType: b.Driver, Config: b.Config, Credentials: creds,
		}); err != nil {
			log.Printf("bootstrap: FAILED to register backend %s (%s): %v", b.ID, b.Driver, err)
			continue
		}
		log.Printf("bootstrap: registered backend %s (%s)", b.ID, b.Driver)
	}
	return cur.Err()
}

