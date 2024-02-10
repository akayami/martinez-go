package martinez_go

import (
	"math"
)

var contourId int

func processPolygon(contourOrHole []Point, isSubject bool, depth int, Q TinyQueue, bbox *[4]float64, isExteriorRing bool) {
	for i := 0; i < len(contourOrHole)-1; i++ {
		s1 := contourOrHole[i]
		s2 := contourOrHole[i+1]
		e1 := &SweepEvent{Point: s1, Left: false, IsSubject: isSubject}
		e2 := &SweepEvent{Point: s2, Left: false, OtherEvent: e1, IsSubject: isSubject}
		e1.OtherEvent = e2

		if s1.X == s2.X && s1.Y == s2.Y {
			continue // Skip collapsed edges
		}

		e1.ContourId = depth
		e2.ContourId = depth
		if !isExteriorRing {
			e1.IsExteriorRing = false
			e2.IsExteriorRing = false
		}
		if CompareEvents(e1, e2) == 1 {
			e2.Left = true
		} else {
			e1.Left = true
		}

		bbox[0] = math.Min(bbox[0], s1.X)
		bbox[1] = math.Min(bbox[1], s1.Y)
		bbox[2] = math.Max(bbox[2], s1.X)
		bbox[3] = math.Max(bbox[3], s1.Y)

		Q.Push(e1)
		Q.Push(e2)
	}
}

func FillQueue(subject, clipping [][][]Point, sbbox, cbbox *[4]float64, operation int) TinyQueue {
	eventQueue := NewTinyQueueDefault(nil, CompareEvents)

	for _, polygonSet := range subject {
		for j, contour := range polygonSet {
			isExteriorRing := j == 0
			if isExteriorRing {
				contourId++
			}
			processPolygon(contour, true, contourId, eventQueue, sbbox, isExteriorRing)
		}
	}

	for _, polygonSet := range clipping {
		for j, contour := range polygonSet {
			isExteriorRing := j == 0
			if operation == Difference {
				isExteriorRing = false
			}
			if isExteriorRing {
				contourId++
			}
			processPolygon(contour, false, contourId, eventQueue, cbbox, isExteriorRing)
		}
	}

	return eventQueue
}
