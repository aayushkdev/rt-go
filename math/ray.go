package math

type Ray struct {
	Origin    Point3
	Direction Vec3
}

func NewRay(origin Point3, direction Vec3) Ray {
	return Ray{Origin: origin, Direction: direction}
}

func (r Ray) At(t float64) Point3 {
	return r.Origin.Add(r.Direction.Mul(t))
}
