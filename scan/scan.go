package scan

import (
	"database/sql"
	"errors"
)

var (
	// PreAllocCollection can be used to override the default capacity of collections returned by CollectRows, CollectRowsKV, CollectRow
	PreAllocCollection uint = 64
)

type (
	// Row is an abstraction for sql.Rows that can be used only for scanning.
	Row interface {
		Scan(...interface{}) error
	}

	// RowCols is the same as Row, but column names are also available
	RowCols interface {
		Row
		Columns() ([]string, error)
	}

	// RowColsCollector should return concrete typed value scanned from row.
	// Basic implementations are:
	//  - StructTagCollectorRowCols to scan row into struct via its field tags "scan"
	//  - StructPosCollectorRowCols to scan row into struct via its field position
	//  - DirectCollectorRowCols to scan row into T directly (primitives and sql.Scanner implementations)
	RowColsCollector[T any] func(row RowCols) (T, error)

	// RowColsCollectorKV should return concrete typed key and value scanned from row.
	// Basic implementations are:
	//  - StructTagCollectorRowColsKV to scan row into struct via its field tags "scan". Key will be used in scan directly
	//  - StructPosCollectorRowColsKV to scan row into struct via its field position. Key will be used in scan directly
	//  - DirectCollectorRowColsKV to scan row into T directly (primitives and sql.Scanner implementations). Both key and value will be used in scan directly
	RowColsCollectorKV[K comparable, V any] func(row RowCols) (K, V, error)

	// RowCollector should return concrete typed value scanned from row, without any information about columns.
	// Basic implementations are:
	//  - StructPosCollectorRow to scan row into struct via its field position
	//  - DirectCollectorRow to scan row into T directly (primitives and sql.Scanner implementations)
	RowCollector[T any] func(row Row) (T, error)
)

// CollectRows can be used for easy scanning collection of T. Note, that it does NOT explicitly close sql.Rows.
func CollectRows[T any](rows *sql.Rows, sc RowColsCollector[T]) ([]T, error) {
	items := make([]T, 0, PreAllocCollection)

	for rows.Next() {
		v, err := sc(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// CollectRowsKV can be used for scanning a map of K and V. Note, that it does NOT explicitly close sql.Rows.
func CollectRowsKV[K comparable, V any](rows *sql.Rows, sc RowColsCollectorKV[K, V]) (map[K]V, error) {
	m := make(map[K]V, PreAllocCollection)

	for rows.Next() {
		k, v, err := sc(rows)
		if err != nil {
			return nil, err
		}

		m[k] = v
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// CollectRow can be used for scanning single value T from a single a row.
func CollectRow[T any](row *sql.Row, sc RowCollector[T]) (T, bool, error) {
	v, err := sc(row)
	if err != nil {
		var defaultValue T
		if errors.Is(err, sql.ErrNoRows) {
			return defaultValue, false, nil
		}
		return defaultValue, false, err
	}

	return v, true, nil
}
