package scene

import (
	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/geometry"
	"github.com/aayushkdev/rt-go/loader"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func DefaultCameraConfig() camera.Config {
	return camera.Config{
		ImageWidth:   800,
		AspectRatio:  16.0 / 9.0,
		VFov:         30,
		LookFrom:     rtmath.NewVec3(0, 1, 2),
		LookAt:       rtmath.NewVec3(0, 0, -2),
		VUp:          rtmath.NewVec3(0, 1, 0),
		DefocusAngle: 0,
		FocusDist:    4.0,
	}
}

func DefaultWorld() geometry.World {
	materialGround := materials.NewLambertian(rtmath.NewVec3(0.8, 0.8, 0.0))
	materialFigurine := materials.NewMetal(rtmath.NewVec3(0.75, 0.75, 0.75), 0.15)
	figurine := loadFigurine(materialFigurine)

	world := geometry.NewWorld(
		figurine,
		geometry.NewSphere(rtmath.NewVec3(0, -100.5, -1), 100, materialGround),
	)

	return geometry.NewWorld(geometry.NewBVH(world.Objects))
}

func loadFigurine(material materials.Material) geometry.Hittable {
	figurine, err := loader.LoadOBJWithMaterials("models/figurine.obj")
	if err != nil {
		return fallbackPyramid(material)
	}

	figurine.FitHeight(1.0)
	min, max, ok := figurine.Bounds()
	if ok {
		figurine.Translate(rtmath.NewVec3(-max.X/2, 0.02-min.Y, -3.1-max.Z/2))
	}
	figurine.BuildBVH()

	return figurine
}

func fallbackPyramid(material materials.Material) geometry.Mesh {
	return geometry.NewMesh(
		material,
		[3]rtmath.Point3{
			rtmath.NewVec3(-0.35, 0.8, -2),
			rtmath.NewVec3(0.35, 0.8, -2),
			rtmath.NewVec3(0, 1.35, -2),
		},
		[3]rtmath.Point3{
			rtmath.NewVec3(0.35, 0.8, -2),
			rtmath.NewVec3(0, 0.8, -2.55),
			rtmath.NewVec3(0, 1.35, -2),
		},
		[3]rtmath.Point3{
			rtmath.NewVec3(0, 0.8, -2.55),
			rtmath.NewVec3(-0.35, 0.8, -2),
			rtmath.NewVec3(0, 1.35, -2),
		},
	)
}
