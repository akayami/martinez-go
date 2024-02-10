package martinez_go

import "github.com/paulmach/orb"

func OrbPointToPoint(point orb.Point) Point {
	return Point{X: point.X(), Y: point.Y()}
}

func OrbPolygonToPolygon(polygon orb.Polygon) [][]Point {
	var result [][]Point

	for _, ring := range polygon {
		var convertedRings []Point
		for _, orbPoint := range ring {
			convertedRings = append(convertedRings, Point{X: orbPoint[0], Y: orbPoint[1]})
		}
		result = append(result, convertedRings)
	}
	return result
}
