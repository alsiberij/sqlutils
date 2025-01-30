package scan

var (
	_ RowColsCollector[interface{}]                = DirectCollectorRowCols
	_ RowColsCollectorKV[interface{}, interface{}] = DirectCollectorRowColsKV
	_ RowCollector[interface{}]                    = DirectCollectorRow
)

// DirectCollectorRowCols simply passes a T pointer into Scan. T can be a primitive type or sql.Scanner implementation.
// Roughly equivalent to .Scan(&T)
func DirectCollectorRowCols[T any](row RowCols) (T, error) {
	var target T

	err := row.Scan(&target)
	if err != nil {
		return target, err
	}

	return target, nil
}

// DirectCollectorRowColsKV simply passes a K pointer and a V pointer into Scan. K and V can be primitive types or sql.Scanner implementations.
// Roughly equivalent to .Scan(&K, &V)
func DirectCollectorRowColsKV[K comparable, V any](row RowCols) (K, V, error) {
	var key K
	var value V

	err := row.Scan(&key, &value)
	if err != nil {
		return key, value, err
	}

	return key, value, nil
}

// DirectCollectorRow simply passes a T pointer into Scan. T can be a primitive type or sql.Scanner implementation.
// Roughly equivalent to .Scan(&T)
func DirectCollectorRow[T any](row Row) (T, error) {
	var target T

	err := row.Scan(&target)
	if err != nil {
		return target, err
	}

	return target, nil
}
