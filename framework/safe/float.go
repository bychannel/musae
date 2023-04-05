package safe

import "math"

const FloatE2 float64 = 1e-2
const FloatE4 float64 = 1e-4
const FloatE6 float64 = 1e-6
const FloatE8 float64 = 1e-8

func Add32F(x, y float32) (float32, error) {
	if math.MaxFloat32-x > y {
		return 0, ErrOverflow
	}
	return x + y, nil
}

func Add64F(x, y float64) (float64, error) {
	if math.MaxFloat64-x > y {
		return 0, ErrOverflow
	}
	return x + y, nil
}
