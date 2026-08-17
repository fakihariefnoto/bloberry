package object

import (
	"context"
	"io"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/storage"
)

// Repository persists object records (metadata) and multipart uploads.
type Repository interface {
	Insert(ctx context.Context, o *domain.Object) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.Object, error)
	GetByName(ctx context.Context, tenantID, folderID, name string) (*domain.Object, error)
	Update(ctx context.Context, o *domain.Object) error
	Delete(ctx context.Context, tenantID, id string) error
	SoftDelete(ctx context.Context, tenantID, id string) error
	ListByFolder(ctx context.Context, tenantID, folderID, backendID string) ([]domain.Object, error)
	// ListActiveByBackend returns active objects on a storage engine (for
	// cross-engine transfers). Empty backendID = all engines.
	ListActiveByBackend(ctx context.Context, tenantID, backendID string) ([]domain.Object, error)
	CountActive(ctx context.Context, tenantID string) (int64, error)
	SumActiveBytes(ctx context.Context, tenantID string) (int64, error)
	GetBackend(ctx context.Context, id string) (*domain.StorageBackend, error)

	InsertMultipart(ctx context.Context, m *domain.MultipartUpload) error
	GetMultipart(ctx context.Context, tenantID, objectID string) (*domain.MultipartUpload, error)
	UpdateMultipartParts(ctx context.Context, tenantID, objectID string, parts []domain.PartRec) error
	DeleteMultipart(ctx context.Context, tenantID, objectID string) error

	// GetTenantBackend returns the tenant's default storage backend record.
	GetTenantBackend(ctx context.Context, tenantID string) (*domain.StorageBackend, error)
}

// Reader is the narrow read interface share depends on.
type Reader interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Object, error)
}

// Writer is the narrow interface job depends on.
type Writer interface {
	GetByID(ctx context.Context, tenantID, id string) (*domain.Object, error)
}

// Registry is the narrow storage-registry interface object depends on.
type Registry interface {
	Get(id string) (storage.Driver, error)
}

type PresignResult struct {
	FileID      string            `json:"file_id"`
	UploadURL   string            `json:"upload_url"`
	Headers     map[string]string `json:"headers,omitempty"`
	ExpiresAt   time.Time         `json:"expires_at"`
	UploadID    string            `json:"upload_id,omitempty"`
	PartSize    int64             `json:"part_size,omitempty"`
	StorageKey  string            `json:"storage_key,omitempty"`
}

type Usecase interface {
	PresignPut(ctx context.Context, tenantID, folderID, name, backendID string, overwrite bool, size int64, contentType string, principalID string) (*PresignResult, error)
	Complete(ctx context.Context, tenantID, fileID, etag string) (*domain.Object, error)
	DirectUpload(ctx context.Context, tenantID, folderID, name, backendID, contentType string, r io.Reader, size int64, principalID string) (*domain.Object, error)
	MultipartInit(ctx context.Context, tenantID, folderID, name, backendID string, overwrite bool, size int64, contentType string, principalID string) (*PresignResult, error)
	MultipartPresignPart(ctx context.Context, tenantID, fileID string, part int) (*PresignResult, error)
	MultipartComplete(ctx context.Context, tenantID, fileID string, parts []storage.Part) (*domain.Object, error)
	MultipartStatus(ctx context.Context, tenantID, fileID string) (*domain.MultipartUpload, error)
	Get(ctx context.Context, tenantID, fileID string) (*domain.Object, error)
	Move(ctx context.Context, tenantID, fileID, targetFolderID string) (*domain.Object, error)
	SetVisibility(ctx context.Context, tenantID, fileID, visibility string) (*domain.Object, error)
	Delete(ctx context.Context, tenantID, fileID string) error
	ListByFolder(ctx context.Context, tenantID, folderID, backendID string) ([]domain.Object, error)
	Download(ctx context.Context, tenantID, fileID string, auditFn func(action string)) (*DownloadResult, error)
}

type DownloadResult struct {
	RedirectURL string            // set when redirecting (presigned)
	Stream      io.ReadCloser     // set when proxying
	Object      *domain.Object
	ContentType string
	Size        int64
}
