package disk_test

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/fakihariefnoto/bloberry/internal/storage/conformance"
	"github.com/fakihariefnoto/bloberry/internal/storage/disk"
)

func TestConformance(t *testing.T) {
	dir := t.TempDir()
	d, err := disk.New(dir, []byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	conformance.Suite(t, d, "disk-test")
}

func TestDiskPresignToken(t *testing.T) {
	dir := t.TempDir()
	d, err := disk.New(dir, []byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	key := "folder/file.txt"
	u, err := d.PresignGet(context.Background(), key, 2*60*1000*1000*1000)
	if err != nil {
		t.Fatal(err)
	}
	// Extract token from the query string of the (relative) presigned URL.
	parsed, err := url.Parse(u.URL)
	if err != nil {
		t.Fatalf("parse presigned url: %v", err)
	}
	token := parsed.Query().Get("token")
	if token == "" {
		t.Fatal("no token in presigned url")
	}

	if err := d.VerifyToken(key, "GET", token); err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}
	if err := d.VerifyToken(key, "PUT", token); err != nil {
		t.Logf("wrong method rejected (expected)")
	} else {
		t.Error("token valid for wrong method")
	}
	if err := d.VerifyToken(key+"x", "GET", token); err == nil {
		t.Error("token valid for wrong key")
	}
}

func TestDiskPutGetDelete(t *testing.T) {
	dir := t.TempDir()
	d, err := disk.New(dir, []byte("s"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	content := []byte("hello disk")
	if err := d.Put(ctx, "a/b.txt", mustReader(content), int64(len(content)), "text/plain"); err != nil {
		t.Fatal(err)
	}
	// file must live under root
	if _, err := os.Stat(filepath.Join(dir, "a", "b.txt")); err != nil {
		t.Fatalf("file not written to root: %v", err)
	}
	rc, info, err := d.Get(ctx, "a/b.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	if info.Size != int64(len(content)) {
		t.Errorf("size = %d", info.Size)
	}
	if err := d.Delete(ctx, []string{"a/b.txt"}); err != nil {
		t.Fatal(err)
	}
}

func TestDiskRejectsTraversal(t *testing.T) {
	dir := t.TempDir()
	d, err := disk.New(dir, []byte("s"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := d.Put(ctx, "../../etc/passwd", mustReader([]byte("x")), 1, ""); err == nil {
		t.Error("traversal key accepted")
	}
}

func mustReader(b []byte) *readCloser {
	return &readCloser{b: b}
}

type readCloser struct{ b []byte }

func (r *readCloser) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

func (r *readCloser) Close() error { return nil }
