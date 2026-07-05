package app

import (
	"github.com/aayushkdev/rt-go/camera"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/scene"
)

type Config struct {
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
