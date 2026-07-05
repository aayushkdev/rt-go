package image

import (
	"fmt"
	"io"
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

func WriteColor(out io.Writer, pixelColor rtmath.Color, samplesPerPixel int) {
	rByte, gByte, bByte := ColorBytes(pixelColor, samplesPerPixel)
	fmt.Fprintf(out, "%d %d %d\n", rByte, gByte, bByte)
}

func ColorBytes(pixelColor rtmath.Color, samplesPerPixel int) (byte, byte, byte) {
	scale := 1.0 / float64(samplesPerPixel)
	intensity := rtmath.NewInterval(0, 0.999)

	r := intensity.Clamp(linearToGamma(pixelColor.X * scale))
	g := intensity.Clamp(linearToGamma(pixelColor.Y * scale))
	b := intensity.Clamp(linearToGamma(pixelColor.Z * scale))

	return byte(255.999 * r), byte(255.999 * g), byte(255.999 * b)
}

func linearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return stdmath.Sqrt(linearComponent)
	}

	return 0
}
