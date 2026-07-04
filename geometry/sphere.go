package geometry

import (
	stdmath "math"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Sphere struct {
	Center   rtmath.Point3
	Radius   float64
	Material materials.Material
}

func NewSphere(center rtmath.Point3, radius float64, material materials.Material) Sphere {
	return Sphere{Center: center, Radius: stdmath.Max(0, radius), Material: material}
}

func (s Sphere) BoundingBox() AABB {
	radius := rtmath.NewVec3(s.Radius, s.Radius, s.Radius)
	return NewAABB(s.Center.Sub(radius), s.Center.Add(radius))
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
		HitInfo: materials.HitInfo{
			T:     root,
			Point: ray.At(root),
		},
		Material: s.Material,
	}
	outwardNormal := record.Point.Sub(s.Center).Div(s.Radius)
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}
