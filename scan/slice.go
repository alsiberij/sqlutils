package scan

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	nullArrayElem = "NULL"
)

var (
	ErrSrcIsNotArray              = errors.New("src is not an array")
	ErrSliceConversionUnsupported = errors.New("slice conversion is unsupported")
	ErrNilUnderlyingSlice         = errors.New("underlying slice is nil")
)

type (
	// Slice is an implementation of sql.Scanner for SQL arrays. Use NewSlice to get a new Slice.
	Slice struct {
		values interface{}
	}
)

// NewSlice returns new Slice with underlying values
func NewSlice[T any](values *[]T) Slice {
	return Slice{
		values: values,
	}
}

func (s Slice) Scan(src any) error {
	var arr string
	var isArr bool

	if arr, isArr = src.(string); !(isArr && len(arr) >= 2 && arr[0] == '{' && arr[len(arr)-1] == '}') {
		return ErrSrcIsNotArray
	}

	arr = arr[1 : len(arr)-1]
	if len(arr) == 0 {
		return nil
	}

	if s.values == nil {
		panic(ErrNilUnderlyingSlice)
	}

	arrValues := strings.Split(arr, ",")

	var err error
	switch ptr := s.values.(type) {
	case *[]int64:
		err = convertInt64(arrValues, ptr)
	case *[]*int64:
		err = convertInt64Ptr(arrValues, ptr)
	case *[]int:
		err = convertInt(arrValues, ptr)
	case *[]*int:
		err = convertIntPtr(arrValues, ptr)
	case *[]uint:
		err = convertUint(arrValues, ptr)
	case *[]*uint:
		err = convertUintPtr(arrValues, ptr)
	case *[]uint64:
		err = convertUint64(arrValues, ptr)
	case *[]*uint64:
		err = convertUint64Ptr(arrValues, ptr)
	case *[]float64:
		err = convertFloat64(arrValues, ptr)
	case *[]*float64:
		err = convertFloat64Ptr(arrValues, ptr)
	case *[]bool:
		err = convertBool(arrValues, ptr)
	case *[]*bool:
		err = convertBoolPtr(arrValues, ptr)
	case *[]string:
		err = convertString(arrValues, ptr)
	case *[]*string:
		err = convertStringPtr(arrValues, ptr)
	case *[]time.Time:
		err = convertTime(arrValues, ptr)
	case *[]*time.Time:
		err = convertTimePtr(arrValues, ptr)
	default:
		return ErrSliceConversionUnsupported
	}
	if err != nil {
		return err
	}

	return nil
}

func convertInt64(ss []string, dst *[]int64) error {
	for _, s := range ss {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, v)
	}

	return nil
}

func convertInt64Ptr(ss []string, dst *[]*int64) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, &v)
	}

	return nil
}

func convertInt(ss []string, dst *[]int) error {
	for _, s := range ss {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, int(v))
	}

	return nil
}

func convertIntPtr(ss []string, dst *[]*int) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		vInt := int(v)
		*dst = append(*dst, &vInt)
	}

	return nil
}

func convertUint64(ss []string, dst *[]uint64) error {
	for _, s := range ss {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, v)
	}

	return nil
}

func convertUint64Ptr(ss []string, dst *[]*uint64) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, &v)
	}

	return nil
}

func convertUint(ss []string, dst *[]uint) error {
	for _, s := range ss {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, uint(v))
	}

	return nil
}

func convertUintPtr(ss []string, dst *[]*uint) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		vInt := uint(v)
		*dst = append(*dst, &vInt)
	}

	return nil
}

func convertFloat64(ss []string, dst *[]float64) error {
	for _, s := range ss {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, v)
	}

	return nil
}

func convertFloat64Ptr(ss []string, dst *[]*float64) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*dst = append(*dst, &v)
	}

	return nil
}

func convertBool(ss []string, dst *[]bool) error {
	for _, s := range ss {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*dst = append(*dst, v)
	}

	return nil
}

func convertBoolPtr(ss []string, dst *[]*bool) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*dst = append(*dst, &v)
	}

	return nil
}

func convertString(ss []string, dst *[]string) error {
	for _, s := range ss {
		*dst = append(*dst, s)
	}

	return nil
}

func convertStringPtr(ss []string, dst *[]*string) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		scp := s
		*dst = append(*dst, &scp)
	}

	return nil
}

func convertTime(ss []string, dst *[]time.Time) error {
	for _, s := range ss {
		v, err := time.Parse("2006-01-02 15:04:05.999999999-07", s[1:len(s)-1])
		if err != nil {
			return err
		}
		*dst = append(*dst, v)
	}

	return nil
}

func convertTimePtr(ss []string, dst *[]*time.Time) error {
	for _, s := range ss {
		if s == nullArrayElem {
			*dst = append(*dst, nil)
			continue
		}
		v, err := time.Parse("2006-01-02 15:04:05.999999999-07", s[1:len(s)-1])
		if err != nil {
			return err
		}
		*dst = append(*dst, &v)
	}

	return nil
}
