package logsql

import (
	"context"
	"database/sql/driver"
	"time"
)

type (
	queryTransaction struct {
		logHandler Logger

		connCtx     context.Context
		transaction driver.Tx
	}
)

func (t *queryTransaction) Commit() error {
	t0 := time.Now()

	err := t.transaction.Commit()
	t.logHandler.TxCommit(t.connCtx, err, time.Since(t0))

	return err
}

func (t *queryTransaction) Rollback() error {
	t0 := time.Now()

	err := t.transaction.Rollback()
	t.logHandler.TxRollback(t.connCtx, err, time.Since(t0))

	return err
}
