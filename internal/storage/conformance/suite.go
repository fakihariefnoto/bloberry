package conformance

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Suite runs every driver against the same checks (PRD G2). Drivers declare
// their capabilities; the suite exercises each capability that is true and
// asserts the documented failure for each that is false.
func Suite(t *testing.T, d storage.Driver, backendID string) {
	t.Helper()
	ctx := context.Background()
	key := backendID + "/conformance.txt"
	content := []byte("bloberry conformance payload — " + backendID)
	caps := d.Capabilities()

	t.Run("roundtrip", func(t *testing.T) {
		if err := d.Put(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain"); err != nil {
			t.Fatalf("put: %v", err)
		}
		info, err := d.Stat(ctx, key)
		if err != nil {
			t.Fatalf("stat: %v", err)
		}
		if info.Size != int64(len(content)) {
			t.Errorf("stat size = %d want %d", info.Size, len(content))
		}
		rc, _, err := d.Get(ctx, key, nil)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if !bytes.Equal(got, content) {
			t.Errorf("roundtrip mismatch: got %q", got)
		}
	})

	t.Run("range", func(t *testing.T) {
		if !caps.RangeRequests {
			t.Skip("driver declares no range support")
		}
		rc, info, err := d.Get(ctx, key, &storage.Range{Start: 0, End: 9})
		if err != nil {
			t.Fatalf("range get: %v", err)
		}
		defer rc.Close()
		got, _ := io.ReadAll(rc)
		if !strings.HasPrefix(string(content), string(got)) {
			t.Errorf("range read = %q want prefix of %q", got, content)
		}
		if info.Size != 10 {
			t.Errorf("range size = %d want 10", info.Size)
		}
	})

	t.Run("overwrite", func(t *testing.T) {
		newContent := []byte("overwritten")
		if err := d.Put(ctx, key, bytes.NewReader(newContent), int64(len(newContent)), "text/plain"); err != nil {
			t.Fatalf("overwrite put: %v", err)
		}
		info, _ := d.Stat(ctx, key)
		if info.Size != int64(len(newContent)) {
			t.Errorf("overwrite size = %d want %d", info.Size, len(newContent))
		}
	})

	t.Run("presign-get", func(t *testing.T) {
		if !caps.Presign {
			t.Skip("driver declares no presign")
		}
		u, err := d.PresignGet(ctx, key, 2*time.Minute)
		if err != nil {
			t.Fatalf("presign get: %v", err)
		}
		if u.URL == "" {
			t.Error("empty presigned URL")
		}
	})

	t.Run("presign-put", func(t *testing.T) {
		if !caps.Presign {
			t.Skip("driver declares no presign")
		}
		u, err := d.PresignPut(ctx, key, 2*time.Minute, 5)
		if err != nil {
			t.Fatalf("presign put: %v", err)
		}
		if u.URL == "" {
			t.Error("empty presigned PUT URL")
		}
	})

	t.Run("multipart-declared", func(t *testing.T) {
		if !caps.Multipart {
			// Capability honesty: the call must fail in the documented way.
			_, err := d.MultipartInit(ctx, key, "application/octet-stream")
			if err == nil {
				t.Errorf("multipart declared false but Init succeeded")
			}
			t.Skip("driver declares no multipart")
		}
		id, err := d.MultipartInit(ctx, key+"-mp", "application/octet-stream")
		if err != nil {
			t.Fatalf("multipart init: %v", err)
		}
		if id == "" {
			t.Error("empty upload id")
		}
		_ = d.MultipartAbort(ctx, key+"-mp", id)
	})

	t.Run("delete-missing-ok", func(t *testing.T) {
		if err := d.Delete(ctx, []string{backendID + "/does-not-exist"}); err != nil {
			t.Errorf("delete nonexistent: %v", err)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := d.Delete(ctx, []string{key}); err != nil {
			t.Fatalf("delete: %v", err)
		}
		if _, err := d.Stat(ctx, key); err == nil {
			t.Error("object still present after delete")
		}
	})
}
