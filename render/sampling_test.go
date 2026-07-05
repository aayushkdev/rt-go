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
