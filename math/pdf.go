package math

type PDF interface {
	Value(direction Vec3) float64
	Generate(random *Random) Vec3
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
