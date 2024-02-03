package martinez_go

func OrderEvents(sortedEvents []*SweepEvent) []*SweepEvent {
	var resultEvents []*SweepEvent
	for _, event := range sortedEvents {
		if (event.Left && event.InResult()) || (!event.Left && event.OtherEvent.InResult()) {
			resultEvents = append(resultEvents, event)
		}
	}

	// Sort resultEvents if necessary
	sorted := false
	for !sorted {
		sorted = true
		for i := 0; i < len(resultEvents)-1; i++ {
			if CompareEvents(resultEvents[i], resultEvents[i+1]) == 1 {
				resultEvents[i], resultEvents[i+1] = resultEvents[i+1], resultEvents[i]
				sorted = false
			}
		}
	}

	// Update otherPos
	for i, event := range resultEvents {
		event.OtherPos = i
		if !event.Left {
			event.OtherPos, event.OtherEvent.OtherPos = event.OtherEvent.OtherPos, event.OtherPos
		}
	}

	return resultEvents
}

func NextPos(pos int, resultEvents []*SweepEvent, processed map[int]bool, origPos int) int {
	newPos := pos + 1
	length := len(resultEvents)

	for newPos < length && resultEvents[newPos].Point.Equals(resultEvents[pos].Point) {
		if !processed[newPos] {
			return newPos
		}
		newPos++
	}

	newPos = pos - 1
	for newPos > origPos && processed[newPos] {
		newPos--
	}

	return newPos
}

// 	}
//
// 	return newPos
// }

// InitializeContourFromContext initializes a contour from the sweep event context.
func InitializeContourFromContext(event *SweepEvent, contours []*Contour, contourId int) *Contour {
	contour := NewContour(contourId)
	if event.PrevInResult != nil {
		lowerContourId := event.PrevInResult.OutputContourId
		lowerContour := contours[lowerContourId]
		if lowerContour.HoleOf != nil {
			parentContourId := lowerContour.HoleOf.Id
			contours[parentContourId].HoleIds = append(contours[parentContourId].HoleIds, contourId)
			contour.HoleOf = contours[parentContourId]
			contour.Depth = lowerContour.Depth
		} else {
			lowerContour.HoleIds = append(lowerContour.HoleIds, contourId)
			contour.HoleOf = lowerContour
			contour.Depth = lowerContour.Depth + 1
		}
	} else {
		contour.HoleOf = nil
		contour.Depth = 0
	}
	return contour
}

// ConnectEdges connects edges to form contours.
func ConnectEdges(sortedEvents []*SweepEvent) []*Contour {
	resultEvents := OrderEvents(sortedEvents)
	processed := make(map[int]bool)
	var contours []*Contour

	for i := 0; i < len(resultEvents); i++ {
		if processed[i] {
			continue
		}

		contourId := len(contours)
		contour := InitializeContourFromContext(resultEvents[i], contours, contourId)

		markAsProcessed := func(pos int) {
			processed[pos] = true
			if pos < len(resultEvents) && resultEvents[pos] != nil {
				resultEvents[pos].OutputContourId = contourId
			}
		}

		pos := i
		origPos := i
		contour.Points = append(contour.Points, resultEvents[i].Point)

		for {
			markAsProcessed(pos)

			pos = resultEvents[pos].OtherPos

			markAsProcessed(pos)
			contour.Points = append(contour.Points, resultEvents[pos].Point)

			pos = NextPos(pos, resultEvents, processed, origPos)

			if pos == origPos || pos >= len(resultEvents) || resultEvents[pos] == nil {
				break
			}
		}

		contours = append(contours, contour)
	}

	return contours
}
