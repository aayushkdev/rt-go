package materials

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Metal struct {
	Albedo rtmath.Color
	Fuzz   float64
}

func NewMetal(albedo rtmath.Color, fuzz float64) Metal {
	return Metal{
		Albedo: albedo,
		Fuzz:   stdmath.Min(fuzz, 1),
	}
}

func (m Metal) Scatter(rayIn rtmath.Ray, hit HitInfo) (rtmath.Color, rtmath.Ray, bool) {
	reflected := rtmath.Reflect(rayIn.Direction, hit.Normal).Unit()
	reflected = reflected.Add(rtmath.RandomUnitVector().Mul(m.Fuzz))
	scattered := rtmath.NewRay(hit.Point, reflected)

	return m.Albedo, scattered, rtmath.Dot(scattered.Direction, hit.Normal) > 0
}
