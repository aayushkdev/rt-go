package render

import (
	"fmt"
	stdmath "math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/aayushkdev/rt-go/camera"
	"github.com/aayushkdev/rt-go/geometry"
	rtimage "github.com/aayushkdev/rt-go/image"
	rtmath "github.com/aayushkdev/rt-go/math"
)

type Renderer struct {
	SamplesPerPixel    int
	MaxDepth           int
	Workers            int
	FlushEveryScanline int
	Background         rtmath.Color
	SkyBackground      bool
}

func NewRenderer() Renderer {
	return Renderer{
		SamplesPerPixel:    10,
		MaxDepth:           50,
		Workers:            runtime.NumCPU(),
		FlushEveryScanline: 10,
		SkyBackground:      true,
	}
}

func (r Renderer) Render(cam camera.Camera, world geometry.Hittable, outputPath string) error {
	headerSize, err := createPPM(outputPath, cam)
	if err != nil {
		return err
	}

	rowStride := cam.ImageWidth * 3
	workerCount := r.workerCount()
	printRenderInfo(cam, r, workerCount)
	jobs := make(chan int)
	results := make(chan scanline, workerCount)
	done := make(chan struct{})
	var closeDone sync.Once
	var wg sync.WaitGroup

	for worker := 0; worker < workerCount; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			random := rtmath.NewRandom(time.Now().UnixNano() + int64(worker))
			for {
				select {
				case <-done:
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					row := r.renderScanline(cam, world, j, rowStride, random)
					select {
					case results <- scanline{Index: j, Row: row}:
					case <-done:
						return
					}
				}
			}
		}(worker)
	}

	go func() {
		defer close(jobs)
		for j := 0; j < cam.ImageHeight; j++ {
			select {
			case jobs <- j:
			case <-done:
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	pendingRows := make(map[int][]byte)
	completedRows := 0
	var renderErr error
	progressTicker := time.NewTicker(time.Second)
	defer progressTicker.Stop()
	printProgress(cam.ImageHeight, completedRows, workerCount)

	for results != nil {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			if renderErr != nil {
				continue
			}

			completedRows++
			printProgress(cam.ImageHeight, completedRows, workerCount)

			pendingRows[result.Index] = result.Row
			if r.FlushEveryScanline <= 0 || len(pendingRows) >= r.FlushEveryScanline {
				renderErr = writeRows(outputPath, headerSize, rowStride, pendingRows)
				if renderErr != nil {
					closeDone.Do(func() { close(done) })
					continue
				}
				pendingRows = make(map[int][]byte)
			}
		case <-progressTicker.C:
			printProgress(cam.ImageHeight, completedRows, workerCount)
		}
	}

	if renderErr != nil {
		return renderErr
	}

	if len(pendingRows) > 0 {
		if err := writeRows(outputPath, headerSize, rowStride, pendingRows); err != nil {
			return err
		}
	}

	printRenderStatus("Done.")
	fmt.Fprintln(os.Stderr)

	return nil
}

func printRenderInfo(cam camera.Camera, renderer Renderer, workerCount int) {
	fmt.Fprintf(os.Stderr, "Image: %dx%d\n", cam.ImageWidth, cam.ImageHeight)
	fmt.Fprintf(os.Stderr, "Samples: %d\n", renderer.SamplesPerPixel)
	fmt.Fprintf(os.Stderr, "Max depth: %d\n", renderer.MaxDepth)
	fmt.Fprintf(os.Stderr, "Workers: %d\n", workerCount)
}

func printProgress(totalRows, completedRows, workerCount int) {
	printRenderStatus(fmt.Sprintf("Scanlines remaining: %d", totalRows-completedRows))
}

func printRenderStatus(message string) {
	fmt.Fprintf(os.Stderr, "\r%-90s", message)
}

type scanline struct {
	Index int
	Row   []byte
}

func (r Renderer) workerCount() int {
	if r.Workers > 0 {
		return r.Workers
	}
	return runtime.NumCPU()
}

func (r Renderer) renderScanline(cam camera.Camera, world geometry.Hittable, j, rowStride int, random *rtmath.Random) []byte {
	row := make([]byte, rowStride)
	for i := 0; i < cam.ImageWidth; i++ {
		pixelColor := rtmath.NewVec3(0, 0, 0)
		for sample := 0; sample < r.SamplesPerPixel; sample++ {
			offsetU, offsetV := stratifiedSampleOffset(sample, r.SamplesPerPixel, random)
			ray := cam.RayForPixelSampleRandom(i, j, offsetU, offsetV, random)
			pixelColor = pixelColor.Add(r.rayColor(ray, world, r.MaxDepth, random))
		}
		rByte, gByte, bByte := rtimage.ColorBytes(pixelColor, r.SamplesPerPixel)
		offset := i * 3
		row[offset] = rByte
		row[offset+1] = gByte
		row[offset+2] = bByte
	}

	return row
}

func createPPM(outputPath string, cam camera.Camera) (int, error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", outputPath, err)
	}

	header := fmt.Sprintf("P6\n%d %d\n255\n", cam.ImageWidth, cam.ImageHeight)
	if _, err := file.WriteString(header); err != nil {
		file.Close()
		return 0, fmt.Errorf("write header %s: %w", outputPath, err)
	}

	imageSize := int64(cam.ImageWidth * cam.ImageHeight * 3)
	if err := file.Truncate(int64(len(header)) + imageSize); err != nil {
		file.Close()
		return 0, fmt.Errorf("size %s: %w", outputPath, err)
	}

	if err := file.Close(); err != nil {
		return 0, fmt.Errorf("close %s: %w", outputPath, err)
	}

	return len(header), nil
}

func writeRows(outputPath string, headerSize, rowStride int, rows map[int][]byte) error {
	file, err := os.OpenFile(outputPath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", outputPath, err)
	}

	for rowIndex, row := range rows {
		rowOffset := int64(headerSize + rowIndex*rowStride)
		if _, err := file.WriteAt(row, rowOffset); err != nil {
			file.Close()
			return fmt.Errorf("write scanline %d: %w", rowIndex, err)
		}
	}

	if err := file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("sync %s: %w", outputPath, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", outputPath, err)
	}

	return nil
}

func (r Renderer) rayColor(ray rtmath.Ray, world geometry.Hittable, depth int, random *rtmath.Random) rtmath.Color {
	if depth <= 0 {
		return rtmath.NewVec3(0, 0, 0)
	}

	record, hit := world.Hit(ray, rtmath.NewInterval(0.001, stdmath.Inf(1)))
	if hit {
		emitted := record.Material.Emitted(record.HitInfo)
		scatter, ok := record.Material.Scatter(ray, record.HitInfo, random)
		if !ok {
			return emitted
		}

		if scatter.SkipPDF {
			return emitted.Add(scatter.Attenuation.MulVec(r.rayColor(scatter.Scattered, world, depth-1, random)))
		}

		scattered := rtmath.NewRay(record.Point, scatter.PDF.Generate(random))
		pdfValue := scatter.PDF.Value(scattered.Direction)
		if pdfValue <= 0 {
			return emitted
		}

		scatteringPDF := record.Material.ScatteringPDF(ray, record.HitInfo, scattered)
		scatterColor := r.rayColor(scattered, world, depth-1, random).
			Mul(scatteringPDF / pdfValue).
			MulVec(scatter.Attenuation)

		return emitted.Add(scatterColor)
	}

	if !r.SkyBackground {
		return r.Background
	}

	unitDirection := rtmath.UnitVector(ray.Direction)
	a := 0.5 * (unitDirection.Y + 1.0)
	white := rtmath.NewVec3(1.0, 1.0, 1.0)
	blue := rtmath.NewVec3(0.5, 0.7, 1.0)

	return white.Mul(1.0 - a).Add(blue.Mul(a))
}

func sampleOffset(random *rtmath.Random) float64 {
	return random.Float64() - 0.5
}

func stratifiedSampleOffset(sample, samplesPerPixel int, random *rtmath.Random) (float64, float64) {
	if samplesPerPixel <= 1 {
		return sampleOffset(random), sampleOffset(random)
	}

	columns := int(stdmath.Ceil(stdmath.Sqrt(float64(samplesPerPixel))))
	rows := int(stdmath.Ceil(float64(samplesPerPixel) / float64(columns)))
	x := sample % columns
	y := sample / columns

	offsetU := (float64(x)+random.Float64())/float64(columns) - 0.5
	offsetV := (float64(y)+random.Float64())/float64(rows) - 0.5

	return offsetU, offsetV
}
