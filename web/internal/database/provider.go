package database

import (
	"context"
	"database/sql"
	"quiz_platform/internal/misc/apperrors"
)

type SqlDatabaseProvider interface {
	GetDb() *sql.DB
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

var sqlDatabaseProviderInstance SqlDatabaseProvider

func InitSqlDatabaseProvider(provider SqlDatabaseProvider) error {
	if sqlDatabaseProviderInstance != nil {
		return &apperrors.ErrInternal{
			Message: "Sql Database Provider is already initiated"}
	}
	sqlDatabaseProviderInstance = provider
	return nil
}

func GetSqlDatabaseProvider() (SqlDatabaseProvider, error) {
	if sqlDatabaseProviderInstance == nil {
		return nil, &apperrors.ErrInternal{
			Message: "Sql Database Provider is not initiated"}
	}
	return sqlDatabaseProviderInstance, nil
}
