package disk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Driver stores objects on the local VPS volume. It has no external signer:
// presigned URLs are Bloberry-issued HMAC tokens against the server's own
// /v1/objects/{id}/raw endpoint (domains.md §6.3).
type Driver struct {
	root    string
	baseURL string
	secret  []byte
}

func New(root, baseURL string, secret []byte) (*Driver, error) {
	if err := os.MkdirAll(root, 0o750); err != nil {
		return nil, fmt.Errorf("disk: mkdir root: %w", err)
	}
	return &Driver{root: root, baseURL: baseURL, secret: secret}, nil
}

func (d *Driver) Capabilities() storage.Capabilities {
	return storage.Capabilities{
		Presign:        true, // via Bloberry's own signer
		Multipart:      false,
		RangeRequests:  true,
		StorageClasses: false,
	}
}

func (d *Driver) path(key string) (string, error) {
	// Prevent traversal: keys must be relative, no "..".
	if key == "" || strings.Contains(key, "..") || strings.HasPrefix(key, "/") {
		return "", errors.New("disk: invalid key")
	}
	return filepath.Join(d.root, filepath.FromSlash(key)), nil
}

func (d *Driver) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	p, err := d.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(p), ".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), p)
}

func (d *Driver) Get(ctx context.Context, key string, rng *storage.Range) (io.ReadCloser, *storage.ObjectInfo, error) {
	p, err := d.path(key)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	info := &storage.ObjectInfo{Size: fi.Size(), LastModified: fi.ModTime()}
	if rng != nil {
		if rng.Start < 0 || rng.Start >= fi.Size() {
			f.Close()
			return nil, nil, fmt.Errorf("disk: invalid range start %d size %d", rng.Start, fi.Size())
		}
		end := rng.End
		if end < 0 || end >= fi.Size() {
			end = fi.Size() - 1
		}
		if _, err := f.Seek(rng.Start, io.SeekStart); err != nil {
			f.Close()
			return nil, nil, err
		}
		info.Size = end - rng.Start + 1
		return &limitedReadCloser{rc: f, n: info.Size}, info, nil
	}
	return f, info, nil
}

type limitedReadCloser struct {
	rc io.ReadCloser
	n  int64
}

func (l *limitedReadCloser) Read(p []byte) (int, error) {
	if l.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.n {
		p = p[:l.n]
	}
	n, err := l.rc.Read(p)
	l.n -= int64(n)
	return n, err
}

func (l *limitedReadCloser) Close() error { return l.rc.Close() }

func (d *Driver) Delete(ctx context.Context, keys []string) error {
	for _, key := range keys {
		p, err := d.path(key)
		if err != nil {
			continue
		}
		_ = os.Remove(p)
	}
	return nil
}

func (d *Driver) Stat(ctx context.Context, key string) (*storage.ObjectInfo, error) {
	p, err := d.path(key)
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	return &storage.ObjectInfo{Size: fi.Size(), LastModified: fi.ModTime()}, nil
}

func (d *Driver) HealthCheck(ctx context.Context) error {
	fi, err := os.Stat(d.root)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("disk: root is not a directory")
	}
	return nil
}

// --- Presigner: HMAC token against the server's own raw endpoint. ---

func (d *Driver) PresignGet(ctx context.Context, key string, ttl time.Duration) (*storage.PresignedURL, error) {
	token := d.sign(key, "GET", ttl)
	return &storage.PresignedURL{
		URL:    d.baseURL + "/v1/objects/raw?key=" + url.QueryEscape(key) + "&token=" + token,
		Method: "GET",
	}, nil
}

func (d *Driver) PresignPut(ctx context.Context, key string, ttl time.Duration, size int64) (*storage.PresignedURL, error) {
	token := d.sign(key, "PUT", ttl)
	return &storage.PresignedURL{
		URL:    d.baseURL + "/v1/objects/raw?key=" + url.QueryEscape(key) + "&token=" + token,
		Method: "PUT",
		Headers: map[string]string{
			"X-Bloberry-Size": strconv.FormatInt(size, 10),
		},
	}, nil
}

func (d *Driver) sign(key, method string, ttl time.Duration) string {
	exp := time.Now().Add(ttl).Unix()
	mac := hmac.New(sha256.New, d.secret)
	fmt.Fprintf(mac, "%s\n%s\n%d", key, method, exp)
	return hex.EncodeToString(mac.Sum(nil)) + "." + strconv.FormatInt(exp, 10)
}

// VerifyToken checks an HMAC token against the raw endpoint contract.
func (d *Driver) VerifyToken(key, method, token string) error {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return errors.New("disk: malformed token")
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.New("disk: malformed token")
	}
	if time.Now().Unix() > exp {
		return errors.New("disk: token expired")
	}
	mac := hmac.New(sha256.New, d.secret)
	fmt.Fprintf(mac, "%s\n%s\n%d", key, method, exp)
	want := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(want), []byte(parts[0])) {
		return errors.New("disk: bad signature")
	}
	return nil
}

func (d *Driver) MultipartInit(ctx context.Context, key, contentType string) (string, error) {
	return "", errors.New("disk: multipart not supported")
}

func (d *Driver) MultipartPresignPart(ctx context.Context, key, uploadID string, part int, ttl time.Duration) (*storage.PresignedURL, error) {
	return nil, errors.New("disk: multipart not supported")
}

func (d *Driver) MultipartComplete(ctx context.Context, key, uploadID string, parts []storage.Part) (*storage.ObjectInfo, error) {
	return nil, errors.New("disk: multipart not supported")
}

func (d *Driver) MultipartAbort(ctx context.Context, key, uploadID string) error {
	return errors.New("disk: multipart not supported")
}

var _ storage.Driver = (*Driver)(nil)
