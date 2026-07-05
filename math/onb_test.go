package math

import (
	stdmath "math"
	"testing"
)

func TestONBFromWAlignsLocalZWithNormal(t *testing.T) {
	normal := NewVec3(0, 1, 0)
	onb := NewONBFromW(normal)
	world := onb.Local(NewVec3(0, 0, 1))

	if !near(world.X, normal.X) || !near(world.Y, normal.Y) || !near(world.Z, normal.Z) {
		t.Fatalf("local z = %#v, want %#v", world, normal)
	}
}

func TestCosineDirectionStaysInPositiveZHemisphere(t *testing.T) {
	random := NewRandom(1)
	for i := 0; i < 1000; i++ {
		direction := random.CosineDirection()
		if direction.Z < 0 {
			t.Fatalf("direction.Z = %v, want non-negative", direction.Z)
		}
		if !near(direction.Length(), 1) {
			t.Fatalf("direction length = %v, want 1", direction.Length())
		}
	}
}

func near(a, b float64) bool {
	return stdmath.Abs(a-b) < 0.000001
}
