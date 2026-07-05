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
