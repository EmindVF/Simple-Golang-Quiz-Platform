package repositories

import (
	"context"

	"quiz_platform/internal/database"
	"quiz_platform/internal/misc/apperrors"
	"quiz_platform/internal/models"
)

type SqlActionsRepository struct {
	DBProvider database.SqlDatabaseProvider
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlActionsRepository) GetAllActions(ctx context.Context) ([]*models.UserAction, error) {
	query :=
		`SELECT
		id, user_id, action_description, action_time
		FROM user_actions ORDER BY action_time DESC`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allActions := make([]*models.UserAction, 0)
	for rows.Next() {
		var action models.UserAction
		err = rows.Scan(
			&action.Id, &action.UserId,
			&action.ActionDescription, &action.ActionTime)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allActions = append(allActions, &action)
	}

	return allActions, nil
}
