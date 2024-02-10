package martinez_go

func CompareEvents(e1, e2 *SweepEvent) int {
	p1 := e1.Point
	p2 := e2.Point

	if p1.X > p2.X {
		return 1
	}
	if p1.X < p2.X {
		return -1
	}

	if p1.Y != p2.Y {
		if p1.Y > p2.Y {
			return 1
		}
		return -1
	}

	return SpecialCases(e1, e2, p1, p2)
}

func SpecialCases(e1, e2 *SweepEvent, p1, p2 Point) int {
	if e1.Left != e2.Left {
		if e1.Left {
			return 1
		}
		return -1
	}

	if SignedArea(p1, e1.OtherEvent.Point, e2.OtherEvent.Point) != 0 {
		if !e1.IsBelow(e2.OtherEvent.Point) {
			return 1
		}
		return -1
	}

	if !e1.IsSubject && e2.IsSubject {
		return 1
	}

	return -1
}
