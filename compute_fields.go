package martinez_go

// Operation types
const (
	Intersection = iota
	Union
	Difference
	XOR
)

// Edge types
const (
	Normal = iota
	SameTransition
	DifferentTransition
	NonContributing
)

// ComputeFields computes various fields for a SweepEvent.
func ComputeFields(event, prev *SweepEvent, operation int) {
	if prev == nil {
		event.InOut = false
		event.OtherInOut = true
	} else {
		if event.IsSubject == prev.IsSubject {
			event.InOut = !prev.InOut
			event.OtherInOut = prev.OtherInOut
		} else {
			// previous line segment in sweepline belongs to the clipping polygon
			event.InOut = !prev.OtherInOut
			if prev.IsVertical() {
				event.OtherInOut = !prev.InOut
			} else {
				event.OtherInOut = prev.InOut
			}
			// event.OtherInOut = !prev.IsVertical() && prev.InOut
		}

		if prev != nil {
			if !InResult(prev, operation) || prev.IsVertical() {
				event.PrevInResult = prev.PrevInResult
			} else {
				event.PrevInResult = prev
			}
		}
	}
	inResult := InResult(event, operation)
	if inResult {
		event.ResultTransition = DetermineResultTransition(event, operation)
	} else {
		event.ResultTransition = 0
	}
}

// InResult determines if a SweepEvent is in the result of the Boolean operation.
func InResult(event *SweepEvent, operation int) bool {
	switch event.Type {
	case Normal:
		switch operation {
		case Intersection:
			return !event.OtherInOut
		case Union:
			return event.OtherInOut
		case Difference:
			return (event.IsSubject && event.OtherInOut) || (!event.IsSubject && !event.OtherInOut)
		case XOR:
			return true
		}
	case SameTransition:
		return operation == Intersection || operation == Union
	case DifferentTransition:
		return operation == Difference
	case NonContributing:
		return false
	}
	return false
}

// DetermineResultTransition determines the result transition for a SweepEvent.
func DetermineResultTransition(event *SweepEvent, operation int) int {
	thisIn := !event.InOut
	thatIn := !event.OtherInOut

	var isIn bool
	switch operation {
	case Intersection:
		isIn = thisIn && thatIn
	case Union:
		isIn = thisIn || thatIn
	case XOR:
		isIn = thisIn != thatIn
	case Difference:
		if event.IsSubject {
			isIn = thisIn && !thatIn
		} else {
			isIn = thatIn && !thisIn
		}
	}
	if isIn {
		return 1
	}
	return -1
}
