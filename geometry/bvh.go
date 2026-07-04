package geometry

import (
	"sort"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type BVHNode struct {
	Left  Hittable
	Right Hittable
	Box   AABB
}

func NewBVH(objects []Hittable) Hittable {
	if len(objects) == 0 {
		return EmptyHittable{}
	}

	items := append([]Hittable(nil), objects...)
	return buildBVH(items)
}

func buildBVH(items []Hittable) Hittable {
	if len(items) == 1 {
		return items[0]
	}

	box := items[0].BoundingBox()
	for i := 1; i < len(items); i++ {
		box = SurroundingBox(box, items[i].BoundingBox())
	}

	axis := box.LongestAxis()
	sort.Slice(items, func(i, j int) bool {
		return items[i].BoundingBox().Axis(axis).Min < items[j].BoundingBox().Axis(axis).Min
	})

	mid := len(items) / 2
	left := buildBVH(items[:mid])
	right := buildBVH(items[mid:])

	return BVHNode{
		Left:  left,
		Right: right,
		Box:   SurroundingBox(left.BoundingBox(), right.BoundingBox()),
	}
}

func (b BVHNode) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	if !b.Box.Hit(ray, rayT) {
		return HitRecord{}, false
	}

	leftRecord, hitLeft := b.Left.Hit(ray, rayT)
	if hitLeft {
		rayT.Max = leftRecord.T
	}

	rightRecord, hitRight := b.Right.Hit(ray, rayT)
	if hitRight {
		return rightRecord, true
	}

	return leftRecord, hitLeft
}

func (b BVHNode) BoundingBox() AABB {
	return b.Box
}

type EmptyHittable struct{}

func (e EmptyHittable) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	return HitRecord{}, false
}

func (e EmptyHittable) BoundingBox() AABB {
	return AABB{}
}
