package martinez_go

// DivideSegment divides a SweepEvent at a point and updates the queue.
func DivideSegment(se *SweepEvent, p Point, queue TinyQueue) {
	r := NewSweepEvent(p, false, se, se.IsSubject, Normal)
	l := NewSweepEvent(p, true, se.OtherEvent, se.IsSubject, Normal)

	if Equals(se.Point, se.OtherEvent.Point) {
		// Log a warning about a collapsed segment, if necessary
		// fmt.Println("Warning: Collapsed segment detected", se)
	}
	// Looks like useless statement
	// r.ContourId = se.ContourId
	// l.ContourId = se.ContourId

	// Avoid a rounding error. The left event would be processed after the right event
	if CompareEvents(l, se.OtherEvent) == 1 {
		se.OtherEvent.Left = true
		l.Left = false
	}

	se.OtherEvent.OtherEvent = l
	se.OtherEvent = r

	queue.Push(l)
	queue.Push(r)
}
