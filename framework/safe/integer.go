package safe

import (
	"errors"
	"math"
	"reflect"
)

var (
	ErrOverflow = errors.New("integer overflow")
	ErrOutRange = errors.New("integer out of range")
)

type NumberInt interface {
	int | int8 | int16 | int32 | int64
}

func AddInt[T NumberInt](x, y T) (T, error) {
	if err := checkAddIntRange(x, y); err != nil {
		return T(0), err
	}

	return x + y, nil
}

func checkAddIntRange[T NumberInt](x, y T) error {
	tmp := reflect.Indirect(reflect.ValueOf(x)).Interface()
	switch tmp.(type) {
	case int:
		return checkAddIntRangeByType(x, y, math.MinInt, math.MaxInt)
	case int8:
		return checkAddIntRangeByType(x, y, math.MinInt8, math.MaxInt8)
	case int16:
		return checkAddIntRangeByType(x, y, math.MinInt16, math.MaxInt16)
	case int32:
		return checkAddIntRangeByType(x, y, math.MinInt32, math.MaxInt32)
	case int64:
		return checkAddIntRangeByType(x, y, math.MinInt64, math.MaxInt64)
	}
	return nil
}

func checkAddIntRangeByType[T NumberInt](x, y T, minVal, maxVal int) error {
	if y > 0 {
		if x > T(maxVal)-y {
			return ErrOverflow
		}
	} else {
		if x < T(minVal)-y {
			return ErrOverflow
		}
	}
	return nil
}

func Add32(x, y int32) (int32, error) {
	if y > 0 {
		if x > math.MaxInt32-y {
			return 0, ErrOverflow
		}
	} else {
		if x < math.MinInt32-y {
			return 0, ErrOverflow
		}
	}
	return x + y, nil
}

func Add64(x, y int64) (int64, error) {
	if y > 0 {
		if x > math.MaxInt64-y {
			return 0, ErrOverflow
		}
	} else {
		if x < math.MinInt64-y {
			return 0, ErrOverflow
		}
	}
	return x + y, nil
}

type NumberUint interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func AddUint[T NumberUint](x, y T) (T, error) {
	v := x + y
	if v < x || v < y {
		return 0, ErrOverflow
	}

	return v, nil
}

func Add32U(x, y uint32) (uint32, error) {
	v := x + y
	if v < x || v < y {
		return 0, ErrOverflow
	}
	return v, nil
}

func Add64U(x, y uint64) (uint64, error) {
	v := x + y
	if v < x || v < y {
		return 0, ErrOverflow
	}
	return v, nil
}

func Mul32(x, y int32) (int32, error) {
	if x == -1 && y == math.MinInt32 {
		return 0, ErrOverflow
	}
	if y == -1 && x == math.MinInt32 {
		return 0, ErrOverflow
	}
	if y != 0 {
		if x > math.MaxInt32/y {
			return 0, ErrOverflow
		}
		if x < math.MinInt32/y {
			return 0, ErrOverflow
		}
	}

	return x * y, nil
}

func Mul64(x, y int64) (int64, error) {
	if x == -1 && y == math.MinInt64 {
		return 0, ErrOverflow
	}
	if y == -1 && x == math.MinInt64 {
		return 0, ErrOverflow
	}
	if y != 0 {
		if x > math.MaxInt64/y {
			return 0, ErrOverflow
		}
		if x < math.MinInt64/y {
			return 0, ErrOverflow
		}
	}

	return x * y, nil
}

func MulInt[T NumberInt](x, y T) (T, error) {
	if err := checkMulIntRange(x, y); err != nil {
		return T(0), err
	}

	return x * y, nil
}

func checkMulIntRange[T NumberInt](x, y T) error {
	tmp := reflect.Indirect(reflect.ValueOf(x)).Interface()
	switch tmp.(type) {
	case int:
		return checkMulIntRangeByType(x, y, math.MinInt, math.MaxInt)
	case int8:
		return checkMulIntRangeByType(x, y, math.MinInt8, math.MaxInt8)
	case int16:
		return checkMulIntRangeByType(x, y, math.MinInt16, math.MaxInt16)
	case int32:
		return checkMulIntRangeByType(x, y, math.MinInt32, math.MaxInt32)
	case int64:
		return checkMulIntRangeByType(x, y, math.MinInt64, math.MaxInt64)
	}
	return nil
}

func checkMulIntRangeByType[T NumberInt](x, y T, minVal, maxVal int) error {
	//if x == -1 && y == T(minVal) {
	//	return ErrOverflow
	//}
	//if y == -1 && x == T(minVal) {
	//	return ErrOverflow
	//}
	//if y != 0 {
	//	if x > T(maxVal)/y {
	//		return ErrOverflow
	//	}
	//	if x < T(minVal)/y {
	//		return ErrOverflow
	//	}
	//}
	//
	//return nil

	if x == 0 || y == 0 || x == 1 || y == 1 || x == -1 || y == -1 {
		return nil
	}
	if x == T(minVal) || y == T(minVal) || x == T(maxVal) || y == T(maxVal) {
		return ErrOverflow
	}
	v := x * y
	if v/y != x {
		return ErrOverflow
	}

	return nil
}

func Mul32U(x, y uint32) (uint32, error) {
	if y != 0 {
		if x > math.MaxUint32/y {
			return 0, ErrOverflow
		}
	}
	return x * y, nil
}

func Mul64U(x, y uint64) (uint64, error) {
	if y != 0 {
		if x > math.MaxUint64/y {
			return 0, ErrOverflow
		}
	}
	return x * y, nil
}

func MulUint[T NumberUint](x, y T) (T, error) {
	if err := checkMulUintRange(x, y); err != nil {
		return T(0), err
	}

	return x * y, nil
}

func checkMulUintRange[T NumberUint](x, y T) error {
	tmp := reflect.Indirect(reflect.ValueOf(x)).Interface()
	switch tmp.(type) {
	case uint:
		return checkMulUintRangeByType(x, y, math.MaxUint)
	case uint8:
		return checkMulUintRangeByType(x, y, math.MaxUint8)
	case uint16:
		return checkMulUintRangeByType(x, y, math.MaxUint16)
	case uint32:
		return checkMulUintRangeByType(x, y, math.MaxUint32)
	case uint64:
		return checkMulUintRangeByType(x, y, math.MaxUint64)
	}
	return nil
}

func checkMulUintRangeByType[T NumberUint](x, y T, maxVal uint) error {
	//if y != 0 {
	//	if x > T(maxVal)/y {
	//		return ErrOverflow
	//	}
	//}
	//return nil

	if x <= 1 || y <= 1 {
		return nil
	}
	if x == T(maxVal) || y == T(maxVal) {
		return ErrOverflow
	}
	v := x * y
	if v/y != x {
		return ErrOverflow
	}
	return nil
}

func Int32From64(x int64) (int32, error) {
	if x > math.MaxInt32 || x < math.MinInt32 {
		return 0, ErrOutRange
	}
	return int32(x), nil
}

func Uint32(x uint64) (uint32, error) {
	if x > math.MaxUint32 {
		return 0, ErrOutRange
	}
	return uint32(x), nil
}
