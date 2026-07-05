package geometry

import (
	"testing"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestQuadPDFValue(t *testing.T) {
	quad := NewQuad(
		rtmath.NewVec3(-1, 0, -1),
		rtmath.NewVec3(2, 0, 0),
		rtmath.NewVec3(0, 0, 2),
		materials.NewLambertian(rtmath.NewVec3(1, 1, 1)),
	)

	got := quad.PDFValue(rtmath.NewVec3(0, 1, 0), rtmath.NewVec3(0, -1, 0))
	if !near(got, 0.25) {
		t.Fatalf("quad pdf = %v, want 0.25", got)
	}
}

func TestSpherePDFValue(t *testing.T) {
	sphere := NewSphere(
		rtmath.NewVec3(0, 0, -2),
		1,
		materials.NewLambertian(rtmath.NewVec3(1, 1, 1)),
	)

	got := sphere.PDFValue(rtmath.NewVec3(0, 0, 0), rtmath.NewVec3(0, 0, -1))
	want := 1 / (2 * rtmath.Pi * (1 - 0.8660254037844386))
	if !near(got, want) {
		t.Fatalf("sphere pdf = %v, want %v", got, want)
	}
}

func TestTrianglePDFValue(t *testing.T) {
	triangle := NewTriangle(
		rtmath.NewVec3(-1, 0, -1),
		rtmath.NewVec3(1, 0, -1),
		rtmath.NewVec3(-1, 0, 1),
		materials.NewLambertian(rtmath.NewVec3(1, 1, 1)),
	)

	got := triangle.PDFValue(rtmath.NewVec3(0, 1, 0), rtmath.NewVec3(0, -1, 0))
	if !near(got, 0.5) {
		t.Fatalf("triangle pdf = %v, want 0.5", got)
	}
}

func TestTriangleRandomReturnsDirectionToTriangle(t *testing.T) {
	triangle := NewTriangle(
		rtmath.NewVec3(0, 0, 0),
		rtmath.NewVec3(1, 0, 0),
		rtmath.NewVec3(0, 0, 1),
		materials.NewLambertian(rtmath.NewVec3(1, 1, 1)),
	)

	origin := rtmath.NewVec3(0, 1, 0)
	direction := triangle.Random(origin, rtmath.NewRandom(1))
	point := origin.Add(direction)

	if point.Y != 0 {
		t.Fatalf("sample point y = %v, want 0", point.Y)
	}
	if point.X < 0 || point.Z < 0 || point.X+point.Z > 1 {
		t.Fatalf("sample point %#v is outside triangle", point)
	}
}
