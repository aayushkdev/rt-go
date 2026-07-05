package main

import (
	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/scene"
)

type AppConfig struct {
	OutputPath string
	Camera     camera.Config
	Render     RenderConfig
	Scene      scene.Config
}

type RenderConfig struct {
	SamplesPerPixel int
	MaxDepth        int
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
			Size(800).
			FOV(28).
			From(0, 1.15, 2.4).
			LookAt(0, 0.45, -2.35).
			Focus(4).
			Config(),

		// Render quality settings.
		// SamplesPerPixel reduces noise/aliasing. Higher is cleaner but slower.
		// MaxDepth controls how many times rays can bounce.
		Render: RenderConfig{
			SamplesPerPixel: 100,
			MaxDepth:        20,
		},

		// Object formats:
		// Sphere(x, y, z, radius, material)
		// Triangle(Point(x1,y1,z1), Point(x2,y2,z2), Point(x3,y3,z3), material)
		// Floor(x1, z1, x2, z2, y, material)
		// WallX(x, y1, y2, z1, z2, material)
		// WallZ(z, x1, x2, y1, y2, material)
		// CeilingLight(x1, z1, x2, z2, y, r, g, b)
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
			Floor(-2, -4, 2, 0, 0, Matte(0.75, 0.75, 0.75)),
			Floor(-2, -4, 2, 0, 3, Matte(0.75, 0.75, 0.75)),
			WallZ(-4, -2, 2, 0, 3, Matte(0.75, 0.75, 0.75)),
			WallX(-2, 0, 3, 0, -4, Matte(0.75, 0.15, 0.15)),
			WallX(2, 0, 3, -4, 0, Matte(0.15, 0.55, 0.2)),

			// Bright rectangle light on the ceiling.
			CeilingLight(-0.65, -2.65, 0.65, -1.35, 2.98, 10, 10, 10),

			// Objects inside the room.
			Sphere(-0.75, 0.45, -2.15, 0.45, Matte(0.75, 0.75, 0.75)),
			Sphere(0.55, 0.5, -2.55, 0.5, Metal(0.8, 0.8, 0.8, 0.05)),
			Sphere(0.25, 0.35, -1.45, 0.35, Glass(1.5)),
		),
	}
}
