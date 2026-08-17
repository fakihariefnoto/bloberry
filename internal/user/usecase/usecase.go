package user

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/user"
)

type usecase struct {
	repo user.Repository
}

func NewUsecase(repo user.Repository) user.Usecase {
	return &usecase{repo: repo}
}

var _ user.Usecase = (*usecase)(nil)

func (u *usecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return u.repo.GetByID(ctx, userID)
}

// GetByEmail resolves a user by email (for adding a project member directly).
func (u *usecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.repo.GetByEmail(ctx, email)
}

func (u *usecase) UpdateProfile(ctx context.Context, userID string, displayName, locale *string, notifications *bool) (*domain.User, error) {
	usr, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if displayName != nil {
		usr.DisplayName = *displayName
	}
	if locale != nil {
		usr.Settings.Locale = *locale
	}
	if notifications != nil {
		usr.Settings.NotificationsEnabled = *notifications
	}
	if err := u.repo.Update(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *usecase) ChangePassword(ctx context.Context, userID, current, next string) error {
	usr, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if usr.PasswordHash == nil {
		return httpx.NewError(httpx.ErrBadRequest, 400)
	}
	ok, err := crypto.VerifyPassword(*usr.PasswordHash, current)
	if err != nil {
		return err
	}
	if !ok {
		return httpx.NewError(httpx.ErrInvalidCredentials, 401)
	}
	hash, err := crypto.HashPassword(next)
	if err != nil {
		return err
	}
	usr.PasswordHash = &hash
	return u.repo.Update(ctx, usr)
}


