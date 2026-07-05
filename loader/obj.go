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
	var texCoords []rtmath.Vec3
	var normals []rtmath.Vec3
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
		case "vt":
			texCoord, err := parseOBJTexCoord(fields)
			if err != nil {
				return geometry.Mesh{}, fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			texCoords = append(texCoords, texCoord)
		case "vn":
			normal, err := parseOBJVertex(fields)
			if err != nil {
				return geometry.Mesh{}, fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			normals = append(normals, rtmath.UnitVector(normal))
		case "f":
			if err := addOBJFace(&mesh, vertices, texCoords, normals, fields[1:], currentMaterial); err != nil {
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

func parseOBJTexCoord(fields []string) (rtmath.Vec3, error) {
	if len(fields) < 3 {
		return rtmath.Vec3{}, fmt.Errorf("texture coordinate needs u v")
	}

	u, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return rtmath.Vec3{}, fmt.Errorf("parse texture u: %w", err)
	}
	v, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return rtmath.Vec3{}, fmt.Errorf("parse texture v: %w", err)
	}

	return rtmath.NewVec3(u, v, 0), nil
}

func addOBJFace(mesh *geometry.Mesh, vertices []rtmath.Point3, texCoords, normals []rtmath.Vec3, tokens []string, material materials.Material) error {
	if len(tokens) < 3 {
		return fmt.Errorf("face needs at least 3 vertices")
	}

	face := make([]objFaceVertex, 0, len(tokens))
	for _, token := range tokens {
		vertex, err := parseOBJFaceVertex(token, len(vertices), len(texCoords), len(normals))
		if err != nil {
			return err
		}
		face = append(face, vertex)
	}

	for i := 1; i < len(face)-1; i++ {
		addOBJTriangle(mesh, vertices, texCoords, normals, face[0], face[i], face[i+1], material)
	}

	return nil
}

type objFaceVertex struct {
	VertexIndex   int
	TexCoordIndex int
	NormalIndex   int
	HasTexCoord   bool
	HasNormal     bool
}

func addOBJTriangle(mesh *geometry.Mesh, vertices []rtmath.Point3, texCoords, normals []rtmath.Vec3, a, b, c objFaceVertex, material materials.Material) {
	hasTexCoords := a.HasTexCoord && b.HasTexCoord && c.HasTexCoord
	hasNormals := a.HasNormal && b.HasNormal && c.HasNormal

	if hasTexCoords && hasNormals {
		mesh.AddSmoothTexturedTriangle(
			vertices[a.VertexIndex],
			vertices[b.VertexIndex],
			vertices[c.VertexIndex],
			normals[a.NormalIndex],
			normals[b.NormalIndex],
			normals[c.NormalIndex],
			texCoords[a.TexCoordIndex],
			texCoords[b.TexCoordIndex],
			texCoords[c.TexCoordIndex],
			material,
		)
		return
	}
	if hasTexCoords {
		mesh.AddTexturedTriangle(
			vertices[a.VertexIndex],
			vertices[b.VertexIndex],
			vertices[c.VertexIndex],
			texCoords[a.TexCoordIndex],
			texCoords[b.TexCoordIndex],
			texCoords[c.TexCoordIndex],
			material,
		)
		return
	}
	if a.HasNormal && b.HasNormal && c.HasNormal {
		mesh.AddSmoothTriangle(
			vertices[a.VertexIndex],
			vertices[b.VertexIndex],
			vertices[c.VertexIndex],
			normals[a.NormalIndex],
			normals[b.NormalIndex],
			normals[c.NormalIndex],
			material,
		)
		return
	}

	mesh.AddTriangle(vertices[a.VertexIndex], vertices[b.VertexIndex], vertices[c.VertexIndex], material)
}

func parseOBJFaceVertex(token string, vertexCount, texCoordCount, normalCount int) (objFaceVertex, error) {
	parts := strings.Split(token, "/")
	if len(parts) == 0 || parts[0] == "" {
		return objFaceVertex{}, fmt.Errorf("missing face vertex index")
	}

	vertexIndex, err := parseOBJIndex(parts[0], vertexCount)
	if err != nil {
		return objFaceVertex{}, fmt.Errorf("parse face vertex index: %w", err)
	}

	faceVertex := objFaceVertex{VertexIndex: vertexIndex}
	if len(parts) >= 2 && parts[1] != "" && texCoordCount > 0 {
		texCoordIndex, err := parseOBJIndex(parts[1], texCoordCount)
		if err != nil {
			return objFaceVertex{}, fmt.Errorf("parse face texture index: %w", err)
		}
		faceVertex.TexCoordIndex = texCoordIndex
		faceVertex.HasTexCoord = true
	}
	if len(parts) >= 3 && parts[2] != "" && normalCount > 0 {
		normalIndex, err := parseOBJIndex(parts[2], normalCount)
		if err != nil {
			return objFaceVertex{}, fmt.Errorf("parse face normal index: %w", err)
		}
		faceVertex.NormalIndex = normalIndex
		faceVertex.HasNormal = true
	}

	return faceVertex, nil
}

func parseOBJIndex(text string, count int) (int, error) {
	index, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	if index < 0 {
		index = count + index + 1
	}
	if index <= 0 || index > count {
		return 0, fmt.Errorf("index %d out of range", index)
	}

	return index - 1, nil
}
