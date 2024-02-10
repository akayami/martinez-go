package martinez_go

// Contour represents a contour with a series of points, holes, and other properties.
type Contour struct {
	Points  []Point
	HoleIds []int
	HoleOf  *Contour
	Depth   int
	Id      int
}

// NewContour creates and initializes a new Contour instance.
func NewContour(id int) *Contour {
	return &Contour{
		Points:  []Point{},
		HoleIds: []int{},
		HoleOf:  nil,
		Depth:   0,
		Id:      id,
	}
}

// IsExterior checks if the contour is an exterior contour.
func (c *Contour) IsExterior() bool {
	return c.HoleOf == nil
}
