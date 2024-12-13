package repository

import (
	"context"
	"quiz_platform/internal/models"
)

var (
	UserRepositoryInstance UserRepository
)

type UserRepository interface {
	// May return ErrInternal or ErrInvalidInput on failure.
	AddUser(ctx context.Context, userName string, email string, passwordHash string) error

	// May return ErrInternal or ErrNotFound on failure.
	DeleteUser(ctx context.Context, id int32) error

	// May return ErrInternal or ErrNotFound on failure.
	UpdateUserName(ctx context.Context, id int32, email string) error

	// May return ErrInternal or ErrNotFound on failure.
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetUserById(ctx context.Context, id int32) (*models.User, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetUserPermissions(ctx context.Context, id int32) (int64, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetAllUsers(ctx context.Context) ([]*models.User, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetUsersRoles(ctx context.Context, baseUserIds []int32) ([]int32, []*models.Role, error)
}
