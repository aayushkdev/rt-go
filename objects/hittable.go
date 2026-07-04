package objects

import rtmath "github.com/aayushkdev/rt-go/math"

type HitRecord struct {
	Point     rtmath.Point3
	Normal    rtmath.Vec3
	T         float64
	FrontFace bool
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
	Hit(ray rtmath.Ray, tMin, tMax float64) (HitRecord, bool)
}
