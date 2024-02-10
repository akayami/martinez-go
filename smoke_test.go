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
		name   string
		op     int
		result [][][]Point
	}{
		{"Union:", Union, [][][]Point{{{
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
		{"Intersection", Intersection, [][][]Point{{{
			Point{5, 5},
			Point{10, 5},
			Point{10, 10},
			Point{5, 10},
			Point{5, 5},
		}}}},
		{"Diff", Difference, [][][]Point{{{
			Point{0, 0},
			Point{10, 0},
			Point{10, 5},
			Point{5, 5},
			Point{5, 10},
			Point{0, 10},
			Point{0, 0},
		}}}},
		{"XOR", XOR, [][][]Point{{{
			Point{0, 0},
			Point{10, 0},
			Point{10, 5},
			Point{5, 5},
			Point{5, 10},
			Point{0, 10},
			Point{0, 0},
		}}, {{
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
		t.Run(test.name, func(t *testing.T) {
			res := Compute(subject, clipping, test.op)
			assert.Equal(t, test.result, res)
		})
	}
}
