package repositories

import (
	"context"
	"fmt"
	"strings"

	"quiz_platform/internal/database"
	"quiz_platform/internal/misc/apperrors"
	"quiz_platform/internal/models"
)

type SqlRoleRepository struct {
	DBProvider database.SqlDatabaseProvider
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlRoleRepository) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	query :=
		`SELECT
		id, name, permissions, created_at, updated_at
		FROM roles ORDER BY created_at DESC`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allRoles := make([]*models.Role, 0)
	for rows.Next() {
		var role models.Role
		err = rows.Scan(
			&role.Id, &role.Name, &role.Permissions,
			&role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allRoles = append(allRoles, &role)
	}

	return allRoles, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlRoleRepository) RemoveUserRoles(ctx context.Context, userId int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM user_roles WHERE user_id = $1`,
		userId)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlRoleRepository) AddUserRoles(ctx context.Context, userId int32, rolesIds []int32) error {
	if len(rolesIds) == 0 {
		return nil
	}

	query := strings.Builder{}
	_, err := query.Write(
		[]byte(
			`INSERT INTO user_roles
		(user_id, role_id)
		VALUES `))
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	values := []any{}

	placeHolderArray := make([]any, 2)
	for i := int32(1); i <= 2; i++ {
		placeHolderArray[i-1] = i
	}

	query.Write([]byte(fmt.Sprintf(" ($%d, $%d) ",
		placeHolderArray...)))
	values = append(values,
		userId, rolesIds[0])

	for i := 1; i < len(rolesIds); i++ {
		for j := range 2 {
			placeHolderArray[j] = i*2 + j + 1
		}
		query.Write([]byte(fmt.Sprintf(", ($%d, $%d) ",
			placeHolderArray...)))
		values = append(values,
			userId, rolesIds[i])
	}

	_, err = repo.DBProvider.ExecContext(
		ctx,
		query.String(),
		values...,
	)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}
