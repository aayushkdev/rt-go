package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type Lambertian struct {
	BaseMaterial
	Albedo Texture
}

func NewLambertian(albedo rtmath.Color) Lambertian {
	return NewTexturedLambertian(NewSolidColor(albedo))
}

func NewTexturedLambertian(albedo Texture) Lambertian {
	return Lambertian{Albedo: albedo}
}

func (l Lambertian) Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (ScatterRecord, bool) {
	attenuation := l.Albedo.Value(0, 0, hit.Point)
	return NewPDFScatter(attenuation, rtmath.NewCosinePDF(hit.Normal)), true
}

func (l Lambertian) ScatteringPDF(rayIn rtmath.Ray, hit HitInfo, scattered rtmath.Ray) float64 {
	cosineTheta := rtmath.Dot(hit.Normal, rtmath.UnitVector(scattered.Direction))
	if cosineTheta <= 0 {
		return 0
	}

	return cosineTheta / rtmath.Pi
}
