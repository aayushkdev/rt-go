package image

import (
	stdmath "math"
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestColorBytesReplacesInvalidComponentsWithBlack(t *testing.T) {
	r, g, b := ColorBytes(rtmath.NewVec3(stdmath.NaN(), stdmath.Inf(1), stdmath.Inf(-1)), 1)

	if r != 0 || g != 0 || b != 0 {
		t.Fatalf("invalid components = (%d, %d, %d), want black", r, g, b)
	}
}

func TestColorBytesClampsHugeFiniteComponents(t *testing.T) {
	r, g, b := ColorBytes(rtmath.NewVec3(1e100, 1e100, 1e100), 1)

	if r != 255 || g != 255 || b != 255 {
		t.Fatalf("huge finite components = (%d, %d, %d), want white clamp", r, g, b)
	}
}
