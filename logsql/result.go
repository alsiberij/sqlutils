package logsql

import (
	"context"
	"database/sql/driver"
)

type (
	queryResult struct {
		logHandler       Logger
		queryErrReplacer QueryErrReplacer

		ctx    context.Context
		result driver.Result
	}
)

func (r *queryResult) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *queryResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}
