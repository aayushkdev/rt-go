package math

import stdmath "math"

type Vec3 struct {
	X, Y, Z float64
}

type Point3 = Vec3
type Color = Vec3

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

func (v Vec3) Mul(t float64) Vec3 {
	return Vec3{X: v.X * t, Y: v.Y * t, Z: v.Z * t}
}

func (v Vec3) MulVec(other Vec3) Vec3 {
	return Vec3{X: v.X * other.X, Y: v.Y * other.Y, Z: v.Z * other.Z}
}

func (v Vec3) Div(t float64) Vec3 {
	return v.Mul(1 / t)
}

func (v Vec3) Neg() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func (v Vec3) Length() float64 {
	return stdmath.Sqrt(v.LengthSquared())
}

func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Unit() Vec3 {
	return v.Div(v.Length())
}

func (v Vec3) NearZero() bool {
	const epsilon = 1e-8

	return stdmath.Abs(v.X) < epsilon &&
		stdmath.Abs(v.Y) < epsilon &&
		stdmath.Abs(v.Z) < epsilon
}

func Dot(a, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Cross(a, b Vec3) Vec3 {
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func UnitVector(v Vec3) Vec3 {
	return v.Unit()
}

func Reflect(v, normal Vec3) Vec3 {
	return v.Sub(normal.Mul(2 * Dot(v, normal)))
}
