package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type Lambertian struct {
	BaseMaterial
	Albedo rtmath.Color
}

func NewLambertian(albedo rtmath.Color) Lambertian {
	return Lambertian{Albedo: albedo}
}

func (l Lambertian) Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (ScatterRecord, bool) {
	return NewPDFScatter(l.Albedo, rtmath.NewCosinePDF(hit.Normal)), true
}

func (l Lambertian) ScatteringPDF(rayIn rtmath.Ray, hit HitInfo, scattered rtmath.Ray) float64 {
	cosineTheta := rtmath.Dot(hit.Normal, rtmath.UnitVector(scattered.Direction))
	if cosineTheta <= 0 {
		return 0
	}

	return cosineTheta / rtmath.Pi
}
