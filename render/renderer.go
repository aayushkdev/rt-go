package render

import (
	"fmt"
	stdmath "math"
	"os"

	"github.com/aayushkdev/rt-go/camera"
	rtimage "github.com/aayushkdev/rt-go/image"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/objects"
)

type Renderer struct{}

func NewRenderer() Renderer {
	return Renderer{}
}

func (r Renderer) Render(cam camera.Camera, world objects.Hittable, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outputPath, err)
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n255\n", cam.ImageWidth, cam.ImageHeight)

	for j := 0; j < cam.ImageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", cam.ImageHeight-j)
		for i := 0; i < cam.ImageWidth; i++ {
			ray := cam.RayForPixel(i, j)
			rtimage.WriteColor(file, r.rayColor(ray, world))
		}
	}

	fmt.Fprintln(os.Stderr, "\rDone.                 ")

	return nil
}

func (r Renderer) rayColor(ray rtmath.Ray, world objects.Hittable) rtmath.Color {
	record, hit := world.Hit(ray, rtmath.NewInterval(0, stdmath.Inf(1)))
	if hit {
		normal := record.Normal
		return rtmath.NewVec3(normal.X+1, normal.Y+1, normal.Z+1).Mul(0.5)
	}

	unitDirection := rtmath.UnitVector(ray.Direction)
	a := 0.5 * (unitDirection.Y + 1.0)

	white := rtmath.NewVec3(1.0, 1.0, 1.0)
	blue := rtmath.NewVec3(0.5, 0.7, 1.0)

	return white.Mul(1.0 - a).Add(blue.Mul(a))
}
