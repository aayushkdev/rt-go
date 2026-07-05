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
