package geometry

import (
	"testing"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestMeshPDFValueWeightsTrianglesByArea(t *testing.T) {
	material := materials.NewLambertian(rtmath.NewVec3(0.5, 0.5, 0.5))
	small := NewTriangle(
		rtmath.NewVec3(0, 0, -1),
		rtmath.NewVec3(1, 0, -1),
		rtmath.NewVec3(0, 1, -1),
		material,
	)
	large := NewTriangle(
		rtmath.NewVec3(10, 0, -1),
		rtmath.NewVec3(20, 0, -1),
		rtmath.NewVec3(10, 10, -1),
		material,
	)
	mesh := Mesh{Triangles: []Triangle{small, large}}

	origin := rtmath.NewVec3(0.25, 0.25, 0)
	direction := rtmath.NewVec3(0, 0, -1)
	got := mesh.PDFValue(origin, direction)
	totalArea := small.Area() + large.Area()
	want := (small.Area() / totalArea) * small.PDFValue(origin, direction)

	if !near(got, want) {
		t.Fatalf("mesh pdf = %v, want %v", got, want)
	}
}
