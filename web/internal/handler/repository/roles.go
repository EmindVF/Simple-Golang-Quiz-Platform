package repository

import (
	"context"
	"quiz_platform/internal/models"
)

var (
	RoleRepositoryInstance RoleRepository
)

type RoleRepository interface {
	// May return ErrInternal or ErrNotFound on failure.
	GetAllRoles(ctx context.Context) ([]*models.Role, error)

	// May return ErrInternal or ErrNotFound on failure.
	RemoveUserRoles(ctx context.Context, userId int32) error

	// May return ErrInternal or ErrNotFound on failure.
	AddUserRoles(ctx context.Context, userId int32, roleIds []int32) error
}
