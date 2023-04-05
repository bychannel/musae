package safe

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestAdd32CostTime(t *testing.T) {
	count := 10000000
	var x int32 = 100

	start := time.Now().UnixMilli()
	for i := 0; i < count; i++ {
		_ = reflect.Indirect(reflect.ValueOf(x)).Interface()
	}
	end := time.Now().UnixMilli()

	// 反射10000000次, 耗时:51 毫秒
	fmt.Println(fmt.Sprintf("反射%d次, 耗时:%d 毫秒", count, end-start))
}

func TestAddInt(t *testing.T) {
	//var aInt int = 10000000000
	//var bInt int = 1200000000

	// 范围 [-128, 127]
	//var aInt int8 = 127
	//var bInt int8 = 1

	//// 范围 [-32768, 32767]
	//var aInt int16 = -32767
	//var bInt int16 = -2

	//// 范围 [-2147483648, 2147483647]
	//var aInt int32 = 2147483647
	//var bInt int32 = -2

	// 范围 [-9223372036854775808, 9223372036854775807]
	var aInt int64 = math.MinInt64
	var bInt int64 = -2

	typ := fmt.Sprintf("类型:%T, (%d) + (%d) = ", aInt, aInt, bInt)

	result, err := AddInt(aInt, bInt)
	if err != nil {
		fmt.Println(typ, err.Error())
		return
	}

	fmt.Println(typ, fmt.Sprintf("%d", result))
}

func TestAddUint(t *testing.T) {
	//// 范围 [0, 18446744073709551615]
	//var aInt uint = math.MaxUint
	//var bInt uint = 1

	//// 范围 [0, 255]
	//var aInt uint8 = math.MaxUint8
	//var bInt uint8 = 1

	//// 范围 [0, 65535]
	//var aInt uint16 = math.MaxUint16
	//var bInt uint16 = 1

	// 范围 [0, 4294967295]
	var aInt uint32 = math.MaxUint32
	var bInt uint32 = 1

	//// 范围 [0, 9223372036854775807]
	//var aInt uint64 = math.MaxUint64
	//var bInt uint64 = 1

	typ := fmt.Sprintf("类型:%T, (%d) + (%d) = ", aInt, aInt, bInt)

	result, err := AddUint(aInt, bInt)
	if err != nil {
		fmt.Println(typ, err.Error())
		return
	}

	fmt.Println(typ, fmt.Sprintf("%d", result))
}

func TestMulInt(t *testing.T) {
	// 范围 [-9223372036854775808, 9223372036854775807]
	var aInt int = math.MinInt
	var bInt int = 2

	//// 范围 [-128, 127]
	//var aInt int8 = math.MaxInt8
	//var bInt int8 = 2

	//// 范围 [-32768, 32767]
	//var aInt int16 = math.MinInt16
	//var bInt int16 = -2

	//// 范围 [-2147483648, 2147483647]
	//var aInt int32 = math.MaxInt32
	//var bInt int32 = -2

	//// 范围 [-9223372036854775808, 9223372036854775807]
	//var aInt int64 = math.MinInt64
	//var bInt int64 = -2

	typ := fmt.Sprintf("类型:%T, (%d) * (%d) = ", aInt, aInt, bInt)

	result, err := MulInt(aInt, bInt)
	if err != nil {
		fmt.Println(typ, err.Error())
		return
	}

	fmt.Println(typ, fmt.Sprintf("%d", result))
}

func TestMulUint(t *testing.T) {
	// 范围 [0, 18446744073709551615]
	//var aInt uint = math.MaxUint
	//var bInt uint = 2

	//// 范围 [0, 255]
	//var aInt uint8 = math.MaxUint8
	//var bInt uint8 = 0

	//// 范围 [0, 65535]
	//var aInt uint16 = math.MaxUint16
	//var bInt uint16 = 2

	// 范围 [0, 4294967295]
	var aInt uint32 = math.MaxUint32
	var bInt uint32 = 2

	//// 范围 [0, 9223372036854775807]
	//var aInt uint64 = math.MaxUint64
	//var bInt uint64 = 2

	typ := fmt.Sprintf("类型:%T, (%d) * (%d) = ", aInt, aInt, bInt)

	result, err := MulUint(aInt, bInt)
	if err != nil {
		fmt.Println(typ, err.Error())
		return
	}

	fmt.Println(typ, fmt.Sprintf("%d", result))
}
