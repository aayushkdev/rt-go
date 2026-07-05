package materials

import (
	"fmt"
	stdimage "image"
	_ "image/jpeg"
	_ "image/png"
	stdmath "math"
	"os"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Texture interface {
	Value(u, v float64, point rtmath.Point3) rtmath.Color
}

type SolidColor struct {
	Color rtmath.Color
}

func NewSolidColor(color rtmath.Color) SolidColor {
	return SolidColor{Color: color}
}

func (s SolidColor) Value(u, v float64, point rtmath.Point3) rtmath.Color {
	return s.Color
}

type CheckerTexture struct {
	Scale float64
	Even  Texture
	Odd   Texture
}

func NewCheckerTexture(scale float64, even, odd Texture) CheckerTexture {
	if scale <= 0 {
		scale = 1
	}
	return CheckerTexture{Scale: scale, Even: even, Odd: odd}
}

func (c CheckerTexture) Value(u, v float64, point rtmath.Point3) rtmath.Color {
	x := int(stdmath.Floor(point.X * c.Scale))
	y := int(stdmath.Floor(point.Y * c.Scale))
	z := int(stdmath.Floor(point.Z * c.Scale))
	if (x+y+z)%2 == 0 {
		return c.Even.Value(u, v, point)
	}

	return c.Odd.Value(u, v, point)
}

type ImageTexture struct {
	Width  int
	Height int
	Pixels []rtmath.Color
}

func NewImageTexture(path string) (ImageTexture, error) {
	file, err := os.Open(path)
	if err != nil {
		return ImageTexture{}, fmt.Errorf("open image texture: %w", err)
	}
	defer file.Close()

	decoded, _, err := stdimage.Decode(file)
	if err != nil {
		return ImageTexture{}, fmt.Errorf("decode image texture: %w", err)
	}

	bounds := decoded.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	pixels := make([]rtmath.Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := decoded.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			pixels[y*width+x] = rtmath.NewVec3(
				sRGBToLinear(float64(r)/65535),
				sRGBToLinear(float64(g)/65535),
				sRGBToLinear(float64(b)/65535),
			)
		}
	}

	return ImageTexture{Width: width, Height: height, Pixels: pixels}, nil
}

func (i ImageTexture) Value(u, v float64, point rtmath.Point3) rtmath.Color {
	if i.Width <= 0 || i.Height <= 0 || len(i.Pixels) == 0 {
		return rtmath.NewVec3(1, 0, 1)
	}

	u = rtmath.NewInterval(0, 1).Clamp(u)
	v = 1 - rtmath.NewInterval(0, 1).Clamp(v)

	x := u * float64(i.Width-1)
	y := v * float64(i.Height-1)
	x0 := int(stdmath.Floor(x))
	y0 := int(stdmath.Floor(y))
	x1 := min(x0+1, i.Width-1)
	y1 := min(y0+1, i.Height-1)
	tx := x - float64(x0)
	ty := y - float64(y0)

	top := lerpColor(i.pixel(x0, y0), i.pixel(x1, y0), tx)
	bottom := lerpColor(i.pixel(x0, y1), i.pixel(x1, y1), tx)
	return lerpColor(top, bottom, ty)
}

func (i ImageTexture) pixel(x, y int) rtmath.Color {
	return i.Pixels[y*i.Width+x]
}

func lerpColor(a, b rtmath.Color, t float64) rtmath.Color {
	return a.Mul(1 - t).Add(b.Mul(t))
}

func sRGBToLinear(value float64) float64 {
	if value <= 0.04045 {
		return value / 12.92
	}

	return stdmath.Pow((value+0.055)/1.055, 2.4)
}
