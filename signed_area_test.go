package martinez_go

import "testing"

func TestSignedArea(t *testing.T) {
	tests := []struct {
		p0, p1, p2 Point
		expected   int
		desc       string
	}{
		{Point{0, 0}, Point{0, 1}, Point{1, 1}, -1, "negative area"},
		{Point{0, 1}, Point{0, 0}, Point{1, 0}, 1, "positive area"},
		{Point{0, 0}, Point{1, 1}, Point{2, 2}, 0, "collinear, 0 area"},
		{Point{-1, 0}, Point{2, 3}, Point{0, 1}, 0, "point on segment"},
		{Point{2, 3}, Point{-1, 0}, Point{0, 1}, 0, "point on segment"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if got := SignedArea(test.p0, test.p1, test.p2); got != test.expected {
				t.Errorf("SignedArea(%v, %v, %v) = %v; want %v", test.p0, test.p1, test.p2, got, test.expected)
			}
		})
	}
}
