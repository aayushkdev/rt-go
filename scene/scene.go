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

func buildObject(object Object) geometry.Hittable {
	switch object.Kind {
	case "model":
		return loadModel(object)
	case "sphere":
		return geometry.NewSphere(object.Center, object.Radius, object.Material)
	case "triangle":
		return geometry.NewTriangle(object.A, object.B, object.C, object.Material)
	case "quad":
		return geometry.NewQuad(object.A, object.U, object.V, object.Material)
	default:
		return geometry.EmptyHittable{}
	}
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
