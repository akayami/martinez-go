package martinez_go

// PossibleIntersection checks for a possible intersection between two sweep events.
func PossibleIntersection(se1, se2 *SweepEvent, queue TinyQueue) int {
	inter := SegmentIntersection(se1.Point, se1.OtherEvent.Point, se2.Point, se2.OtherEvent.Point, false)

	nintersections := len(inter)

	if nintersections == 0 {
		return 0 // No intersection
	}

	// Check if the intersection is at the endpoint of both line segments
	if nintersections == 1 && (Equals(se1.Point, se2.Point) || Equals(se1.OtherEvent.Point, se2.OtherEvent.Point)) {
		return 0
	}
	// The line segments associated to se1 and se2 intersect
	if nintersections == 2 && se1.IsSubject == se2.IsSubject {
		return 0 // Overlapping edges from the same polygon
	}

	if nintersections == 1 {
		// if the intersection point is not an endpoint of se1
		if !Equals(se1.Point, inter[0]) && !Equals(se1.OtherEvent.Point, inter[0]) {
			DivideSegment(se1, inter[0], queue)
		}
		if !Equals(se2.Point, inter[0]) && !Equals(se2.OtherEvent.Point, inter[0]) {
			DivideSegment(se2, inter[0], queue)
		}
		return 1
	}

	return processOverlap(se1, se2, queue)
}

// Example implementation of the logic.
func processOverlap(se1, se2 *SweepEvent, queue TinyQueue) int {
	var events []*SweepEvent
	leftCoincide := false
	rightCoincide := false

	if Equals(se1.Point, se2.Point) {
		leftCoincide = true // linked
	} else if CompareEvents(se1, se2) == 1 {
		events = append(events, se2, se1)
	} else {
		events = append(events, se1, se2)
	}

	if Equals(se1.OtherEvent.Point, se2.OtherEvent.Point) {
		rightCoincide = true
	} else if CompareEvents(se1.OtherEvent, se2.OtherEvent) == 1 {
		events = append(events, se2.OtherEvent, se1.OtherEvent)
	} else {
		events = append(events, se1.OtherEvent, se2.OtherEvent)
	}

	if (leftCoincide && rightCoincide) || leftCoincide {
		// both line segments are equal or share the left endpoint
		se2.Type = NonContributing
		if se2.InOut == se1.InOut {
			se1.Type = SameTransition
		} else {
			se1.Type = DifferentTransition
		}

		if leftCoincide && !rightCoincide {
			// Fixes the overlapping self-intersecting polygons issue
			DivideSegment(events[1].OtherEvent, events[0].Point, queue)
		}
		return 2
	}

	if rightCoincide {
		DivideSegment(events[0], events[1].Point, queue)
		return 3
	}

	// No line segment fully includes the other
	if events[0] != events[3].OtherEvent {
		DivideSegment(events[0], events[1].Point, queue)
		DivideSegment(events[1], events[2].Point, queue)
		return 3
	}

	// One line segment fully includes the other
	DivideSegment(events[0], events[1].Point, queue)
	DivideSegment(events[3].OtherEvent, events[2].Point, queue)

	return 3
}
