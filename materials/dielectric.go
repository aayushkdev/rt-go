package materials

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Dielectric struct {
	BaseMaterial
	RefractionIndex float64
	Tint            Texture
	Roughness       float64
}

func NewDielectric(refractionIndex float64) Dielectric {
	return NewTintedDielectric(refractionIndex, NewSolidColor(rtmath.NewVec3(1, 1, 1)))
}

func NewTintedDielectric(refractionIndex float64, tint Texture) Dielectric {
	return NewRoughDielectric(refractionIndex, tint, 0)
}

func NewRoughDielectric(refractionIndex float64, tint Texture, roughness float64) Dielectric {
	return Dielectric{
		RefractionIndex: refractionIndex,
		Tint:            tint,
		Roughness:       stdmath.Min(stdmath.Max(roughness, 0), 1),
	}
}

func (d Dielectric) Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (ScatterRecord, bool) {
	attenuation := d.Tint.Value(hit.U, hit.V, hit.Point)
	refractionRatio := d.RefractionIndex
	if hit.FrontFace {
		refractionRatio = 1.0 / d.RefractionIndex
	}

	unitDirection := rtmath.UnitVector(rayIn.Direction)
	cosTheta := stdmath.Min(rtmath.Dot(unitDirection.Neg(), hit.Normal), 1.0)
	sinTheta := stdmath.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1.0
	direction := rtmath.Refract(unitDirection, hit.Normal, refractionRatio)
	if cannotRefract || reflectance(cosTheta, refractionRatio) > random.Float64() {
		direction = rtmath.Reflect(unitDirection, hit.Normal)
	}
	if d.Roughness > 0 {
		direction = roughenDirection(direction, d.Roughness, random)
	}

	return NewSpecularScatter(attenuation, rtmath.NewRay(hit.Point, direction)), true
}

func roughenDirection(direction rtmath.Vec3, roughness float64, random *rtmath.Random) rtmath.Vec3 {
	coneCosine := stdmath.Cos(roughness * rtmath.Pi / 2)
	z := random.Float64Range(coneCosine, 1)
	phi := 2 * rtmath.Pi * random.Float64()
	r := stdmath.Sqrt(1 - z*z)
	local := rtmath.NewVec3(stdmath.Cos(phi)*r, stdmath.Sin(phi)*r, z)

	return rtmath.NewONBFromW(direction).Local(local).Unit()
}

func reflectance(cosine, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0

	return r0 + (1-r0)*stdmath.Pow(1-cosine, 5)
}
