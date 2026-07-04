package scene

import (
	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/objects"
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

func DefaultWorld() objects.World {
	materialGround := materials.NewLambertian(rtmath.NewVec3(0.8, 0.8, 0.0))
	materialCenter := materials.NewLambertian(rtmath.NewVec3(0.1, 0.2, 0.5))
	materialLeft := materials.NewMetal(rtmath.NewVec3(0.75, 0.75, 0.75), 0)
	materialRight := materials.NewDielectric(1.5)

	return objects.NewWorld(
		objects.NewSphere(rtmath.NewVec3(0, 0, -2), 0.5, materialLeft),
		objects.NewSphere(rtmath.NewVec3(-1, 0, -2), 0.5, materialCenter),
		objects.NewSphere(rtmath.NewVec3(1, 0, -2), 0.5, materialRight),
		objects.NewSphere(rtmath.NewVec3(0, -100.5, -1), 100, materialGround),
	)
}
