package logsql

type (
	// Config must contain Logger for logging sql.DB events.
	// If Qer is nil, NoOpQueryErrReplacer will be used
	Config struct {
		Qer        QueryErrReplacer
		LogHandler Logger
	}
)

// Validate returns error if some fields are nil
func (c Config) Validate() error {
	if c.LogHandler == nil {
		return ErrNilLogHandler
	}
	if c.Qer == nil {
		c.Qer = NoOpQueryErrReplacer
	}

	return nil
}
