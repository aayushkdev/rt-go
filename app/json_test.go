package app

import (
	stdimage "image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestExampleSceneFilesParse(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Fatal(err)
		}
	}()

	files, err := filepath.Glob(filepath.Join("examples", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no example scene files found")
	}

	for _, file := range files {
		name := filepath.ToSlash(file)
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ParseJSONConfig(data); err != nil {
				t.Fatalf("parse %s: %v", name, err)
			}
		})
	}
}

func TestParseJSONConfigDoesNotUseDefaultScene(t *testing.T) {
	config, err := ParseJSONConfig([]byte(`{
		"output": "custom.ppm",
		"camera": {
			"size": 200,
			"aspect": 1,
			"fov": 40,
			"from": [0, 1, 3],
			"look_at": [0, 1, -1],
			"focus": 4
		},
		"render": {
			"samples": 10,
			"max_depth": 5,
			"workers": 0,
			"flush_every_scanline": 10,
			"background": [0, 0, 0],
			"sky": false,
			"sampling_target_weight": 0.5
		},
		"objects": [
			{
				"type": "sphere",
				"center": [0, 0.5, -1],
				"radius": 0.5,
				"material": { "type": "matte", "color": [0.7, 0.7, 0.7] }
			}
		]
	}`))
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	if config.OutputPath != "custom.ppm" {
		t.Fatalf("output = %q, want custom.ppm", config.OutputPath)
	}
	if len(config.Scene.Objects) != 1 {
		t.Fatalf("object count = %d, want 1", len(config.Scene.Objects))
	}
}

func TestJSONSamplingRules(t *testing.T) {
	config, err := ParseJSONConfig([]byte(`{
		"output": "sample.ppm",
		"camera": {
			"size": 200,
			"aspect": 1,
			"fov": 40,
			"from": [0, 1, 3],
			"look_at": [0, 1, -1],
			"focus": 4
		},
		"render": {
			"samples": 10,
			"max_depth": 5,
			"workers": 0,
			"flush_every_scanline": 10,
			"background": [0, 0, 0],
			"sky": false,
			"sampling_target_weight": 0.5
		},
		"objects": [
			{
				"type": "sphere",
				"center": [0, 0.5, -1],
				"radius": 0.5,
				"material": { "type": "glass", "refraction": 1.5 }
			},
			{
				"type": "sphere",
				"center": [1, 0.5, -1],
				"radius": 0.5,
				"sample": false,
				"material": { "type": "glass", "refraction": 1.5 }
			},
			{
				"type": "sphere",
				"center": [-1, 0.5, -1],
				"radius": 0.5,
				"sample": false,
				"material": { "type": "light", "color": [4, 4, 4] }
			}
		]
	}`))
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	objects := config.Scene.Objects
	if !objects[0].Sample {
		t.Fatalf("glass object should auto-sample")
	}
	if objects[1].Sample {
		t.Fatalf("sample:false should disable glass sampling")
	}
	if !objects[2].Light || !objects[2].Sample {
		t.Fatalf("light object should always be light and sampled")
	}
}

func TestParseJSONConfigSupportsSolidTexture(t *testing.T) {
	config, err := ParseJSONConfig([]byte(`{
		"output": "texture.ppm",
		"camera": {
			"size": 200,
			"aspect": 1,
			"fov": 40,
			"from": [0, 1, 3],
			"look_at": [0, 1, -1],
			"focus": 4
		},
		"render": {
			"samples": 10,
			"max_depth": 5,
			"workers": 0,
			"flush_every_scanline": 10,
			"background": [0, 0, 0],
			"sky": false,
			"sampling_target_weight": 0.5
		},
		"objects": [
			{
				"type": "sphere",
				"center": [0, 0.5, -1],
				"radius": 0.5,
				"material": {
					"type": "matte",
					"texture": { "type": "solid", "color": [0.2, 0.4, 0.8] }
				}
			}
		]
	}`))
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	material, ok := config.Scene.Objects[0].Material.(materials.Lambertian)
	if !ok {
		t.Fatalf("material = %T, want materials.Lambertian", config.Scene.Objects[0].Material)
	}
	color := material.Albedo.Value(0, 0, rtmath.Point3{})
	if color != rtmath.NewVec3(0.2, 0.4, 0.8) {
		t.Fatalf("texture color = %#v", color)
	}
}

func TestTextureFromJSONSupportsChecker(t *testing.T) {
	texture, err := textureFromJSON(fileTexture{
		Type:  "checker",
		Scale: 1,
		Even:  &fileTexture{Type: "solid", Color: []float64{1, 1, 1}},
		Odd:   &fileTexture{Type: "solid", Color: []float64{0, 0, 0}},
	}, nil)
	if err != nil {
		t.Fatalf("build texture: %v", err)
	}

	if _, ok := texture.(materials.CheckerTexture); !ok {
		t.Fatalf("texture = %T, want materials.CheckerTexture", texture)
	}
}

func TestTextureFromJSONSupportsImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "texture.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 25, G: 50, B: 100, A: 255})
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	texture, err := textureFromJSON(fileTexture{Type: "image", Path: path}, nil)
	if err != nil {
		t.Fatalf("build texture: %v", err)
	}

	if _, ok := texture.(materials.ImageTexture); !ok {
		t.Fatalf("texture = %T, want materials.ImageTexture", texture)
	}
}

func TestTextureFromJSONSupportsNoise(t *testing.T) {
	texture, err := textureFromJSON(fileTexture{
		Type:  "noise",
		Scale: 3,
		Color: []float64{0.8, 0.7, 0.6},
	}, nil)
	if err != nil {
		t.Fatalf("build texture: %v", err)
	}

	if _, ok := texture.(materials.NoiseTexture); !ok {
		t.Fatalf("texture = %T, want materials.NoiseTexture", texture)
	}
}
