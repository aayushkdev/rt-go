package image

import (
	"fmt"
	"io"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func WriteColor(out io.Writer, pixelColor rtmath.Color) {
	r := pixelColor.X
	g := pixelColor.Y
	b := pixelColor.Z

	rByte := int(255.999 * r)
	gByte := int(255.999 * g)
	bByte := int(255.999 * b)

	fmt.Fprintf(out, "%d %d %d\n", rByte, gByte, bByte)
}
