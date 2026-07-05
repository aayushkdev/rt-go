package geometry

import (
	stdmath "math"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Quad struct {
	Q        rtmath.Point3
	U        rtmath.Vec3
	V        rtmath.Vec3
	Normal   rtmath.Vec3
	D        float64
	W        rtmath.Vec3
	Box      AABB
	Material materials.Material
}

func NewQuad(q rtmath.Point3, u, v rtmath.Vec3, material materials.Material) Quad {
	normal := rtmath.Cross(u, v)
	unitNormal := rtmath.UnitVector(normal)

	return Quad{
		Q:        q,
		U:        u,
		V:        v,
		Normal:   unitNormal,
		D:        rtmath.Dot(unitNormal, q),
		W:        normal.Div(rtmath.Dot(normal, normal)),
		Box:      SurroundingBox(NewAABB(q, q.Add(u).Add(v)), NewAABB(q.Add(u), q.Add(v))),
		Material: material,
	}
}

func (q Quad) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	const epsilon = 1e-8

	denominator := rtmath.Dot(q.Normal, ray.Direction)
	if stdmath.Abs(denominator) < epsilon {
		return HitRecord{}, false
	}

	t := (q.D - rtmath.Dot(q.Normal, ray.Origin)) / denominator
	if !rayT.Contains(t) {
		return HitRecord{}, false
	}

	intersection := ray.At(t)
	planarHit := intersection.Sub(q.Q)
	alpha := rtmath.Dot(q.W, rtmath.Cross(planarHit, q.V))
	beta := rtmath.Dot(q.W, rtmath.Cross(q.U, planarHit))
	if !rtmath.NewInterval(0, 1).Contains(alpha) || !rtmath.NewInterval(0, 1).Contains(beta) {
		return HitRecord{}, false
	}

	record := HitRecord{
		HitInfo: materials.HitInfo{
			T:     t,
			Point: intersection,
		},
		Material: q.Material,
	}
	record.SetFaceNormal(ray, q.Normal)

	return record, true
}

func (q Quad) BoundingBox() AABB {
	return q.Box
}

func (q Quad) PDFValue(origin rtmath.Point3, direction rtmath.Vec3) float64 {
	record, hit := q.Hit(rtmath.NewRay(origin, direction), rtmath.NewInterval(0.001, stdmath.Inf(1)))
	if !hit {
		return 0
	}

	area := rtmath.Cross(q.U, q.V).Length()
	distanceSquared := record.T * record.T * direction.LengthSquared()
	cosine := stdmath.Abs(rtmath.Dot(direction, q.Normal) / direction.Length())
	if cosine == 0 {
		return 0
	}

	return distanceSquared / (cosine * area)
}

func (q Quad) Random(origin rtmath.Point3, random *rtmath.Random) rtmath.Vec3 {
	point := q.Q.
		Add(q.U.Mul(random.Float64())).
		Add(q.V.Mul(random.Float64()))

	return point.Sub(origin)
}
