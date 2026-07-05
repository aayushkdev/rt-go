package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type Texture interface {
	Value(u, v float64, point rtmath.Point3) rtmath.Color
}

type SolidColor struct {
	Color rtmath.Color
}

func NewSolidColor(color rtmath.Color) SolidColor {
	return SolidColor{Color: color}
}

func (s SolidColor) Value(u, v float64, point rtmath.Point3) rtmath.Color {
	return s.Color
}
