package usecase

import (
	"context"
	"strings"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/folder"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
)

type usecase struct {
	repo folder.Repository
}

func NewUsecase(repo folder.Repository) folder.Usecase {
	return &usecase{repo: repo}
}

var _ folder.Usecase = (*usecase)(nil)

func (u *usecase) Create(ctx context.Context, tenantID, parentID, name string) (*domain.Folder, error) {
	if name == "" {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	var parent *domain.Folder
	if parentID != "" {
		p, err := u.repo.GetByID(ctx, tenantID, parentID)
		if err != nil {
			return nil, err
		}
		parent = p
	}
	f := &domain.Folder{
		TenantID: tenantID,
		Name:     name,
	}
	if parent != nil {
		f.ParentID = &parent.ID
		f.Ancestors = append(append([]string{}, parent.Ancestors...), parent.ID)
		f.Depth = parent.Depth + 1
		f.Path = join(parent.Path, name)
	} else {
		f.Ancestors = []string{}
		f.Depth = 0
		f.Path = "/" + name
	}
	if err := u.repo.Create(ctx, f); err != nil {
		return nil, httpx.NewError(httpx.ErrNameConflict, 409)
	}
	return f, nil
}

func (u *usecase) Get(ctx context.Context, tenantID, id string) (*domain.Folder, error) {
	return u.repo.GetByID(ctx, tenantID, id)
}

func (u *usecase) GetRoot(ctx context.Context, tenantID string) (*domain.Folder, error) {
	return u.repo.GetByPath(ctx, tenantID, "/")
}

func (u *usecase) Rename(ctx context.Context, tenantID, id, name string) (*domain.Folder, error) {
	f, err := u.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	f.Name = name
	if f.ParentID == nil {
		f.Path = "/" + name
	} else {
		parent, err := u.repo.GetByID(ctx, tenantID, *f.ParentID)
		if err == nil {
			f.Path = join(parent.Path, name)
		}
	}
	if err := u.repo.Update(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (u *usecase) Move(ctx context.Context, tenantID, id, targetParentID string) (*domain.Folder, error) {
	f, err := u.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	target, err := u.repo.GetByID(ctx, tenantID, targetParentID)
	if err != nil {
		return nil, err
	}
	// cycle prevention: target must not be f or a descendant of f
	if target.ID == f.ID {
		return nil, httpx.NewError(httpx.ErrFolderCycle, 422)
	}
	for _, a := range target.Ancestors {
		if a == f.ID {
			return nil, httpx.NewError(httpx.ErrFolderCycle, 422)
		}
	}
	f.ParentID = &target.ID
	f.Ancestors = append(append([]string{}, target.Ancestors...), target.ID)
	f.Depth = target.Depth + 1
	f.Path = join(target.Path, f.Name)
	if err := u.repo.Update(ctx, f); err != nil {
		return nil, err
	}
	// rewrite descendants (path + ancestors)
	desc, err := u.repo.Descendants(ctx, tenantID, id)
	if err == nil {
		for i := range desc {
			d := &desc[i]
			d.Ancestors = append(append([]string{}, f.Ancestors...), f.ID)
			_ = u.repo.Update(ctx, d)
		}
	}
	return f, nil
}

func (u *usecase) Delete(ctx context.Context, tenantID, id string) error {
	return u.repo.Delete(ctx, tenantID, id)
}

func (u *usecase) ListChildren(ctx context.Context, tenantID string, parentID *string) ([]domain.Folder, error) {
	return u.repo.ListChildren(ctx, tenantID, parentID)
}

func join(base, name string) string {
	if base == "" || base == "/" {
		return "/" + strings.TrimPrefix(name, "/")
	}
	return strings.TrimSuffix(base, "/") + "/" + name
}


