package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Driver serves AWS S3, Cloudflare R2 (endpoint override), MinIO, B2,
// Spaces and Wasabi. R2 capability flags differ — it is not fully S3
// compatible (TRD R2).
type Driver struct {
	client       *s3.Client
	bucket       string
	prefix       string
	capabilities storage.Capabilities
}

type Options struct {
	Endpoint        string
	Region          string
	Bucket          string
	Prefix          string
	AccessKeyID     string
	SecretAccessKey string
	R2              bool // capability deltas for Cloudflare R2
	UsePathStyle    bool
}

func New(opts Options) (*Driver, error) {
	region := orDefault(opts.Region, "us-east-1")
	if opts.R2 {
		// R2 requires the "auto" region in SigV4; a fixed region breaks the
		// signature on presigned requests.
		region = "auto"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(opts.AccessKeyID, opts.SecretAccessKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("s3: %w", err)
	}
	// R2 does not support virtual-hosted-style URLs — presigned URLs must be
	// path-style against the account endpoint (<bucket>/<key>), otherwise the
	// signature is built for the wrong host and the provider returns 403.
	usePathStyle := opts.UsePathStyle || opts.R2
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if opts.Endpoint != "" {
			o.BaseEndpoint = aws.String(opts.Endpoint)
			o.UsePathStyle = usePathStyle
		}
	})
	caps := storage.Capabilities{
		Presign:          true,
		Multipart:        true,
		MinPartSize:      5 * 1024 * 1024,
		MaxPartCount:     10000,
		StorageClasses:   true,
		ServerSideCopy:   true,
		RangeRequests:    true,
		ObjectAttributes: true,
	}
	if opts.R2 {
		caps.StorageClasses = false
		caps.ObjectAttributes = false
		caps.MinPartSize = 10 * 1024 * 1024
	}
	return &Driver{client: client, bucket: opts.Bucket, prefix: opts.Prefix, capabilities: caps}, nil
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func (d *Driver) key(k string) string { return d.prefix + k }
func (d *Driver) Capabilities() storage.Capabilities {
	return d.capabilities
}

func (d *Driver) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	input := &s3.PutObjectInput{Bucket: aws.String(d.bucket), Key: aws.String(d.key(key)), Body: r}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	_, err := d.client.PutObject(ctx, input)
	return wrapErr(err)
}

func (d *Driver) Get(ctx context.Context, key string, rng *storage.Range) (io.ReadCloser, *storage.ObjectInfo, error) {
	input := &s3.GetObjectInput{Bucket: aws.String(d.bucket), Key: aws.String(d.key(key))}
	if rng != nil {
		input.Range = aws.String(fmt.Sprintf("bytes=%d-%d", rng.Start, rng.End))
	}
	out, err := d.client.GetObject(ctx, input)
	if err != nil {
		return nil, nil, wrapErr(err)
	}
	info := &storage.ObjectInfo{Size: *out.ContentLength, ETag: *out.ETag, LastModified: time.Now()}
	if out.ContentType != nil {
		info.ContentType = *out.ContentType
	}
	if out.LastModified != nil {
		info.LastModified = *out.LastModified
	}
	return out.Body, info, nil
}

func (d *Driver) Delete(ctx context.Context, keys []string) error {
	for i := 0; i < len(keys); i += 1000 {
		batch := keys[i:min(i+1000, len(keys))]
		objects := make([]types.ObjectIdentifier, 0, len(batch))
		for _, k := range batch {
			objects = append(objects, types.ObjectIdentifier{Key: aws.String(d.key(k))})
		}
		_, err := d.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(d.bucket),
			Delete: &types.Delete{Objects: objects},
		})
		if err != nil {
			return wrapErr(err)
		}
	}
	return nil
}

func (d *Driver) Stat(ctx context.Context, key string) (*storage.ObjectInfo, error) {
	head := &s3.HeadObjectInput{Bucket: aws.String(d.bucket), Key: aws.String(d.key(key))}
	out, err := d.client.HeadObject(ctx, head)
	if err != nil {
		return nil, wrapErr(err)
	}
	info := &storage.ObjectInfo{Size: *out.ContentLength, ETag: *out.ETag, LastModified: time.Now()}
	if out.LastModified != nil {
		info.LastModified = *out.LastModified
	}
	if out.ContentType != nil {
		info.ContentType = *out.ContentType
	}
	return info, nil
}

func (d *Driver) HealthCheck(ctx context.Context) error {
	_, err := d.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(d.bucket)})
	return wrapErr(err)
}

func (d *Driver) PresignGet(ctx context.Context, key string, ttl time.Duration) (*storage.PresignedURL, error) {
	input := &s3.GetObjectInput{Bucket: aws.String(d.bucket), Key: aws.String(d.key(key))}
	return presignGet(ctx, d.client, input, ttl)
}

func (d *Driver) PresignPut(ctx context.Context, key string, ttl time.Duration, size int64) (*storage.PresignedURL, error) {
	input := &s3.PutObjectInput{Bucket: aws.String(d.bucket), Key: aws.String(d.key(key))}
	url, err := s3.NewPresignClient(d.client).PresignPutObject(ctx, input, s3.WithPresignExpires(ttl))
	if err != nil {
		return nil, wrapErr(err)
	}
	return &storage.PresignedURL{URL: url.URL, Method: "PUT", Headers: map[string]string{"Content-Length": strconv.FormatInt(size, 10)}}, nil
}

func presignGet(ctx context.Context, client *s3.Client, input *s3.GetObjectInput, ttl time.Duration) (*storage.PresignedURL, error) {
	url, err := s3.NewPresignClient(client).PresignGetObject(ctx, input, s3.WithPresignExpires(ttl))
	if err != nil {
		return nil, wrapErr(err)
	}
	return &storage.PresignedURL{URL: url.URL, Method: "GET"}, nil
}

func (d *Driver) MultipartInit(ctx context.Context, key, contentType string) (string, error) {
	input := &s3.CreateMultipartUploadInput{Bucket: aws.String(d.bucket), Key: aws.String(d.key(key))}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	out, err := d.client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return "", wrapErr(err)
	}
	return *out.UploadId, nil
}

func (d *Driver) MultipartPresignPart(ctx context.Context, key, uploadID string, part int, ttl time.Duration) (*storage.PresignedURL, error) {
	input := &s3.UploadPartInput{
		Bucket:     aws.String(d.bucket),
		Key:        aws.String(d.key(key)),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int32(int32(part)),
	}
	url, err := s3.NewPresignClient(d.client).PresignUploadPart(ctx, input, s3.WithPresignExpires(ttl))
	if err != nil {
		return nil, wrapErr(err)
	}
	return &storage.PresignedURL{URL: url.URL, Method: "PUT"}, nil
}

func (d *Driver) MultipartComplete(ctx context.Context, key, uploadID string, parts []storage.Part) (*storage.ObjectInfo, error) {
	completed := make([]types.CompletedPart, 0, len(parts))
	for _, p := range parts {
		completed = append(completed, types.CompletedPart{
			ETag:       aws.String(strings.Trim(p.ETag, `"`)),
			PartNumber: aws.Int32(int32(p.Number)),
		})
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(d.bucket),
		Key:      aws.String(d.key(key)),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{Parts: completed},
	}
	out, err := d.client.CompleteMultipartUpload(ctx, input)
	if err != nil {
		return nil, wrapErr(err)
	}
	return &storage.ObjectInfo{Size: 0, ETag: *out.ETag}, nil
}

func (d *Driver) MultipartAbort(ctx context.Context, key, uploadID string) error {
	_, err := d.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(d.bucket),
		Key:      aws.String(d.key(key)),
		UploadId: aws.String(uploadID),
	})
	return wrapErr(err)
}

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	return errors.New(err.Error())
}

var _ storage.Driver = (*Driver)(nil)
