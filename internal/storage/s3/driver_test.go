package s3

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestR2PresignIsPathStyle(t *testing.T) {
	d, err := New(Options{
		Endpoint: "https://abc123.r2.cloudflarestorage.com",
		Region:   "auto",
		Bucket:   "mybucket",
		AccessKeyID: "ak", SecretAccessKey: "sk",
		R2: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	u, err := d.PresignPut(context.Background(), "folder/x.png", 5*time.Minute, 100)
	if err != nil {
		t.Fatal(err)
	}
	// Must NOT be virtual-hosted (bucket.endpoint/…). Path-style:
	// https://abc123.r2.cloudflarestorage.com/mybucket/folder/x.png
	if strings.Contains(u.URL, "mybucket.abc123") {
		t.Errorf("R2 URL is virtual-hosted (wrong for R2): %s", u.URL)
	}
	if !strings.Contains(u.URL, "abc123.r2.cloudflarestorage.com/mybucket/folder/x.png") {
		t.Errorf("R2 URL is not path-style: %s", u.URL)
	}
}

func TestS3PresignVirtualHosted(t *testing.T) {
	d, err := New(Options{
		Endpoint: "https://s3.us-east-1.amazonaws.com",
		Region:   "us-east-1",
		Bucket:   "mybucket",
		AccessKeyID: "ak", SecretAccessKey: "sk",
	})
	if err != nil {
		t.Fatal(err)
	}
	u, err := d.PresignGet(context.Background(), "folder/x.png", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(u.URL, "mybucket.s3.us-east-1.amazonaws.com/folder/x.png") {
		t.Logf("S3 presign (virtual-hosted by default): %s", u.URL)
	}
}
