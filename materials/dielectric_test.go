package materials

import (
	stdmath "math"
	"testing"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func TestRoughenDirectionStaysInsideRoughnessCone(t *testing.T) {
	direction := rtmath.NewVec3(0, 0, -1)
	roughness := 0.25
	minCosine := stdmath.Cos(roughness * rtmath.Pi / 2)
	random := rtmath.NewRandom(1)

	for range 100 {
		roughDirection := roughenDirection(direction, roughness, random)
		cosine := rtmath.Dot(roughDirection, direction)
		if cosine < minCosine-1e-12 {
			t.Fatalf("rough direction cosine = %v, want at least %v", cosine, minCosine)
		}
	}
}
