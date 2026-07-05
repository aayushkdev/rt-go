package geometry

import rtmath "github.com/aayushkdev/rt-go/math"

type Translate struct {
	Object Hittable
	Offset rtmath.Vec3
	Box    AABB
}

func NewTranslate(object Hittable, offset rtmath.Vec3) Translate {
	return Translate{
		Object: object,
		Offset: offset,
		Box:    object.BoundingBox().Translate(offset),
	}
}

func (t Translate) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	movedRay := rtmath.NewRay(ray.Origin.Sub(t.Offset), ray.Direction)
	record, hit := t.Object.Hit(movedRay, rayT)
	if !hit {
		return HitRecord{}, false
	}

	record.Point = record.Point.Add(t.Offset)
	return record, true
}

func (t Translate) BoundingBox() AABB {
	return t.Box
}

func (t Translate) PDFValue(origin rtmath.Point3, direction rtmath.Vec3) float64 {
	sampler, ok := t.Object.(Sampler)
	if !ok {
		return 0
	}

	return sampler.PDFValue(origin.Sub(t.Offset), direction)
}

func (t Translate) Random(origin rtmath.Point3, random *rtmath.Random) rtmath.Vec3 {
	sampler, ok := t.Object.(Sampler)
	if !ok {
		return rtmath.NewVec3(1, 0, 0)
	}

	return sampler.Random(origin.Sub(t.Offset), random)
}
