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
	Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (ScatterRecord, bool)
	ScatteringPDF(rayIn rtmath.Ray, hit HitInfo, scattered rtmath.Ray) float64
}

type ScatterRecord struct {
	Attenuation rtmath.Color
	Scattered   rtmath.Ray
	PDF         rtmath.PDF
	SkipPDF     bool
}

type BaseMaterial struct{}

func (b BaseMaterial) Emitted(hit HitInfo) rtmath.Color {
	return rtmath.NewVec3(0, 0, 0)
}

func (b BaseMaterial) ScatteringPDF(rayIn rtmath.Ray, hit HitInfo, scattered rtmath.Ray) float64 {
	return 0
}
