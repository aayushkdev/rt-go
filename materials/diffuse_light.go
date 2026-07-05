package materials

import rtmath "github.com/aayushkdev/rt-go/math"

type DiffuseLight struct {
	Emit rtmath.Color
}

func NewDiffuseLight(emit rtmath.Color) DiffuseLight {
	return DiffuseLight{Emit: emit}
}

func (d DiffuseLight) Emitted(hit HitInfo) rtmath.Color {
	return d.Emit
}

func (d DiffuseLight) Scatter(rayIn rtmath.Ray, hit HitInfo) (rtmath.Color, rtmath.Ray, bool) {
	return rtmath.Color{}, rtmath.Ray{}, false
}
