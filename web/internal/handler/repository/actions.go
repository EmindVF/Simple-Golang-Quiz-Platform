package repository

import (
	"context"
	"quiz_platform/internal/models"
)

var (
	ActionsRepositoryInstance ActionsRepository
)

type ActionsRepository interface {
	// May return ErrInternal or ErrNotFound on failure.
	GetAllActions(ctx context.Context) ([]*models.UserAction, error)
}
