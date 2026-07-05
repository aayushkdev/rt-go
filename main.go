package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aayushkdev/rt-go/app"
)

func main() {
	renderPath := flag.String("r", "", "render the scene JSON file")
	renderPathLong := flag.String("render", "", "render the scene JSON file")
	viewerPath := flag.String("v", "", "start the viewer with the scene JSON file")
	viewerPathLong := flag.String("viewer", "", "start the viewer with the scene JSON file")
	flag.Usage = printUsage
	flag.Parse()

	renderScenePath := firstNonEmpty(*renderPath, *renderPathLong)
	viewerScenePath := firstNonEmpty(*viewerPath, *viewerPathLong)
	if renderScenePath == "" && viewerScenePath == "" {
		printUsage()
		os.Exit(1)
	}
	if renderScenePath != "" && viewerScenePath != "" {
		fmt.Fprintln(os.Stderr, "choose either render or viewer, not both")
		printUsage()
		os.Exit(1)
	}
	if flag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument: %s\n", flag.Arg(0))
		printUsage()
		os.Exit(1)
	}

	scenePath := firstNonEmpty(renderScenePath, viewerScenePath)
	config, err := app.LoadJSONConfig(scenePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load scene failed: %v\n", err)
		os.Exit(1)
	}

	if renderScenePath != "" {
		if err := app.RenderScene(config); err != nil {
			fmt.Fprintf(os.Stderr, "render failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	app.RunViewer(config)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  go run . -r <scene.json>")
	fmt.Fprintln(os.Stderr, "  go run . -v <scene.json>")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "flags:")
	fmt.Fprintln(os.Stderr, "  -r, --render  render a scene to its configured output file")
	fmt.Fprintln(os.Stderr, "  -v, --viewer  start the browser viewer for a scene")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
