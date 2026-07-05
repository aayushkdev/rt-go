package math

import (
	stdmath "math"
	stdrand "math/rand"
)

type Random struct {
	source *stdrand.Rand
}

func NewRandom(seed int64) *Random {
	return &Random{source: stdrand.New(stdrand.NewSource(seed))}
}

func (r *Random) Float64() float64 {
	if r == nil {
		return RandomFloat64()
	}
	return r.source.Float64()
}

func (r *Random) Float64Range(min, max float64) float64 {
	return min + (max-min)*r.Float64()
}

func (r *Random) Vec3Range(min, max float64) Vec3 {
	return NewVec3(
		r.Float64Range(min, max),
		r.Float64Range(min, max),
		r.Float64Range(min, max),
	)
}

func (r *Random) UnitVector() Vec3 {
	for {
		p := r.Vec3Range(-1, 1)
		lengthSquared := p.LengthSquared()
		if 1e-160 < lengthSquared && lengthSquared <= 1 {
			return p.Unit()
		}
	}
}

func (r *Random) InUnitDisk() Vec3 {
	for {
		p := NewVec3(r.Float64Range(-1, 1), r.Float64Range(-1, 1), 0)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func (r *Random) CosineDirection() Vec3 {
	r1 := r.Float64()
	r2 := r.Float64()
	phi := 2 * Pi * r1
	x := stdmath.Cos(phi) * stdmath.Sqrt(r2)
	y := stdmath.Sin(phi) * stdmath.Sqrt(r2)
	z := stdmath.Sqrt(1 - r2)

	return NewVec3(x, y, z)
}

func (r *Random) ToSphere(radius, distanceSquared float64) Vec3 {
	r1 := r.Float64()
	r2 := r.Float64()
	z := 1 + r2*(stdmath.Sqrt(1-radius*radius/distanceSquared)-1)
	phi := 2 * Pi * r1
	x := stdmath.Cos(phi) * stdmath.Sqrt(1-z*z)
	y := stdmath.Sin(phi) * stdmath.Sqrt(1-z*z)

	return NewVec3(x, y, z)
}

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
