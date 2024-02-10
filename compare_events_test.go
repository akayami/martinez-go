package martinez_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueue(t *testing.T) {

	t.Run("first group", func(t *testing.T) {
		t.Run("test 1", func(t *testing.T) {
			queue := NewTinyQueueDefault(nil, CompareEvents)

			// Events with different X and Y values
			e1 := &SweepEvent{Point: Point{X: 1.0, Y: 1.0}}
			e2 := &SweepEvent{Point: Point{X: 2.0, Y: 2.0}}

			// Events with the same X but different Y values
			e3 := &SweepEvent{Point: Point{X: 2.0, Y: 1.5}}

			// Add events to the queue
			queue.Push(e3) // Adding out of order to test the sorting
			queue.Push(e1)
			queue.Push(e2)

			// Pop events and test the order based on CompareEvents logic
			if pop1 := queue.Pop(); pop1 != e1 {
				t.Errorf("Expected first pop to be e1, got %+v", pop1)
			}
			if pop2 := queue.Pop(); pop2 != e3 {
				t.Errorf("Expected second pop to be e3 (same X as e2 but lower Y), got %+v", pop2)
			}
			if pop3 := queue.Pop(); pop3 != e2 {
				t.Errorf("Expected third pop to be e2, got %+v", pop3)
			}
		})

		t.Run("queue should process least (by x) sweep event first", func(t *testing.T) {
			queue := NewTinyQueueDefault(nil, CompareEvents)
			e1 := &SweepEvent{Point: Point{0.0, 0.0}}
			e2 := &SweepEvent{Point: Point{0.5, 0.5}}

			queue.Push(e1)
			queue.Push(e2)

			if pop1 := queue.Pop(); pop1 != e1 {
				t.Errorf("Expected first pop to be e1, got %+v", pop1)
			}

			if pop2 := queue.Pop(); pop2 != e2 {
				t.Errorf("Expected second pop to be e2, got %+v", pop2)
			}
		})

		t.Run("queue should process lest(by y) sweep event first", func(t *testing.T) {
			queue := NewTinyQueueDefault(nil, CompareEvents)
			e1 := &SweepEvent{Point: Point{0.0, 0.0}}
			e2 := &SweepEvent{Point: Point{0.0, 0.5}}

			queue.Push(e1)
			queue.Push(e2)

			if pop1 := queue.Pop(); pop1 != e1 {
				t.Errorf("Expected first pop to be e1 based on y coordinate, got %+v", pop1)
			}

			if pop2 := queue.Pop(); pop2 != e2 {
				t.Errorf("Expected second pop to be e2 based on y coordinate, got %+v", pop2)
			}
		})

		t.Run("queue should pop least(by left prop) sweep event first", func(t *testing.T) {
			assert := assert.New(t)
			queue := NewTinyQueueDefault(nil, CompareEvents)

			e1 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: true}
			e2 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: false}

			queue.Push(e1)
			queue.Push(e2)

			assert.Equal(e2, queue.Pop(), "Expected e2 to be popped first due to its left property being false")
			assert.Equal(e1, queue.Pop(), "Expected e1 to be popped second due to its left property being true")
		})
	})
	t.Run("sweep event comparison x coordinates", func(t *testing.T) {
		e1 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}}
		e2 := &SweepEvent{Point: Point{X: 0.5, Y: 0.5}}

		if result := CompareEvents(e1, e2); result != -1 {
			t.Errorf("Expected comparison of e1 with e2 to be -1, got %d", result)
		}

		if result := CompareEvents(e2, e1); result != 1 {
			t.Errorf("Expected comparison of e2 with e1 to be 1, got %d", result)
		}
	})

	t.Run("sweep event comparison y coordinates", func(t *testing.T) {
		e1 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}}
		e2 := &SweepEvent{Point: Point{X: 0.0, Y: 0.5}}

		// Expect e1 to be considered "less than" e2 because its Y value is smaller
		if result := CompareEvents(e1, e2); result != -1 {
			t.Errorf("Expected comparison of e1 with e2 (by Y) to be -1, got %d", result)
		}

		// Expect e2 to be considered "greater than" e1 because its Y value is larger
		if result := CompareEvents(e2, e1); result != 1 {
			t.Errorf("Expected comparison of e2 with e1 (by Y) to be 1, got %d", result)
		}
	})

	t.Run("sweep event comparison not left first", func(t *testing.T) {
		e1 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: true}
		e2 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: false}

		assert.Equal(t, 1, CompareEvents(e1, e2))
		assert.Equal(t, -1, CompareEvents(e2, e1))
	})

	t.Run("sweep event comparison shared start point not collinear edges", func(t *testing.T) {
		e1 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: true, OtherEvent: &SweepEvent{Point: Point{1, 1}, Left: false}}
		e2 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: true, OtherEvent: &SweepEvent{Point: Point{2, 3}, Left: false}}

		assert.Equal(t, -1, CompareEvents(e1, e2))
		assert.Equal(t, 1, CompareEvents(e2, e1))
	})

	t.Run("sweep event comparison collinear edges", func(t *testing.T) {
		e1 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: true, OtherEvent: &SweepEvent{Point: Point{1, 1}, Left: true}, IsSubject: true}
		e2 := &SweepEvent{Point: Point{X: 0.0, Y: 0.0}, Left: true, OtherEvent: &SweepEvent{Point: Point{2, 2}, Left: false}, IsSubject: false}

		assert.Equal(t, -1, CompareEvents(e1, e2))
		assert.Equal(t, 1, CompareEvents(e2, e1))
	})
}
