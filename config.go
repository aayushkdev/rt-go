package main

import (
	"github.com/aayushkdev/rt-go/camera"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/scene"
)

type AppConfig struct {
	OutputPath string
	Camera     camera.Config
	Render     RenderConfig
	Scene      scene.Config
}

type RenderConfig struct {
	SamplesPerPixel      int
	MaxDepth             int
	Workers              int
	FlushEveryScanline   int
	Background           rtmath.Color
	SkyBackground        bool
	SamplingTargetWeight float64
}

// DefaultConfig is the main place to edit the render.
// Run with `go run .` so Go includes every file in this package.
// `go run main.go` only compiles main.go and will not see this config.
func DefaultConfig() AppConfig {
	return AppConfig{
		// OutputPath is where the normal renderer writes the final image.
		OutputPath: "image.ppm",

		// Change these values to move the normal image render camera:
		// Size controls output width, FOV controls zoom, From is the camera
		// position, LookAt is what it points at, and Focus/Defocus control
		// depth of field.
		Camera: Camera().
			Size(600).
			Aspect(1).
			FOV(40).
			From(0, 1.5, 4.0).
			LookAt(0, 1.5, -2.0).
			Focus(6).
			Config(),

		// Render quality settings.
		// SamplesPerPixel reduces noise/aliasing. Higher is cleaner but slower.
		// MaxDepth controls how many times rays can bounce.
		// Workers controls CPU goroutines. Use 0 to use all CPU cores.
		// FlushEveryScanline controls how often image.ppm is synced while rendering.
		// SamplingTargetWeight is the chance to sample lights/important objects
		// instead of the material PDF. Higher usually reduces small-light grain.
		Render: RenderConfig{
			SamplesPerPixel:      10,
			MaxDepth:             20,
			Workers:              0,
			FlushEveryScanline:   10,
			Background:           Point(0, 0, 0),
			SkyBackground:        false,
			SamplingTargetWeight: 0.999999,
		},

		// Object formats:
		// Sphere(x, y, z, radius, material)
		// Box(Point(minX,minY,minZ), Point(maxX,maxY,maxZ), material)
		// Triangle(Point(x1,y1,z1), Point(x2,y2,z2), Point(x3,y3,z3), material)
		// Floor(x1, z1, x2, z2, y, material)
		// WallX(x, y1, y2, z1, z2, material)
		// WallZ(z, x1, x2, y1, y2, material)
		// CeilingLight(x1, z1, x2, z2, y, r, g, b)
		// Translate(object, x, y, z)
		// RotateY(object, angleDegrees)
		// AsSampleTarget(object) makes a non-light object important to sample.
		// Model(path).WithHeight(height).At(x, y, z)
		//
		// Material formats:
		// Matte(r, g, b)
		// Metal(r, g, b, fuzz)
		// Glass(refractionIndex)
		// Light(r, g, b)
		//
		// Keep the camera outside reflective/glass objects, or the render can go black/slow.
		Scene: Scene(
			// Cornell-style room.
			Floor(-2, -4, 2, 0, 0, Matte(0.73, 0.73, 0.73)),
			Floor(-2, -4, 2, 0, 3, Matte(0.73, 0.73, 0.73)),
			WallZ(-4, -2, 2, 0, 3, Matte(0.73, 0.73, 0.73)),
			WallX(-2, 0, 3, 0, -4, Matte(0.12, 0.45, 0.15)),
			WallX(2, 0, 3, -4, 0, Matte(0.65, 0.05, 0.05)),

			// Bright rectangle light on the ceiling.
			CeilingLight(-0.55, -2.45, 0.55, -1.55, 2.98, 15, 15, 15),

			// Boxes inside the room.
			Translate(
				RotateY(
					Box(Point(-0.45, 0, -0.45), Point(0.45, 1.75, 0.45), Matte(0.73, 0.73, 0.73)),
					15,
				),
				-0.65, 0, -2.45,
			),
			Translate(
				RotateY(
					Box(Point(-0.45, 0, -0.45), Point(0.45, 0.9, 0.45), Matte(0.73, 0.73, 0.73)),
					-18,
				),
				0.75, 0, -1.75,
			),
		),
	}
}
