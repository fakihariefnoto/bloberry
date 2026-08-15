package gcs

import (
	"context"
		"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	gstorage "cloud.google.com/go/storage"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Driver serves Google Cloud Storage. Credential shape differs from the
// other drivers: a service-account signer, not a static key pair (TRD R1).
type Driver struct {
	client *gstorage.Client
	bucket string
	prefix string
	saJSON []byte
}

type Options struct {
	Bucket             string
	Prefix             string
	ServiceAccountJSON string
}

func New(opts Options) (*Driver, error) {
	var clientOpts []option.ClientOption
	if opts.ServiceAccountJSON != "" {
		clientOpts = append(clientOpts, option.WithCredentialsJSON([]byte(opts.ServiceAccountJSON)))
	}
	client, err := gstorage.NewClient(context.Background(), clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("gcs: %w", err)
	}
	return &Driver{client: client, bucket: opts.Bucket, prefix: opts.Prefix, saJSON: []byte(opts.ServiceAccountJSON)}, nil
}

func (d *Driver) key(k string) string { return d.prefix + k }

func (d *Driver) Capabilities() storage.Capabilities {
	return storage.Capabilities{
		Presign:          true, // service-account signer
		Multipart:        false, // GCS has no server-side multipart; resume via XML-compat (out of v1 scope)
		MinPartSize:      0,
		MaxPartCount:     0,
		StorageClasses:   true,
		ServerSideCopy:   true,
		RangeRequests:    true,
		ObjectAttributes: true,
	}
}

func (d *Driver) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	w := d.client.Bucket(d.bucket).Object(d.key(key)).NewWriter(ctx)
	if contentType != "" {
		w.ContentType = contentType
	}
	if _, err := io.Copy(w, r); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func (d *Driver) Get(ctx context.Context, key string, rng *storage.Range) (io.ReadCloser, *storage.ObjectInfo, error) {
	obj := d.client.Bucket(d.bucket).Object(d.key(key))
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, nil, err
	}
	r, err := obj.NewRangeReader(ctx, rng.Start, rng.End-rng.Start+1)
	if err != nil {
		return nil, nil, err
	}
	info := &storage.ObjectInfo{Size: r.Remain(), LastModified: attrs.Updated}
	return r, info, nil
}

func (d *Driver) Delete(ctx context.Context, keys []string) error {
	for _, k := range keys {
		if err := d.client.Bucket(d.bucket).Object(d.key(k)).Delete(ctx); err != nil && err != gstorage.ErrObjectNotExist {
			return err
		}
	}
	return nil
}

func (d *Driver) Stat(ctx context.Context, key string) (*storage.ObjectInfo, error) {
	attrs, err := d.client.Bucket(d.bucket).Object(d.key(key)).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return &storage.ObjectInfo{Size: attrs.Size, LastModified: attrs.Updated, ContentType: attrs.ContentType}, nil
}

func (d *Driver) HealthCheck(ctx context.Context) error {
	it := d.client.Bucket(d.bucket).Objects(ctx, nil)
	_, err := it.Next()
	if err == iterator.Done {
		return nil
	}
	return err
}

func (d *Driver) PresignGet(ctx context.Context, key string, ttl time.Duration) (*storage.PresignedURL, error) {
	opts := &gstorage.SignedURLOptions{Method: "GET", Expires: time.Now().Add(ttl)}
	if len(d.saJSON) > 0 {
		accessID, pk, err := d.signer()
		if err != nil {
			return nil, err
		}
		opts.GoogleAccessID = accessID
		opts.PrivateKey = pk
	}
	u, err := d.client.Bucket(d.bucket).SignedURL(d.key(key), opts)
	if err != nil {
		return nil, err
	}
	return &storage.PresignedURL{URL: u, Method: "GET"}, nil
}

func (d *Driver) PresignPut(ctx context.Context, key string, ttl time.Duration, size int64) (*storage.PresignedURL, error) {
	opts := &gstorage.SignedURLOptions{Method: "PUT", Expires: time.Now().Add(ttl)}
	if len(d.saJSON) > 0 {
		accessID, pk, err := d.signer()
		if err != nil {
			return nil, err
		}
		opts.GoogleAccessID = accessID
		opts.PrivateKey = pk
	}
	u, err := d.client.Bucket(d.bucket).SignedURL(d.key(key), opts)
	if err != nil {
		return nil, err
	}
	return &storage.PresignedURL{URL: u, Method: "PUT"}, nil
}

// signer extracts the GoogleAccessID + PEM private key from the service-account
// JSON for presigning (the IAM signBlob path is deferred — the key-file case
// is the common one).
func (d *Driver) signer() (string, []byte, error) {
	var sa struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	if err := json.Unmarshal(d.saJSON, &sa); err != nil {
		return "", nil, fmt.Errorf("gcs: parse service account: %w", err)
	}
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sa.PrivateKey)); err != nil {
		return "", nil, fmt.Errorf("gcs: parse private key: %w", err)
	}
	return sa.ClientEmail, []byte(sa.PrivateKey), nil
}

func (d *Driver) MultipartInit(ctx context.Context, key, contentType string) (string, error) {
	return "", errors.New("gcs: multipart not supported in v1")
}

func (d *Driver) MultipartPresignPart(ctx context.Context, key, uploadID string, part int, ttl time.Duration) (*storage.PresignedURL, error) {
	return nil, errors.New("gcs: multipart not supported in v1")
}

func (d *Driver) MultipartComplete(ctx context.Context, key, uploadID string, parts []storage.Part) (*storage.ObjectInfo, error) {
	return nil, errors.New("gcs: multipart not supported in v1")
}

func (d *Driver) MultipartAbort(ctx context.Context, key, uploadID string) error {
	return errors.New("gcs: multipart not supported in v1")
}

var _ storage.Driver = (*Driver)(nil)
