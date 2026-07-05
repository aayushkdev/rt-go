package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/scene"
)

type fileConfig struct {
	Output  string       `json:"output"`
	Camera  fileCamera   `json:"camera"`
	Render  fileRender   `json:"render"`
	Objects []fileObject `json:"objects"`
}

type fileCamera struct {
	Size    int       `json:"size"`
	Aspect  float64   `json:"aspect"`
	FOV     float64   `json:"fov"`
	From    []float64 `json:"from"`
	LookAt  []float64 `json:"look_at"`
	Focus   float64   `json:"focus"`
	Defocus float64   `json:"defocus"`
}

type fileRender struct {
	Samples              int       `json:"samples"`
	MaxDepth             int       `json:"max_depth"`
	Workers              int       `json:"workers"`
	FlushEveryScanline   int       `json:"flush_every_scanline"`
	Background           []float64 `json:"background"`
	Sky                  bool      `json:"sky"`
	SamplingTargetWeight float64   `json:"sampling_target_weight"`
}

type fileMaterial struct {
	Type       string    `json:"type"`
	Color      []float64 `json:"color"`
	Fuzz       float64   `json:"fuzz"`
	Refraction float64   `json:"refraction"`
}

type fileObject struct {
	Type     string       `json:"type"`
	Material fileMaterial `json:"material"`
	Light    bool         `json:"light"`
	Sample   *bool        `json:"sample"`
	RotateY  float64      `json:"rotate_y"`
	Move     []float64    `json:"move"`

	Path     string    `json:"path"`
	Height   float64   `json:"height"`
	Position []float64 `json:"position"`

	Center []float64 `json:"center"`
	Radius float64   `json:"radius"`
	Min    []float64 `json:"min"`
	Max    []float64 `json:"max"`
	A      []float64 `json:"a"`
	B      []float64 `json:"b"`
	C      []float64 `json:"c"`
	Q      []float64 `json:"q"`
	U      []float64 `json:"u"`
	V      []float64 `json:"v"`
	Color  []float64 `json:"color"`

	X  float64 `json:"x"`
	X1 float64 `json:"x1"`
	X2 float64 `json:"x2"`
	Y  float64 `json:"y"`
	Y1 float64 `json:"y1"`
	Y2 float64 `json:"y2"`
	Z  float64 `json:"z"`
	Z1 float64 `json:"z1"`
	Z2 float64 `json:"z2"`
}

func LoadJSONConfig(path string) (AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("read scene file: %w", err)
	}

	return ParseJSONConfig(data)
}

func ParseJSONConfig(data []byte) (AppConfig, error) {
	var file fileConfig
	if err := json.Unmarshal(data, &file); err != nil {
		return AppConfig{}, fmt.Errorf("parse scene json: %w", err)
	}
	if file.Output == "" {
		return AppConfig{}, fmt.Errorf("output is required")
	}
	if len(file.Objects) == 0 {
		return AppConfig{}, fmt.Errorf("objects must contain at least one object")
	}

	cameraConfig, err := cameraFromJSON(file.Camera)
	if err != nil {
		return AppConfig{}, fmt.Errorf("camera: %w", err)
	}

	renderConfig, err := renderFromJSON(file.Render)
	if err != nil {
		return AppConfig{}, fmt.Errorf("render: %w", err)
	}

	objects := make([]scene.Object, 0, len(file.Objects))
	for index, object := range file.Objects {
		built, err := objectFromJSON(object)
		if err != nil {
			return AppConfig{}, fmt.Errorf("object %d: %w", index, err)
		}
		objects = append(objects, built)
	}

	return AppConfig{
		OutputPath: file.Output,
		Camera:     cameraConfig,
		Render:     renderConfig,
		Scene:      scene.New(objects...),
	}, nil
}

func cameraFromJSON(cameraFile fileCamera) (camera.Config, error) {
	from, err := requiredVec(cameraFile.From, "from")
	if err != nil {
		return camera.Config{}, err
	}
	lookAt, err := requiredVec(cameraFile.LookAt, "look_at")
	if err != nil {
		return camera.Config{}, err
	}
	if cameraFile.Size <= 0 {
		return camera.Config{}, fmt.Errorf("size must be greater than zero")
	}
	if cameraFile.Aspect <= 0 {
		return camera.Config{}, fmt.Errorf("aspect must be greater than zero")
	}
	if cameraFile.FOV <= 0 {
		return camera.Config{}, fmt.Errorf("fov must be greater than zero")
	}
	if cameraFile.Focus <= 0 {
		return camera.Config{}, fmt.Errorf("focus must be greater than zero")
	}

	return Camera().
		Size(cameraFile.Size).
		Aspect(cameraFile.Aspect).
		FOV(cameraFile.FOV).
		From(from.X, from.Y, from.Z).
		LookAt(lookAt.X, lookAt.Y, lookAt.Z).
		Focus(cameraFile.Focus).
		Defocus(cameraFile.Defocus).
		Config(), nil
}

func renderFromJSON(render fileRender) (RenderConfig, error) {
	background, err := requiredVec(render.Background, "background")
	if err != nil {
		return RenderConfig{}, err
	}
	if render.Samples <= 0 {
		return RenderConfig{}, fmt.Errorf("samples must be greater than zero")
	}
	if render.MaxDepth <= 0 {
		return RenderConfig{}, fmt.Errorf("max_depth must be greater than zero")
	}
	if render.SamplingTargetWeight < 0 || render.SamplingTargetWeight > 1 {
		return RenderConfig{}, fmt.Errorf("sampling_target_weight must be between 0 and 1")
	}

	return RenderConfig{
		SamplesPerPixel:      render.Samples,
		MaxDepth:             render.MaxDepth,
		Workers:              render.Workers,
		FlushEveryScanline:   render.FlushEveryScanline,
		Background:           background,
		SkyBackground:        render.Sky,
		SamplingTargetWeight: render.SamplingTargetWeight,
	}, nil
}

func objectFromJSON(object fileObject) (scene.Object, error) {
	var material materials.Material
	materialType := ""
	if object.Type != "ceiling_light" && object.Type != "model" {
		var err error
		material, materialType, err = materialFromJSON(object.Material)
		if err != nil {
			return scene.Object{}, err
		}
	}

	var built scene.Object
	switch object.Type {
	case "sphere":
		center, err := requiredVec(object.Center, "center")
		if err != nil {
			return built, err
		}
		if object.Radius <= 0 {
			return built, fmt.Errorf("radius must be greater than zero")
		}
		built = Sphere(center.X, center.Y, center.Z, object.Radius, material)
	case "box":
		min, err := requiredVec(object.Min, "min")
		if err != nil {
			return built, err
		}
		max, err := requiredVec(object.Max, "max")
		if err != nil {
			return built, err
		}
		built = Box(min, max, material)
	case "triangle":
		a, err := requiredVec(object.A, "a")
		if err != nil {
			return built, err
		}
		b, err := requiredVec(object.B, "b")
		if err != nil {
			return built, err
		}
		c, err := requiredVec(object.C, "c")
		if err != nil {
			return built, err
		}
		built = Triangle(a, b, c, material)
	case "quad":
		q, err := requiredVec(object.Q, "q")
		if err != nil {
			return built, err
		}
		u, err := requiredVec(object.U, "u")
		if err != nil {
			return built, err
		}
		v, err := requiredVec(object.V, "v")
		if err != nil {
			return built, err
		}
		built = Quad(q, u, v, material)
	case "floor":
		built = Floor(object.X1, object.Z1, object.X2, object.Z2, object.Y, material)
	case "wall_x":
		built = WallX(object.X, object.Y1, object.Y2, object.Z1, object.Z2, material)
	case "wall_z":
		built = WallZ(object.Z, object.X1, object.X2, object.Y1, object.Y2, material)
	case "ceiling_light":
		color, err := requiredVec(object.Color, "color")
		if err != nil {
			return built, err
		}
		built = CeilingLight(object.X1, object.Z1, object.X2, object.Z2, object.Y, color.X, color.Y, color.Z)
		materialType = "light"
	case "model":
		if object.Path == "" {
			return built, fmt.Errorf("path is required")
		}
		built = Model(object.Path)
		if object.Height > 0 {
			built = built.WithHeight(object.Height)
		}
		if len(object.Position) > 0 {
			position, err := requiredVec(object.Position, "position")
			if err != nil {
				return built, err
			}
			built = built.At(position.X, position.Y, position.Z)
		}
	default:
		return built, fmt.Errorf("unknown object type %q", object.Type)
	}

	if object.RotateY != 0 {
		built = RotateY(built, object.RotateY)
	}
	if len(object.Move) > 0 {
		move, err := requiredVec(object.Move, "move")
		if err != nil {
			return built, err
		}
		built = Translate(built, move.X, move.Y, move.Z)
	}

	shouldSample := materialType == "glass" || materialType == "light" || object.Light
	if object.Sample != nil {
		shouldSample = *object.Sample
	}
	if object.Light || materialType == "light" {
		return AsLight(built), nil
	}
	if shouldSample {
		return AsSampleTarget(built), nil
	}

	return built, nil
}

func materialFromJSON(material fileMaterial) (materials.Material, string, error) {
	switch material.Type {
	case "matte":
		color, err := requiredVec(material.Color, "color")
		if err != nil {
			return nil, "", err
		}
		return Matte(color.X, color.Y, color.Z), "matte", nil
	case "metal":
		color, err := requiredVec(material.Color, "color")
		if err != nil {
			return nil, "", err
		}
		return Metal(color.X, color.Y, color.Z, material.Fuzz), "metal", nil
	case "glass":
		if material.Refraction <= 0 {
			return nil, "", fmt.Errorf("refraction must be greater than zero")
		}
		return Glass(material.Refraction), "glass", nil
	case "light":
		color, err := requiredVec(material.Color, "color")
		if err != nil {
			return nil, "", err
		}
		return Light(color.X, color.Y, color.Z), "light", nil
	default:
		return nil, "", fmt.Errorf("unknown material type %q", material.Type)
	}
}

func requiredVec(values []float64, name string) (rtmath.Vec3, error) {
	if len(values) != 3 {
		return rtmath.Vec3{}, fmt.Errorf("%s must have exactly 3 numbers", name)
	}

	return rtmath.NewVec3(values[0], values[1], values[2]), nil
}
