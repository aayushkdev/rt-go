package geometry

import (
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type HitRecord struct {
	materials.HitInfo
	Material materials.Material
}

func (h *HitRecord) SetFaceNormal(ray rtmath.Ray, outwardNormal rtmath.Vec3) {
	h.FrontFace = rtmath.Dot(ray.Direction, outwardNormal) < 0
	if h.FrontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.Neg()
	}
}

type Hittable interface {
	Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool)
	BoundingBox() AABB
}
