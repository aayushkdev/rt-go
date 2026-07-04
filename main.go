package main

import (
	"fmt"
	"os"

	"github.com/aayushkdev/rt-go/camera"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/objects"
	"github.com/aayushkdev/rt-go/render"
)

func main() {
	cam := camera.New(400, 16.0/9.0)
	world := objects.NewWorld(
		objects.NewSphere(rtmath.NewVec3(0, 0, -1), 0.5),
		objects.NewSphere(rtmath.NewVec3(0, -100.5, -1), 100),
	)

	renderer := render.NewRenderer()
	if err := renderer.Render(cam, world, "image.ppm"); err != nil {
		fmt.Fprintf(os.Stderr, "render failed: %v\n", err)
		os.Exit(1)
	}
}
