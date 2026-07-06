package materials

import (
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestPlasticWithNoSpecularUsesDiffusePDF(t *testing.T) {
	plastic := NewPlastic(rtmath.NewVec3(0.8, 0.1, 0.05), 0, 0.2)
	record, ok := plastic.Scatter(
		rtmath.NewRay(rtmath.Point3{}, rtmath.NewVec3(0, 0, -1)),
		HitInfo{Point: rtmath.Point3{}, Normal: rtmath.NewVec3(0, 0, 1), FrontFace: true},
		rtmath.NewRandom(1),
	)

	if !ok {
		t.Fatal("scatter failed")
	}
	if record.SkipPDF {
		t.Fatal("plastic with no specular should use diffuse PDF scatter")
	}
	if record.Attenuation != rtmath.NewVec3(0.8, 0.1, 0.05) {
		t.Fatalf("attenuation = %#v", record.Attenuation)
	}
}

func TestPlasticWithFullSpecularUsesReflection(t *testing.T) {
	plastic := NewPlastic(rtmath.NewVec3(0.8, 0.1, 0.05), 1, 0)
	record, ok := plastic.Scatter(
		rtmath.NewRay(rtmath.Point3{}, rtmath.NewVec3(0, 0, -1)),
		HitInfo{Point: rtmath.Point3{}, Normal: rtmath.NewVec3(0, 0, 1), FrontFace: true},
		rtmath.NewRandom(1),
	)

	if !ok {
		t.Fatal("scatter failed")
	}
	if !record.SkipPDF {
		t.Fatal("plastic with full specular should use specular scatter")
	}
	if record.Attenuation != rtmath.NewVec3(1, 1, 1) {
		t.Fatalf("specular attenuation = %#v", record.Attenuation)
	}
	if record.Scattered.Direction != rtmath.NewVec3(0, 0, 1) {
		t.Fatalf("scattered direction = %#v", record.Scattered.Direction)
	}
}
