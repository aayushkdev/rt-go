package render

import (
	"fmt"
	stdmath "math"
	"os"

	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/geometry"
	rtimage "github.com/aayushkdev/rt-go/image"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Renderer struct {
	SamplesPerPixel    int
	MaxDepth           int
	FlushEveryScanline int
	Background         rtmath.Color
	SkyBackground      bool
}

func NewRenderer() Renderer {
	return Renderer{
		SamplesPerPixel:    10,
		MaxDepth:           50,
		FlushEveryScanline: 10,
		SkyBackground:      true,
	}
}

func (r Renderer) Render(cam camera.Camera, world geometry.Hittable, outputPath string) error {
	headerSize, err := createPPM(outputPath, cam)
	if err != nil {
		return err
	}

	rowStride := cam.ImageWidth * 3
	pendingRows := make(map[int][]byte)

	for j := 0; j < cam.ImageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", cam.ImageHeight-j)
		row := make([]byte, rowStride)
		for i := 0; i < cam.ImageWidth; i++ {
			pixelColor := rtmath.NewVec3(0, 0, 0)
			for sample := 0; sample < r.SamplesPerPixel; sample++ {
				ray := cam.RayForPixelSample(i, j, sampleOffset(), sampleOffset())
				pixelColor = pixelColor.Add(r.rayColor(ray, world, r.MaxDepth))
			}
			rByte, gByte, bByte := rtimage.ColorBytes(pixelColor, r.SamplesPerPixel)
			offset := i * 3
			row[offset] = rByte
			row[offset+1] = gByte
			row[offset+2] = bByte
		}

		pendingRows[j] = row
		if r.FlushEveryScanline <= 0 || (j+1)%r.FlushEveryScanline == 0 {
			if err := writeRows(outputPath, headerSize, rowStride, pendingRows); err != nil {
				return err
			}
			pendingRows = make(map[int][]byte)
		}
	}

	if len(pendingRows) > 0 {
		if err := writeRows(outputPath, headerSize, rowStride, pendingRows); err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stderr, "\rDone.                 ")

	return nil
}

func createPPM(outputPath string, cam camera.Camera) (int, error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", outputPath, err)
	}

	header := fmt.Sprintf("P6\n%d %d\n255\n", cam.ImageWidth, cam.ImageHeight)
	if _, err := file.WriteString(header); err != nil {
		file.Close()
		return 0, fmt.Errorf("write header %s: %w", outputPath, err)
	}

	imageSize := int64(cam.ImageWidth * cam.ImageHeight * 3)
	if err := file.Truncate(int64(len(header)) + imageSize); err != nil {
		file.Close()
		return 0, fmt.Errorf("size %s: %w", outputPath, err)
	}

	if err := file.Close(); err != nil {
		return 0, fmt.Errorf("close %s: %w", outputPath, err)
	}

	return len(header), nil
}

func writeRows(outputPath string, headerSize, rowStride int, rows map[int][]byte) error {
	file, err := os.OpenFile(outputPath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", outputPath, err)
	}

	for rowIndex, row := range rows {
		rowOffset := int64(headerSize + rowIndex*rowStride)
		if _, err := file.WriteAt(row, rowOffset); err != nil {
			file.Close()
			return fmt.Errorf("write scanline %d: %w", rowIndex, err)
		}
	}

	if err := file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("sync %s: %w", outputPath, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", outputPath, err)
	}

	return nil
}

func (r Renderer) rayColor(ray rtmath.Ray, world geometry.Hittable, depth int) rtmath.Color {
	if depth <= 0 {
		return rtmath.NewVec3(0, 0, 0)
	}

	record, hit := world.Hit(ray, rtmath.NewInterval(0.001, stdmath.Inf(1)))
	if hit {
		emitted := record.Material.Emitted(record.HitInfo)
		attenuation, scattered, ok := record.Material.Scatter(ray, record.HitInfo)
		if !ok {
			return emitted
		}

		return emitted.Add(attenuation.MulVec(r.rayColor(scattered, world, depth-1)))
	}

	if !r.SkyBackground {
		return r.Background
	}

	unitDirection := rtmath.UnitVector(ray.Direction)
	a := 0.5 * (unitDirection.Y + 1.0)
	white := rtmath.NewVec3(1.0, 1.0, 1.0)
	blue := rtmath.NewVec3(0.5, 0.7, 1.0)

	return white.Mul(1.0 - a).Add(blue.Mul(a))
}

func sampleOffset() float64 {
	return rtmath.RandomFloat64() - 0.5
}
