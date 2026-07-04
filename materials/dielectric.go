package materials

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Dielectric struct {
	RefractionIndex float64
}

func NewDielectric(refractionIndex float64) Dielectric {
	return Dielectric{RefractionIndex: refractionIndex}
}

func (d Dielectric) Scatter(rayIn rtmath.Ray, hit HitInfo) (rtmath.Color, rtmath.Ray, bool) {
	attenuation := rtmath.NewVec3(1, 1, 1)
	refractionRatio := d.RefractionIndex
	if hit.FrontFace {
		refractionRatio = 1.0 / d.RefractionIndex
	}

	unitDirection := rtmath.UnitVector(rayIn.Direction)
	cosTheta := stdmath.Min(rtmath.Dot(unitDirection.Neg(), hit.Normal), 1.0)
	sinTheta := stdmath.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1.0
	direction := rtmath.Refract(unitDirection, hit.Normal, refractionRatio)
	if cannotRefract || reflectance(cosTheta, refractionRatio) > rtmath.RandomFloat64() {
		direction = rtmath.Reflect(unitDirection, hit.Normal)
	}

	return attenuation, rtmath.NewRay(hit.Point, direction), true
}

func reflectance(cosine, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0

	return r0 + (1-r0)*stdmath.Pow(1-cosine, 5)
}
