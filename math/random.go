package math

import stdrand "math/rand"

func RandomFloat64() float64 {
	return stdrand.Float64()
}

func RandomFloat64Range(min, max float64) float64 {
	return min + (max-min)*RandomFloat64()
}
