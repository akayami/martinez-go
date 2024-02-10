package martinez_go

type SweepEvent struct {
	Point            Point
	Left             bool
	OtherEvent       *SweepEvent
	IsSubject        bool
	Type             int
	InOut            bool
	OtherInOut       bool
	PrevInResult     *SweepEvent
	ResultTransition int
	OtherPos         int
	OutputContourId  int
	IsExteriorRing   bool
	ContourId        int
}

func NewSweepEvent(point Point, left bool, otherEvent *SweepEvent, isSubject bool, edgeType int) *SweepEvent {
	return &SweepEvent{
		Point:            point,
		Left:             left,
		OtherEvent:       otherEvent,
		IsSubject:        isSubject,
		Type:             edgeType,
		InOut:            false,
		OtherInOut:       false,
		PrevInResult:     nil,
		ResultTransition: 0,
		OtherPos:         -1,
		OutputContourId:  -1,
		IsExteriorRing:   true,
		ContourId:        -1,
	}
}

func (e *SweepEvent) IsBelow(p Point) bool {
	p0 := e.Point
	p1 := e.OtherEvent.Point
	if e.Left {
		return (p0.X-p.X)*(p1.Y-p.Y)-(p1.X-p.X)*(p0.Y-p.Y) > 0
	} else {
		return (p1.X-p.X)*(p0.Y-p.Y)-(p0.X-p.X)*(p1.Y-p.Y) > 0
	}
}

func (e *SweepEvent) IsAbove(p Point) bool {
	return !e.IsBelow(p)
}

func (e *SweepEvent) IsVertical() bool {
	return e.Point.X == e.OtherEvent.Point.X
}

func (e *SweepEvent) InResult() bool {
	return e.ResultTransition != 0
}

func (e *SweepEvent) Clone() *SweepEvent {
	clone := NewSweepEvent(e.Point, e.Left, e.OtherEvent, e.IsSubject, e.Type)
	clone.PrevInResult = e.PrevInResult
	clone.IsExteriorRing = e.IsExteriorRing
	clone.InOut = e.InOut
	clone.OtherInOut = e.OtherInOut
	clone.ResultTransition = e.ResultTransition
	clone.OutputContourId = e.OutputContourId
	return clone
}
