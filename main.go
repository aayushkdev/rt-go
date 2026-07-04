package main

import (
	"fmt"
	"os"

	"github.com/aayushkdev/rt-go/camera"
	rtimage "github.com/aayushkdev/rt-go/image"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func rayColor(ray rtmath.Ray) rtmath.Color {
	unitDirection := rtmath.UnitVector(ray.Direction)
	a := 0.5 * (unitDirection.Y + 1.0)

	white := rtmath.NewVec3(1.0, 1.0, 1.0)
	blue := rtmath.NewVec3(0.5, 0.7, 1.0)

	return white.Mul(1.0 - a).Add(blue.Mul(a))
}

func main() {
	cam := camera.New(400, 16.0/9.0)

	file, err := os.Create("image.ppm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create image.ppm: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n255\n", cam.ImageWidth, cam.ImageHeight)

	for j := 0; j < cam.ImageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", cam.ImageHeight-j)
		for i := 0; i < cam.ImageWidth; i++ {
			ray := cam.RayForPixel(i, j)
			rtimage.WriteColor(file, rayColor(ray))
		}
	}

	fmt.Fprintln(os.Stderr, "\rDone.                 ")
}
