package math

import stdrand "math/rand"

func RandomFloat64() float64 {
	return stdrand.Float64()
}

func RandomFloat64Range(min, max float64) float64 {
	return min + (max-min)*RandomFloat64()
}

func RandomVec3() Vec3 {
	return NewVec3(RandomFloat64(), RandomFloat64(), RandomFloat64())
}

func RandomVec3Range(min, max float64) Vec3 {
	return NewVec3(
		RandomFloat64Range(min, max),
		RandomFloat64Range(min, max),
		RandomFloat64Range(min, max),
	)
}

func RandomUnitVector() Vec3 {
	for {
		p := RandomVec3Range(-1, 1)
		lengthSquared := p.LengthSquared()
		if 1e-160 < lengthSquared && lengthSquared <= 1 {
			return p.Unit()
		}
	}
}

func RandomInUnitDisk() Vec3 {
	for {
		p := NewVec3(RandomFloat64Range(-1, 1), RandomFloat64Range(-1, 1), 0)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}
