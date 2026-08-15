package azblob

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

		"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	azservice "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"

	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Driver serves Azure Blob Storage. Container-scoped, SharedKey or AAD auth,
// block-blob staging for multipart (domains.md §6.2).
type Driver struct {
	client    *container.Client
	prefix    string
	account   string
	container string
	sharedKey string
	service   string
}

type Options struct {
	Endpoint    string
	AccountName string
	Container   string
	Prefix      string
	SharedKey   string
}

func New(opts Options) (*Driver, error) {
	serviceURL := opts.Endpoint
	if serviceURL == "" {
		serviceURL = "https://" + opts.AccountName + ".blob.core.windows.net"
	}

	var c *container.Client
	if opts.SharedKey != "" {
		cred, err := azblob.NewSharedKeyCredential(opts.AccountName, opts.SharedKey)
		if err != nil {
			return nil, fmt.Errorf("azblob: %w", err)
		}
		svc, err := azservice.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
		if err != nil {
			return nil, fmt.Errorf("azblob: %w", err)
		}
		c = svc.NewContainerClient(opts.Container)
	} else {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, fmt.Errorf("azblob: %w", err)
		}
		svc, err := azservice.NewClient(serviceURL, cred, nil)
		if err != nil {
			return nil, fmt.Errorf("azblob: %w", err)
		}
		c = svc.NewContainerClient(opts.Container)
	}
	return &Driver{client: c, prefix: opts.Prefix, account: opts.AccountName, container: opts.Container, sharedKey: opts.SharedKey, service: serviceURL}, nil
}

func (d *Driver) key(k string) string { return d.prefix + k }

func (d *Driver) Capabilities() storage.Capabilities {
	return storage.Capabilities{
		Presign:          true, // SharedKey/SAS
		Multipart:        true, // block-blob staging
		MinPartSize:      0,
		MaxPartCount:     50000,
		StorageClasses:   true,
		ServerSideCopy:   true,
		RangeRequests:    true,
		ObjectAttributes: true,
	}
}

func (d *Driver) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	bb := d.client.NewBlockBlobClient(d.key(key))
	opts := &blockblob.UploadStreamOptions{}
	if contentType != "" {
		opts.HTTPHeaders = &blob.HTTPHeaders{BlobContentType: to.Ptr(contentType)}
	}
	_, err := bb.UploadStream(ctx, r, opts)
	return err
}

func (d *Driver) Get(ctx context.Context, key string, rng *storage.Range) (io.ReadCloser, *storage.ObjectInfo, error) {
	bb := d.client.NewBlockBlobClient(d.key(key))
	opts := &blob.DownloadStreamOptions{}
	if rng != nil {
		opts.Range = blob.HTTPRange{Offset: rng.Start, Count: rng.End - rng.Start + 1}
	}
	resp, err := bb.DownloadStream(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	info := &storage.ObjectInfo{LastModified: time.Now()}
	if resp.ContentLength != nil {
		info.Size = *resp.ContentLength
	}
	if resp.ContentType != nil {
		info.ContentType = *resp.ContentType
	}
	return resp.Body, info, nil
}

func (d *Driver) Delete(ctx context.Context, keys []string) error {
	for _, k := range keys {
		bc := d.client.NewBlobClient(d.key(k))
		if _, err := bc.Delete(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}

func (d *Driver) Stat(ctx context.Context, key string) (*storage.ObjectInfo, error) {
	bb := d.client.NewBlockBlobClient(d.key(key))
	resp, err := bb.GetProperties(ctx, nil)
	if err != nil {
		return nil, err
	}
	info := &storage.ObjectInfo{LastModified: time.Now()}
	if resp.ContentLength != nil {
		info.Size = *resp.ContentLength
	}
	if resp.ContentType != nil {
		info.ContentType = *resp.ContentType
	}
	return info, nil
}

func (d *Driver) HealthCheck(ctx context.Context) error {
	_, err := d.client.GetProperties(ctx, nil)
	return err
}

func (d *Driver) PresignGet(ctx context.Context, key string, ttl time.Duration) (*storage.PresignedURL, error) {
	return d.sasURL(key, ttl, sas.BlobPermissions{Read: true})
}

func (d *Driver) PresignPut(ctx context.Context, key string, ttl time.Duration, size int64) (*storage.PresignedURL, error) {
	u, err := d.sasURL(key, ttl, sas.BlobPermissions{Write: true, Create: true})
	if err != nil {
		return nil, err
	}
	u.Headers = map[string]string{"Content-Length": strconv.FormatInt(size, 10)}
	return u, nil
}

func (d *Driver) sasURL(key string, ttl time.Duration, perms sas.BlobPermissions) (*storage.PresignedURL, error) {
	if d.sharedKey == "" {
		return nil, errors.New("azblob: presign requires SharedKey (AAD path deferred)")
	}
	cred, err := azblob.NewSharedKeyCredential(d.account, d.sharedKey)
	if err != nil {
		return nil, err
	}
	bb := d.client.NewBlockBlobClient(d.key(key))
	values := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().Add(-time.Minute),
		ExpiryTime:    time.Now().Add(ttl),
		Permissions:   perms.String(),
		ContainerName: d.container,
		BlobName:      d.key(key),
	}
	qp, err := values.SignWithSharedKey(cred)
	if err != nil {
		return nil, err
	}
	return &storage.PresignedURL{URL: bb.URL() + "?" + qp.Encode(), Method: "GET"}, nil
}

func (d *Driver) MultipartInit(ctx context.Context, key, contentType string) (string, error) {
	// Azure uses block IDs, not upload IDs — return the blob key as a handle.
	return d.key(key), nil
}

func (d *Driver) MultipartPresignPart(ctx context.Context, key, uploadID string, part int, ttl time.Duration) (*storage.PresignedURL, error) {
	return nil, errors.New("azblob: presigned part URLs unsupported (use SDK block staging)")
}

func (d *Driver) MultipartComplete(ctx context.Context, key, uploadID string, parts []storage.Part) (*storage.ObjectInfo, error) {
	return nil, errors.New("azblob: use block staging via Put")
}

func (d *Driver) MultipartAbort(ctx context.Context, key, uploadID string) error {
	return nil
}

var _ storage.Driver = (*Driver)(nil)
