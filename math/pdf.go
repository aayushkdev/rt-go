package math

type PDF interface {
	Value(direction Vec3) float64
	Generate(random *Random) Vec3
}

type SpherePDF struct{}

func NewSpherePDF() SpherePDF {
	return SpherePDF{}
}

func (p SpherePDF) Value(direction Vec3) float64 {
	return 1 / (4 * Pi)
}

func (p SpherePDF) Generate(random *Random) Vec3 {
	return random.UnitVector()
}

type CosinePDF struct {
	UVW ONB
}

func NewCosinePDF(w Vec3) CosinePDF {
	return CosinePDF{UVW: NewONBFromW(w)}
}

func (p CosinePDF) Value(direction Vec3) float64 {
	cosineTheta := Dot(UnitVector(direction), p.UVW.W)
	if cosineTheta <= 0 {
		return 0
	}

	return cosineTheta / Pi
}

func (p CosinePDF) Generate(random *Random) Vec3 {
	return p.UVW.Local(random.CosineDirection())
}

type MixturePDF struct {
	A       PDF
	B       PDF
	WeightA float64
}

func NewMixturePDF(a, b PDF) MixturePDF {
	return NewWeightedMixturePDF(a, b, 0.5)
}

func NewWeightedMixturePDF(a, b PDF, weightA float64) MixturePDF {
	if weightA < 0 {
		weightA = 0
	}
	if weightA > 1 {
		weightA = 1
	}

	return MixturePDF{A: a, B: b, WeightA: weightA}
}

func (p MixturePDF) Value(direction Vec3) float64 {
	return p.WeightA*p.A.Value(direction) + (1-p.WeightA)*p.B.Value(direction)
}

func (p MixturePDF) Generate(random *Random) Vec3 {
	if random.Float64() < p.WeightA {
		return p.A.Generate(random)
	}

	return p.B.Generate(random)
}
