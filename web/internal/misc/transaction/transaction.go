package transaction

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type TransactionManager interface {
	Run(ctx context.Context, callback func(ctx context.Context) error) error
}

type txKey string

var ctxWithTx = txKey("tx")

type SQLTransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *SQLTransactionManager {
	return &SQLTransactionManager{
		db: db,
	}
}

func (manager *SQLTransactionManager) Run(
	ctx context.Context,
	callback func(ctx context.Context) error,
) (rErr error) {
	tx, err := manager.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if rErr != nil {
			rErr = multierr.Combine(rErr, errors.WithStack(tx.Rollback()))
		}
	}()
	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				rErr = e
			} else {
				rErr = errors.Errorf("%s", rec)
			}
		}
	}()

	if err = callback(insertTxIntoContext(ctx, tx)); err != nil {
		return err
	}
	return errors.WithStack(tx.Commit())
}

func ExtractTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx := ctx.Value(ctxWithTx)
	if t, ok := tx.(*sql.Tx); ok {
		return t, true
	}
	return nil, false
}

func insertTxIntoContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxWithTx, tx)
}
