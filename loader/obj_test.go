package loader

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

func TestLoadOBJTriangulatesFace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "quad.obj")
	err := os.WriteFile(path, []byte(`
v 0 0 0
v 1 0 0
v 1 1 0
v 0 1 0
f 1 2 3 4
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJ(path, materials.NewLambertian(rtmath.NewVec3(1, 0, 0)))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(mesh.Triangles), 2; got != want {
		t.Fatalf("len(mesh.Triangles) = %d, want %d", got, want)
	}
}

func TestLoadOBJSupportsSlashFaceTokens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "triangle.obj")
	err := os.WriteFile(path, []byte(`
v 0 0 0
v 1 0 0
v 0 1 0
f 1/1/1 2/2/2 3/3/3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJ(path, materials.NewLambertian(rtmath.NewVec3(1, 0, 0)))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(mesh.Triangles), 1; got != want {
		t.Fatalf("len(mesh.Triangles) = %d, want %d", got, want)
	}
}

func TestLoadOBJUsesVertexNormalsForSmoothTriangles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "triangle.obj")
	err := os.WriteFile(path, []byte(`
v 0 0 0
v 1 0 0
v 0 1 0
vn 0 0 1
vn 0 1 0
vn 1 0 0
f 1//1 2//2 3//3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJ(path, materials.NewLambertian(rtmath.NewVec3(1, 0, 0)))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(mesh.Triangles), 1; got != want {
		t.Fatalf("len(mesh.Triangles) = %d, want %d", got, want)
	}
	triangle := mesh.Triangles[0]
	if !triangle.Smooth {
		t.Fatal("triangle.Smooth = false, want true")
	}
	if triangle.NormalA != rtmath.NewVec3(0, 0, 1) {
		t.Fatalf("triangle.NormalA = %#v", triangle.NormalA)
	}
	if triangle.NormalB != rtmath.NewVec3(0, 1, 0) {
		t.Fatalf("triangle.NormalB = %#v", triangle.NormalB)
	}
	if triangle.NormalC != rtmath.NewVec3(1, 0, 0) {
		t.Fatalf("triangle.NormalC = %#v", triangle.NormalC)
	}
}

func TestLoadOBJUsesTextureCoordinates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "triangle.obj")
	err := os.WriteFile(path, []byte(`
v 0 0 0
v 1 0 0
v 0 1 0
vt 0.2 0.3
vt 0.8 0.3
vt 0.2 0.9
f 1/1 2/2 3/3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJ(path, materials.NewLambertian(rtmath.NewVec3(1, 0, 0)))
	if err != nil {
		t.Fatal(err)
	}

	triangle := mesh.Triangles[0]
	if !triangle.HasUV {
		t.Fatal("triangle.HasUV = false, want true")
	}
	if triangle.UVA != rtmath.NewVec3(0.2, 0.3, 0) {
		t.Fatalf("triangle.UVA = %#v", triangle.UVA)
	}
	if triangle.UVB != rtmath.NewVec3(0.8, 0.3, 0) {
		t.Fatalf("triangle.UVB = %#v", triangle.UVB)
	}
	if triangle.UVC != rtmath.NewVec3(0.2, 0.9, 0) {
		t.Fatalf("triangle.UVC = %#v", triangle.UVC)
	}
}

func TestLoadOBJWithMaterialsUsesMTLColor(t *testing.T) {
	dir := t.TempDir()
	objPath := filepath.Join(dir, "model.obj")
	mtlPath := filepath.Join(dir, "model.mtl")

	err := os.WriteFile(objPath, []byte(`
mtllib model.mtl
v 0 0 0
v 1 0 0
v 0 1 0
usemtl blue
f 1 2 3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(mtlPath, []byte(`
newmtl blue
Kd 0.1 0.2 0.9
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJWithMaterials(objPath)
	if err != nil {
		t.Fatal(err)
	}

	material, ok := mesh.Triangles[0].Material.(materials.Lambertian)
	if !ok {
		t.Fatalf("material = %T, want materials.Lambertian", mesh.Triangles[0].Material)
	}
	albedo := material.Albedo.Value(0, 0, rtmath.Point3{})
	if albedo != rtmath.NewVec3(0.1, 0.2, 0.9) {
		t.Fatalf("material albedo = %#v", albedo)
	}
}

func TestLoadOBJWithMaterialsUsesDiffuseTextureMap(t *testing.T) {
	dir := t.TempDir()
	objPath := filepath.Join(dir, "model.obj")
	mtlPath := filepath.Join(dir, "model.mtl")
	texturePath := filepath.Join(dir, "diffuse.png")

	if err := writeTestPNG(texturePath, color.RGBA{R: 255, G: 128, B: 0, A: 255}); err != nil {
		t.Fatal(err)
	}
	err := os.WriteFile(objPath, []byte(`
mtllib model.mtl
v 0 0 0
v 1 0 0
v 0 1 0
vt 0 0
vt 1 0
vt 0 1
usemtl mapped
f 1/1 2/2 3/3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(mtlPath, []byte(`
newmtl mapped
Kd 0.1 0.1 0.1
map_Kd diffuse.png
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJWithMaterials(objPath)
	if err != nil {
		t.Fatal(err)
	}

	material, ok := mesh.Triangles[0].Material.(materials.Lambertian)
	if !ok {
		t.Fatalf("material = %T, want materials.Lambertian", mesh.Triangles[0].Material)
	}
	if _, ok := material.Albedo.(materials.ImageTexture); !ok {
		t.Fatalf("material.Albedo = %T, want materials.ImageTexture", material.Albedo)
	}
}

func TestLoadOBJWithMaterialsMapsSpecularToMetal(t *testing.T) {
	dir := t.TempDir()
	objPath := filepath.Join(dir, "model.obj")
	mtlPath := filepath.Join(dir, "model.mtl")

	err := os.WriteFile(objPath, []byte(`
mtllib model.mtl
v 0 0 0
v 1 0 0
v 0 1 0
usemtl shiny
f 1 2 3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(mtlPath, []byte(`
newmtl shiny
Kd 0.6 0.5 0.1
Ks 0.5 0.5 0.5
Ns 10
illum 2
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJWithMaterials(objPath)
	if err != nil {
		t.Fatal(err)
	}

	material, ok := mesh.Triangles[0].Material.(materials.Metal)
	if !ok {
		t.Fatalf("material = %T, want materials.Metal", mesh.Triangles[0].Material)
	}
	albedo := material.Albedo.Value(0, 0, rtmath.Point3{})
	if albedo != rtmath.NewVec3(0.6, 0.5, 0.1) {
		t.Fatalf("material albedo = %#v", albedo)
	}
	if material.Fuzz != 0.5 {
		t.Fatalf("material.Fuzz = %v, want 0.5", material.Fuzz)
	}
}

func TestLoadOBJWithMaterialsUsesDiffuseTextureMapForMetal(t *testing.T) {
	dir := t.TempDir()
	objPath := filepath.Join(dir, "model.obj")
	mtlPath := filepath.Join(dir, "model.mtl")
	texturePath := filepath.Join(dir, "metal.png")

	if err := writeTestPNG(texturePath, color.RGBA{R: 190, G: 190, B: 210, A: 255}); err != nil {
		t.Fatal(err)
	}
	err := os.WriteFile(objPath, []byte(`
mtllib model.mtl
v 0 0 0
v 1 0 0
v 0 1 0
vt 0 0
vt 1 0
vt 0 1
usemtl mapped_metal
f 1/1 2/2 3/3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(mtlPath, []byte(`
newmtl mapped_metal
Kd 0.5 0.5 0.5
Ks 0.7 0.7 0.7
Ns 50
illum 2
map_Kd metal.png
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJWithMaterials(objPath)
	if err != nil {
		t.Fatal(err)
	}

	material, ok := mesh.Triangles[0].Material.(materials.Metal)
	if !ok {
		t.Fatalf("material = %T, want materials.Metal", mesh.Triangles[0].Material)
	}
	if _, ok := material.Albedo.(materials.ImageTexture); !ok {
		t.Fatalf("material.Albedo = %T, want materials.ImageTexture", material.Albedo)
	}
}

func writeTestPNG(path string, pixel color.RGBA) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	image := stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1))
	image.Set(0, 0, pixel)

	return png.Encode(file, image)
}

func TestLoadOBJWithMaterialsMapsTransparencyToDielectric(t *testing.T) {
	dir := t.TempDir()
	objPath := filepath.Join(dir, "model.obj")
	mtlPath := filepath.Join(dir, "model.mtl")

	err := os.WriteFile(objPath, []byte(`
mtllib model.mtl
v 0 0 0
v 1 0 0
v 0 1 0
usemtl glass
f 1 2 3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(mtlPath, []byte(`
newmtl glass
Kd 0.9 0.9 1.0
d 0.4
Ni 1.45
illum 4
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJWithMaterials(objPath)
	if err != nil {
		t.Fatal(err)
	}

	material, ok := mesh.Triangles[0].Material.(materials.Dielectric)
	if !ok {
		t.Fatalf("material = %T, want materials.Dielectric", mesh.Triangles[0].Material)
	}
	if material.RefractionIndex != 1.45 {
		t.Fatalf("material.RefractionIndex = %v, want 1.45", material.RefractionIndex)
	}
	tint := material.Tint.Value(0, 0, rtmath.Point3{})
	if tint != rtmath.NewVec3(0.9, 0.9, 1.0) {
		t.Fatalf("material tint = %#v", tint)
	}
}

func TestLoadOBJWithMaterialsFallsBackToRedWhenMTLMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "model.obj")
	err := os.WriteFile(path, []byte(`
mtllib missing.mtl
v 0 0 0
v 1 0 0
v 0 1 0
usemtl missing
f 1 2 3
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	mesh, err := LoadOBJWithMaterials(path)
	if err != nil {
		t.Fatal(err)
	}

	material, ok := mesh.Triangles[0].Material.(materials.Metal)
	if !ok {
		t.Fatalf("material = %T, want materials.Metal", mesh.Triangles[0].Material)
	}
	albedo := material.Albedo.Value(0, 0, rtmath.Point3{})
	if albedo != rtmath.NewVec3(0.8, 0.8, 0.8) {
		t.Fatalf("material albedo = %#v", albedo)
	}
	if material.Fuzz != 0 {
		t.Fatalf("material.Fuzz = %v, want 0", material.Fuzz)
	}
}
