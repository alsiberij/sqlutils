# Description
Module provides:
- Package `logsql` that contains a wrap for logging `database/sql` events. Example:
```go
package main

import (
	"database/sql"
	"database/sql/driver"
	"github.com/alsiberij/sqlutils/logsql"
	"log/slog"
)

type (
	MySlogLogger struct {
		logger *slog.Logger
	}
)

// Implement logsql.Logger interface!

func main() {
	var l *slog.Logger
	var dr driver.Driver
	var dsn string

	// Init l, dr, dsn

	loggedConnector := logsql.NewConnectorFromDriver(dr, dsn, logsql.Config{
		LogHandler: MySlogLogger{
			logger: l,
		},
	})

	db := sql.OpenDB(loggedConnector)
	// Everything you do with db will be logged via logsql.Logger interface
}
```
- Package `scan` that has useful generic utils for handling `sql.Rows`. Example:
```go
package main

import (
	"context"
	"database/sql"
	"github.com/alsiberij/sqlutils/scan"
)

type (
	User struct {
		Id    string `scan:"id"`
		Name  string `scan:"name"`
		Email string `scan:"email"`
	}
)

func getUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scan.CollectRows(rows, scan.StructTagCollectorRowCols[User])
}

func getUsersMap(ctx context.Context, db *sql.DB) (map[int]User, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scan.CollectRowsKV(rows, scan.StructTagCollectorRowColsKV[int, User])
}

func getUser(ctx context.Context, db *sql.DB, id int) (User, bool, error) {
	row := db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id = $1`, id)
	if row.Err() != nil {
		return User{}, false, row.Err()
	}

	return scan.CollectRow(row, scan.StructPosCollectorRow[User])
}
```
See more info in concrete types.

