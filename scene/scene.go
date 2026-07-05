package scene

import (
	"github.com/aayushkdev/rt-go/geometry"
	"github.com/aayushkdev/rt-go/loader"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func BuildWorld(config Config) geometry.World {
	world := geometry.NewWorld()
	for _, object := range config.Objects {
		world.Add(buildObject(object))
	}

	return geometry.NewWorld(geometry.NewBVH(world.Objects))
}

func BuildLights(config Config) geometry.World {
	lights := geometry.NewWorld()
	for _, object := range config.Objects {
		if object.Light {
			lights.Add(buildObject(object))
		}
	}

	return lights
}

func buildObject(object Object) geometry.Hittable {
	var hittable geometry.Hittable

	switch object.Kind {
	case "model":
		hittable = loadModel(object)
	case "sphere":
		hittable = geometry.NewSphere(object.Center, object.Radius, object.Material)
	case "box":
		hittable = geometry.NewBox(object.Min, object.Max, object.Material)
	case "triangle":
		hittable = geometry.NewTriangle(object.A, object.B, object.C, object.Material)
	case "quad":
		hittable = geometry.NewQuad(object.A, object.U, object.V, object.Material)
	default:
		hittable = geometry.EmptyHittable{}
	}

	if object.RotateY != 0 {
		hittable = geometry.NewRotateY(hittable, object.RotateY)
	}
	if !object.Offset.NearZero() {
		hittable = geometry.NewTranslate(hittable, object.Offset)
	}

	return hittable
}

func loadModel(object Object) geometry.Hittable {
	figurine, err := loader.LoadOBJWithMaterials(object.Path)
	if err != nil {
		return fallbackPyramid(materials.NewMetal(rtmath.NewVec3(0.75, 0.75, 0.75), 0.15))
	}

	figurine.FitHeight(object.Height)
	min, max, ok := figurine.Bounds()
	if ok {
		figurine.Translate(rtmath.NewVec3(
			object.Position.X-max.X/2,
			object.Position.Y-min.Y,
			object.Position.Z-max.Z/2,
		))
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
