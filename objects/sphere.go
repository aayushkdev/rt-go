package objects

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Sphere struct {
	Center rtmath.Point3
	Radius float64
}

func NewSphere(center rtmath.Point3, radius float64) Sphere {
	return Sphere{Center: center, Radius: stdmath.Max(0, radius)}
}

func (s Sphere) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	oc := s.Center.Sub(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := rtmath.Dot(ray.Direction, oc)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := h*h - a*c

	if discriminant < 0 {
		return HitRecord{}, false
	}

	sqrtd := stdmath.Sqrt(discriminant)
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return HitRecord{}, false
		}
	}

	record := HitRecord{
		T:     root,
		Point: ray.At(root),
	}
	outwardNormal := record.Point.Sub(s.Center).Div(s.Radius)
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}
