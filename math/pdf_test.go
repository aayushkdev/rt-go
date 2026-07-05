package math

import "testing"

func TestSpherePDFValue(t *testing.T) {
	pdf := NewSpherePDF()

	if got := pdf.Value(NewVec3(1, 0, 0)); !near(got, 1/(4*Pi)) {
		t.Fatalf("sphere pdf = %v, want %v", got, 1/(4*Pi))
	}
}

func TestSpherePDFGenerateReturnsUnitVector(t *testing.T) {
	random := NewRandom(1)
	pdf := NewSpherePDF()

	for i := 0; i < 1000; i++ {
		direction := pdf.Generate(random)
		if !near(direction.Length(), 1) {
			t.Fatalf("direction length = %v, want 1", direction.Length())
		}
	}
}

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

type constantPDF struct {
	value     float64
	direction Vec3
}

func (p constantPDF) Value(direction Vec3) float64 {
	return p.value
}

func (p constantPDF) Generate(random *Random) Vec3 {
	return p.direction
}

func TestMixturePDFValueAveragesInputs(t *testing.T) {
	pdf := NewMixturePDF(
		constantPDF{value: 2},
		constantPDF{value: 6},
	)

	if got := pdf.Value(NewVec3(1, 0, 0)); got != 4 {
		t.Fatalf("mixture pdf = %v, want 4", got)
	}
}
