package loader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func loadMTL(path string, materialsByName map[string]materials.Material) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open mtl %s: %w", path, err)
	}
	defer file.Close()

	current := newMTLMaterial("")
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "newmtl":
			current.store(materialsByName)
			if len(fields) >= 2 {
				current = newMTLMaterial(fields[1])
			}
		case "Kd":
			color, err := parseMTLColor(fields)
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.diffuse = color
			current.hasDiffuse = true
		case "Ks":
			color, err := parseMTLColor(fields)
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.specular = color
			current.hasSpecular = true
		case "Ns":
			value, err := parseMTLFloat(fields, "Ns")
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.shininess = value
			current.hasShininess = true
		case "Ni":
			value, err := parseMTLFloat(fields, "Ni")
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.refractionIndex = value
			current.hasRefractionIndex = true
		case "d":
			value, err := parseMTLFloat(fields, "d")
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.alpha = value
			current.hasAlpha = true
		case "Tr":
			value, err := parseMTLFloat(fields, "Tr")
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.alpha = 1 - value
			current.hasAlpha = true
		case "illum":
			value, err := parseMTLInt(fields, "illum")
			if err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			current.illumination = value
			current.hasIllumination = true
		case "map_Kd":
			if len(fields) >= 2 {
				current.diffuseMapPath = filepath.Join(filepath.Dir(path), fields[len(fields)-1])
			}
		}
	}
	current.store(materialsByName)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read mtl %s: %w", path, err)
	}

	return nil
}

type mtlMaterial struct {
	name               string
	diffuse            rtmath.Color
	specular           rtmath.Color
	shininess          float64
	refractionIndex    float64
	alpha              float64
	illumination       int
	diffuseMapPath     string
	hasDiffuse         bool
	hasSpecular        bool
	hasShininess       bool
	hasRefractionIndex bool
	hasAlpha           bool
	hasIllumination    bool
}

func newMTLMaterial(name string) mtlMaterial {
	return mtlMaterial{
		name:            name,
		diffuse:         rtmath.NewVec3(0.8, 0.8, 0.8),
		refractionIndex: 1.5,
		alpha:           1,
	}
}

func (m mtlMaterial) store(materialsByName map[string]materials.Material) {
	if m.name == "" {
		return
	}

	if m.isGlass() {
		if texture, ok := m.diffuseTexture(); ok {
			materialsByName[m.name] = materials.NewRoughDielectric(m.refractionIndex, texture, m.glassRoughness())
			return
		}
		if m.hasDiffuse {
			materialsByName[m.name] = materials.NewRoughDielectric(m.refractionIndex, materials.NewSolidColor(m.diffuse), m.glassRoughness())
			return
		}
		materialsByName[m.name] = materials.NewRoughDielectric(m.refractionIndex, materials.NewSolidColor(rtmath.NewVec3(1, 1, 1)), m.glassRoughness())
		return
	}
	if m.isMetal() {
		if texture, ok := m.diffuseTexture(); ok {
			materialsByName[m.name] = materials.NewTexturedMetal(texture, m.metalFuzz())
			return
		}
		materialsByName[m.name] = materials.NewMetal(m.diffuse, m.metalFuzz())
		return
	}
	if texture, ok := m.diffuseTexture(); ok {
		materialsByName[m.name] = materials.NewTexturedLambertian(texture)
		return
	}

	materialsByName[m.name] = materials.NewLambertian(m.diffuse)
}

func (m mtlMaterial) diffuseTexture() (materials.Texture, bool) {
	if m.diffuseMapPath == "" {
		return nil, false
	}

	texture, err := materials.NewImageTexture(m.diffuseMapPath)
	return texture, err == nil
}

func (m mtlMaterial) isGlass() bool {
	if m.hasAlpha && m.alpha < 0.95 {
		return true
	}
	return m.hasIllumination && (m.illumination == 4 || m.illumination == 6 || m.illumination == 7 || m.illumination == 9)
}

func (m mtlMaterial) isMetal() bool {
	if !m.hasSpecular {
		return false
	}
	specularStrength := maxColorComponent(m.specular)
	return specularStrength > 0.05 && (!m.hasIllumination || m.illumination >= 2)
}

func (m mtlMaterial) metalFuzz() float64 {
	if !m.hasShininess {
		return 0.2
	}
	return 1 / (1 + m.shininess/10)
}

func (m mtlMaterial) glassRoughness() float64 {
	if !m.hasShininess {
		return 0
	}
	return 1 / (1 + m.shininess/10)
}

func parseMTLColor(fields []string) (rtmath.Color, error) {
	if len(fields) < 4 {
		return rtmath.Color{}, fmt.Errorf("%s needs r g b", fields[0])
	}

	r, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return rtmath.Color{}, fmt.Errorf("parse Kd r: %w", err)
	}
	g, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return rtmath.Color{}, fmt.Errorf("parse Kd g: %w", err)
	}
	b, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return rtmath.Color{}, fmt.Errorf("parse Kd b: %w", err)
	}

	return rtmath.NewVec3(r, g, b), nil
}

func parseMTLFloat(fields []string, name string) (float64, error) {
	if len(fields) < 2 {
		return 0, fmt.Errorf("%s needs a value", name)
	}

	value, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}

	return value, nil
}

func parseMTLInt(fields []string, name string) (int, error) {
	if len(fields) < 2 {
		return 0, fmt.Errorf("%s needs a value", name)
	}

	value, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}

	return value, nil
}

func maxColorComponent(color rtmath.Color) float64 {
	return max(color.X, max(color.Y, color.Z))
}
