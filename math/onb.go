package math

import stdmath "math"

type ONB struct {
	U Vec3
	V Vec3
	W Vec3
}

func NewONBFromW(n Vec3) ONB {
	w := UnitVector(n)
	a := NewVec3(1, 0, 0)
	if stdmath.Abs(w.X) > 0.9 {
		a = NewVec3(0, 1, 0)
	}

	v := UnitVector(Cross(w, a))
	u := Cross(w, v)

	return ONB{U: u, V: v, W: w}
}

func (o ONB) Local(a Vec3) Vec3 {
	return o.U.Mul(a.X).Add(o.V.Mul(a.Y)).Add(o.W.Mul(a.Z))
}
