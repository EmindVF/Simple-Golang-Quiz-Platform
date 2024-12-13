package repository

import (
	"context"
	"quiz_platform/internal/models"
)

var (
	NewsRepositoryInstance NewsRepository
)

type NewsRepository interface {
	// May return ErrInternal or ErrInvalidInput on failure.
	AddNews(ctx context.Context, title string, text string, author_id int32) (int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	EditNews(ctx context.Context, id int32, title string, text string) error

	// May return ErrInternal or ErrNotFound on failure.
	DeleteNews(ctx context.Context, id int32) error

	// May return ErrInternal or ErrNotFound on failure.
	GetAllNews(ctx context.Context) ([]*models.News, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetNewsById(ctx context.Context, id int32) (*models.News, error)
}
