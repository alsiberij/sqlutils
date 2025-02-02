package logsql

import (
	"context"
	"database/sql/driver"
	"time"
)

// NewConnectorFromConnector returns new [driver.Connector] based on existing connector. Panics if [Config.Validate]
// returns non nil error
func NewConnectorFromConnector(connector driver.Connector, cfg Config) driver.Connector {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	if cfg.Qer == nil {
		cfg.Qer = NoOpQueryErrReplacer
	}

	return &connectorFromConnector{
		logHandler:       cfg.LogHandler,
		queryErrReplacer: cfg.Qer,
		connector:        connector,
	}
}

type (
	connectorFromConnector struct {
		logHandler       Logger
		queryErrReplacer QueryErrReplacer

		connector driver.Connector
	}
)

func (c *connectorFromConnector) Connect(ctx context.Context) (driver.Conn, error) {
	t0 := time.Now()

	conn, err := c.connector.Connect(ctx)
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

func (c *connectorFromConnector) Driver() driver.Driver {
	return c.connector.Driver()
}
