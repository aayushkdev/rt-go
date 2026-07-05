package render

import (
	"image/color"
	stdmath "math"
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestRGBAColorReplacesInvalidComponentsWithBlack(t *testing.T) {
	got := rgbaColor(rtmath.NewVec3(stdmath.NaN(), stdmath.Inf(1), stdmath.Inf(-1)), 1)
	want := color.RGBA{A: 255}

	if got != want {
		t.Fatalf("invalid components = %v, want %v", got, want)
	}
}
