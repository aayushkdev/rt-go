package math

import "testing"

func TestCosinePDFValue(t *testing.T) {
	pdf := NewCosinePDF(NewVec3(0, 1, 0))

	if got := pdf.Value(NewVec3(0, 2, 0)); !near(got, 1/Pi) {
		t.Fatalf("pdf straight up = %v, want %v", got, 1/Pi)
	}
	if got := pdf.Value(NewVec3(1, 0, 0)); got != 0 {
		t.Fatalf("pdf tangent = %v, want 0", got)
	}
	if got := pdf.Value(NewVec3(0, -1, 0)); got != 0 {
		t.Fatalf("pdf below surface = %v, want 0", got)
	}
}

func TestCosinePDFGenerateUsesPositiveHemisphere(t *testing.T) {
	random := NewRandom(1)
	pdf := NewCosinePDF(NewVec3(0, 1, 0))

	for i := 0; i < 1000; i++ {
		direction := pdf.Generate(random)
		if Dot(UnitVector(direction), NewVec3(0, 1, 0)) < 0 {
			t.Fatalf("generated direction %#v is below the surface", direction)
		}
	}
}
