
package crypto

import (
	"strings"
	"testing"
)

func TestNewStorageID(t *testing.T) {
	cases := map[string]string{
		"s3": "S3-", "r2": "R2-", "oss": "OSS", "gcs": "GCS", "azblob": "AZB", "disk": "DSK", "weird": "STO",
	}
	for driver, wantPrefix := range cases {
		id := NewStorageID(driver)
		if len(id) != 20 {
			t.Errorf("%s: got length %d, want 20", driver, len(id))
		}
		if !strings.HasPrefix(id, wantPrefix) {
			t.Errorf("%s: got prefix %q, want %q", driver, id[:3], wantPrefix)
		}
	}
	// uniqueness
	a, b := NewStorageID("disk"), NewStorageID("disk")
	if a == b {
		t.Error("two storage ids collided")
	}
}
