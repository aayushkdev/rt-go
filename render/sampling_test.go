package render

import (
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestStratifiedSampleOffsetStaysInsidePixel(t *testing.T) {
	random := rtmath.NewRandom(1)

	for sample := 0; sample < 100; sample++ {
		u, v := stratifiedSampleOffset(sample, 100, random)
		if u < -0.5 || u >= 0.5 {
			t.Fatalf("u offset = %v, want [-0.5, 0.5)", u)
		}
		if v < -0.5 || v >= 0.5 {
			t.Fatalf("v offset = %v, want [-0.5, 0.5)", v)
		}
	}
}

func TestStratifiedSampleOffsetSpreadsSamplesAcrossGrid(t *testing.T) {
	random := rtmath.NewRandom(1)

	firstU, firstV := stratifiedSampleOffset(0, 4, random)
	lastU, lastV := stratifiedSampleOffset(3, 4, random)

	if firstU >= 0 || firstV >= 0 {
		t.Fatalf("first sample = (%v, %v), want lower-left pixel cell", firstU, firstV)
	}
	if lastU < 0 || lastV < 0 {
		t.Fatalf("last sample = (%v, %v), want upper-right pixel cell", lastU, lastV)
	}
}

func TestRussianRouletteSurvivalUsesStrongestColorChannel(t *testing.T) {
	survival := russianRouletteSurvival(rtmath.NewVec3(0.2, 0.7, 0.4))
	if survival != 0.7 {
		t.Fatalf("survival = %v, want 0.7", survival)
	}
}

func TestRussianRouletteSurvivalIsClamped(t *testing.T) {
	low := russianRouletteSurvival(rtmath.NewVec3(0.001, 0.002, 0.003))
	if low != 0.05 {
		t.Fatalf("low survival = %v, want 0.05", low)
	}

	high := russianRouletteSurvival(rtmath.NewVec3(2, 1, 0.5))
	if high != 0.95 {
		t.Fatalf("high survival = %v, want 0.95", high)
	}
}

func TestRendererUsesRussianRouletteAfterFiveBounces(t *testing.T) {
	renderer := Renderer{MaxDepth: 10}

	if renderer.useRussianRoulette(6) {
		t.Fatal("russian roulette should not run before five completed bounces")
	}
	if !renderer.useRussianRoulette(5) {
		t.Fatal("russian roulette should run after five completed bounces")
	}
}
