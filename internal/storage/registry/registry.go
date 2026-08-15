package registry

import (
	"errors"
	"sync"

	"github.com/fakihariefnoto/bloberry/internal/storage"
	"github.com/fakihariefnoto/bloberry/internal/storage/azblob"
	"github.com/fakihariefnoto/bloberry/internal/storage/disk"
	"github.com/fakihariefnoto/bloberry/internal/storage/gcs"
	"github.com/fakihariefnoto/bloberry/internal/storage/oss"
	"github.com/fakihariefnoto/bloberry/internal/storage/s3"
)

// BackendRecord is the subset of a storage_backends document the registry
// needs to construct a driver. Credentials arrive already decrypted.
type BackendRecord struct {
	ID           string
	DriverType   string // s3 | r2 | oss | gcs | azblob | disk
	Config       map[string]interface{}
	Credentials  map[string]interface{}
	Capabilities storage.Capabilities
}

// Registry maps backend_id → constructed storage.Driver. Drivers are built once and
// reused (credentials decrypted only in memory at construction, TRD R7).
type Registry struct {
	mu      sync.RWMutex
	drivers map[string]storage.Driver
	factory Factory
}

type Factory func(record BackendRecord) (storage.Driver, error)

func NewRegistry(f Factory) *Registry {
	return &Registry{
		drivers: map[string]storage.Driver{},
		factory: f,
	}
}

// Register builds (or rebuilds) the driver for a backend.
func (r *Registry) Register(record BackendRecord) (storage.Driver, error) {
	d, err := r.factory(record)
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	r.drivers[record.ID] = d
	r.mu.Unlock()
	return d, nil
}

func (r *Registry) Get(id string) (storage.Driver, error) {
	r.mu.RLock()
	d, ok := r.drivers[id]
	r.mu.RUnlock()
	if !ok {
		return nil, errors.New("storage: backend not constructed")
	}
	return d, nil
}

// All returns every constructed driver (backendID → driver). Used by the disk
// driver's raw HMAC endpoint to find the right driver to verify/stream.
func (r *Registry) All() map[string]storage.Driver {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]storage.Driver, len(r.drivers))
	for id, d := range r.drivers {
		out[id] = d
	}
	return out
}

func (r *Registry) Remove(id string) {
	r.mu.Lock()
	delete(r.drivers, id)
	r.mu.Unlock()
}

// DefaultFactory builds the concrete driver for a record.
// secret is the disk-driver HMAC signer key (domains.md §6.3).
func DefaultFactory(secret []byte) Factory {
	return func(record BackendRecord) (storage.Driver, error) {
		switch record.DriverType {
		case "disk":
			root := str(record.Config, "root", "/var/lib/bloberry/objects")
			return disk.New(root, secret)
		case "s3", "r2":
			return s3.New(s3.Options{
				Endpoint:        str(record.Config, "endpoint", ""),
				Region:          str(record.Config, "region", ""),
				Bucket:          str(record.Config, "bucket", ""),
				Prefix:          str(record.Config, "prefix", ""),
				AccessKeyID:     str(record.Credentials, "access_key_id", ""),
				SecretAccessKey: str(record.Credentials, "secret_access_key", ""),
				R2:              record.DriverType == "r2",
				UsePathStyle:    boolv(record.Config, "use_path_style", false),
			})
		case "oss":
			return oss.New(oss.Options{
				Endpoint:   str(record.Config, "endpoint", ""),
				Bucket:     str(record.Config, "bucket", ""),
				Prefix:     str(record.Config, "prefix", ""),
				AccessKeyID: str(record.Credentials, "access_key_id", ""),
				SecretAccessKey: str(record.Credentials, "secret_access_key", ""),
			})
		case "gcs":
			return gcs.New(gcs.Options{
				Bucket:             str(record.Config, "bucket", ""),
				Prefix:             str(record.Config, "prefix", ""),
				ServiceAccountJSON: str(record.Credentials, "service_account_json", ""),
			})
		case "azblob":
			return azblob.New(azblob.Options{
				AccountName:   str(record.Config, "account_name", ""),
				Container:     str(record.Config, "container", ""),
				Prefix:        str(record.Config, "prefix", ""),
				SharedKey:     str(record.Credentials, "shared_key", ""),
				Endpoint:      str(record.Config, "endpoint", ""),
			})
		default:
			return nil, errors.New("storage: unknown driver " + record.DriverType)
		}
	}
}

func str(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

func boolv(m map[string]interface{}, key string, def bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}
