package geometry

import (
	stdmath "math"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Triangle struct {
	A        rtmath.Point3
	B        rtmath.Point3
	C        rtmath.Point3
	Material materials.Material
}

func NewTriangle(a, b, c rtmath.Point3, material materials.Material) Triangle {
	return Triangle{A: a, B: b, C: c, Material: material}
}

func (t Triangle) BoundingBox() AABB {
	return SurroundingBox(NewAABB(t.A, t.B), NewAABB(t.A, t.C))
}

func (t Triangle) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	const epsilon = 1e-8

	edgeAB := t.B.Sub(t.A)
	edgeAC := t.C.Sub(t.A)
	pvec := rtmath.Cross(ray.Direction, edgeAC)
	determinant := rtmath.Dot(edgeAB, pvec)
	if stdmath.Abs(determinant) < epsilon {
		return HitRecord{}, false
	}

	invDeterminant := 1.0 / determinant
	tvec := ray.Origin.Sub(t.A)
	u := rtmath.Dot(tvec, pvec) * invDeterminant
	if u < 0 || u > 1 {
		return HitRecord{}, false
	}

	qvec := rtmath.Cross(tvec, edgeAB)
	v := rtmath.Dot(ray.Direction, qvec) * invDeterminant
	if v < 0 || u+v > 1 {
		return HitRecord{}, false
	}

	hitT := rtmath.Dot(edgeAC, qvec) * invDeterminant
	if !rayT.Surrounds(hitT) {
		return HitRecord{}, false
	}

	outwardNormal := rtmath.UnitVector(rtmath.Cross(edgeAB, edgeAC))
	record := HitRecord{
		HitInfo: materials.HitInfo{
			T:     hitT,
			Point: ray.At(hitT),
		},
		Material: t.Material,
	}
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}
