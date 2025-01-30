package logsql

type (
	// QueryErrReplacer can be used for replacing actual errors from a driver ("Query", "QueryContext", "Exec", "ExecContext").
	// QueryErrReplacer must return non-nil error if new error should substitute the original one, nil otherwise.
	// Such substitution is useful in case of checking raw SQL errors.
	//
	// Example of handling SQL specific 42P07 error:
	//
	// 1. SQL driver query error which underlying driver usually returns
	//  type SqlError struct { Code string }
	//  func (e SqlError) Error() string {
	//	    return e.Code
	//  }
	//
	// 2. Your own specific error
	//  var errDuplicateTable = errors.New("duplicate table")
	//
	// 3. QueryErrReplacer function. If err is not a SqlError, don't substitute anything. If it is expected error 42P07
	// return your own errDuplicateTable
	//  func myErrReplacer(err error) error {
	//	    var pgErr SqlError
	//	    if !errors.As(err, &pgErr) {
	//          return nil
	//	    }
	//
	//	    switch pgErr.Code {
	//      case "42P07":
	//          return errDuplicateTable
	//      default:
	//          return nil
	//      }
	//  }
	//
	// 4. Handling. Thus, your own error can be handled in a specific way.
	// Moreover, if Logger is properly configured, no ERROR entry will appear in logs because this error was expected
	//
	//  err := db.ExecContext(ctx, query)
	//  if err != nil {
	//	    if errors.Is(err, errDuplicateTable) {
	//          // ...
	//	    }
	//	    return err
	//  }
	QueryErrReplacer func(err error) error
)

var (
	_ QueryErrReplacer = NoOpQueryErrReplacer
)

// NoOpQueryErrReplacer will not replace anything
func NoOpQueryErrReplacer(_ error) error {
	return nil
}
