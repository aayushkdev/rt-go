package geometry

import (
	stdmath "math"
	"testing"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestRotateYMovesRayIntoObjectSpace(t *testing.T) {
	material := materials.NewLambertian(rtmath.NewVec3(1, 1, 1))
	sphere := NewSphere(rtmath.NewVec3(0, 0, -1), 0.5, material)
	rotated := RotateY(sphere, 90)

	ray := rtmath.NewRay(rtmath.NewVec3(-2, 0, 0), rtmath.NewVec3(1, 0, 0))
	record, hit := rotated.Hit(ray, rtmath.NewInterval(0.001, 100))
	if !hit {
		t.Fatal("rotated sphere was not hit")
	}
	if !near(record.Point.X, -1.5) || !near(record.Point.Y, 0) || !near(record.Point.Z, 0) {
		t.Fatalf("hit point = %#v, want approximately (-1.5, 0, 0)", record.Point)
	}
}

func TestRotateXMovesRayIntoObjectSpace(t *testing.T) {
	material := materials.NewLambertian(rtmath.NewVec3(1, 1, 1))
	sphere := NewSphere(rtmath.NewVec3(0, 0, -1), 0.5, material)
	rotated := RotateX(sphere, 90)

	ray := rtmath.NewRay(rtmath.NewVec3(0, 2, 0), rtmath.NewVec3(0, -1, 0))
	record, hit := rotated.Hit(ray, rtmath.NewInterval(0.001, 100))
	if !hit {
		t.Fatal("rotated sphere was not hit")
	}
	if !near(record.Point.X, 0) || !near(record.Point.Y, 1.5) || !near(record.Point.Z, 0) {
		t.Fatalf("hit point = %#v, want approximately (0, 1.5, 0)", record.Point)
	}
}

func TestRotateZMovesRayIntoObjectSpace(t *testing.T) {
	material := materials.NewLambertian(rtmath.NewVec3(1, 1, 1))
	sphere := NewSphere(rtmath.NewVec3(1, 0, 0), 0.5, material)
	rotated := RotateZ(sphere, 90)

	ray := rtmath.NewRay(rtmath.NewVec3(0, 2, 0), rtmath.NewVec3(0, -1, 0))
	record, hit := rotated.Hit(ray, rtmath.NewInterval(0.001, 100))
	if !hit {
		t.Fatal("rotated sphere was not hit")
	}
	if !near(record.Point.X, 0) || !near(record.Point.Y, 1.5) || !near(record.Point.Z, 0) {
		t.Fatalf("hit point = %#v, want approximately (0, 1.5, 0)", record.Point)
	}
}

func TestRotateYBoundingBox(t *testing.T) {
	material := materials.NewLambertian(rtmath.NewVec3(1, 1, 1))
	box := NewBox(rtmath.NewVec3(0, 0, -2), rtmath.NewVec3(1, 1, 0), material)
	rotated := RotateY(box, 90)
	bounds := rotated.BoundingBox()

	if !near(bounds.X.Min, -2) || !near(bounds.X.Max, 0) {
		t.Fatalf("X bounds = [%v, %v], want approximately [-2, 0]", bounds.X.Min, bounds.X.Max)
	}
	if !near(bounds.Z.Min, -1) || !near(bounds.Z.Max, 0) {
		t.Fatalf("Z bounds = [%v, %v], want approximately [-1, 0]", bounds.Z.Min, bounds.Z.Max)
	}
}

func near(a, b float64) bool {
	return stdmath.Abs(a-b) < 0.001
}
