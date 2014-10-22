package perspective

import (
	"image"
	"math"
)

type wave struct {
	w   int              // Width of the visualization
	h   int              // Height of the visualization
	tA  float64          // Lower limit of time range to be visualized
	tΩ  float64          // Upper limit of time range to be visualized
	vis *image.RGBA      // Visualization canvas
	x   int              // Current x-position for recording events
	p   []EventDataPoint // Passing event data points in current x-position
	f   []EventDataPoint // Failing event data points in current x-position
}

// NewWave returns a wave-visualization generator.
func NewWave(width int, height int, minTime int, maxTime int) *wave {
	return &wave{
		width,
		height,
		float64(minTime),
		float64(maxTime),
		initializeVisualization(width, height),
		0,
		[]EventDataPoint{},
		[]EventDataPoint{}}
}

// Record accepts an EventDataPoint and plots it onto the visualization.
//
// NOTE: Event input is expected to be received in chronological order. If
//       it is not received in chronological order, the graph will not be
//       rendered properly (with the severity of the issue being dependent
//       upon the degree of deviation between the input order and the ideal
//       chronologically-sorted input.
func (v *wave) Record(e EventDataPoint) {
	pʹ := make([]EventDataPoint, 0, len(v.p)+64)
	fʹ := make([]EventDataPoint, 0, len(v.f)+64)
	for _, p := range v.p {
		if p.Start+p.Run > e.Start {
			pʹ = append(pʹ, p)
		}
	}
	for _, f := range v.f {
		if f.Start+f.Run > e.Start {
			fʹ = append(fʹ, f)
		}
	}
	v.p = pʹ
	v.f = fʹ
	if e.Status == 0 {
		v.p = append(v.p, e)
	} else {
		v.f = append(v.f, e)
	}
	t := float64(e.Start)
	xʹ := int(float64(v.w) * (t - v.tA) / (v.tΩ - v.tA))
	for xʹ > v.x {
		v.x++
		yP := 0
		yF := 0
		for i := 0; i < len(v.p); i++ {
			var (
				p    = v.p[len(v.p)-i-1]
				prog = float64(e.Start-p.Start) / float64(p.Run+1)
				rg16 = uint32(math.Min(maxC16, float64(bg<<8+maxC16*prog/4)))
				b16  = uint32(math.Min(maxC16, float64(bg<<8+maxC16*prog)))
				yPʹ  = yP + 1
			)
			for ; yP < yPʹ; yP++ {
				plot(v.vis, v.x, v.h/2-yP, rg16, rg16, b16)
			}
		}
		for i := 0; i < len(v.f); i++ {
			var (
				f    = v.f[len(v.f)-i-1]
				prog = float64(e.Start-f.Start) / float64(f.Run+1)
				r16  = uint32(math.Min(maxC16, float64(bg<<8+maxC16*prog)))
				gb16 = uint32(math.Min(maxC16, float64(bg<<8+maxC16*prog/4)))
				yFʹ  = yF + 1
			)
			for ; yF < yFʹ; yF++ {
				plot(v.vis, v.x, v.h/2+yF, r16, gb16, gb16)
			}
		}
	}
}

// Render returns the visualization constructed from all previously-recorded
// data points.
func (v *wave) Render() *image.RGBA {
	return v.vis
}
