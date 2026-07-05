package geometry

import (
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func NewBox(min, max rtmath.Point3, material materials.Material) World {
	dx := rtmath.NewVec3(max.X-min.X, 0, 0)
	dy := rtmath.NewVec3(0, max.Y-min.Y, 0)
	dz := rtmath.NewVec3(0, 0, max.Z-min.Z)

	return NewWorld(
		NewQuad(rtmath.NewVec3(min.X, min.Y, max.Z), dx, dy, material),
		NewQuad(rtmath.NewVec3(max.X, min.Y, min.Z), dx.Neg(), dy, material),
		NewQuad(rtmath.NewVec3(min.X, max.Y, max.Z), dx, dz.Neg(), material),
		NewQuad(rtmath.NewVec3(min.X, min.Y, min.Z), dx, dz, material),
		NewQuad(rtmath.NewVec3(min.X, min.Y, min.Z), dz, dy, material),
		NewQuad(rtmath.NewVec3(max.X, min.Y, max.Z), dz.Neg(), dy, material),
	)
}
