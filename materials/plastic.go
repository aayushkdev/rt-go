package materials

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Plastic struct {
	BaseMaterial
	Albedo           Texture
	SpecularStrength float64
	Roughness        float64
}

func NewPlastic(albedo rtmath.Color, specularStrength, roughness float64) Plastic {
	return NewTexturedPlastic(NewSolidColor(albedo), specularStrength, roughness)
}

func NewTexturedPlastic(albedo Texture, specularStrength, roughness float64) Plastic {
	return Plastic{
		Albedo:           albedo,
		SpecularStrength: clamp01(specularStrength),
		Roughness:        clamp01(roughness),
	}
}

func (p Plastic) Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (ScatterRecord, bool) {
	if random.Float64() < p.SpecularStrength {
		reflected := rtmath.Reflect(rayIn.Direction.Unit(), hit.Normal)
		if p.Roughness > 0 {
			reflected = roughenDirection(reflected, p.Roughness, random)
		}
		scattered := rtmath.NewRay(hit.Point, reflected)

		return NewSpecularScatter(rtmath.NewVec3(1, 1, 1), scattered), rtmath.Dot(scattered.Direction, hit.Normal) > 0
	}

	attenuation := p.Albedo.Value(hit.U, hit.V, hit.Point)
	return NewPDFScatter(attenuation, rtmath.NewCosinePDF(hit.Normal)), true
}

func (p Plastic) ScatteringPDF(rayIn rtmath.Ray, hit HitInfo, scattered rtmath.Ray) float64 {
	cosineTheta := rtmath.Dot(hit.Normal, rtmath.UnitVector(scattered.Direction))
	if cosineTheta <= 0 {
		return 0
	}

	return cosineTheta / rtmath.Pi
}

func clamp01(value float64) float64 {
	return stdmath.Min(stdmath.Max(value, 0), 1)
}
