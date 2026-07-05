package materials

import (
	stdimage "image"
	"image/color"
	"image/png"
	stdmath "math"
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
	if !nearColor(left, rtmath.NewVec3(0.75, 0, 0.25)) {
		t.Fatalf("left color = %#v", left)
	}

	right := texture.Value(0.75, 0.5, rtmath.Point3{})
	if !nearColor(right, rtmath.NewVec3(0.25, 0, 0.75)) {
		t.Fatalf("right color = %#v", right)
	}
}

func TestImageTextureConvertsSRGBToLinear(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gray.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 128, G: 128, B: 128, A: 255})
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

	color := texture.Value(0, 0, rtmath.Point3{})
	want := rtmath.NewVec3(0.2158605, 0.2158605, 0.2158605)
	if !nearColor(color, want) {
		t.Fatalf("linear color = %#v, want %#v", color, want)
	}
}

func TestTextureTransformScalesOffsetsAndRotatesUV(t *testing.T) {
	source := recordingTexture{}
	transform := NewTextureTransform(source, 2, 3, 0.1, 0.2, 90)

	color := transform.Value(0.25, 0.5, rtmath.Point3{})
	want := rtmath.NewVec3(-0.4, 0.7, 0)
	if !nearColor(color, want) {
		t.Fatalf("transformed uv = %#v, want %#v", color, want)
	}
}

func TestNoiseTextureReturnsBoundedColor(t *testing.T) {
	texture := NewNoiseTexture(4, rtmath.NewVec3(0.8, 0.6, 0.4))
	color := texture.Value(0, 0, rtmath.NewVec3(0.25, 0.5, 0.75))

	if color.X < 0 || color.X > 0.8 ||
		color.Y < 0 || color.Y > 0.6 ||
		color.Z < 0 || color.Z > 0.4 {
		t.Fatalf("noise color out of bounds: %#v", color)
	}
}

type recordingTexture struct{}

func (r recordingTexture) Value(u, v float64, point rtmath.Point3) rtmath.Color {
	return rtmath.NewVec3(u, v, 0)
}

func nearColor(a, b rtmath.Color) bool {
	const epsilon = 1e-5
	return stdmath.Abs(a.X-b.X) < epsilon &&
		stdmath.Abs(a.Y-b.Y) < epsilon &&
		stdmath.Abs(a.Z-b.Z) < epsilon
}
