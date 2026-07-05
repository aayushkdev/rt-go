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

	r := intensity.Clamp(linearToGamma(cleanComponent(pixelColor.X * scale)))
	g := intensity.Clamp(linearToGamma(cleanComponent(pixelColor.Y * scale)))
	b := intensity.Clamp(linearToGamma(cleanComponent(pixelColor.Z * scale)))

	return byte(255.999 * r), byte(255.999 * g), byte(255.999 * b)
}

func cleanComponent(component float64) float64 {
	if stdmath.IsNaN(component) || stdmath.IsInf(component, 0) {
		return 0
	}

	return component
}

func linearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return stdmath.Sqrt(linearComponent)
	}

	return 0
}
