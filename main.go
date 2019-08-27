package main

import (
	"math/rand"

	"github.com/deadsy/sdfx/sdf"
)

type series []float64

const spacing = 5.0
const depth = 1
const numPoints = 7
const numSeries = 30
const depthOverlap = 0.1
const padding = 2
const border = 5

func main() {

	s := sdf.Transform3D(
		sdf.Box3D(sdf.V3{
			X: (numPoints + padding+1) * spacing + (border/ 2),
			Y: 2,
			Z: numSeries * (depth + depth*depthOverlap) + (border/ 2),
		}, 0),
		sdf.Translate3d(sdf.V3{
			X: ((numPoints + padding+1) * spacing) / 2,
			Y: 0,
			Z: (numSeries * (depth + depth*depthOverlap)) / 2,
		}),
	)

	for i := 0.0; i < numSeries; i++ {
		data := randSeries(numPoints)
		s2 := series3d(data)

		s2 = sdf.Transform3D(
			s2,
			sdf.Translate3d(
				sdf.V3{
					X: 0,
					Y: 0,
					Z: depth*i,
				},
			),
		)
		s = sdf.Union3D(s, s2)
	}
	sdf.RenderSTL(s, 200, "result.stl")
}

func series3d(s series) sdf.SDF3 {
	b := sdf.NewBezier()

	curPos := 0.0
	for _, p := range zeroSeries(s) {
		b.Add(curPos, p)
		curPos += spacing
	}
	b.Close()

	p := b.Polygon()
	s0 := sdf.Polygon2D(p.Vertices())

	return sdf.Extrude3D(s0, depth+depth*depthOverlap)
}

func randSeries(size int) series {
	s := series{}
	for i := 0; i < size; i++ {
		s = append(s, 10*rand.Float64())
	}
	return s
}

func zeroSeries(s series) series {
	// ensure series starts and ends at 0
	for i := 0; i < padding; i++ {
		s = append(series{0}, s...)
		s = append(s, 0)
	}
	return s
}
