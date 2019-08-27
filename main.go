package main

import (
	"encoding/csv"
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/deadsy/sdfx/sdf"
)

type series []float64

const spacing = 2.0
const depth = 0.5
const depthOverlap = 0.1
const padding = 2
const border = 1

func main() {
	render(seriesFromZXYCSV(os.Stdin))
}

func render(series []series) {

	numSeries := float64(len(series))
	if numSeries == 0 {
		panic("no data")
	}
	numPoints := float64(len(series[0]))

	// base
	s := sdf.Transform3D(
		sdf.Box3D(sdf.V3{
			X: ((float64(len(series[0])) + padding) * spacing) + (border * 2),
			Y: 0.5,
			Z: (numSeries * (depth + (depth * depthOverlap))) + (border * 2),
		}, 0),
		sdf.Translate3d(sdf.V3{
			X: ((numPoints + padding + 1) * spacing) / 2,
			Y: 0,
			Z: ((numSeries * (depth + (depth * depthOverlap))) / 2) - border,
		}),
	)
	for k, data := range series {

		scaleSeries(data)
		spew.Dump(data)

		s2 := sdf.Transform3D(
			series3d(data, false),
			sdf.Translate3d(
				sdf.V3{
					X: 0,
					Y: 0,
					Z: depth * float64(k),
				},
			),
		)
		s = sdf.Union3D(s, s2)
	}
	sdf.RenderSTL(s, 250, "result.stl")
}

func series3d(s series, smooth bool) sdf.SDF3 {
	b := sdf.NewBezier()

	curPos := 0.0
	for _, p := range zeroSeries(s) {
		if smooth {
			b.Add(curPos, p).HandleFwd(sdf.DtoR(15), 2)
		} else {
			b.Add(curPos, p)
		}
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

// convert a csv with the columns as Z, X, Y e.g. week, day, count
// data must be sorted by Z, X
// Y must be parsable as a float
func seriesFromZXYCSV(r io.Reader) []series {

	csvr := csv.NewReader(r)

	allSeries := []series{}

	var curZ string
	var curSeries series
	for {
		rowStrings, err := csvr.Read()
		if err == io.EOF {
			if curSeries != nil && len(curSeries) > 0 {
				allSeries = append(allSeries, curSeries)
			}
			break
		}
		if err != nil {
			panic(err)
		}
		if curSeries == nil {
			curSeries = series{}
		}
		if rowStrings[0] != curZ {
			if len(curSeries) > 0 {
				allSeries = append(allSeries, curSeries)
			}
			curSeries = series{}
			curZ = rowStrings[0]
			continue
		}

		val, err := strconv.ParseFloat(rowStrings[2], 64)
		if err != nil {
			panic(err)
		}
		curSeries = append(curSeries, val)
	}

	return allSeries
}

func scaleSeries(s series) {
	total := 0.0
	for _, v := range s {
		total += v
	}
	for k, v := range s {
		s[k] = v / total * 30
	}
}
