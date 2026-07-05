package geometry

import rtmath "github.com/aayushkdev/rt-go/math"

type HittablePDF struct {
	Object Sampler
	Origin rtmath.Point3
}

func NewHittablePDF(object Sampler, origin rtmath.Point3) HittablePDF {
	return HittablePDF{Object: object, Origin: origin}
}

func (p HittablePDF) Value(direction rtmath.Vec3) float64 {
	return p.Object.PDFValue(p.Origin, direction)
}

func (p HittablePDF) Generate(random *rtmath.Random) rtmath.Vec3 {
	return p.Object.Random(p.Origin, random)
}
