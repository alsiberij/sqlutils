package logsql

import (
	"database/sql/driver"
)

func driverValuesToNamed(args []driver.Value) []driver.NamedValue {
	result := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		result[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   arg,
		}
	}
	return result
}

func driverNamedToValues(args []driver.NamedValue) []driver.Value {
	result := make([]driver.Value, len(args))
	for i, arg := range args {
		result[i] = arg.Value
	}
	return result
}
