package martinez_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntryPointArea(t *testing.T) {

	subject := [][][]Point{{{
		Point{0, 0},
		Point{0, 10},
		Point{10, 10},
		Point{10, 0},
		Point{0, 0},
	}}}

	clipping := [][][]Point{{{
		Point{5, 5},
		Point{5, 15},
		Point{15, 15},
		Point{15, 5},
		Point{5, 5},
	}}}

	tests := []struct {
		op     int
		result [][][]Point
	}{
		{Union, [][][]Point{{{
			Point{0, 0},
			Point{10, 0},
			Point{10, 5},
			Point{15, 5},
			Point{15, 15},
			Point{5, 15},
			Point{5, 10},
			Point{0, 10},
			Point{0, 0},
		}}}},
		{Intersection, [][][]Point{{{
			Point{5, 5},
			Point{10, 5},
			Point{10, 10},
			Point{5, 10},
			Point{5, 5},
		}}}},
		{Difference, [][][]Point{{{
			Point{0, 0},
			Point{10, 0},
			Point{10, 5},
			Point{5, 5},
			Point{5, 10},
			Point{0, 10},
			Point{0, 0},
		}}}},
		{XOR, [][][]Point{{{
			Point{0, 0},
			Point{10, 0},
			Point{10, 5},
			Point{5, 5},
			Point{5, 10},
			Point{0, 10},
			Point{0, 0},
		}, {
			Point{5, 10},
			Point{10, 10},
			Point{10, 5},
			Point{15, 5},
			Point{15, 15},
			Point{5, 15},
			Point{5, 10},
		}}}},
	}

	for _, test := range tests {
		res := Boolean(subject, clipping, test.op)
		assert.Equal(t, test.result, res)
	}
}
