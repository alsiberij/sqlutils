package logsql

import "errors"

var (
	ErrNilLogHandler       = errors.New("log handler is nil")
	ErrUnsupportedByDriver = errors.New("unsupported by underlying driver")
)
