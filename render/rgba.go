package render

import (
	stdimage "image"
	"image/color"
	stdmath "math"
	"sync"
	"time"

	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/geometry"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func (r Renderer) RenderRGBA(cam camera.Camera, world geometry.Hittable) *stdimage.RGBA {
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, cam.ImageWidth, cam.ImageHeight))
	workerCount := r.workerCount()
	jobs := make(chan int)
	results := make(chan rgbaScanline, workerCount)
	var wg sync.WaitGroup

	for worker := 0; worker < workerCount; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			random := rtmath.NewRandom(time.Now().UnixNano() + int64(worker))
			for j := range jobs {
				row := make([]color.RGBA, cam.ImageWidth)
				for i := 0; i < cam.ImageWidth; i++ {
					pixelColor := rtmath.NewVec3(0, 0, 0)
					for sample := 0; sample < r.SamplesPerPixel; sample++ {
						offsetU, offsetV := stratifiedSampleOffset(sample, r.SamplesPerPixel, random)
						ray := cam.RayForPixelSampleRandom(i, j, offsetU, offsetV, random)
						pixelColor = pixelColor.Add(r.rayColor(ray, world, r.MaxDepth, random))
					}
					row[i] = rgbaColor(pixelColor, r.SamplesPerPixel)
				}
				results <- rgbaScanline{Index: j, Row: row}
			}
		}(worker)
	}

	go func() {
		defer close(jobs)
		for j := 0; j < cam.ImageHeight; j++ {
			jobs <- j
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		for i, pixel := range result.Row {
			img.SetRGBA(i, result.Index, pixel)
		}
	}

	return img
}

type rgbaScanline struct {
	Index int
	Row   []color.RGBA
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
