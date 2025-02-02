package logsql

import (
	"context"
	"database/sql/driver"
	"time"
)

type (
	queryStatement struct {
		logHandler       Logger
		queryErrReplacer QueryErrReplacer

		connCtx   context.Context
		query     string
		statement driver.Stmt
	}
)

func (s *queryStatement) Close() error {
	t0 := time.Now()

	err := s.statement.Close()
	s.logHandler.ClosePreparedStatement(s.connCtx, s.query, err, time.Since(t0))

	return err
}

func (s *queryStatement) NumInput() int {
	return s.statement.NumInput()
}

func (s *queryStatement) Exec(args []driver.Value) (driver.Result, error) {
	t0 := time.Now()

	result, err := s.statement.Exec(args)
	var replacedErr error
	if err != nil {
		replacedErr = s.queryErrReplacer(err)
	}
	s.logHandler.ExecPreparedStatement(s.connCtx, s.query, driverValuesToNamed(args), replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return result, nil
}

func (s *queryStatement) Query(args []driver.Value) (driver.Rows, error) {
	t0 := time.Now()

	rows, err := s.statement.Query(args)
	var replacedErr error
	if err != nil {
		replacedErr = s.queryErrReplacer(err)
	}
	s.logHandler.QueryPreparedStatement(s.connCtx, s.query, driverValuesToNamed(args), replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return rows, nil
}

func (s *queryStatement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	stExecerCtx, ok := s.statement.(driver.StmtExecContext)
	if !ok {
		return s.Exec(driverNamedToValues(args))
	}

	t0 := time.Now()

	result, err := stExecerCtx.ExecContext(ctx, args)
	var replacedErr error
	if err != nil {
		replacedErr = s.queryErrReplacer(err)
	}
	s.logHandler.ExecPreparedStatement(ctx, s.query, args, replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return &queryResult{
		ctx:    ctx,
		result: result,
	}, nil
}

func (s *queryStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	stQueryerCtx, ok := s.statement.(driver.StmtQueryContext)
	if !ok {
		return s.Query(driverNamedToValues(args))
	}

	t0 := time.Now()

	rows, err := stQueryerCtx.QueryContext(ctx, args)
	var replacedErr error
	if err != nil {
		replacedErr = s.queryErrReplacer(err)
	}
	s.logHandler.QueryPreparedStatement(ctx, s.query, args, replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return &queryRows{
		logHandler: s.logHandler,
		connCtx:    ctx,
		rows:       rows,
	}, nil
}

func (s *queryStatement) CheckNamedValue(value *driver.NamedValue) error {
	connValueChecker, ok := s.statement.(driver.NamedValueChecker)
	if !ok {
		return driver.ErrSkip
	}

	return connValueChecker.CheckNamedValue(value)
}

func (s *queryStatement) ColumnConverter(idx int) driver.ValueConverter {
	stColumnConverter, ok := s.statement.(driver.ColumnConverter)
	if !ok {
		return driver.DefaultParameterConverter
	}

	return stColumnConverter.ColumnConverter(idx)
}
