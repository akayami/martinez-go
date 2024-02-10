package martinez_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompareSegments(t *testing.T) {

	t.Run("not collinear", func(t *testing.T) {

		t.Run("shared left point", func(t *testing.T) {
			tree := NewSplayTree(CompareSegments)
			pt := Point{X: 0.0, Y: 0.0}
			se1 := NewSweepEvent(pt, true, NewSweepEvent(Point{X: 1, Y: 1}, false, nil, false, 0), false, 0)
			se2 := NewSweepEvent(pt, true, NewSweepEvent(Point{X: 2, Y: 3}, false, nil, false, 0), false, 0)

			tree.Insert(NewSplayTreeNode(se1))
			tree.Insert(NewSplayTreeNode(se2))

			assert.Equal(t, Point{X: 2, Y: 3}, tree.MaxNode(tree.Root).Key.OtherEvent.Point)
			assert.Equal(t, Point{X: 1, Y: 1}, tree.MinNode(tree.Root).Key.OtherEvent.Point)
		})

		t.Run("different left point", func(t *testing.T) {
			tree := NewSplayTree(CompareSegments)
			se1 := NewSweepEvent(Point{X: 0.0, Y: 1.0}, true, NewSweepEvent(Point{X: 1, Y: 1}, false, nil, false, 0), true, 0)
			se2 := NewSweepEvent(Point{X: 0.0, Y: 2.0}, true, NewSweepEvent(Point{X: 2, Y: 3}, false, nil, false, 0), true, 0)

			tree.Insert(NewSplayTreeNode(se1))
			tree.Insert(NewSplayTreeNode(se2))

			assert.Equal(t, Point{X: 1, Y: 1}, tree.MinNode(tree.Root).Key.OtherEvent.Point)
			assert.Equal(t, Point{X: 2, Y: 3}, tree.MaxNode(tree.Root).Key.OtherEvent.Point)
		})

		t.Run("events order in sweep line", func(t *testing.T) {
			se1 := NewSweepEvent(Point{0, 1}, true, NewSweepEvent(Point{2, 1}, false, nil, false, 0), false, 0)
			se2 := NewSweepEvent(Point{-1, 0}, true, NewSweepEvent(Point{2, 3}, false, nil, false, 0), false, 0)
			se3 := NewSweepEvent(Point{0, 1}, true, NewSweepEvent(Point{3, 4}, false, nil, false, 0), false, 0)
			se4 := NewSweepEvent(Point{-1, 0}, true, NewSweepEvent(Point{3, 1}, false, nil, false, 0), false, 0)

			assert.Equal(t, 1, CompareEvents(se1, se2))
			assert.False(t, se2.IsBelow(se1.Point))
			assert.True(t, se2.IsAbove(se1.Point))

			assert.Equal(t, -1, CompareSegments(se1, se2), "compare segments")
			assert.Equal(t, 1, CompareSegments(se2, se1), "compare segments inverted")

			assert.Equal(t, 1, CompareEvents(se3, se4))
			assert.False(t, se4.IsAbove(se3.Point))
		})

		t.Run("first point is below", func(t *testing.T) {
			se2 := NewSweepEvent(Point{0, 1}, true, NewSweepEvent(Point{2, 1}, false, nil, false, 0), false, 0)
			se1 := NewSweepEvent(Point{-1, 0}, true, NewSweepEvent(Point{2, 3}, false, nil, false, 0), false, 0)

			assert.False(t, se1.IsBelow(se2.Point))

			assert.Equal(t, 1, CompareSegments(se1, se2), "compare segments")
		})

	})

	t.Run("collinear segments", func(t *testing.T) {
		se1 := NewSweepEvent(Point{1, 1}, true, NewSweepEvent(Point{5, 1}, false, nil, false, 0), true, 0)
		se2 := NewSweepEvent(Point{2, 1}, true, NewSweepEvent(Point{3, 1}, false, nil, false, 0), false, 0)
		assert.NotEqual(t, se1.IsSubject, se2.IsSubject)
		assert.Equal(t, -1, CompareSegments(se1, se2))
	})

	t.Run("collinear shared left point", func(t *testing.T) {
		pt := Point{0, 1}
		se1 := NewSweepEvent(pt, true, NewSweepEvent(Point{5, 1}, false, nil, false, 0), false, 0)
		se2 := NewSweepEvent(pt, true, NewSweepEvent(Point{3, 1}, false, nil, false, 0), false, 0)

		se1.ContourId = 1
		se2.ContourId = 2

		assert.Equal(t, se1.IsSubject, se2.IsSubject)
		assert.Equal(t, se1.Point, se2.Point)

		assert.Equal(t, -1, CompareSegments(se1, se2))

		se1.ContourId = 2
		se2.ContourId = 1

		assert.Equal(t, 1, CompareSegments(se1, se2))
	})
	t.Run("collinear same polygon different left points", func(t *testing.T) {
		se1 := NewSweepEvent(Point{1, 1}, true, NewSweepEvent(Point{5, 1}, false, nil, false, 0), true, 0)
		se2 := NewSweepEvent(Point{2, 1}, true, NewSweepEvent(Point{3, 1}, false, nil, false, 0), true, 0)

		assert.Equal(t, se1.IsSubject, se2.IsSubject)
		assert.NotEqual(t, se1.Point, se2.Point)
		assert.Equal(t, -1, CompareSegments(se1, se2))
		assert.Equal(t, 1, CompareSegments(se2, se1))
	})
}
