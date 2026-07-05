package materials

import (
	stdimage "image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestCheckerTextureAlternatesByPosition(t *testing.T) {
	checker := NewCheckerTexture(
		1,
		NewSolidColor(rtmath.NewVec3(1, 1, 1)),
		NewSolidColor(rtmath.NewVec3(0, 0, 0)),
	)

	even := checker.Value(0, 0, rtmath.NewVec3(0.2, 0.2, 0.2))
	if even != rtmath.NewVec3(1, 1, 1) {
		t.Fatalf("even checker color = %#v", even)
	}

	odd := checker.Value(0, 0, rtmath.NewVec3(1.2, 0.2, 0.2))
	if odd != rtmath.NewVec3(0, 0, 0) {
		t.Fatalf("odd checker color = %#v", odd)
	}
}

func TestImageTextureLoadsImageFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "texture.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{B: 255, A: 255})
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	texture, err := NewImageTexture(path)
	if err != nil {
		t.Fatal(err)
	}

	left := texture.Value(0.25, 0.5, rtmath.Point3{})
	if left != rtmath.NewVec3(1, 0, 0) {
		t.Fatalf("left color = %#v", left)
	}

	right := texture.Value(0.75, 0.5, rtmath.Point3{})
	if right != rtmath.NewVec3(0, 0, 1) {
		t.Fatalf("right color = %#v", right)
	}
}
