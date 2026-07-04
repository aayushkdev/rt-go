package camera

import rtmath "github.com/aayushkdev/rt-go/math"

type Camera struct {
	ImageWidth  int
	ImageHeight int

	Center      rtmath.Point3
	Pixel00     rtmath.Point3
	PixelDeltaU rtmath.Vec3
	PixelDeltaV rtmath.Vec3
}

func New(imageWidth int, aspectRatio float64) Camera {
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)
	cameraCenter := rtmath.NewVec3(0, 0, 0)

	viewportU := rtmath.NewVec3(viewportWidth, 0, 0)
	viewportV := rtmath.NewVec3(0, -viewportHeight, 0)

	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	viewportUpperLeft := cameraCenter.
		Sub(rtmath.NewVec3(0, 0, focalLength)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2))
	pixel00 := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).Mul(0.5))

	return Camera{
		ImageWidth:  imageWidth,
		ImageHeight: imageHeight,
		Center:      cameraCenter,
		Pixel00:     pixel00,
		PixelDeltaU: pixelDeltaU,
		PixelDeltaV: pixelDeltaV,
	}
}

func (c Camera) RayForPixel(i, j int) rtmath.Ray {
	pixelCenter := c.Pixel00.
		Add(c.PixelDeltaU.Mul(float64(i))).
		Add(c.PixelDeltaV.Mul(float64(j)))
	rayDirection := pixelCenter.Sub(c.Center)

	return rtmath.NewRay(c.Center, rayDirection)
}
