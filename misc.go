package martinez_go

import (
	"math"
)

// EMPTY slice to represent an empty result
var EMPTY = [][][]Point{}

// TrivialOperation checks for trivial cases based on the operation type.
func TrivialOperation(subject, clipping [][][]Point, operation int) [][][]Point {
	if len(subject)*len(clipping) == 0 {
		switch operation {
		case Intersection:
			return EMPTY
		case Difference:
			return subject
		case Union, XOR:
			if len(subject) == 0 {
				return clipping
			}
			return subject
		}
	}
	return nil
}

// CompareBBoxes compares bounding boxes for trivial solutions.
func CompareBBoxes(subject, clipping [][][]Point, sbbox, cbbox [4]float64, operation int) [][][]Point {
	if sbbox[0] > cbbox[2] || cbbox[0] > sbbox[2] || sbbox[1] > cbbox[3] || cbbox[1] > sbbox[3] {
		switch operation {
		case Intersection:
			return EMPTY
		case Difference:
			return subject
		case Union, XOR:
			// Concatenate subject and clipping polygons
			result := append(subject, clipping...)
			return result
		}
	}
	return nil
}

// // Boolean performs the Boolean operation on two polygon sets.
func Boolean(subject, clipping [][][]Point, operation int) [][][]Point {
	// Ensure subject and clipping are in the correct format
	// This step is specific to JavaScript's handling of arrays and might not be necessary in Go

	trivial := TrivialOperation(subject, clipping, operation)
	if trivial != nil {
		return trivial
	}

	sbbox := [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	cbbox := [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}

	// eventQueue := NewTinyQueue()
	eventQueue := FillQueue(subject, clipping, &sbbox, &cbbox, operation)

	trivial = CompareBBoxes(subject, clipping, sbbox, cbbox, operation)
	if trivial != nil {
		return trivial
	}

	sortedEvents := SubdivideSegments(eventQueue, subject, clipping, &sbbox, &cbbox, operation)

	contours := ConnectEdges(sortedEvents)

	var polygons [][][]Point
	for _, contour := range contours {
		if contour.IsExterior() {
			rings := [][]Point{contour.Points}
			for _, holeId := range contour.HoleIds {
				rings = append(rings, contours[holeId].Points)
			}
			polygons = append(polygons, rings)
		}
	}

	return polygons
}
