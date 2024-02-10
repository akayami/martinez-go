package martinez_go

type CompareSeg func(le1, le2 *SweepEvent) int

// CompareSegments compares two SweepEvent instances and returns an integer indicating their order.
func CompareSegments(le1, le2 *SweepEvent) int {
	if le1 == le2 {
		return 0
	}

	// Check if segments are not collinear
	if SignedArea(le1.Point, le1.OtherEvent.Point, le2.Point) != 0 ||
		SignedArea(le1.Point, le1.OtherEvent.Point, le2.OtherEvent.Point) != 0 {

		// If they share their left endpoint, use the right endpoint to sort
		if Equals(le1.Point, le2.Point) {
			if le1.IsBelow(le2.OtherEvent.Point) {
				return -1
			}
			return 1
		}

		// Different left endpoint: use the left endpoint to sort
		if le1.Point.X == le2.Point.X {
			if le1.Point.Y < le2.Point.Y {
				return -1
			}
			return 1
		}

		// Check the insertion order in the sweep line
		if CompareEvents(le1, le2) == 1 {
			if le2.IsAbove(le1.Point) {
				return -1
			}
			return 1
		}

		// The line segment associated with le2 inserted after le1
		if le1.IsBelow(le2.Point) {
			return -1
		}
		return 1
	}

	// Check if segments belong to the same polygon
	if le1.IsSubject == le2.IsSubject {
		p1, p2 := le1.Point, le2.Point
		if p1.X == p2.X && p1.Y == p2.Y {
			p1, p2 = le1.OtherEvent.Point, le2.OtherEvent.Point
			if p1.X == p2.X && p1.Y == p2.Y {
				return 0
			}
			if le1.ContourId > le2.ContourId {
				return 1
			}
			return -1
		}
	} else {
		if le1.IsSubject {
			return -1
		}
		return 1
	}
	return CompareEvents(le1, le2)
}
