package example

import (
	"fmt"
	"github.com/akayami/martinez-go"
)

func main() {

	type Point = martinez_go.Point

	subject := [][][]Point{{{
		Point{0, 0},
		Point{0, 10},
		Point{10, 10},
		Point{10, 0},
		Point{0, 0},
	}}}

	clipping := [][][]martinez_go.Point{{{
		Point{5, 5},
		Point{5, 15},
		Point{15, 15},
		Point{15, 5},
		Point{5, 5},
	}}}

	res := martinez_go.Compute(subject, clipping, martinez_go.Intersection)

	fmt.Println(res)

}
