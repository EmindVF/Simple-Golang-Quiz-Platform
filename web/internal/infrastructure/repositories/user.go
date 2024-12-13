package repositories

import (
	"context"
	"database/sql"

	"quiz_platform/internal/database"
	"quiz_platform/internal/misc/apperrors"
	"quiz_platform/internal/models"

	"github.com/lib/pq"
)

type SqlUserRepository struct {
	DBProvider database.SqlDatabaseProvider
}

// May return ErrInternal or ErrInvalidInput on failure.
func (repo *SqlUserRepository) AddUser(ctx context.Context, username string, email string, passwordHash string) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`INSERT INTO users 
		(username, email, password_hash) 
		VALUES ($1, $2, $3)`,
		username, email, passwordHash)

	if err == sql.ErrConnDone {
		return &apperrors.ErrInternal{Message: "connection is done"}
	} else if err != nil {
		return &apperrors.ErrInvalidInput{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrInvalidInput on failure.
func (repo *SqlUserRepository) UpdateUserName(ctx context.Context, userId int32, userName string) error {
	res, err := repo.DBProvider.ExecContext(
		ctx,
		`UPDATE users SET
		username = $1
		WHERE id = $2`,
		userName, userId)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}
	if rowsAffected == 0 {
		return &apperrors.ErrNotFound{Message: "user not found"}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT id, username, email, password_hash, created_at, updated_at
		FROM users 
		WHERE email = $1 `,
		email).Scan(
		&user.Id, &user.UserName, &user.Email,
		&user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return user, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlUserRepository) GetUserById(ctx context.Context, id int32) (*models.User, error) {
	user := &models.User{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT id, username, email, password_hash, created_at, updated_at
		FROM users 
		WHERE id = $1 `,
		id).Scan(
		&user.Id, &user.UserName, &user.Email,
		&user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return user, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlUserRepository) GetUserPermissions(ctx context.Context, id int32) (int64, error) {
	permissions := int64(0)
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT 
			BIT_OR(permissions)
		FROM 
			user_roles ur
		JOIN 
			roles r ON ur.role_id = r.id
		WHERE 
			ur.user_id = $1
		HAVING 
			COUNT(permissions) > 0;`,
		id).Scan(
		&permissions)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return permissions, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query :=
		`SELECT
		id, username, email, password_hash, created_at, updated_at
		FROM users ORDER BY created_at DESC`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allUsers := make([]*models.User, 0)
	for rows.Next() {
		var news models.User
		err = rows.Scan(
			&news.Id, &news.UserName, &news.Email, &news.PasswordHash,
			&news.CreatedAt, &news.UpdatedAt)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allUsers = append(allUsers, &news)
	}

	return allUsers, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlUserRepository) GetUsersRoles(ctx context.Context, baseUserIds []int32) ([]int32, []*models.Role, error) {
	query :=
		`SELECT 
			ur.user_id, ur.role_id, r.name, r.permissions, r.created_at, r.updated_at
		FROM 
			user_roles ur
		JOIN 
			roles r ON ur.role_id = r.id
		WHERE 
			ur.user_id = ANY($1)`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
		pq.Array(baseUserIds),
	)
	if err != nil {
		return nil, nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	userIds := make([]int32, 0)
	roles := make([]*models.Role, 0)
	for rows.Next() {
		var userId int32
		var role models.Role
		err = rows.Scan(
			&userId, &role.Id, &role.Name, &role.Permissions,
			&role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		userIds = append(userIds, userId)
		roles = append(roles, &role)
	}

	return userIds, roles, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlUserRepository) DeleteUser(ctx context.Context, id int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM users WHERE
		id = $1`, id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}
