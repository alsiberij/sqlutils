package logsql

import (
	"context"
	"database/sql/driver"
	"time"
)

type (
	connection struct {
		logHandler       Logger
		queryErrReplacer QueryErrReplacer

		conn driver.Conn
	}
)

func (c *connection) Prepare(query string) (driver.Stmt, error) {
	t0 := time.Now()

	stmt, err := c.conn.Prepare(query)
	c.logHandler.PrepareStatement(context.Background(), query, err, time.Since(t0))
	if err != nil {
		return nil, err
	}

	return &queryStatement{
		logHandler:       c.logHandler,
		queryErrReplacer: c.queryErrReplacer,
		connCtx:          context.Background(),
		query:            query,
		statement:        stmt,
	}, nil
}

func (c *connection) Close() error {
	t0 := time.Now()

	err := c.conn.Close()
	c.logHandler.ConnClose(context.Background(), err, time.Since(t0))

	return err
}

func (c *connection) Begin() (driver.Tx, error) {
	t0 := time.Now()

	tx, err := c.conn.Begin()
	c.logHandler.TxBegin(context.Background(), err, time.Since(t0))
	if err != nil {
		return nil, err
	}

	return &queryTransaction{
		logHandler:  c.logHandler,
		connCtx:     context.Background(),
		transaction: tx,
	}, err
}

func (c *connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	connBeginTx, ok := c.conn.(driver.ConnBeginTx)
	if !ok {
		return c.Begin()
	}

	t0 := time.Now()

	tx, err := connBeginTx.BeginTx(ctx, opts)
	c.logHandler.TxBegin(ctx, err, time.Since(t0))
	if err != nil {
		return nil, err
	}

	return &queryTransaction{
		logHandler:  c.logHandler,
		connCtx:     ctx,
		transaction: tx,
	}, nil
}

func (c *connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	connPrepareTx, ok := c.conn.(driver.ConnPrepareContext)
	if !ok {
		return c.Prepare(query)
	}

	t0 := time.Now()

	stmt, err := connPrepareTx.PrepareContext(ctx, query)
	c.logHandler.PrepareStatement(ctx, query, err, time.Since(t0))
	if err != nil {
		return nil, err
	}

	return &queryStatement{
		logHandler:       c.logHandler,
		queryErrReplacer: c.queryErrReplacer,
		connCtx:          ctx,
		query:            query,
		statement:        stmt,
	}, nil
}

func (c *connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	connExecerCtx, ok := c.conn.(driver.ExecerContext)
	if !ok {
		return c.Exec(query, driverNamedToValues(args))
	}

	t0 := time.Now()

	result, err := connExecerCtx.ExecContext(ctx, query, args)
	var replacedErr error
	if err != nil {
		replacedErr = c.queryErrReplacer(err)
	}
	c.logHandler.Exec(ctx, query, args, replacedErr, err, time.Since(t0))

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

func (c *connection) Exec(query string, args []driver.Value) (driver.Result, error) {
	connExecer, ok := c.conn.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	t0 := time.Now()

	result, err := connExecer.Exec(query, args)
	var replacedErr error
	if err != nil {
		replacedErr = c.queryErrReplacer(err)
	}
	c.logHandler.Exec(context.Background(), query, driverValuesToNamed(args), replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return &queryResult{
		ctx:    context.Background(),
		result: result,
	}, nil
}

func (c *connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	connQueryerCtx, ok := c.conn.(driver.QueryerContext)
	if !ok {
		return c.Query(query, driverNamedToValues(args))
	}

	t0 := time.Now()

	rows, err := connQueryerCtx.QueryContext(ctx, query, args)
	var replacedErr error
	if err != nil {
		replacedErr = c.queryErrReplacer(err)
	}
	c.logHandler.Query(context.Background(), query, args, replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return &queryRows{
		logHandler: c.logHandler,
		connCtx:    ctx,
		rows:       rows,
	}, nil
}

func (c *connection) Query(query string, args []driver.Value) (driver.Rows, error) {
	connQueryer, ok := c.conn.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	t0 := time.Now()

	rows, err := connQueryer.Query(query, args)
	var replacedErr error
	if err != nil {
		replacedErr = c.queryErrReplacer(err)
	}
	c.logHandler.Query(context.Background(), query, driverValuesToNamed(args), replacedErr, err, time.Since(t0))

	if err != nil {
		if replacedErr != nil {
			return nil, replacedErr
		}
		return nil, err
	}

	return &queryRows{
		logHandler: c.logHandler,
		connCtx:    context.Background(),
		rows:       rows,
	}, nil
}

func (c *connection) Ping(ctx context.Context) error {
	connPinger, ok := c.conn.(driver.Pinger)
	if !ok {
		return ErrUnsupportedByDriver
	}

	t0 := time.Now()

	err := connPinger.Ping(ctx)
	c.logHandler.Ping(ctx, err, time.Since(t0))

	return err
}

func (c *connection) ResetSession(ctx context.Context) error {
	connSessionResetter, ok := c.conn.(driver.SessionResetter)
	if !ok {
		return nil
	}

	return connSessionResetter.ResetSession(ctx)
}

func (c *connection) IsValid() bool {
	connValidator, ok := c.conn.(driver.Validator)
	if !ok {
		return true
	}

	return connValidator.IsValid()
}

func (c *connection) CheckNamedValue(value *driver.NamedValue) error {
	connValueChecker, ok := c.conn.(driver.NamedValueChecker)
	if !ok {
		return driver.ErrSkip
	}

	return connValueChecker.CheckNamedValue(value)
}
