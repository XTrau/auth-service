package repositories

import (
	"context"
	"database/sql"

	errs "github.com/XTrau/auth-service/internal/errors"
)

var txKey string = "tx"

func GetTxFromContext(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)

	if !ok || tx == nil {
		return nil, errs.ErrTxIsNil
	}

	return tx, nil
}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
