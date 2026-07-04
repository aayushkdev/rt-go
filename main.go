package main

import (
	"fmt"
	"os"

	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/objects"
	"github.com/aayushkdev/rt-go/render"
)

func main() {
	cam := camera.New(camera.Config{
		ImageWidth:  800,
		AspectRatio: 16.0 / 9.0,
		VFov:        20,
		LookFrom:    rtmath.NewVec3(-2, 2, 1),
		LookAt:      rtmath.NewVec3(0, 0, -1),
		VUp:         rtmath.NewVec3(0, 1, 0),
	})
	materialGround := materials.NewLambertian(rtmath.NewVec3(0.8, 0.8, 0.0))
	materialCenter := materials.NewLambertian(rtmath.NewVec3(0.1, 0.2, 0.5))
	materialLeft := materials.NewMetal(rtmath.NewVec3(0.8, 0.6, 0.2), 0)
	materialRight := materials.NewDielectric(1.5)

	world := objects.NewWorld(
		objects.NewSphere(rtmath.NewVec3(0, 0, -2), 0.5, materialLeft),
		objects.NewSphere(rtmath.NewVec3(-1, 0, -2), 0.5, materialCenter),
		objects.NewSphere(rtmath.NewVec3(1, 0, -2), 0.5, materialRight),
		objects.NewSphere(rtmath.NewVec3(0, -100.5, -1), 100, materialGround),
	)

	renderer := render.NewRenderer()
	if err := renderer.Render(cam, world, "image.ppm"); err != nil {
		fmt.Fprintf(os.Stderr, "render failed: %v\n", err)
		os.Exit(1)
	}
}
