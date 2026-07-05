package main

import (
	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/scene"
)

// Camera starts a camera config builder with sensible defaults.
type CameraBuilder struct {
	config camera.Config
}

func Camera() CameraBuilder {
	return CameraBuilder{
		config: camera.Config{
			ImageWidth:   800,
			AspectRatio:  16.0 / 9.0,
			VFov:         30,
			LookFrom:     rtmath.NewVec3(0, 1, 2),
			LookAt:       rtmath.NewVec3(0, 0, -2),
			VUp:          rtmath.NewVec3(0, 1, 0),
			DefocusAngle: 0,
			FocusDist:    4,
		},
	}
}

func (c CameraBuilder) Size(width int) CameraBuilder {
	c.config.ImageWidth = width
	return c
}

func (c CameraBuilder) FOV(vfov float64) CameraBuilder {
	c.config.VFov = vfov
	return c
}

func (c CameraBuilder) From(x, y, z float64) CameraBuilder {
	c.config.LookFrom = rtmath.NewVec3(x, y, z)
	return c
}

func (c CameraBuilder) LookAt(x, y, z float64) CameraBuilder {
	c.config.LookAt = rtmath.NewVec3(x, y, z)
	return c
}

func (c CameraBuilder) Focus(distance float64) CameraBuilder {
	c.config.FocusDist = distance
	return c
}

func (c CameraBuilder) Defocus(angle float64) CameraBuilder {
	c.config.DefocusAngle = angle
	return c
}

func (c CameraBuilder) Config() camera.Config {
	return c.config
}

// Scene groups all objects that should be rendered.
func Scene(objects ...scene.Object) scene.Config {
	return scene.New(objects...)
}

// Model loads an OBJ file. Use WithHeight and At to size/place it.
func Model(path string) scene.Object {
	return scene.Model(path)
}

// Sphere creates a sphere at x,y,z with radius and material.
func Sphere(x, y, z, radius float64, material materials.Material) scene.Object {
	return scene.Sphere(x, y, z, radius, material)
}

// Triangle creates a single triangle from three points and one material.
func Triangle(a, b, c rtmath.Point3, material materials.Material) scene.Object {
	return scene.Triangle(a, b, c, material)
}

// Point is a short helper for 3D positions.
func Point(x, y, z float64) rtmath.Point3 {
	return rtmath.NewVec3(x, y, z)
}

// Matte creates a diffuse material.
func Matte(r, g, b float64) materials.Material {
	return materials.NewLambertian(rtmath.NewVec3(r, g, b))
}

// Metal creates a reflective material. Lower fuzz is sharper.
func Metal(r, g, b, fuzz float64) materials.Material {
	return materials.NewMetal(rtmath.NewVec3(r, g, b), fuzz)
}

// Glass creates a transparent/refractive material.
func Glass(refractionIndex float64) materials.Material {
	return materials.NewDielectric(refractionIndex)
}
