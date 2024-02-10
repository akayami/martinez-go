# martinez-go
Martinez-Rueda polygon clipping algorithm for golang.

Follows https://github.com/w8r/martinez and https://github.com/21re/rust-geo-booleanop as reference. 

# Usage
```go
package main

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

	clipping := [][][]Point{{{
		Point{5, 5},
		Point{5, 15},
		Point{15, 15},
		Point{15, 5},
		Point{5, 5},
	}}}

	res := martinez_go.Compute(subject, clipping, martinez_go.Intersection)

	fmt.Println(res)

}
```

```shell
go run ./main.go 
[[[{5 5} {10 5} {10 10} {5 10} {5 5}]]]
```

# Authors
* [Tomasz Rakowski](https://github.com/akayami)


### Based on
* [A Javascript implementation](https://github.com/w8r/martinez/)
* [A new algorithm for computing Boolean operations on polygons](http://www.sciencedirect.com/science/article/pii/S0965997813000379) (2008, 2013) by Francisco Martinez, Antonio Jesus Rueda, Francisco Ramon Feito (and its C++ code)
