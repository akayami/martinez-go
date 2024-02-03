package martinez_go

// Equals checks if two points are equal.
func Equals(p1, p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}
