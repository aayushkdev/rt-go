package camera

import (
	stdmath "math"

	rtmath "github.com/aayushkdev/rt-go/math"
)

type Config struct {
	ImageWidth   int
	AspectRatio  float64
	VFov         float64
	LookFrom     rtmath.Point3
	LookAt       rtmath.Point3
	VUp          rtmath.Vec3
	DefocusAngle float64
	FocusDist    float64
}

type Camera struct {
	ImageWidth  int
	ImageHeight int

	Center      rtmath.Point3
	Pixel00     rtmath.Point3
	PixelDeltaU rtmath.Vec3
	PixelDeltaV rtmath.Vec3
	DefocusU    rtmath.Vec3
	DefocusV    rtmath.Vec3
	Defocus     bool
}

func New(config Config) Camera {
	if config.ImageWidth == 0 {
		config.ImageWidth = 100
	}
	if config.AspectRatio == 0 {
		config.AspectRatio = 1
	}
	if config.VFov == 0 {
		config.VFov = 90
	}
	if config.VUp.NearZero() {
		config.VUp = rtmath.NewVec3(0, 1, 0)
	}
	if config.FocusDist == 0 {
		config.FocusDist = config.LookFrom.Sub(config.LookAt).Length()
	}

	imageHeight := int(float64(config.ImageWidth) / config.AspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

	cameraCenter := config.LookFrom
	theta := degreesToRadians(config.VFov)
	h := stdmath.Tan(theta / 2)
	viewportHeight := 2 * h * config.FocusDist
	viewportWidth := viewportHeight * float64(config.ImageWidth) / float64(imageHeight)

	w := rtmath.UnitVector(config.LookFrom.Sub(config.LookAt))
	u := rtmath.UnitVector(rtmath.Cross(config.VUp, w))
	v := rtmath.Cross(w, u)

	viewportU := u.Mul(viewportWidth)
	viewportV := v.Neg().Mul(viewportHeight)

	pixelDeltaU := viewportU.Div(float64(config.ImageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	viewportUpperLeft := cameraCenter.
		Sub(w.Mul(config.FocusDist)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2))
	pixel00 := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).Mul(0.5))
	defocusRadius := config.FocusDist * stdmath.Tan(degreesToRadians(config.DefocusAngle/2))

	return Camera{
		ImageWidth:  config.ImageWidth,
		ImageHeight: imageHeight,
		Center:      cameraCenter,
		Pixel00:     pixel00,
		PixelDeltaU: pixelDeltaU,
		PixelDeltaV: pixelDeltaV,
		DefocusU:    u.Mul(defocusRadius),
		DefocusV:    v.Mul(defocusRadius),
		Defocus:     config.DefocusAngle > 0,
	}
}

func (c Camera) RayForPixel(i, j int) rtmath.Ray {
	return c.RayForPixelSample(i, j, 0, 0)
}

func (c Camera) RayForPixelSample(i, j int, offsetU, offsetV float64) rtmath.Ray {
	return c.RayForPixelSampleRandom(i, j, offsetU, offsetV, nil)
}

func (c Camera) RayForPixelSampleRandom(i, j int, offsetU, offsetV float64, random *rtmath.Random) rtmath.Ray {
	pixelCenter := c.Pixel00.
		Add(c.PixelDeltaU.Mul(float64(i) + offsetU)).
		Add(c.PixelDeltaV.Mul(float64(j) + offsetV))
	rayOrigin := c.rayOrigin(random)
	rayDirection := pixelCenter.Sub(rayOrigin)

	return rtmath.NewRay(rayOrigin, rayDirection)
}

func degreesToRadians(degrees float64) float64 {
	return degrees * rtmath.Pi / 180
}

func (c Camera) rayOrigin(random *rtmath.Random) rtmath.Point3 {
	if !c.Defocus {
		return c.Center
	}

	p := random.InUnitDisk()
	return c.Center.Add(c.DefocusU.Mul(p.X)).Add(c.DefocusV.Mul(p.Y))
}
