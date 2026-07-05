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
	record.U, record.V = sphereUV(outwardNormal)
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}

func sphereUV(point rtmath.Point3) (float64, float64) {
	theta := stdmath.Acos(-point.Y)
	phi := stdmath.Atan2(-point.Z, point.X) + rtmath.Pi

	return phi / (2 * rtmath.Pi), theta / rtmath.Pi
}

func (s Sphere) PDFValue(origin rtmath.Point3, direction rtmath.Vec3) float64 {
	_, hit := s.Hit(rtmath.NewRay(origin, direction), rtmath.NewInterval(0.001, stdmath.Inf(1)))
	if !hit {
		return 0
	}

	distanceSquared := s.Center.Sub(origin).LengthSquared()
	if distanceSquared <= s.Radius*s.Radius {
		return 1 / (4 * rtmath.Pi)
	}

	cosThetaMax := stdmath.Sqrt(1 - s.Radius*s.Radius/distanceSquared)
	solidAngle := 2 * rtmath.Pi * (1 - cosThetaMax)

	return 1 / solidAngle
}

func (s Sphere) Random(origin rtmath.Point3, random *rtmath.Random) rtmath.Vec3 {
	direction := s.Center.Sub(origin)
	distanceSquared := direction.LengthSquared()
	if distanceSquared <= s.Radius*s.Radius {
		return random.UnitVector()
	}

	uvw := rtmath.NewONBFromW(direction)
	return uvw.Local(random.ToSphere(s.Radius, distanceSquared))
}
