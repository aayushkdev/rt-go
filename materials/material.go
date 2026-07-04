package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type HitInfo struct {
	Point     rtmath.Point3
	Normal    rtmath.Vec3
	T         float64
	FrontFace bool
}

type Material interface {
	Scatter(rayIn rtmath.Ray, hit HitInfo) (rtmath.Color, rtmath.Ray, bool)
}
