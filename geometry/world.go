package geometry

import rtmath "github.com/aayushkdev/rt-go/math"

type World struct {
	Objects []Hittable
}

func NewWorld(objects ...Hittable) World {
	world := World{}
	for _, object := range objects {
		world.Add(object)
	}

	return world
}

func (w *World) Clear() {
	w.Objects = nil
}

func (w *World) Add(object Hittable) {
	w.Objects = append(w.Objects, object)
}

func (w World) Hit(ray rtmath.Ray, rayT rtmath.Interval) (HitRecord, bool) {
	closest := rayT.Max
	hitAnything := false
	closestRecord := HitRecord{}

	for _, object := range w.Objects {
		record, hit := object.Hit(ray, rtmath.NewInterval(rayT.Min, closest))
		if hit {
			hitAnything = true
			closest = record.T
			closestRecord = record
		}
	}

	return closestRecord, hitAnything
}

func (w World) BoundingBox() AABB {
	if len(w.Objects) == 0 {
		return AABB{}
	}

	box := w.Objects[0].BoundingBox()
	for i := 1; i < len(w.Objects); i++ {
		box = SurroundingBox(box, w.Objects[i].BoundingBox())
	}

	return box
}

func (w World) PDFValue(origin rtmath.Point3, direction rtmath.Vec3) float64 {
	samplers := w.samplers()
	if len(samplers) == 0 {
		return 0
	}

	weight := 1.0 / float64(len(samplers))
	sum := 0.0
	for _, sampler := range samplers {
		sum += weight * sampler.PDFValue(origin, direction)
	}

	return sum
}

func (w World) Random(origin rtmath.Point3, random *rtmath.Random) rtmath.Vec3 {
	samplers := w.samplers()
	if len(samplers) == 0 {
		return rtmath.NewVec3(1, 0, 0)
	}

	index := int(random.Float64() * float64(len(samplers)))
	if index >= len(samplers) {
		index = len(samplers) - 1
	}

	return samplers[index].Random(origin, random)
}

func (w World) samplers() []Sampler {
	samplers := []Sampler{}
	for _, object := range w.Objects {
		sampler, ok := object.(Sampler)
		if ok {
			samplers = append(samplers, sampler)
		}
	}

	return samplers
}
