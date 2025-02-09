package logsql

type (
	// Config must contain Logger for logging sql.DB events.
	// If Qer is nil, NoOpQueryErrReplacer will be used
	Config struct {
		Qer        QueryErrReplacer
		LogHandler Logger
	}
)

// Validate returns error if config is not ready to use
func (c Config) Validate() error {
	if c.LogHandler == nil {
		return ErrNilLogHandler
	}

	return nil
}
