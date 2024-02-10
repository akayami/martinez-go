package martinez_go

import "math"

func SegmentIntersection(a1, a2, b1, b2 Point, noEndpointTouch bool) []Point {
	va := Point{X: a2.X - a1.X, Y: a2.Y - a1.Y}
	vb := Point{X: b2.X - b1.X, Y: b2.Y - b1.Y}
	e := Point{X: b1.X - a1.X, Y: b1.Y - a1.Y}

	kross := crossProduct(va, vb)
	sqrKross := kross * kross
	sqrLenA := dotProduct(va, va)

	// Check for line intersection
	if sqrKross > 0 {
		s := crossProduct(e, vb) / kross
		if s < 0 || s > 1 {
			return nil
		}
		t := crossProduct(e, va) / kross
		if t < 0 || t > 1 {
			return nil
		}
		if s == 0 || s == 1 {
			if noEndpointTouch {
				return nil
			}
			return []Point{toPoint(a1, s, va)}
		}
		if t == 0 || t == 1 {
			if noEndpointTouch {
				return nil
			}
			return []Point{toPoint(b1, t, vb)}
		}
		return []Point{toPoint(a1, s, va)}
	}

	// Checking for parallel lines
	kross = crossProduct(e, va)
	sqrKross = kross * kross

	if sqrKross > 0 {
		return nil
	}

	sa := dotProduct(va, e) / sqrLenA
	sb := sa + dotProduct(va, vb)/sqrLenA
	smin := math.Min(sa, sb)
	smax := math.Max(sa, sb)

	if smin <= 1 && smax >= 0 {
		if smin == 1 {
			if noEndpointTouch {
				return nil
			}
			return []Point{toPoint(a1, math.Max(smin, 0), va)}
		}
		if smax == 0 {
			if noEndpointTouch {
				return nil
			}
			return []Point{toPoint(a1, math.Min(smax, 1), va)}
		}
		if noEndpointTouch && smin == 0 && smax == 1 {
			return nil
		}
		return []Point{toPoint(a1, math.Max(smin, 0), va), toPoint(a1, math.Min(smax, 1), va)}
	}

	return nil
}

// Helper function to calculate cross product of two vectors
func crossProduct(v1, v2 Point) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Helper function to calculate dot product of two vectors
func dotProduct(v1, v2 Point) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Function to convert from parameter space to point
func toPoint(p Point, s float64, d Point) Point {
	return Point{
		X: p.X + s*d.X,
		Y: p.Y + s*d.Y,
	}
}
