package scan

import (
	"fmt"
	"reflect"
)

const (
	StructTag = "scan"
)

var (
	_ RowColsCollector[interface{}]                = StructTagCollectorRowCols
	_ RowColsCollectorKV[interface{}, interface{}] = StructTagCollectorRowColsKV
)

// StructTagCollectorRowCols maps returned columns on T struct fields by "scan" tag and passes those fields in Scan. If no
// field is found for column error will be returned.
// Roughly equivalent to .Scan(&T.FieldForColumn1, &T.FieldForColumn2, &T.FieldForColumn3, ...)
func StructTagCollectorRowCols[T any](row RowCols) (T, error) {
	var target T

	cols, err := row.Columns()
	if err != nil {
		return target, err
	}

	val := reflect.ValueOf(&target)
	typ := reflect.TypeOf(target)

	if typ.Kind() != reflect.Struct {
		return target, ErrStructRequired
	}

	destinations := make([]interface{}, len(cols))

	for i := range destinations {
		var found bool
		for j := range typ.NumField() {
			if tag, ok := typ.Field(j).Tag.Lookup(StructTag); ok && tag != "" && tag == cols[i] {
				destinations[i] = reflect.Indirect(val).Field(j).Addr().Interface()
				found = true
				break
			}
		}
		if !found {
			return target, fmt.Errorf("destination for #%d column %s not found", i, cols[i])
		}
	}

	err = row.Scan(destinations...)
	if err != nil {
		return target, err
	}

	return target, nil
}

// StructTagCollectorRowColsKV maps returned columns on T struct fields by "scan" tag and passes those fields in Scan.
// The first column is directly scanned in K. If no field is found for column error will be returned.
// Roughly equivalent to .Scan(&K, &V.FieldForColumn2, &V.FieldForColumn3, &V.FieldForColumn4, ...)
func StructTagCollectorRowColsKV[K comparable, V any](row RowCols) (K, V, error) {
	var key K
	var value V

	cols, err := row.Columns()
	if err != nil {
		return key, value, err
	}

	vaRf := reflect.ValueOf(&value)
	typ := reflect.TypeOf(value)

	if typ.Kind() != reflect.Struct {
		return key, value, ErrStructRequired
	}

	destinations := make([]interface{}, len(cols))

	destinations[0] = &key
	for i, col := range cols[1:] {
		var found bool
		for j := range typ.NumField() {
			if tag, ok := typ.Field(j).Tag.Lookup(StructTag); ok && tag != "" && tag == col {
				destinations[i+1] = reflect.Indirect(vaRf).Field(j).Addr().Interface()
				found = true
				break
			}
		}
		if !found {
			return key, value, fmt.Errorf("destination for #%d column %s not found", i, cols[i])
		}
	}

	err = row.Scan(destinations...)
	if err != nil {
		return key, value, err
	}

	return key, value, nil
}
