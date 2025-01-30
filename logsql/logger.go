package logsql

import (
	"context"
	"database/sql/driver"
	"time"
)

type (
	// Logger is used to log specific cases internally in [sql.DB]
	Logger interface {
		Connect(ctx context.Context, err error, dt time.Duration)
		ConnClose(ctx context.Context, err error, dt time.Duration)

		TxBegin(ctx context.Context, err error, dt time.Duration)
		TxCommit(ctx context.Context, err error, dt time.Duration)
		TxRollback(ctx context.Context, err error, dt time.Duration)

		Exec(ctx context.Context, query string, args []driver.NamedValue, replacedErr error, err error, dt time.Duration)
		Query(ctx context.Context, query string, args []driver.NamedValue, replacedErr error, err error, dt time.Duration)

		NamedValueCheck(ctx context.Context, err error, dt time.Duration)
		ValidateConn(ctx context.Context, isValid bool, dt time.Duration)

		Ping(ctx context.Context, err error, dt time.Duration)
		ResetSession(ctx context.Context, err error, dt time.Duration)

		RowsClose(ctx context.Context, err error, dt time.Duration)
		// RowsNext can receive [io.EOF] as err in the end of scanning
		RowsNext(ctx context.Context, dest []driver.Value, err error, dt time.Duration)

		PreparingStatement(ctx context.Context, query string, err error, dt time.Duration)
		ClosePreparedStatement(ctx context.Context, query string, err error, dt time.Duration)
		ExecPreparedStatement(ctx context.Context, query string, args []driver.NamedValue, replacedErr error, err error, dt time.Duration)
		QueryPreparedStatement(ctx context.Context, query string, args []driver.NamedValue, replacedErr error, err error, dt time.Duration)
	}
)
