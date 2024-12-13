package repositories

import (
	"context"
	"database/sql"

	"quiz_platform/internal/database"
	"quiz_platform/internal/misc/apperrors"
	"quiz_platform/internal/models"
)

type SqlNewsRepository struct {
	DBProvider database.SqlDatabaseProvider
}

// May return ErrInternal or ErrInvalidInput on failure.
func (repo *SqlNewsRepository) AddNews(ctx context.Context, title string, text string, authorId int32) (int32, error) {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`INSERT INTO news 
		(title, news_text, author_id) 
		VALUES ($1, $2, $3) RETURNING id`,
		title, text, authorId).Scan(
		&id)

	if err == sql.ErrConnDone {
		return 0, &apperrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return 0, &apperrors.ErrInvalidInput{Message: err.Error()}
	}

	return id, nil
}

// May return ErrInternal or ErrInvalidInput on failure.
func (repo *SqlNewsRepository) EditNews(ctx context.Context, id int32, title string, text string) error {
	res, err := repo.DBProvider.ExecContext(
		ctx,
		`UPDATE news SET
		title = $1, news_text = $2
		WHERE id = $3`,
		title, text, id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}
	if rowsAffected == 0 {
		return &apperrors.ErrNotFound{Message: "news not found"}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlNewsRepository) GetNewsById(ctx context.Context, id int32) (*models.News, error) {
	news := &models.News{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT id, title, news_text, author_id, created_at, updated_at
		FROM news 
		WHERE id = $1 `,
		id).Scan(
		&news.Id, &news.Title, &news.NewsText,
		&news.AuthorId, &news.CreatedAt, &news.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return news, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlNewsRepository) GetAllNews(ctx context.Context) ([]*models.News, error) {
	query :=
		`SELECT
		id, title, news_text, author_id, created_at, updated_at
		FROM news ORDER BY created_at DESC`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allNews := make([]*models.News, 0)
	for rows.Next() {
		var news models.News
		err = rows.Scan(
			&news.Id, &news.Title, &news.NewsText, &news.AuthorId,
			&news.CreatedAt, &news.UpdatedAt)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allNews = append(allNews, &news)
	}

	return allNews, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlNewsRepository) DeleteNews(ctx context.Context, id int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM news WHERE
		id = $1`, id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}
