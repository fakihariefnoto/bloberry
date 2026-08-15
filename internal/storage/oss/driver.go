package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Driver serves Alibaba OSS. Its own signature version, separate SDK
// (domains.md §6.2).
type Driver struct {
	client oss.Client
	bucket oss.Bucket
	prefix string
}

type Options struct {
	Endpoint        string
	Bucket          string
	Prefix          string
	AccessKeyID     string
	SecretAccessKey string
}

func New(opts Options) (*Driver, error) {
	client, err := oss.New(opts.Endpoint, opts.AccessKeyID, opts.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("oss: %w", err)
	}
	bucket, err := client.Bucket(opts.Bucket)
	if err != nil {
		return nil, fmt.Errorf("oss: %w", err)
	}
	return &Driver{client: *client, bucket: *bucket, prefix: opts.Prefix}, nil
}

func (d *Driver) key(k string) string { return d.prefix + k }

func (d *Driver) Capabilities() storage.Capabilities {
	return storage.Capabilities{
		Presign:          true,
		Multipart:        true,
		MinPartSize:      100 * 1024,
		MaxPartCount:     10000,
		StorageClasses:   true,
		ServerSideCopy:   true,
		RangeRequests:    true,
		ObjectAttributes: true,
	}
}

func (d *Driver) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	var opts []oss.Option
	if contentType != "" {
		opts = append(opts, oss.ContentType(contentType))
	}
	return d.bucket.PutObject(d.key(key), r, opts...)
}

func (d *Driver) Get(ctx context.Context, key string, rng *storage.Range) (io.ReadCloser, *storage.ObjectInfo, error) {
	var opts []oss.Option
	if rng != nil {
		opts = append(opts, oss.Range(rng.Start, rng.End))
	}
	body, err := d.bucket.GetObject(d.key(key), opts...)
	if err != nil {
		return nil, nil, err
	}
	info := &storage.ObjectInfo{LastModified: time.Now()}
	return body, info, nil
}

func (d *Driver) Delete(ctx context.Context, keys []string) error {
	for _, k := range keys {
		if err := d.bucket.DeleteObject(d.key(k)); err != nil {
			return err
		}
	}
	return nil
}

func (d *Driver) Stat(ctx context.Context, key string) (*storage.ObjectInfo, error) {
	props, err := d.bucket.GetObjectMeta(d.key(key))
	if err != nil {
		return nil, err
	}
	size, _ := strconv.ParseInt(props.Get("Content-Length"), 10, 64)
	return &storage.ObjectInfo{Size: size, LastModified: time.Now()}, nil
}

func (d *Driver) HealthCheck(ctx context.Context) error {
	_, err := d.client.GetBucketInfo(d.bucket.BucketName)
	return err
}

func (d *Driver) PresignGet(ctx context.Context, key string, ttl time.Duration) (*storage.PresignedURL, error) {
	u, err := d.bucket.SignURL(d.key(key), oss.HTTPGet, int64(ttl.Seconds()))
	if err != nil {
		return nil, err
	}
	return &storage.PresignedURL{URL: u, Method: "GET"}, nil
}

func (d *Driver) PresignPut(ctx context.Context, key string, ttl time.Duration, size int64) (*storage.PresignedURL, error) {
	u, err := d.bucket.SignURL(d.key(key), oss.HTTPPut, int64(ttl.Seconds()))
	if err != nil {
		return nil, err
	}
	return &storage.PresignedURL{URL: u, Method: "PUT", Headers: map[string]string{"Content-Length": strconv.FormatInt(size, 10)}}, nil
}

func (d *Driver) MultipartInit(ctx context.Context, key, contentType string) (string, error) {
	init, err := d.bucket.InitiateMultipartUpload(d.key(key), oss.ContentType(contentType))
	if err != nil {
		return "", err
	}
	return init.UploadID, nil
}

func (d *Driver) MultipartPresignPart(ctx context.Context, key, uploadID string, part int, ttl time.Duration) (*storage.PresignedURL, error) {
	// OSS v3 has no presigned-part URL helper; the blob stores the uploadID
	// in the object record, and parts are staged via UploadPart. Presigned
	// part URLs are the S3 shape; return a descriptive error.
	return nil, errors.New("oss: presigned part URLs unsupported (stage parts via UploadPart)")
}

func (d *Driver) MultipartComplete(ctx context.Context, key, uploadID string, parts []storage.Part) (*storage.ObjectInfo, error) {
	imur := oss.InitiateMultipartUploadResult{Bucket: d.bucket.BucketName, Key: d.key(key), UploadID: uploadID}
	completed := make([]oss.UploadPart, 0, len(parts))
	for _, p := range parts {
		completed = append(completed, oss.UploadPart{PartNumber: p.Number, ETag: p.ETag})
	}
	_, err := d.bucket.CompleteMultipartUpload(imur, completed)
	if err != nil {
		return nil, err
	}
	return &storage.ObjectInfo{Size: 0}, nil
}

func (d *Driver) MultipartAbort(ctx context.Context, key, uploadID string) error {
	imur := oss.InitiateMultipartUploadResult{Bucket: d.bucket.BucketName, Key: d.key(key), UploadID: uploadID}
	return d.bucket.AbortMultipartUpload(imur)
}

var _ storage.Driver = (*Driver)(nil)
var _ = errors.New
