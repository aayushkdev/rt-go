package scene

import (
	"testing"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestBuildSampleTargetsIncludesLightsAndExplicitTargets(t *testing.T) {
	matte := materials.NewLambertian(rtmath.NewVec3(0.5, 0.5, 0.5))
	light := materials.NewDiffuseLight(rtmath.NewVec3(4, 4, 4))
	config := New(
		Sphere(0, 0, -1, 0.5, matte),
		AsSampleTarget(Sphere(1, 0, -1, 0.5, matte)),
		AsLight(Quad(rtmath.NewVec3(-1, 1, -2), rtmath.NewVec3(2, 0, 0), rtmath.NewVec3(0, 0, 2), light)),
	)

	targets := BuildSampleTargets(config)
	if len(targets.Objects) != 2 {
		t.Fatalf("sample target count = %d, want 2", len(targets.Objects))
	}
}

func TestAsLightAlsoMarksObjectAsSampleTarget(t *testing.T) {
	light := AsLight(Sphere(0, 0, -1, 0.5, materials.NewDiffuseLight(rtmath.NewVec3(4, 4, 4))))

	if !light.Light || !light.Sample {
		t.Fatalf("light flags = Light:%v Sample:%v, want both true", light.Light, light.Sample)
	}
}
