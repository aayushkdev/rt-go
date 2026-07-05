package scene

import (
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Config struct {
	Objects []Object
}

type Object struct {
	Kind     string
	Material materials.Material

	Path     string
	Height   float64
	Position rtmath.Point3

	Center rtmath.Point3
	Radius float64

	A rtmath.Point3
	B rtmath.Point3
	C rtmath.Point3
	U rtmath.Vec3
	V rtmath.Vec3
}

func New(objects ...Object) Config {
	return Config{Objects: objects}
}

func Model(path string) Object {
	return Object{Kind: "model", Path: path, Height: 1}
}

func Sphere(x, y, z, radius float64, material materials.Material) Object {
	return Object{
		Kind:     "sphere",
		Center:   rtmath.NewVec3(x, y, z),
		Radius:   radius,
		Material: material,
	}
}

func Triangle(a, b, c rtmath.Point3, material materials.Material) Object {
	return Object{Kind: "triangle", A: a, B: b, C: c, Material: material}
}

func Quad(q, u, v rtmath.Vec3, material materials.Material) Object {
	return Object{Kind: "quad", A: q, U: u, V: v, Material: material}
}

func (o Object) At(x, y, z float64) Object {
	o.Position = rtmath.NewVec3(x, y, z)
	return o
}

func (o Object) WithHeight(height float64) Object {
	o.Height = height
	return o
}
