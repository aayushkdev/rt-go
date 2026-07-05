package materials

import (
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestDiffuseLightEmitsOnlyFromFrontFace(t *testing.T) {
	light := NewDiffuseLight(rtmath.NewVec3(4, 5, 6))

	front := light.Emitted(HitInfo{FrontFace: true})
	if front != (rtmath.NewVec3(4, 5, 6)) {
		t.Fatalf("front emission = %v, want %v", front, rtmath.NewVec3(4, 5, 6))
	}

	back := light.Emitted(HitInfo{FrontFace: false})
	if back != (rtmath.NewVec3(0, 0, 0)) {
		t.Fatalf("back emission = %v, want black", back)
	}
}
