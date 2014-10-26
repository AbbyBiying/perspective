// Perspective: Graphing library for quality control in event-driven systems

// Copyright (C) 2014 Christian Paro <christian.paro@gmail.com>,
//                                   <cparo@digitalocean.com>

// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License version 2 as published by the
// Free Software Foundation.

// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
// details.

// You should have received a copy of the GNU General Public License along with
// this program. If not, see <http://www.gnu.org/licenses/>.

package perspective

import (
	"image"
	"math"
)

type sweep struct {
	w     int         // Width of the visualization
	h     int         // Height of the visualization
	vis   *image.RGBA // Visualization canvas
	tA    float64     // Lower limit of time range to be visualized
	tΩ    float64     // Upper limit of time range to be visualized
	yLog2 float64     // Number of pixels over which elapsed times double
	cΔ    float64     // Increment for color channel value increases
}

// NewSweep returns an sweep-visualization generator.
func NewSweep(
	width int,
	height int,
	minTime int,
	maxTime int,
	yLog2 float64,
	colorSteps int,
	xGrid int) Visualizer {

	return (&sweep{
		width,
		height,
		initializeVisualization(width, height),
		float64(minTime),
		float64(maxTime),
		float64(yLog2),
		saturated / float64(colorSteps)}).drawGrid(xGrid)
}

// Record accepts an EventDataPoint and plots it onto the visualization.
func (v *sweep) Record(e EventDataPoint) {

	tMin := float64(e.Start)
	tMax := float64(e.Start + e.Run)
	y := v.h / 2

	// Each event is drawn as an arc tracing its time of existance, with the
	// x-axis representing absolute time and the y-axis being a logarithmic
	// representation of time elapsed since the event was started. Since
	// recorded events may collide in space with other recorded events in the
	// visualization, we use a color progression to indicate the density of
	// events in a given pixel of the visualization. This requires that we take
	// into account the existing color of the point on the canvas to which the
	// event will be plotted and calculate its new color as a function of its
	// existing color.
	for t := tMin; t <= tMax; t++ {
		x := int(float64(v.w) * (t - v.tA) / (v.tΩ - v.tA))
		yMin := v.h/2 - int(v.yLog2*(math.Log2(math.Max(1, t-tMin))))
		for yʹ := y; yʹ > yMin; yʹ-- {
			y = yʹ
			if e.Status == 0 {
				// Successes are plotted above the center line and allowed to
				// desaturate in high-density regions for reasons of aesthetics
				// and additional expressive range.
				c := getRGBA(v.vis, x, y)
				c.R = uint8(math.Min(saturated, float64(c.R)+v.cΔ/4))
				c.G = uint8(math.Min(saturated, float64(c.G)+v.cΔ/4))
				c.B = uint8(math.Min(saturated, float64(c.B)+v.cΔ))
			} else {
				// Failures are plotted below the center line and kept saturated
				// to make them more visible and for the perceptual advantage of
				// keeping them all red, all the time to clearly convey that
				// they are an indication of something gone wrong.
				c := getRGBA(v.vis, x, v.h-y)
				c.R = uint8(math.Min(saturated, float64(c.R)+v.cΔ))
			}
		}
	}
}

// Render returns the visualization constructed from all previously-recorded
// data points.
func (v *sweep) Render() image.Image {
	return v.vis
}

func (v *sweep) drawGrid(xGrid int) *sweep {

	// Draw vertical grid lines, if vertical divisions were specified
	if xGrid > 0 {
		for x := 0; x < v.w; x = x + v.w/xGrid {
			drawXGridLine(v.vis, x)
		}
	}

	// Draw horizontal grid lines on each doubling of the run time in seconds
	for y := v.h / 2; y < v.h; y = y + int(float64(v.h)/v.yLog2) {
		drawYGridLine(v.vis, y)
		drawYGridLine(v.vis, v.h-y)
	}

	// Draw a line up top, for the sake of tidy appearance
	drawYGridLine(v.vis, 0)

	// Return the seep visualization struct, so this can be conveniently
	// used in the visualization's constructor.
	return v
}
