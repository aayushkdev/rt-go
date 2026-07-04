package loader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aayushkdev/rt-go/geometry"
	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

func LoadOBJWithMaterials(path string) (geometry.Mesh, error) {
	return loadOBJ(path, defaultOBJMaterial(), true)
}

func LoadOBJ(path string, material materials.Material) (geometry.Mesh, error) {
	return loadOBJ(path, material, false)
}

func loadOBJ(path string, material materials.Material, useMaterials bool) (geometry.Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return geometry.Mesh{}, fmt.Errorf("open obj %s: %w", path, err)
	}
	defer file.Close()

	var vertices []rtmath.Point3
	mesh := geometry.Mesh{}
	materialsByName := map[string]materials.Material{}
	currentMaterial := material
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
		case "mtllib":
			if useMaterials {
				loadOBJMaterialLibraries(path, fields[1:], materialsByName)
			}
		case "usemtl":
			if useMaterials {
				currentMaterial = materialForOBJName(materialsByName, fields, material)
			}
		case "v":
			vertex, err := parseOBJVertex(fields)
			if err != nil {
				return geometry.Mesh{}, fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			vertices = append(vertices, vertex)
		case "f":
			if err := addOBJFace(&mesh, vertices, fields[1:], currentMaterial); err != nil {
				return geometry.Mesh{}, fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return geometry.Mesh{}, fmt.Errorf("read obj %s: %w", path, err)
	}

	mesh.BuildBVH()
	return mesh, nil
}

func defaultOBJMaterial() materials.Material {
	return materials.NewMetal(rtmath.NewVec3(0.8, 0.8, 0.8), 0)
}

func loadOBJMaterialLibraries(objPath string, names []string, materialsByName map[string]materials.Material) {
	objDir := filepath.Dir(objPath)
	for _, name := range names {
		mtlPath := filepath.Join(objDir, name)
		_ = loadMTL(mtlPath, materialsByName)
	}
}

func materialForOBJName(materialsByName map[string]materials.Material, fields []string, fallback materials.Material) materials.Material {
	if len(fields) < 2 {
		return fallback
	}
	material, ok := materialsByName[fields[1]]
	if !ok {
		return fallback
	}

	return material
}

func parseOBJVertex(fields []string) (rtmath.Point3, error) {
	if len(fields) < 4 {
		return rtmath.Point3{}, fmt.Errorf("vertex needs x y z")
	}

	x, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return rtmath.Point3{}, fmt.Errorf("parse vertex x: %w", err)
	}
	y, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return rtmath.Point3{}, fmt.Errorf("parse vertex y: %w", err)
	}
	z, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return rtmath.Point3{}, fmt.Errorf("parse vertex z: %w", err)
	}

	return rtmath.NewVec3(x, y, z), nil
}

func addOBJFace(mesh *geometry.Mesh, vertices []rtmath.Point3, tokens []string, material materials.Material) error {
	if len(tokens) < 3 {
		return fmt.Errorf("face needs at least 3 vertices")
	}

	face := make([]rtmath.Point3, 0, len(tokens))
	for _, token := range tokens {
		index, err := parseOBJFaceIndex(token, len(vertices))
		if err != nil {
			return err
		}
		face = append(face, vertices[index])
	}

	for i := 1; i < len(face)-1; i++ {
		mesh.AddTriangle(face[0], face[i], face[i+1], material)
	}

	return nil
}

func parseOBJFaceIndex(token string, vertexCount int) (int, error) {
	vertexIndexText := strings.Split(token, "/")[0]
	if vertexIndexText == "" {
		return 0, fmt.Errorf("missing face vertex index")
	}

	index, err := strconv.Atoi(vertexIndexText)
	if err != nil {
		return 0, fmt.Errorf("parse face vertex index: %w", err)
	}
	if index < 0 {
		index = vertexCount + index + 1
	}
	if index <= 0 || index > vertexCount {
		return 0, fmt.Errorf("face vertex index %d out of range", index)
	}

	return index - 1, nil
}
