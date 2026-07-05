package geometry

import (
	"testing"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestTranslateMovesHitPointAndBoundingBox(t *testing.T) {
	material := materials.NewLambertian(rtmath.NewVec3(1, 1, 1))
	sphere := NewSphere(rtmath.NewVec3(0, 0, -1), 0.5, material)
	translated := NewTranslate(sphere, rtmath.NewVec3(1, 0, 0))

	ray := rtmath.NewRay(rtmath.NewVec3(1, 0, 0), rtmath.NewVec3(0, 0, -1))
	record, hit := translated.Hit(ray, rtmath.NewInterval(0.001, 100))
	if !hit {
		t.Fatal("translated sphere was not hit")
	}
	if record.Point != rtmath.NewVec3(1, 0, -0.5) {
		t.Fatalf("hit point = %#v, want %#v", record.Point, rtmath.NewVec3(1, 0, -0.5))
	}

	box := translated.BoundingBox()
	if box.X.Min != 0.5 || box.X.Max != 1.5 {
		t.Fatalf("translated X bounds = [%v, %v], want [0.5, 1.5]", box.X.Min, box.X.Max)
	}
}
