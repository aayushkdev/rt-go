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
			SamplesPerPixel: 200,
			MaxDepth:        20,
		},

		// Add or remove objects here.
		// Sphere(x, y, z, radius, material) places a sphere. Keep the camera
		// outside reflective/glass objects, or the render can go black/slow.
		// Model(path).WithHeight(...).At(...) loads and places an OBJ.
		// Triangle(Point(...), Point(...), Point(...), material) adds a triangle.
		Scene: Scene(
			// Large sphere used as the ground.
			Sphere(0, -100.5, -1, 100, Matte(0.8, 0.8, 0.0)),

			// Spheres with different materials.
			Sphere(-1.6, 0.35, -2.2, 0.35, Matte(0.8, 0.2, 0.2)),
			Sphere(-0.55, 0.45, -2.0, 0.45, Metal(0.8, 0.8, 0.8, 0)),
			Sphere(0.55, 0.45, -2.0, 0.45, Metal(0.8, 0.6, 0.2, 0.35)),
			Sphere(1.6, 0.45, -2.2, 0.45, Glass(1.5)),

			// Small matte balls around the main objects.
			Sphere(-2.2, 0.12, -1.45, 0.12, Matte(0.8, 0.3, 0.3)),
			Sphere(-1.75, 0.10, -1.2, 0.10, Matte(0.2, 0.5, 0.9)),
			Sphere(-1.15, 0.11, -1.35, 0.11, Matte(0.9, 0.7, 0.2)),
			Sphere(-0.45, 0.10, -1.25, 0.10, Matte(0.3, 0.8, 0.4)),
			Sphere(0.25, 0.12, -1.3, 0.12, Matte(0.7, 0.3, 0.9)),
			Sphere(0.9, 0.10, -1.2, 0.10, Matte(0.9, 0.4, 0.2)),
			Sphere(1.55, 0.11, -1.35, 0.11, Matte(0.2, 0.8, 0.8)),
			Sphere(2.15, 0.12, -1.5, 0.12, Matte(0.8, 0.8, 0.3)),
			Sphere(-2.35, 0.13, -2.45, 0.13, Matte(0.5, 0.3, 0.8)),
			Sphere(-1.85, 0.10, -2.85, 0.10, Matte(0.8, 0.2, 0.5)),
			Sphere(-1.25, 0.12, -2.95, 0.12, Matte(0.2, 0.7, 0.3)),
			Sphere(-0.65, 0.11, -2.75, 0.11, Matte(0.9, 0.8, 0.5)),
			Sphere(0.0, 0.13, -2.95, 0.13, Matte(0.4, 0.6, 0.9)),
			Sphere(0.65, 0.11, -2.75, 0.11, Matte(0.9, 0.5, 0.4)),
			Sphere(1.25, 0.12, -2.95, 0.12, Matte(0.5, 0.9, 0.5)),
			Sphere(1.85, 0.10, -2.85, 0.10, Matte(0.6, 0.4, 0.9)),
			Sphere(2.35, 0.13, -2.45, 0.13, Matte(0.9, 0.6, 0.2)),
			Sphere(-0.95, 0.09, -3.45, 0.09, Matte(0.3, 0.7, 0.9)),
			Sphere(0.0, 0.10, -3.55, 0.10, Matte(0.8, 0.3, 0.6)),
			Sphere(0.95, 0.09, -3.45, 0.09, Matte(0.4, 0.9, 0.7)),

			// A couple of manual triangles.
			Triangle(
				Point(-0.95, 0.05, -1.45),
				Point(-0.45, 0.05, -1.45),
				Point(-0.7, 0.5, -1.45),
				Matte(0.1, 0.4, 0.9),
			),
			Triangle(
				Point(0.5, 0.05, -1.4),
				Point(1.0, 0.05, -1.4),
				Point(0.75, 0.48, -1.4),
				Metal(0.9, 0.9, 0.9, 0.1),
			),
		),
	}
}
