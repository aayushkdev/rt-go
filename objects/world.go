package objects

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
