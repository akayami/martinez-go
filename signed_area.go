package martinez_go

// Orient2d calculates the orientation of three points.
func Orient2d(p0, p1, p2 Point) float64 {
	return (p2.X-p0.X)*(p1.Y-p0.Y) - (p1.X-p0.X)*(p2.Y-p0.Y)
}

// SignedArea determines the orientation of the triangle formed by three points.
// It returns -1 for counterclockwise, 1 for clockwise, and 0 for collinear points.
func SignedArea(p0, p1, p2 Point) int {
	res := Orient2d(p0, p1, p2)
	if res > 0 {
		return -1
	}
	if res < 0 {
		return 1
	}
	return 0
}
