package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type Lambertian struct {
	BaseMaterial
	Albedo rtmath.Color
}

func NewLambertian(albedo rtmath.Color) Lambertian {
	return Lambertian{Albedo: albedo}
}

func (l Lambertian) Scatter(rayIn rtmath.Ray, hit HitInfo, random *rtmath.Random) (rtmath.Color, rtmath.Ray, bool) {
	scatterDirection := hit.Normal.Add(random.UnitVector())
	if scatterDirection.NearZero() {
		scatterDirection = hit.Normal
	}

	return l.Albedo, rtmath.NewRay(hit.Point, scatterDirection), true
}
