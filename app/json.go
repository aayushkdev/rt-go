package app

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
	Type       string      `json:"type"`
	Color      []float64   `json:"color"`
	Texture    fileTexture `json:"texture"`
	Tint       fileTexture `json:"tint"`
	Fuzz       float64     `json:"fuzz"`
	Roughness  float64     `json:"roughness"`
	Refraction float64     `json:"refraction"`
}

type fileTexture struct {
	Type     string       `json:"type"`
	Color    []float64    `json:"color"`
	Path     string       `json:"path"`
	Scale    float64      `json:"scale"`
	UVScale  []float64    `json:"uv_scale"`
	UVOffset []float64    `json:"uv_offset"`
	UVRotate float64      `json:"uv_rotate"`
	Even     *fileTexture `json:"even"`
	Odd      *fileTexture `json:"odd"`
}

type fileObject struct {
	Type     string       `json:"type"`
	Material fileMaterial `json:"material"`
	Light    bool         `json:"light"`
	Sample   *bool        `json:"sample"`
	RotateX  float64      `json:"rotate_x"`
	RotateY  float64      `json:"rotate_y"`
	RotateZ  float64      `json:"rotate_z"`
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

func LoadJSONConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read scene file: %w", err)
	}

	return ParseJSONConfig(data)
}

func ParseJSONConfig(data []byte) (Config, error) {
	var file fileConfig
	if err := json.Unmarshal(data, &file); err != nil {
		return Config{}, fmt.Errorf("parse scene json: %w", err)
	}
	if file.Output == "" {
		return Config{}, fmt.Errorf("output is required")
	}
	if len(file.Objects) == 0 {
		return Config{}, fmt.Errorf("objects must contain at least one object")
	}

	cameraConfig, err := cameraFromJSON(file.Camera)
	if err != nil {
		return Config{}, fmt.Errorf("camera: %w", err)
	}

	renderConfig, err := renderFromJSON(file.Render)
	if err != nil {
		return Config{}, fmt.Errorf("render: %w", err)
	}

	objects := make([]scene.Object, 0, len(file.Objects))
	for index, object := range file.Objects {
		built, err := objectFromJSON(object)
		if err != nil {
			return Config{}, fmt.Errorf("object %d: %w", index, err)
		}
		objects = append(objects, built)
	}

	return Config{
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

	return camera.Config{
		ImageWidth:   cameraFile.Size,
		AspectRatio:  cameraFile.Aspect,
		VFov:         cameraFile.FOV,
		LookFrom:     from,
		LookAt:       lookAt,
		VUp:          rtmath.NewVec3(0, 1, 0),
		DefocusAngle: cameraFile.Defocus,
		FocusDist:    cameraFile.Focus,
	}, nil
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
		built = scene.Sphere(center.X, center.Y, center.Z, object.Radius, material)
	case "box":
		min, err := requiredVec(object.Min, "min")
		if err != nil {
			return built, err
		}
		max, err := requiredVec(object.Max, "max")
		if err != nil {
			return built, err
		}
		built = scene.Box(min, max, material)
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
		built = scene.Triangle(a, b, c, material)
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
		built = scene.Quad(q, u, v, material)
	case "floor":
		built = floorObject(object.X1, object.Z1, object.X2, object.Z2, object.Y, material)
	case "wall_x":
		built = wallXObject(object.X, object.Y1, object.Y2, object.Z1, object.Z2, material)
	case "wall_z":
		built = wallZObject(object.Z, object.X1, object.X2, object.Y1, object.Y2, material)
	case "ceiling_light":
		color, err := requiredVec(object.Color, "color")
		if err != nil {
			return built, err
		}
		built = ceilingLightObject(object.X1, object.Z1, object.X2, object.Z2, object.Y, color)
		materialType = "light"
	case "model":
		if object.Path == "" {
			return built, fmt.Errorf("path is required")
		}
		built = scene.Model(object.Path)
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

	built.RotateX = object.RotateX
	built.RotateY = object.RotateY
	built.RotateZ = object.RotateZ

	if len(object.Move) > 0 {
		move, err := requiredVec(object.Move, "move")
		if err != nil {
			return built, err
		}
		built = scene.Translate(built, move.X, move.Y, move.Z)
	}

	shouldSample := materialType == "glass" || materialType == "light" || object.Light
	if object.Sample != nil {
		shouldSample = *object.Sample
	}
	if object.Light || materialType == "light" {
		return scene.AsLight(built), nil
	}
	if shouldSample {
		return scene.AsSampleTarget(built), nil
	}

	return built, nil
}

func materialFromJSON(material fileMaterial) (materials.Material, string, error) {
	switch material.Type {
	case "matte":
		texture, err := textureFromJSON(material.Texture, material.Color)
		if err != nil {
			return nil, "", err
		}
		return materials.NewTexturedLambertian(texture), "matte", nil
	case "metal":
		texture, err := textureFromJSON(material.Texture, material.Color)
		if err != nil {
			return nil, "", err
		}
		return materials.NewTexturedMetal(texture, material.Fuzz), "metal", nil
	case "glass":
		if material.Refraction <= 0 {
			return nil, "", fmt.Errorf("refraction must be greater than zero")
		}
		if material.Roughness < 0 || material.Roughness > 1 {
			return nil, "", fmt.Errorf("roughness must be between 0 and 1")
		}
		tint := materials.Texture(materials.NewSolidColor(rtmath.NewVec3(1, 1, 1)))
		if textureConfigured(material.Tint) {
			var err error
			tint, err = textureFromJSON(material.Tint, material.Tint.Color)
			if err != nil {
				return nil, "", fmt.Errorf("tint: %w", err)
			}
		}
		return materials.NewRoughDielectric(material.Refraction, tint, material.Roughness), "glass", nil
	case "light":
		color, err := requiredVec(material.Color, "color")
		if err != nil {
			return nil, "", err
		}
		return materials.NewDiffuseLight(color), "light", nil
	default:
		return nil, "", fmt.Errorf("unknown material type %q", material.Type)
	}
}

func textureConfigured(texture fileTexture) bool {
	return texture.Type != "" ||
		len(texture.Color) > 0 ||
		texture.Path != "" ||
		texture.Scale != 0 ||
		len(texture.UVScale) > 0 ||
		len(texture.UVOffset) > 0 ||
		texture.UVRotate != 0 ||
		texture.Even != nil ||
		texture.Odd != nil
}

func textureFromJSON(texture fileTexture, fallbackColor []float64) (materials.Texture, error) {
	if texture.Type == "" {
		color, err := requiredVec(fallbackColor, "color")
		if err != nil {
			return nil, err
		}
		return materials.NewSolidColor(color), nil
	}

	switch texture.Type {
	case "solid":
		color, err := requiredVec(texture.Color, "texture.color")
		if err != nil {
			return nil, err
		}
		return applyTextureTransform(materials.NewSolidColor(color), texture)
	case "checker":
		if texture.Even == nil {
			return nil, fmt.Errorf("texture.even is required for checker textures")
		}
		if texture.Odd == nil {
			return nil, fmt.Errorf("texture.odd is required for checker textures")
		}
		even, err := textureFromJSON(*texture.Even, nil)
		if err != nil {
			return nil, fmt.Errorf("texture.even: %w", err)
		}
		odd, err := textureFromJSON(*texture.Odd, nil)
		if err != nil {
			return nil, fmt.Errorf("texture.odd: %w", err)
		}
		return applyTextureTransform(materials.NewCheckerTexture(texture.Scale, even, odd), texture)
	case "image":
		if texture.Path == "" {
			return nil, fmt.Errorf("texture.path is required for image textures")
		}
		imageTexture, err := materials.NewImageTexture(texture.Path)
		if err != nil {
			return nil, err
		}
		return applyTextureTransform(imageTexture, texture)
	case "noise":
		color := rtmath.NewVec3(1, 1, 1)
		if len(texture.Color) > 0 {
			parsedColor, err := requiredVec(texture.Color, "texture.color")
			if err != nil {
				return nil, err
			}
			color = parsedColor
		}
		return applyTextureTransform(materials.NewNoiseTexture(texture.Scale, color), texture)
	default:
		return nil, fmt.Errorf("unknown texture type %q", texture.Type)
	}
}

func applyTextureTransform(texture materials.Texture, config fileTexture) (materials.Texture, error) {
	scaleU, scaleV := 1.0, 1.0
	if len(config.UVScale) > 0 {
		if len(config.UVScale) != 2 {
			return nil, fmt.Errorf("texture.uv_scale must have exactly 2 numbers")
		}
		scaleU, scaleV = config.UVScale[0], config.UVScale[1]
	}

	offsetU, offsetV := 0.0, 0.0
	if len(config.UVOffset) > 0 {
		if len(config.UVOffset) != 2 {
			return nil, fmt.Errorf("texture.uv_offset must have exactly 2 numbers")
		}
		offsetU, offsetV = config.UVOffset[0], config.UVOffset[1]
	}

	if scaleU == 1 && scaleV == 1 && offsetU == 0 && offsetV == 0 && config.UVRotate == 0 {
		return texture, nil
	}

	return materials.NewTextureTransform(texture, scaleU, scaleV, offsetU, offsetV, config.UVRotate), nil
}

func floorObject(x1, z1, x2, z2, y float64, material materials.Material) scene.Object {
	return scene.Quad(
		rtmath.NewVec3(x1, y, z1),
		rtmath.NewVec3(x2-x1, 0, 0),
		rtmath.NewVec3(0, 0, z2-z1),
		material,
	)
}

func wallXObject(x, y1, y2, z1, z2 float64, material materials.Material) scene.Object {
	return scene.Quad(
		rtmath.NewVec3(x, y1, z1),
		rtmath.NewVec3(0, 0, z2-z1),
		rtmath.NewVec3(0, y2-y1, 0),
		material,
	)
}

func wallZObject(z, x1, x2, y1, y2 float64, material materials.Material) scene.Object {
	return scene.Quad(
		rtmath.NewVec3(x1, y1, z),
		rtmath.NewVec3(x2-x1, 0, 0),
		rtmath.NewVec3(0, y2-y1, 0),
		material,
	)
}

func ceilingLightObject(x1, z1, x2, z2, y float64, color rtmath.Color) scene.Object {
	return scene.AsLight(floorObject(x1, z1, x2, z2, y, materials.NewDiffuseLight(color)))
}

func requiredVec(values []float64, name string) (rtmath.Vec3, error) {
	if len(values) != 3 {
		return rtmath.Vec3{}, fmt.Errorf("%s must have exactly 3 numbers", name)
	}

	return rtmath.NewVec3(values[0], values[1], values[2]), nil
}
