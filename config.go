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
			FOV(44).
			From(0, 1.35, 4.2).
			LookAt(0, 1.15, -1.6).
			Focus(5.8).
			Config(),

		// Render quality settings.
		// SamplesPerPixel reduces noise/aliasing. Higher is cleaner but slower.
		// MaxDepth controls how many times rays can bounce.
		// Workers controls CPU goroutines. Use 0 to use all CPU cores.
		// FlushEveryScanline controls how often image.ppm is synced while rendering.
		// SamplingTargetWeight is the chance to sample lights/important objects
		// instead of the material PDF. Higher usually reduces small-light grain.
		Render: RenderConfig{
			SamplesPerPixel:      1000,
			MaxDepth:             20,
			Workers:              0,
			FlushEveryScanline:   10,
			Background:           Point(0, 0, 0),
			SkyBackground:        false,
			SamplingTargetWeight: 0.5,
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
			// Material gallery room.
			Floor(-3.0, -4.8, 3.0, 3.6, 0, Matte(0.68, 0.70, 0.66)),
			Floor(-3.0, -4.8, 3.0, 0.9, 3, Matte(0.58, 0.60, 0.62)),
			WallZ(-4.8, -3.0, 3.0, 0, 3, Matte(0.58, 0.60, 0.62)),
			WallX(-3.0, 0, 3, 0.9, -4.8, Matte(0.52, 0.56, 0.58)),
			WallX(3.0, 0, 3, -4.8, 0.9, Matte(0.60, 0.56, 0.52)),

			// Centered soft white key light.
			CeilingLight(-0.65, -2.1, 0.65, -1.1, 2.98, 13, 13, 13),

			// Matte color panels attached to the back wall for visible refraction.
			Quad(Point(-0.78, 0.68, -4.79), Point(0.26, 0, 0), Point(0, 1.05, 0), Matte(0.75, 0.20, 0.16)),
			Quad(Point(-0.39, 0.68, -4.79), Point(0.26, 0, 0), Point(0, 1.05, 0), Matte(0.92, 0.72, 0.18)),
			Quad(Point(0.00, 0.68, -4.79), Point(0.26, 0, 0), Point(0, 1.05, 0), Matte(0.18, 0.42, 0.78)),

			// Matte pedestal.
			Translate(
				RotateY(
					Box(Point(-0.72, 0, -0.42), Point(0.72, 0.34, 0.42), Matte(0.72, 0.70, 0.64)),
					-8,
				),
				0, 0, -1.6,
			),

			// Main materials.
			AsSampleTarget(Sphere(0, 0.82, -1.6, 0.48, Glass(1.5))),
			Sphere(-1.18, 0.38, -1.55, 0.38, Metal(0.88, 0.86, 0.80, 0.05)),
			Sphere(1.18, 0.34, -1.63, 0.34, Matte(0.15, 0.38, 0.85)),

			// Small diffuse color swatches.
			Sphere(-1.35, 0.16, -0.95, 0.16, Matte(0.9, 0.14, 0.10)),
			Sphere(-0.95, 0.14, -0.82, 0.14, Matte(0.95, 0.72, 0.12)),
			Sphere(0.98, 0.13, -0.82, 0.13, Matte(0.15, 0.75, 0.35)),
			Sphere(1.34, 0.15, -0.96, 0.15, Metal(0.75, 0.80, 0.92, 0.35)),
		),
	}
}
