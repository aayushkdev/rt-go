package geometry

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type AABB struct {
	X rtmath.Interval
	Y rtmath.Interval
	Z rtmath.Interval
}

func NewAABB(a, b rtmath.Point3) AABB {
	return AABB{
		X: rtmath.NewInterval(stdmath.Min(a.X, b.X), stdmath.Max(a.X, b.X)),
		Y: rtmath.NewInterval(stdmath.Min(a.Y, b.Y), stdmath.Max(a.Y, b.Y)),
		Z: rtmath.NewInterval(stdmath.Min(a.Z, b.Z), stdmath.Max(a.Z, b.Z)),
	}.Pad()
}

func SurroundingBox(a, b AABB) AABB {
	return AABB{
		X: rtmath.NewInterval(stdmath.Min(a.X.Min, b.X.Min), stdmath.Max(a.X.Max, b.X.Max)),
		Y: rtmath.NewInterval(stdmath.Min(a.Y.Min, b.Y.Min), stdmath.Max(a.Y.Max, b.Y.Max)),
		Z: rtmath.NewInterval(stdmath.Min(a.Z.Min, b.Z.Min), stdmath.Max(a.Z.Max, b.Z.Max)),
	}
}

func (box AABB) Translate(offset rtmath.Vec3) AABB {
	return AABB{
		X: rtmath.NewInterval(box.X.Min+offset.X, box.X.Max+offset.X),
		Y: rtmath.NewInterval(box.Y.Min+offset.Y, box.Y.Max+offset.Y),
		Z: rtmath.NewInterval(box.Z.Min+offset.Z, box.Z.Max+offset.Z),
	}
}

func (box AABB) Hit(ray rtmath.Ray, rayT rtmath.Interval) bool {
	for axis := 0; axis < 3; axis++ {
		axisInterval := box.Axis(axis)
		origin := axisValue(ray.Origin, axis)
		direction := axisValue(ray.Direction, axis)
		if direction == 0 {
			if !axisInterval.Contains(origin) {
				return false
			}
			continue
		}

		invDirection := 1.0 / direction
		t0 := (axisInterval.Min - origin) * invDirection
		t1 := (axisInterval.Max - origin) * invDirection
		if invDirection < 0 {
			t0, t1 = t1, t0
		}

		rayT.Min = stdmath.Max(t0, rayT.Min)
		rayT.Max = stdmath.Min(t1, rayT.Max)
		if rayT.Max <= rayT.Min {
			return false
		}
	}

	return true
}

func (box AABB) Axis(axis int) rtmath.Interval {
	switch axis {
	case 0:
		return box.X
	case 1:
		return box.Y
	default:
		return box.Z
	}
}

func (box AABB) LongestAxis() int {
	if box.X.Size() > box.Y.Size() {
		if box.X.Size() > box.Z.Size() {
			return 0
		}
		return 2
	}
	if box.Y.Size() > box.Z.Size() {
		return 1
	}
	return 2
}

func (box AABB) Pad() AABB {
	const delta = 0.0001

	if box.X.Size() < delta {
		box.X = expandInterval(box.X, delta)
	}
	if box.Y.Size() < delta {
		box.Y = expandInterval(box.Y, delta)
	}
	if box.Z.Size() < delta {
		box.Z = expandInterval(box.Z, delta)
	}

	return box
}

func expandInterval(interval rtmath.Interval, delta float64) rtmath.Interval {
	padding := delta / 2
	return rtmath.NewInterval(interval.Min-padding, interval.Max+padding)
}

func axisValue(v rtmath.Vec3, axis int) float64 {
	switch axis {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		return v.Z
	}
}
