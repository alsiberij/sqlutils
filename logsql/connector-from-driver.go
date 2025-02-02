package logsql

import (
	"context"
	"database/sql/driver"
	"time"
)

// NewConnectorFromDriver returns new [driver.Connector] based on an existing driver and DSN. Panics if [Config.Validate]
// returns non-nil error
func NewConnectorFromDriver(d driver.Driver, dsn string, cfg Config) driver.Connector {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	if cfg.Qer == nil {
		cfg.Qer = NoOpQueryErrReplacer
	}

	return &connectorFromDriver{
		logHandler:       cfg.LogHandler,
		queryErrReplacer: cfg.Qer,
		drv:              d,
		dsn:              dsn,
	}
}

type (
	connectorFromDriver struct {
		logHandler       Logger
		queryErrReplacer QueryErrReplacer

		drv driver.Driver
		dsn string
	}
)

func (c *connectorFromDriver) Connect(ctx context.Context) (driver.Conn, error) {
	t0 := time.Now()

	conn, err := c.drv.Open(c.dsn)
	c.logHandler.Connect(ctx, err, time.Since(t0))
	if err != nil {
		return nil, err
	}

	return &connection{
		logHandler:       c.logHandler,
		queryErrReplacer: c.queryErrReplacer,
		conn:             conn,
	}, nil
}

func (c *connectorFromDriver) Driver() driver.Driver {
	return c.drv
}
