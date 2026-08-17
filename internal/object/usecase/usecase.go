package usecase

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/object"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/storage"
)

type usecase struct {
	repo       object.Repository
	registry   object.Registry
	quota      quotaChecker
	folders    folderReader
	maxSize    int64
	partSize   int64
	baseURL    string
	rawSecret  []byte
}

type quotaChecker interface {
	CheckQuota(ctx context.Context, tenantID string, addBytes, addObjects int64) error
	IncrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
	DecrementUsed(ctx context.Context, tenantID string, bytes, objects int64) error
}

type folderReader interface {
	Get(ctx context.Context, tenantID, id string) (*domain.Folder, error)
	GetRoot(ctx context.Context, tenantID string) (*domain.Folder, error)
}

// resolveFolderID maps the pseudo-id "root" to the tenant's root folder.
func resolveFolderID(ctx context.Context, f folderReader, tenantID, folderID string) string {
	if folderID == "" || folderID == "root" {
		if root, err := f.GetRoot(ctx, tenantID); err == nil {
			return root.ID
		}
		return folderID
	}
	return folderID
}

type Deps struct {
	Repo      object.Repository
	Registry  object.Registry
	Quota     quotaChecker
	Folders   folderReader
	MaxSize   int64
	PartSize  int64
	BaseURL   string
	RawSecret []byte
}

func NewUsecase(d Deps) object.Usecase {
	return &usecase{
		repo: d.Repo, registry: d.Registry, quota: d.Quota, folders: d.Folders,
		maxSize: d.MaxSize, partSize: d.PartSize, baseURL: d.BaseURL, rawSecret: d.RawSecret,
	}
}

var _ object.Usecase = (*usecase)(nil)

func (u *usecase) backendFor(ctx context.Context, tenantID, backendID string) (*domain.StorageBackend, storage.Driver, error) {
	var be *domain.StorageBackend
	var err error
	if backendID != "" {
		be, err = u.repo.GetBackend(ctx, backendID)
	} else {
		be, err = u.repo.GetTenantBackend(ctx, tenantID)
	}
	if err != nil {
		return nil, nil, httpx.NewErrorContent(httpx.ErrBackendUnreachable, 502, "storage engine lookup failed: "+err.Error())
	}
	drv, err := u.registry.Get(be.ID)
	if err != nil {
		return nil, nil, httpx.NewErrorContent(httpx.ErrBackendUnreachable, 502, "storage engine driver not ready: "+err.Error())
	}
	return be, drv, nil
}

func (u *usecase) PresignPut(ctx context.Context, tenantID, folderID, name, backendID string, overwrite bool, size int64, contentType string, principalID string) (*object.PresignResult, error) {
	if size > u.maxSize {
		return nil, httpx.NewError(httpx.ErrPayloadTooLarge, 413)
	}
	folder, err := u.folders.Get(ctx, tenantID, resolveFolderID(ctx, u.folders, tenantID, folderID))
	if err != nil {
		return nil, err
	}
	if err := u.quota.CheckQuota(ctx, tenantID, size, 1); err != nil {
		return nil, err
	}
	be, drv, err := u.backendFor(ctx, tenantID, backendID)
	if err != nil {
		return nil, err
	}
	key := storageKey(folder.Path, name)
	// name conflict: replace (overwrite) or 409
	if existing, err := u.repo.GetByName(ctx, tenantID, folder.ID, name); err == nil {
		if !overwrite {
			return nil, httpx.NewError(httpx.ErrNameConflict, 409)
		}
		// delete the old blob + record, then upload fresh under the same name
		oldDrv, derr := u.registry.Get(existing.BackendID)
		if derr == nil {
			_ = oldDrv.Delete(ctx, []string{existing.StorageKey})
		}
		_ = u.quota.DecrementUsed(ctx, tenantID, existing.SizeBytes, 1)
		_ = u.repo.Delete(ctx, tenantID, existing.ID)
	}
	// pending record first (two-phase write, ADR-5)
	obj := &domain.Object{
		ID: crypto.NewID(), TenantID: tenantID, FolderID: folder.ID,
		Ancestors: folder.Ancestors, Name: name, BackendID: be.ID,
		StorageKey: key, State: "pending", SizeBytes: size, ContentType: contentType,
		Visibility: "private", UploadedBy: principalID,
	}
	if err := u.repo.Insert(ctx, obj); err != nil {
		return nil, httpx.NewError(httpx.ErrNameConflict, 409)
	}
	pu, err := drv.PresignPut(ctx, key, 5*time.Minute, size)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	return &object.PresignResult{
		FileID: obj.ID, UploadURL: pu.URL, Headers: pu.Headers,
		ExpiresAt: time.Now().Add(5 * time.Minute), StorageKey: key,
	}, nil
}

func (u *usecase) Complete(ctx context.Context, tenantID, fileID, etag string) (*domain.Object, error) {
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return nil, err
	}
	if obj.State == "active" {
		return obj, nil
	}
	drv, err := u.registry.Get(obj.BackendID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	info, err := drv.Stat(ctx, obj.StorageKey)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	obj.State = "active"
	obj.SizeBytes = info.Size
	obj.UpdatedAt = time.Now().UTC()
	if err := u.repo.Update(ctx, obj); err != nil {
		return nil, err
	}
	if err := u.quota.IncrementUsed(ctx, tenantID, info.Size, 1); err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *usecase) DirectUpload(ctx context.Context, tenantID, folderID, name, backendID, contentType string, r io.Reader, size int64, principalID string) (*domain.Object, error) {
	// proxy path: bytes through Bloberry
	res, err := u.PresignPut(ctx, tenantID, folderID, name, backendID, false, size, contentType, principalID)
	if err != nil {
		return nil, err
	}
	be, drv, err := u.backendFor(ctx, tenantID, backendID)
	if err != nil {
		return nil, err
	}
	if err := drv.Put(ctx, res.StorageKey, r, size, contentType); err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	obj, err := u.repo.GetByID(ctx, tenantID, res.FileID)
	if err != nil {
		return nil, err
	}
	obj.State = "active"
	obj.BackendID = be.ID
	obj.SizeBytes = size
	obj.ContentType = contentType
	obj.UpdatedAt = time.Now().UTC()
	if err := u.repo.Update(ctx, obj); err != nil {
		return nil, err
	}
	if err := u.quota.IncrementUsed(ctx, tenantID, size, 1); err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *usecase) MultipartInit(ctx context.Context, tenantID, folderID, name, backendID string, overwrite bool, size int64, contentType string, principalID string) (*object.PresignResult, error) {
	if size > u.maxSize {
		return nil, httpx.NewError(httpx.ErrPayloadTooLarge, 413)
	}
	folder, err := u.folders.Get(ctx, tenantID, resolveFolderID(ctx, u.folders, tenantID, folderID))
	if err != nil {
		return nil, err
	}
	if err := u.quota.CheckQuota(ctx, tenantID, size, 1); err != nil {
		return nil, err
	}
	be, drv, err := u.backendFor(ctx, tenantID, backendID)
	if err != nil {
		return nil, err
	}
	if !drv.Capabilities().Multipart {
		return nil, httpx.NewErrorContent(httpx.ErrBadRequest, 400, "driver has no multipart; use presigned-PUT")
	}
	key := storageKey(folder.Path, name)
	// name conflict: replace (overwrite) or 409
	if existing, err := u.repo.GetByName(ctx, tenantID, folder.ID, name); err == nil {
		if !overwrite {
			return nil, httpx.NewError(httpx.ErrNameConflict, 409)
		}
		oldDrv, derr := u.registry.Get(existing.BackendID)
		if derr == nil {
			_ = oldDrv.Delete(ctx, []string{existing.StorageKey})
		}
		_ = u.quota.DecrementUsed(ctx, tenantID, existing.SizeBytes, 1)
		_ = u.repo.Delete(ctx, tenantID, existing.ID)
	}
	uploadID, err := drv.MultipartInit(ctx, key, contentType)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	obj := &domain.Object{
		ID: crypto.NewID(), TenantID: tenantID, FolderID: folder.ID,
		Ancestors: folder.Ancestors, Name: name, BackendID: be.ID,
		StorageKey: key, State: "pending", SizeBytes: size, ContentType: contentType,
		Visibility: "private", UploadedBy: principalID,
	}
	if err := u.repo.Insert(ctx, obj); err != nil {
		return nil, httpx.NewError(httpx.ErrNameConflict, 409)
	}
	mp := &domain.MultipartUpload{
		ID: crypto.NewID(), ObjectID: obj.ID, TenantID: tenantID,
		ProviderUploadID: uploadID, PartSizeBytes: u.partSize,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := u.repo.InsertMultipart(ctx, mp); err != nil {
		return nil, err
	}
	return &object.PresignResult{
		FileID: obj.ID, UploadID: uploadID, PartSize: u.partSize,
		ExpiresAt: time.Now().Add(24 * time.Hour), StorageKey: key,
	}, nil
}

func (u *usecase) MultipartPresignPart(ctx context.Context, tenantID, fileID string, part int) (*object.PresignResult, error) {
	mp, err := u.repo.GetMultipart(ctx, tenantID, fileID)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return nil, err
	}
	drv, err := u.registry.Get(obj.BackendID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	pu, err := drv.MultipartPresignPart(ctx, obj.StorageKey, mp.ProviderUploadID, part, 5*time.Minute)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	return &object.PresignResult{UploadURL: pu.URL, Headers: pu.Headers, ExpiresAt: time.Now().Add(5 * time.Minute)}, nil
}

func (u *usecase) MultipartComplete(ctx context.Context, tenantID, fileID string, parts []storage.Part) (*domain.Object, error) {
	mp, err := u.repo.GetMultipart(ctx, tenantID, fileID)
	if err != nil {
		return nil, httpx.ErrResourceNotFound
	}
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return nil, err
	}
	drv, err := u.registry.Get(obj.BackendID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	info, err := drv.MultipartComplete(ctx, obj.StorageKey, mp.ProviderUploadID, parts)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	obj.State = "active"
	if info != nil {
		obj.SizeBytes = info.Size
	}
	obj.UpdatedAt = time.Now().UTC()
	if err := u.repo.Update(ctx, obj); err != nil {
		return nil, err
	}
	if err := u.quota.IncrementUsed(ctx, tenantID, obj.SizeBytes, 1); err != nil {
		return nil, err
	}
	_ = u.repo.DeleteMultipart(ctx, tenantID, fileID)
	return obj, nil
}

func (u *usecase) MultipartStatus(ctx context.Context, tenantID, fileID string) (*domain.MultipartUpload, error) {
	return u.repo.GetMultipart(ctx, tenantID, fileID)
}

func (u *usecase) Get(ctx context.Context, tenantID, fileID string) (*domain.Object, error) {
	return u.repo.GetByID(ctx, tenantID, fileID)
}

func (u *usecase) Move(ctx context.Context, tenantID, fileID, targetFolderID string) (*domain.Object, error) {
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return nil, err
	}
	target, err := u.folders.Get(ctx, tenantID, resolveFolderID(ctx, u.folders, tenantID, targetFolderID))
	if err != nil {
		return nil, err
	}
	obj.FolderID = target.ID
	obj.Ancestors = target.Ancestors
	obj.UpdatedAt = time.Now().UTC()
	if err := u.repo.Update(ctx, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *usecase) SetVisibility(ctx context.Context, tenantID, fileID, visibility string) (*domain.Object, error) {
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return nil, err
	}
	if visibility != "private" && visibility != "public" {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	obj.Visibility = visibility
	obj.UpdatedAt = time.Now().UTC()
	if err := u.repo.Update(ctx, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *usecase) Delete(ctx context.Context, tenantID, fileID string) error {
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return err
	}
	drv, err := u.registry.Get(obj.BackendID)
	if err == nil {
		_ = drv.Delete(ctx, []string{obj.StorageKey})
	}
	obj.DeletedAt = nowPtr()
	if err := u.repo.Update(ctx, obj); err != nil {
		return err
	}
	return u.quota.DecrementUsed(ctx, tenantID, obj.SizeBytes, 1)
}

func (u *usecase) ListByFolder(ctx context.Context, tenantID, folderID, backendID string) ([]domain.Object, error) {
	return u.repo.ListByFolder(ctx, tenantID, folderID, backendID)
}

func (u *usecase) Download(ctx context.Context, tenantID, fileID string, auditFn func(action string)) (*object.DownloadResult, error) {
	obj, err := u.repo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		return nil, err
	}
	drv, err := u.registry.Get(obj.BackendID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	if drv.Capabilities().Presign {
		pu, err := drv.PresignGet(ctx, obj.StorageKey, 5*time.Minute)
		if err != nil {
			return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
		}
		if auditFn != nil {
			auditFn("object.link_issued")
		}
		return &object.DownloadResult{RedirectURL: pu.URL, Object: obj, ContentType: obj.ContentType, Size: obj.SizeBytes}, nil
	}
	// proxy path
	rc, info, err := drv.Get(ctx, obj.StorageKey, nil)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrBackendUnreachable, 502)
	}
	if auditFn != nil {
		auditFn("object.read")
	}
	size := obj.SizeBytes
	if info != nil {
		size = info.Size
	}
	return &object.DownloadResult{Stream: rc, Object: obj, ContentType: obj.ContentType, Size: size}, nil
}

func nowPtr() *time.Time {
	t := time.Now().UTC()
	return &t
}



func storageKey(folderPath, name string) string {
	p := strings.Trim(strings.TrimSpace(folderPath), "/")
	if p == "" {
		return name
	}
	return p + "/" + name
}
