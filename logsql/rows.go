package logsql

import (
	"context"
	"database/sql/driver"
	"io"
	"reflect"
	"time"
)

type (
	queryRows struct {
		logHandler Logger

		connCtx context.Context
		rows    driver.Rows
	}
)

func (r *queryRows) Columns() []string {
	return r.rows.Columns()
}

func (r *queryRows) Close() error {
	t0 := time.Now()

	err := r.rows.Close()
	r.logHandler.RowsClose(r.connCtx, err, time.Since(t0))

	return err
}

func (r *queryRows) Next(dest []driver.Value) error {
	t0 := time.Now()

	err := r.rows.Next(dest)
	r.logHandler.RowsNext(r.connCtx, dest, err, time.Since(t0))

	return err
}

func (r *queryRows) HasNextResultSet() bool {
	rs, ok := r.rows.(driver.RowsNextResultSet)
	if !ok {
		return false
	}

	return rs.HasNextResultSet()
}

func (r *queryRows) NextResultSet() error {
	rs, ok := r.rows.(driver.RowsNextResultSet)
	if !ok {
		return io.EOF
	}

	return rs.NextResultSet()
}

func (r *queryRows) ColumnTypeScanType(index int) reflect.Type {
	rs, ok := r.rows.(driver.RowsColumnTypeScanType)
	if !ok {
		return reflect.TypeFor[any]()
	}

	return rs.ColumnTypeScanType(index)
}

func (r *queryRows) ColumnTypeDatabaseTypeName(index int) string {
	rs, ok := r.rows.(driver.RowsColumnTypeDatabaseTypeName)
	if !ok {
		return ""
	}

	return rs.ColumnTypeDatabaseTypeName(index)
}

func (r *queryRows) ColumnTypeLength(index int) (int64, bool) {
	rs, ok := r.rows.(driver.RowsColumnTypeLength)
	if !ok {
		return 0, false
	}

	return rs.ColumnTypeLength(index)
}

func (r *queryRows) ColumnTypeNullable(index int) (bool, bool) {
	rs, ok := r.rows.(driver.RowsColumnTypeNullable)
	if !ok {
		return false, false
	}

	return rs.ColumnTypeNullable(index)
}

func (r *queryRows) ColumnTypePrecisionScale(index int) (int64, int64, bool) {
	rs, ok := r.rows.(driver.RowsColumnTypePrecisionScale)
	if !ok {
		return 0, 0, false
	}

	return rs.ColumnTypePrecisionScale(index)
}
