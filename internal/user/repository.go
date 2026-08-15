package user

import (
	"context"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

// Repository is the persistence boundary for user data.
type Repository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Insert(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error
	UpdateLastLogin(ctx context.Context, id string) error
}

// Reader is the narrow read interface other domains depend on.
type Reader interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// Usecase is the user domain service.
type Usecase interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, displayName, locale *string, notifications *bool) (*domain.User, error)
	ChangePassword(ctx context.Context, userID, current, next string) error
}
