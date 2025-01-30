package scan

import (
	"reflect"
)

var (
	_ RowColsCollector[interface{}]                = StructPosCollectorRowCols
	_ RowCollector[interface{}]                    = StructPosCollectorRow
	_ RowColsCollectorKV[interface{}, interface{}] = StructPosCollectorRowColsKV
)

// StructPosCollectorRowCols traverses struct T and uses its field as destinations for Scan. Order of struct fields must
// be the same as order of returned columns. Roughly equivalent to .Scan(&T.Field1, &T.Field2, &T.Field3, ...)
func StructPosCollectorRowCols[T any](row RowCols) (T, error) {
	var target T

	valRf := reflect.ValueOf(&target)

	if valRf.Elem().Kind() != reflect.Struct {
		return target, ErrStructRequired
	}

	destinations := make([]interface{}, valRf.Elem().Type().NumField())
	for i := range destinations {
		destinations[i] = reflect.Indirect(valRf).Field(i).Addr().Interface()
	}

	err := row.Scan(destinations...)
	if err != nil {
		return target, err
	}

	return target, nil
}

// StructPosCollectorRow traverses struct T and uses its field as destinations for Scan. Order of struct fields must
// be the same as order of returned columns. Roughly equivalent to .Scan(&T.Field1, &T.Field2, &T.Field3, ...)
func StructPosCollectorRow[T any](row Row) (T, error) {
	var target T

	valRf := reflect.ValueOf(&target)

	if valRf.Elem().Kind() != reflect.Struct {
		return target, ErrStructRequired
	}

	destinations := make([]interface{}, valRf.Elem().Type().NumField())
	for i := range destinations {
		destinations[i] = reflect.Indirect(valRf).Field(i).Addr().Interface()
	}

	err := row.Scan(destinations...)
	if err != nil {
		return target, err
	}

	return target, nil
}

// StructPosCollectorRowColsKV traverses struct V and uses its field as destinations for Scan. The first column is directly scanned in K,
// Order of struct fields must be the same as order of returned columns.
// Roughly equivalent to .Scan(&K, &V.Field1, &V.Field2, &V.Field3, ...)
func StructPosCollectorRowColsKV[K comparable, V any](row RowCols) (K, V, error) {
	var key K
	var value V

	valRf := reflect.ValueOf(&value)

	if valRf.Elem().Kind() != reflect.Struct {
		return key, value, ErrStructRequired
	}

	destinations := make([]interface{}, valRf.Elem().Type().NumField()+1)
	destinations[0] = &key

	for i := range destinations[1:] {
		destinations[i+1] = reflect.Indirect(valRf).Field(i).Addr().Interface()
	}

	err := row.Scan(destinations...)
	if err != nil {
		return key, value, err
	}

	return key, value, nil
}
