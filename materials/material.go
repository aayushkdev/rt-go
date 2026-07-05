package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type HitInfo struct {
	Point     rtmath.Point3
	Normal    rtmath.Vec3
	T         float64
	FrontFace bool
}

type Material interface {
	Emitted(hit HitInfo) rtmath.Color
	Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (rtmath.Color, rtmath.Ray, bool)
}

type BaseMaterial struct{}

func (b BaseMaterial) Emitted(hit HitInfo) rtmath.Color {
	return rtmath.NewVec3(0, 0, 0)
}
