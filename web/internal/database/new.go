package database

import (
	"context"
	"fmt"

	"quiz_platform/internal/misc/config"
	"quiz_platform/internal/misc/transaction"

	"database/sql"

	_ "github.com/lib/pq"
)

type postgresSqlDatabaseProvider struct {
	Db *sql.DB
}

func (p *postgresSqlDatabaseProvider) GetDb() *sql.DB {
	return p.Db
}

func (db *postgresSqlDatabaseProvider) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	tx, ok := transaction.ExtractTxFromContext(ctx)
	if !ok {
		return db.Db.QueryRowContext(ctx, query, args...)
	}
	return tx.QueryRowContext(ctx, query, args...)
}

func (db *postgresSqlDatabaseProvider) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	tx, ok := transaction.ExtractTxFromContext(ctx)
	if !ok {
		return db.Db.QueryContext(ctx, query, args...)
	}
	return tx.QueryContext(ctx, query, args...)
}

func (db *postgresSqlDatabaseProvider) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, ok := transaction.ExtractTxFromContext(ctx)
	if !ok {
		return db.Db.ExecContext(ctx, query, args...)
	}
	return tx.ExecContext(ctx, query, args...)
}

func (db *postgresSqlDatabaseProvider) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	tx, ok := transaction.ExtractTxFromContext(ctx)
	if !ok {
		return db.Db.PrepareContext(ctx, query)
	}
	return tx.PrepareContext(ctx, query)
}

func NewPostgresSqlDatabaseProvider() SqlDatabaseProvider {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.GlobalConfig.Database.Host,
		config.GlobalConfig.Database.User,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.DBName,
		config.GlobalConfig.Database.Port,
		config.GlobalConfig.Database.SSLMode,
		config.GlobalConfig.Database.TimeZone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect Postgres database.")
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("ping to database failed: %v", err))
	}

	//initializePostgresDatabase(db)

	return &postgresSqlDatabaseProvider{Db: db}
}
