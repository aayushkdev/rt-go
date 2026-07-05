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
	NormalA  rtmath.Vec3
	NormalB  rtmath.Vec3
	NormalC  rtmath.Vec3
	Smooth   bool
	Material materials.Material
}

func NewTriangle(a, b, c rtmath.Point3, material materials.Material) Triangle {
	return Triangle{A: a, B: b, C: c, Material: material}
}

func NewSmoothTriangle(a, b, c rtmath.Point3, normalA, normalB, normalC rtmath.Vec3, material materials.Material) Triangle {
	return Triangle{
		A:        a,
		B:        b,
		C:        c,
		NormalA:  rtmath.UnitVector(normalA),
		NormalB:  rtmath.UnitVector(normalB),
		NormalC:  rtmath.UnitVector(normalC),
		Smooth:   true,
		Material: material,
	}
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

	outwardNormal := t.normalAt(u, v, edgeAB, edgeAC)
	record := HitRecord{
		HitInfo: materials.HitInfo{
			T:     hitT,
			Point: ray.At(hitT),
			U:     u,
			V:     v,
		},
		Material: t.Material,
	}
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}

func (t Triangle) PDFValue(origin rtmath.Point3, direction rtmath.Vec3) float64 {
	record, hit := t.Hit(rtmath.NewRay(origin, direction), rtmath.NewInterval(0.001, stdmath.Inf(1)))
	if !hit {
		return 0
	}

	area := t.Area()
	if area == 0 {
		return 0
	}

	normal := rtmath.UnitVector(rtmath.Cross(t.B.Sub(t.A), t.C.Sub(t.A)))
	distanceSquared := record.T * record.T * direction.LengthSquared()
	cosine := stdmath.Abs(rtmath.Dot(direction, normal) / direction.Length())
	if cosine == 0 {
		return 0
	}

	return distanceSquared / (cosine * area)
}

func (t Triangle) Random(origin rtmath.Point3, random *rtmath.Random) rtmath.Vec3 {
	r1 := random.Float64()
	r2 := random.Float64()
	if r1+r2 > 1 {
		r1 = 1 - r1
		r2 = 1 - r2
	}

	point := t.A.
		Add(t.B.Sub(t.A).Mul(r1)).
		Add(t.C.Sub(t.A).Mul(r2))

	return point.Sub(origin)
}

func (t Triangle) Area() float64 {
	return rtmath.Cross(t.B.Sub(t.A), t.C.Sub(t.A)).Length() / 2
}

func (t Triangle) normalAt(u, v float64, edgeAB, edgeAC rtmath.Vec3) rtmath.Vec3 {
	if !t.Smooth {
		return rtmath.UnitVector(rtmath.Cross(edgeAB, edgeAC))
	}

	w := 1 - u - v
	normal := t.NormalA.Mul(w).
		Add(t.NormalB.Mul(u)).
		Add(t.NormalC.Mul(v))
	if normal.NearZero() {
		return rtmath.UnitVector(rtmath.Cross(edgeAB, edgeAC))
	}

	return rtmath.UnitVector(normal)
}
