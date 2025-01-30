// Package logsql provides wrapper for [database/sql] with Logger interface. To create new logged *sql.DB use either
// NewConnectorFromDriver or NewConnectorFromConnector to retrieve driver.Connector and pass it to sql.OpenDB.
package logsql
