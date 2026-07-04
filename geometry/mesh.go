package geometry

import (
	stdmath "math"

	"github.com/aayushkdev/rt-go/materials"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Mesh struct {
	Triangles []Triangle
	Tree      Hittable
}

func NewMesh(material materials.Material, faces ...[3]rtmath.Point3) Mesh {
	mesh := Mesh{}
	for _, face := range faces {
		mesh.AddTriangle(face[0], face[1], face[2], material)
	}
	mesh.BuildBVH()

	return mesh
}

func (m *Mesh) AddTriangle(a, b, c rtmath.Point3, material materials.Material) {
	m.Triangles = append(m.Triangles, NewTriangle(a, b, c, material))
	m.Tree = nil
}

func (m Mesh) Bounds() (rtmath.Point3, rtmath.Point3, bool) {
	if len(m.Triangles) == 0 {
		return rtmath.Point3{}, rtmath.Point3{}, false
	}

	min := rtmath.NewVec3(stdmath.Inf(1), stdmath.Inf(1), stdmath.Inf(1))
	max := rtmath.NewVec3(stdmath.Inf(-1), stdmath.Inf(-1), stdmath.Inf(-1))
	for _, triangle := range m.Triangles {
		min, max = includePoint(min, max, triangle.A)
		min, max = includePoint(min, max, triangle.B)
		min, max = includePoint(min, max, triangle.C)
	}

	return min, max, true
}

func (m *Mesh) Scale(factor float64) {
	for i := range m.Triangles {
		m.Triangles[i].A = m.Triangles[i].A.Mul(factor)
		m.Triangles[i].B = m.Triangles[i].B.Mul(factor)
		m.Triangles[i].C = m.Triangles[i].C.Mul(factor)
	}
	m.Tree = nil
}

func (m *Mesh) Translate(offset rtmath.Vec3) {
	for i := range m.Triangles {
		m.Triangles[i].A = m.Triangles[i].A.Add(offset)
		m.Triangles[i].B = m.Triangles[i].B.Add(offset)
		m.Triangles[i].C = m.Triangles[i].C.Add(offset)
	}
	m.Tree = nil
}

func (m *Mesh) FitHeight(height float64) {
	min, max, ok := m.Bounds()
	if !ok {
		return
	}

	currentHeight := max.Y - min.Y
	if currentHeight <= 0 {
		return
	}

	m.Translate(min.Neg())
	m.Scale(height / currentHeight)
}

func (m *Mesh) BuildBVH() {
	items := make([]Hittable, len(m.Triangles))
	for i := range m.Triangles {
		items[i] = m.Triangles[i]
	}
	m.Tree = NewBVH(items)
}

func (m Mesh) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	if m.Tree != nil {
		return m.Tree.Hit(ray, rayT)
	}

	closest := rayT.Max
	hitAnything := false
	closestRecord := HitRecord{}

	for _, triangle := range m.Triangles {
		record, hit := triangle.Hit(ray, rtmath.NewInterval(rayT.Min, closest))
		if hit {
			hitAnything = true
			closest = record.T
			closestRecord = record
		}
	}

	return closestRecord, hitAnything
}

func (m Mesh) BoundingBox() AABB {
	if m.Tree != nil {
		return m.Tree.BoundingBox()
	}
	if len(m.Triangles) == 0 {
		return AABB{}
	}

	box := m.Triangles[0].BoundingBox()
	for i := 1; i < len(m.Triangles); i++ {
		box = SurroundingBox(box, m.Triangles[i].BoundingBox())
	}

	return box
}

func includePoint(min, max, point rtmath.Point3) (rtmath.Point3, rtmath.Point3) {
	return rtmath.NewVec3(
			stdmath.Min(min.X, point.X),
			stdmath.Min(min.Y, point.Y),
			stdmath.Min(min.Z, point.Z),
		), rtmath.NewVec3(
			stdmath.Max(max.X, point.X),
			stdmath.Max(max.Y, point.Y),
			stdmath.Max(max.Z, point.Z),
		)
}
