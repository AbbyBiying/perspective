package main

import (
	"flag"
	"log"
	"os"
	"perspective"
	"perspective/feeds"
	"time"
)

// Variables for command-line option flags.
var (
	errorClassConf string  // Optional conf file for error classification.
	typeFilter     int     // Event type code to filter for, if non-negative.
	tA             int     // Lower limit of time range to be visualized.
	tΩ             int     // Upper limit of time range to be visualized.
	xGrid          int     // Number of horizontal grid divisions.
	yLog2          float64 // Number of pixels over which elapsed times double.
	w              int     // Visualization width, in pixels.
	h              int     // Visualization height, in pixels.
	colorSteps     int     // The number of color steps before saturation.
)

// Variables for fixed-position command-line arguments.
var (
	action string // Indication of what type of visualization to generate.
	iPath  string // Filesystem path for input.
	oPath  string // Filesystem path for output.
)

func convertCommaSeparatedToBinary() {
	feeds.ConvertCSVToBinary(iPath, oPath, tA, tΩ, typeFilter, errorClassConf)
}

func generateErrorStackVisualization() {
	v := perspective.NewErrorStack(w, h)
	generateVisualization(v)
}

func generateHistogramVisualization() {
	v := perspective.NewHistogram(w, h, yLog2)
	generateVisualization(v)
}

func generateRollingStackVisualization() {
	v := perspective.NewRollingStack(w, h, tA, tΩ)
	generateVisualization(v)
}

func generateScatterVisualization() {
	v := perspective.NewScatter(w, h, tΩ, tA, yLog2, colorSteps, xGrid)
	generateVisualization(v)
}

func generateSortedWaveVisualization() {
	v := perspective.NewSortedWave(w, h, tA, tΩ)
	generateVisualization(v)
}

func generateStatusStackVisualization() {
	v := perspective.NewStatusStack(w, h)
	generateVisualization(v)
}

func generateSweepVisualization() {
	v := perspective.NewSweep(w, h, tA, tΩ, yLog2, colorSteps, xGrid)
	generateVisualization(v)
}

func generateWaveVisualization() {
	v := perspective.NewWave(w, h, tA, tΩ)
	generateVisualization(v)
}

func generateVisualization(v perspective.Visualizer) {
	feeds.GeneratePNGFromBinLog(iPath, oPath, tA, tΩ, typeFilter, v)
}

func main() {

	flag.StringVar(
		&errorClassConf,
		"error-reason-filter",
		"",
		"Error reason filter congfiguration.")

	flag.IntVar(
		&typeFilter,
		"event-type-id",
		-1,
		"Event type ID to filter for.")

	flag.IntVar(
		&tA,
		"min-time",
		0,
		"Least recent time to show, expressed as seconds in Unix epoch time.")

	flag.IntVar(
		&tΩ,
		"max-time",
		int(time.Now().Unix()),
		"Most recent time to show, expressed as seconds in Unix epoch time.")

	flag.IntVar(
		&xGrid,
		"x-grid",
		0,
		"Number of divisions to be separated with vertical grid lines.")

	flag.Float64Var(
		&yLog2,
		"run-time-scale",
		16,
		"Pixels along y-axis for every doubling in seconds of run time.")

	flag.IntVar(
		&w,
		"width",
		256,
		"Width of the rendered graph, in pixels.")

	flag.IntVar(
		&h,
		"height",
		128,
		"Height of the rendered graph, in pixels.")

	flag.IntVar(
		&colorSteps,
		"color-steps",
		1,
		"Number of color steps to use in rendering before clipping.")

	flag.Parse()

	if flag.NArg() != 3 {
		log.Println("Incorrect argument count.")
		os.Exit(1)
	}

	action = flag.Arg(0)
	iPath = flag.Arg(1)
	oPath = flag.Arg(2)

	if action == "csv-convert" {
		convertCommaSeparatedToBinary()
	} else if action == "vis-error-stack" {
		generateErrorStackVisualization()
	} else if action == "vis-histogram" {
		generateHistogramVisualization()
	} else if action == "vis-rolling-stack" {
		generateRollingStackVisualization()
	} else if action == "vis-scatter" {
		generateScatterVisualization()
	} else if action == "vis-status-stack" {
		generateStatusStackVisualization()
	} else if action == "vis-sweep" {
		generateSweepVisualization()
	} else if action == "vis-wave" {
		generateWaveVisualization()
	} else if action == "vis-wave-sorted" {
		generateSortedWaveVisualization()
	} else {
		log.Println("Unrecognized action.")
		os.Exit(1)
	}
}
