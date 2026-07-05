package geometry

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type rotationAxis int

const (
	rotateX rotationAxis = iota
	rotateY
	rotateZ
)

type Rotate struct {
	Object Hittable
	Axis   rotationAxis
	Sin    float64
	Cos    float64
	Box    AABB
}

func RotateX(object Hittable, angle float64) Rotate {
	return newRotate(object, angle, rotateX)
}

func RotateY(object Hittable, angle float64) Rotate {
	return newRotate(object, angle, rotateY)
}

func RotateZ(object Hittable, angle float64) Rotate {
	return newRotate(object, angle, rotateZ)
}

func newRotate(object Hittable, angle float64, axis rotationAxis) Rotate {
	radians := angle * rtmath.Pi / 180
	rotation := Rotate{
		Object: object,
		Axis:   axis,
		Sin:    stdmath.Sin(radians),
		Cos:    stdmath.Cos(radians),
	}
	rotation.Box = rotation.rotatedBox(object.BoundingBox())

	return rotation
}

func (r Rotate) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	origin := r.rotateIntoObject(ray.Origin)
	direction := r.rotateIntoObject(ray.Direction)
	rotatedRay := rtmath.NewRay(origin, direction)

	record, hit := r.Object.Hit(rotatedRay, rayT)
	if !hit {
		return HitRecord{}, false
	}

	record.Point = r.rotateFromObject(record.Point)
	record.Normal = r.rotateFromObject(record.Normal)
	record.SetFaceNormal(ray, record.Normal)

	return record, true
}

func (r Rotate) BoundingBox() AABB {
	return r.Box
}

func (r Rotate) PDFValue(origin rtmath.Point3, direction rtmath.Vec3) float64 {
	sampler, ok := r.Object.(Sampler)
	if !ok {
		return 0
	}

	return sampler.PDFValue(r.rotateIntoObject(origin), r.rotateIntoObject(direction))
}

func (r Rotate) Random(origin rtmath.Point3, random *rtmath.Random) rtmath.Vec3 {
	sampler, ok := r.Object.(Sampler)
	if !ok {
		return rtmath.NewVec3(1, 0, 0)
	}

	direction := sampler.Random(r.rotateIntoObject(origin), random)
	return r.rotateFromObject(direction)
}

func (r Rotate) rotateIntoObject(v rtmath.Vec3) rtmath.Vec3 {
	switch r.Axis {
	case rotateX:
		return rtmath.NewVec3(
			v.X,
			r.Cos*v.Y+r.Sin*v.Z,
			-r.Sin*v.Y+r.Cos*v.Z,
		)
	case rotateY:
		return rtmath.NewVec3(
			r.Cos*v.X-r.Sin*v.Z,
			v.Y,
			r.Sin*v.X+r.Cos*v.Z,
		)
	case rotateZ:
		return rtmath.NewVec3(
			r.Cos*v.X+r.Sin*v.Y,
			-r.Sin*v.X+r.Cos*v.Y,
			v.Z,
		)
	default:
		return v
	}
}

func (r Rotate) rotateFromObject(v rtmath.Vec3) rtmath.Vec3 {
	switch r.Axis {
	case rotateX:
		return rtmath.NewVec3(
			v.X,
			r.Cos*v.Y-r.Sin*v.Z,
			r.Sin*v.Y+r.Cos*v.Z,
		)
	case rotateY:
		return rtmath.NewVec3(
			r.Cos*v.X+r.Sin*v.Z,
			v.Y,
			-r.Sin*v.X+r.Cos*v.Z,
		)
	case rotateZ:
		return rtmath.NewVec3(
			r.Cos*v.X-r.Sin*v.Y,
			r.Sin*v.X+r.Cos*v.Y,
			v.Z,
		)
	default:
		return v
	}
}

func (r Rotate) rotatedBox(box AABB) AABB {
	min := rtmath.NewVec3(stdmath.Inf(1), stdmath.Inf(1), stdmath.Inf(1))
	max := rtmath.NewVec3(stdmath.Inf(-1), stdmath.Inf(-1), stdmath.Inf(-1))

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := chooseBoxSide(box.X, i)
				y := chooseBoxSide(box.Y, j)
				z := chooseBoxSide(box.Z, k)
				corner := r.rotateFromObject(rtmath.NewVec3(x, y, z))

				min.X = stdmath.Min(min.X, corner.X)
				min.Y = stdmath.Min(min.Y, corner.Y)
				min.Z = stdmath.Min(min.Z, corner.Z)
				max.X = stdmath.Max(max.X, corner.X)
				max.Y = stdmath.Max(max.Y, corner.Y)
				max.Z = stdmath.Max(max.Z, corner.Z)
			}
		}
	}

	return NewAABB(min, max)
}

func chooseBoxSide(interval rtmath.Interval, side int) float64 {
	if side == 0 {
		return interval.Min
	}
	return interval.Max
}
