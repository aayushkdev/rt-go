package math

import stdmath "math"

type Interval struct {
	Min float64
	Max float64
}

func NewInterval(min, max float64) Interval {
	return Interval{Min: min, Max: max}
}

func EmptyInterval() Interval {
	return Interval{Min: stdmath.Inf(1), Max: stdmath.Inf(-1)}
}

func UniverseInterval() Interval {
	return Interval{Min: stdmath.Inf(-1), Max: stdmath.Inf(1)}
}

func (i Interval) Size() float64 {
	return i.Max - i.Min
}

func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}
