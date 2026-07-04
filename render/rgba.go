package render

import (
	stdimage "image"
	"image/color"
	stdmath "math"

	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/geometry"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func (r Renderer) RenderRGBA(cam camera.Camera, world geometry.Hittable) *stdimage.RGBA {
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, cam.ImageWidth, cam.ImageHeight))

	for j := 0; j < cam.ImageHeight; j++ {
		for i := 0; i < cam.ImageWidth; i++ {
			pixelColor := rtmath.NewVec3(0, 0, 0)
			for sample := 0; sample < r.SamplesPerPixel; sample++ {
				ray := cam.RayForPixelSample(i, j, sampleOffset(), sampleOffset())
				pixelColor = pixelColor.Add(r.rayColor(ray, world, r.MaxDepth))
			}
			img.SetRGBA(i, j, rgbaColor(pixelColor, r.SamplesPerPixel))
		}
	}

	return img
}

func rgbaColor(pixelColor rtmath.Color, samplesPerPixel int) color.RGBA {
	scale := 1.0 / float64(samplesPerPixel)
	intensity := rtmath.NewInterval(0, 0.999)

	r := intensity.Clamp(linearToGamma(pixelColor.X * scale))
	g := intensity.Clamp(linearToGamma(pixelColor.Y * scale))
	b := intensity.Clamp(linearToGamma(pixelColor.Z * scale))

	return color.RGBA{
		R: uint8(255.999 * r),
		G: uint8(255.999 * g),
		B: uint8(255.999 * b),
		A: 255,
	}
}

func linearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return stdmath.Sqrt(linearComponent)
	}

	return 0
}
