package martinez_go

// Point represents a 2D point with X and Y coordinates.
type Point struct {
	X, Y float64
}

func (f *Point) Equals(p Point) bool {
	if f.X == p.X {
		if f.Y == p.Y {
			return true
		}
	}
	return false
}
