package martinez_go

import "math"

// CrossProduct calculates the cross product of two vectors.
func CrossProduct(a, b Point) float64 {
	return a.X*b.Y - a.Y*b.X
}

// DotProduct calculates the dot product of two vectors.
func DotProduct(a, b Point) float64 {
	return a.X*b.X + a.Y*b.Y
}

// ToPoint converts a parametric line representation back to a point.
func ToPoint(p Point, s float64, d Point) Point {
	return Point{X: p.X + s*d.X, Y: p.Y + s*d.Y}
}

// Intersection finds the intersection between two line segments.
func SegmentIntersection(a1, a2, b1, b2 Point, noEndpointTouch bool) ([]Point, int) {
	va := Point{X: a2.X - a1.X, Y: a2.Y - a1.Y}
	vb := Point{X: b2.X - b1.X, Y: b2.Y - b1.Y}
	e := Point{X: b1.X - a1.X, Y: b1.Y - a1.Y}

	kross := CrossProduct(va, vb)
	sqrKross := kross * kross
	sqrLenA := DotProduct(va, va)

	if sqrKross > 0 {
		s := CrossProduct(e, vb) / kross
		t := CrossProduct(e, va) / kross
		if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
			if (s == 0 || s == 1 || t == 0 || t == 1) && noEndpointTouch {
				return nil, 0
			}
			return []Point{ToPoint(a1, s, va)}, 1
		}
	} else {
		// Check for overlap between segments
		kross = CrossProduct(e, va)
		sqrKross = kross * kross
		if sqrKross > 0 {
			return nil, 0 // Lines are parallel but not the same
		}
		sa := DotProduct(va, e) / sqrLenA
		sb := sa + DotProduct(va, vb)/sqrLenA
		smin := math.Min(sa, sb)
		smax := math.Max(sa, sb)
		if smin <= 1 && smax >= 0 {
			if (smin == 1 || smax == 0) && noEndpointTouch {
				return nil, 0
			}
			return []Point{ToPoint(a1, math.Max(smin, 0), va), ToPoint(a1, math.Min(smax, 1), va)}, 2
		}
	}
	return nil, 0
}
