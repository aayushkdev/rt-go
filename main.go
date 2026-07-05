package main

import (
	_ "embed"
	"fmt"
	"image/png"
	"log"
	stdmath "math"
	"net/http"
	"os"
	"strconv"

	"github.com/aayushkdev/rt-go/camera"
	rtmath "github.com/aayushkdev/rt-go/math"
	"github.com/aayushkdev/rt-go/render"
	"github.com/aayushkdev/rt-go/scene"
)

//go:embed viewer.html
var viewerHTML string

const (
	spawnX     = 0.0
	spawnY     = 1.0
	spawnZ     = 2.0
	spawnYaw   = -stdmath.Pi / 2
	spawnPitch = -0.25
	spawnFocus = 4.0
	spawnDOF   = false
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "viewer" {
		runViewer()
		return
	}

	config := DefaultConfig()
	world := scene.BuildWorld(config.Scene)
	cam := camera.New(config.Camera)
	renderer := render.NewRenderer()
	renderer.SamplesPerPixel = config.Render.SamplesPerPixel
	renderer.MaxDepth = config.Render.MaxDepth
	if err := renderer.Render(cam, world, config.OutputPath); err != nil {
		fmt.Fprintf(os.Stderr, "render failed: %v\n", err)
		os.Exit(1)
	}
}

func runViewer() {
	config := DefaultConfig()
	world := scene.BuildWorld(config.Scene)
	renderer := render.NewRenderer()
	renderer.SamplesPerPixel = 1
	renderer.MaxDepth = 6

	http.HandleFunc("/", serveViewer)
	http.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		config := configFromRequest(r)
		img := renderer.RenderRGBA(camera.New(config), world)

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-store")
		if err := png.Encode(w, img); err != nil {
			log.Printf("encode png: %v", err)
		}
	})

	addr := "localhost:8080"
	fmt.Printf("viewer running at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func configFromRequest(r *http.Request) camera.Config {
	base := DefaultConfig().Camera
	width := queryInt(r, "width", 480)
	height := queryInt(r, "height", 270)
	if height < 1 {
		height = 1
	}

	position := rtmath.NewVec3(
		queryFloat(r, "px", spawnX),
		queryFloat(r, "py", spawnY),
		queryFloat(r, "pz", spawnZ),
	)
	yaw := queryFloat(r, "yaw", spawnYaw)
	pitch := queryFloat(r, "pitch", spawnPitch)
	direction := viewDirection(yaw, pitch)

	base.ImageWidth = width
	base.AspectRatio = float64(width) / float64(height)
	base.LookFrom = position
	base.LookAt = position.Add(direction)
	base.VUp = rtmath.NewVec3(0, 1, 0)
	base.FocusDist = stdmath.Max(0.1, queryFloat(r, "focus", spawnFocus))
	if queryBool(r, "dof", spawnDOF) {
		base.DefocusAngle = 2.0
	} else {
		base.DefocusAngle = 0
	}

	return base
}

func viewDirection(yaw, pitch float64) rtmath.Vec3 {
	cp := stdmath.Cos(pitch)
	return rtmath.UnitVector(rtmath.NewVec3(
		cp*stdmath.Cos(yaw),
		stdmath.Sin(pitch),
		cp*stdmath.Sin(yaw),
	))
}

func queryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return fallback
	}
	return value
}

func queryFloat(r *http.Request, key string, fallback float64) float64 {
	value, err := strconv.ParseFloat(r.URL.Query().Get(key), 64)
	if err != nil {
		return fallback
	}
	return value
}

func queryBool(r *http.Request, key string, fallback bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	return value == "1" || value == "true"
}

func serveViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, viewerHTML)
}
