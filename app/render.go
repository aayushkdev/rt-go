package app

import (
	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/geometry"
	"github.com/aayushkdev/rt-go/render"
	"github.com/aayushkdev/rt-go/scene"
)

func RenderScene(config Config) error {
	world := scene.BuildWorld(config.Scene)
	samplingTargets := scene.BuildSampleTargets(config.Scene)
	cam := camera.New(config.Camera)
	renderer := rendererFromConfig(config, samplingTargets)

	return renderer.Render(cam, world, config.OutputPath)
}

func rendererFromConfig(config Config, samplingTargets geometry.World) render.Renderer {
	renderer := render.NewRenderer()
	renderer.SamplesPerPixel = config.Render.SamplesPerPixel
	renderer.MaxDepth = config.Render.MaxDepth
	renderer.Workers = config.Render.Workers
	renderer.FlushEveryScanline = config.Render.FlushEveryScanline
	renderer.Background = config.Render.Background
	renderer.SkyBackground = config.Render.SkyBackground
	renderer.SamplingTargetWeight = config.Render.SamplingTargetWeight
	if len(samplingTargets.Objects) > 0 {
		renderer.SamplingTargets = samplingTargets
	}

	return renderer
}
