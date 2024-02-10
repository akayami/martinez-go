package martinez_go

import "testing"

func TestSweepEventIsBelow(t *testing.T) {
	s1 := NewSweepEvent(Point{0, 0}, true, NewSweepEvent(Point{1, 1}, false, nil, false, 0), false, 0)
	s2 := NewSweepEvent(Point{0, 1}, false, NewSweepEvent(Point{0, 0}, false, nil, false, 0), false, 0)

	if !s1.IsBelow(Point{0, 1}) {
		t.Errorf("s1.IsBelow(Point{0, 1}) = false; want true")
	}
	if !s1.IsBelow(Point{1, 2}) {
		t.Errorf("s1.IsBelow(Point{1, 2}) = false; want true")
	}
	if s1.IsBelow(Point{0, 0}) {
		t.Errorf("s1.IsBelow(Point{0, 0}) = true; want false")
	}
	if s1.IsBelow(Point{5, -1}) {
		t.Errorf("s1.IsBelow(Point{5, -1}) = true; want false")
	}

	if s2.IsBelow(Point{0, 1}) {
		t.Errorf("s2.IsBelow(Point{0, 1}) = true; want false")
	}
	if s2.IsBelow(Point{1, 2}) {
		t.Errorf("s2.IsBelow(Point{1, 2}) = true; want false")
	}
	if s2.IsBelow(Point{0, 0}) {
		t.Errorf("s2.IsBelow(Point{0, 0}) = true; want false")
	}
	if s2.IsBelow(Point{5, -1}) {
		t.Errorf("s2.IsBelow(Point{5, -1}) = true; want false")
	}
}

func TestSweepEventIsAbove(t *testing.T) {
	s1 := NewSweepEvent(Point{0, 0}, true, NewSweepEvent(Point{1, 1}, false, nil, false, 0), false, 0)
	s2 := NewSweepEvent(Point{0, 1}, false, NewSweepEvent(Point{0, 0}, false, nil, false, 0), false, 0)

	if s1.IsAbove(Point{0, 1}) {
		t.Errorf("s1.IsAbove(Point{0, 1}) = true; want false")
	}
	if s1.IsAbove(Point{1, 2}) {
		t.Errorf("s1.IsAbove(Point{1, 2}) = true; want false")
	}
	if !s1.IsAbove(Point{0, 0}) {
		t.Errorf("s1.IsAbove(Point{0, 0}) = false; want true")
	}
	if !s1.IsAbove(Point{5, -1}) {
		t.Errorf("s1.IsAbove(Point{5, -1}) = false; want true")
	}

	if !s2.IsAbove(Point{0, 1}) {
		t.Errorf("s2.IsAbove(Point{0, 1}) = false; want true")
	}
	if !s2.IsAbove(Point{1, 2}) {
		t.Errorf("s2.IsAbove(Point{1, 2}) = false; want true")
	}
	if !s2.IsAbove(Point{0, 0}) {
		t.Errorf("s2.IsAbove(Point{0, 0}) = false; want true")
	}
	if !s2.IsAbove(Point{5, -1}) {
		t.Errorf("s2.IsAbove(Point{5, -1}) = false; want true")
	}
}

func TestSweepEventIsVertical(t *testing.T) {
	if !NewSweepEvent(Point{0, 0}, true, NewSweepEvent(Point{0, 1}, false, nil, false, 0), false, 0).IsVertical() {
		t.Errorf("IsVertical() = false; want true for vertical line")
	}
	if NewSweepEvent(Point{0, 0}, true, NewSweepEvent(Point{0.0001, 1}, false, nil, false, 0), false, 0).IsVertical() {
		t.Errorf("IsVertical() = true; want false for non-vertical line")
	}
}
