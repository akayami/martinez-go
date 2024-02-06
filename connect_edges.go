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

func InitializeContourFromContext(event *SweepEvent, contours []*Contour, contourId int) *Contour {
	contour := &Contour{}

	if event.PrevInResult != nil {
		prevInResult := event.PrevInResult
		lowerContourId := prevInResult.OutputContourId
		lowerResultTransition := prevInResult.ResultTransition

		if lowerResultTransition > 0 {
			lowerContour := contours[lowerContourId]
			if lowerContour.HoleOf != nil {
				// The lower contour is a hole => Connect the new contour as a hole to its parent,
				// and use the same depth.
				parentContourId := lowerContour.HoleOf.Id
				contours[parentContourId].HoleIds = append(contours[parentContourId].HoleIds, contourId)
				contour.HoleOf = contours[parentContourId]
				contour.Depth = lowerContour.Depth
			} else {
				// The lower contour is an exterior contour => Connect the new contour as a hole,
				// and increment depth.
				contours[lowerContourId].HoleIds = append(contours[lowerContourId].HoleIds, contourId)
				contour.HoleOf = lowerContour
				contour.Depth = lowerContour.Depth + 1
			}
		} else {
			// We are outside => this contour is an exterior contour of same depth.
			contour.HoleOf = nil
			contour.Depth = contours[lowerContourId].Depth
		}
	} else {
		// There is no lower/previous contour => this contour is an exterior contour of depth 0.
		contour.HoleOf = nil
		contour.Depth = 0
	}

	contour.Id = contourId
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
