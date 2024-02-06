package martinez_go

import (
	"math"
	"testing"
)

// Helper function to check if slices of Point are equal, considering floating point imprecision.
func pointsEqual(a, b []Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i].X-b[i].X) > 1e-9 || math.Abs(a[i].Y-b[i].Y) > 1e-9 {
			return false
		}
	}
	return true
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		a1, a2, b1, b2  Point
		noEndpointTouch bool
		want            []Point
	}{
		{Point{0, 0}, Point{1, 1}, Point{1, 0}, Point{2, 2}, false, nil},
		{Point{0, 0}, Point{1, 1}, Point{1, 0}, Point{10, 2}, false, nil},
		{Point{2, 2}, Point{3, 3}, Point{0, 6}, Point{2, 4}, false, nil},
		{Point{0, 0}, Point{1, 1}, Point{1, 0}, Point{0, 1}, false, []Point{{0.5, 0.5}}},
		{Point{0, 0}, Point{1, 1}, Point{0, 1}, Point{0, 0}, false, []Point{{0, 0}}},
		{Point{0, 0}, Point{1, 1}, Point{0, 1}, Point{1, 1}, false, []Point{{1, 1}}},
		{Point{0, 0}, Point{1, 1}, Point{0.5, 0.5}, Point{1, 0}, false, []Point{{0.5, 0.5}}},
		{Point{0, 0}, Point{10, 10}, Point{1, 1}, Point{5, 5}, false, []Point{{1, 1}, {5, 5}}},
		// Add other tests as needed
	}

	for _, tt := range tests {
		got := SegmentIntersection(tt.a1, tt.a2, tt.b1, tt.b2, tt.noEndpointTouch)
		if !pointsEqual(got, tt.want) {
			t.Errorf("PossibleIntersection(%v, %v, %v, %v, %v) = %v, want %v", tt.a1, tt.a2, tt.b1, tt.b2, tt.noEndpointTouch, got, tt.want)
		}
	}
}
