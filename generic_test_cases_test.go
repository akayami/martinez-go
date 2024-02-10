package martinez_go

import (
	"errors"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneric(t *testing.T) {

	type Feature struct {
		Geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Operation string `json:"operation"`
		} `json:"properties"`
	}

	type GeoJSON struct {
		Features []Feature `json:"features"`
	}

	onlyRun := []string{}
	skipRun := []string{}

	t.Run("Generic cases", func(t *testing.T) {

		onlyRun = []string{
			// "hourglasses.geojson",
			// "issue76.geojson",
			// "touching_boxes.geojson",
		}

		skipRun = []string{
			// "issue76.geojson",
		}

		caseDir := "./test/genericTestCases" // Adjust path as needed
		files, err := filepath.Glob(filepath.Join(caseDir, "*.geojson"))

		fmt.Println(files)

		if err != nil {
			t.Fatalf("Failed to list test case files: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("No test cases found, this must not happen")
		}

		for _, file := range files {
			if len(onlyRun) > 0 && !match(onlyRun, file) {
				continue
			}

			if len(skipRun) > 0 && match(skipRun, file) {
				continue
			}

			data, err := load(file)
			assert.Nil(t, err)

			p1Geometry := TranslateMultiPolygon(data[0].Polygon)
			p2Geometry := TranslateMultiPolygon(data[1].Polygon)
			// p2Geometry := geojson.Features[1].Geometry.(orb.Polygon)

			// if geojson.Features[0].Geometry.Bound().

			expectedResults := data[2:]

			t.Run(fmt.Sprintf("Case: %s", file), func(t *testing.T) {
				for _, result := range expectedResults {
					if op, ok := result.Properties["operation"]; ok {
						t.Run(fmt.Sprintf("Operation: %s", op), func(t *testing.T) {
							realOp := getOp(op)
							var res [][][]Point
							if realOp >= 0 {
								res = Compute(p1Geometry, p2Geometry, realOp)
							} else if realOp == -1 {
								res = Compute(p2Geometry, p1Geometry, Difference)
							}
							expected := TranslateMultiPolygon(result.Polygon)
							assert.Equal(t, expected, res)
						})
					} else {
						fmt.Println("Missing operation value")
					}
				}
			})
		}
	})
}

func match(list []string, name string) bool {
	for _, needle := range list {
		if strings.Contains(strings.ToLower(name), strings.ToLower(needle)) {
			return true
		}
	}
	return false
}

// func extractExpectedResults([]orb.MultiPolygon)

type TestFeature struct {
	Polygon    orb.MultiPolygon
	Properties map[string]string
}

func TranslateMultiPolygon(multiPoly orb.MultiPolygon) [][][]Point {
	var result [][][]Point // Initialize the result slice

	for _, poly := range multiPoly { // Iterate over each polygon in the MultiPolygon
		var r [][]Point
		for _, ring := range poly { // Iterate over each ring in the polygon
			var ringPoints []Point   // To store converted points of the current ring
			for _, p := range ring { // Iterate over each point in the ring
				// Convert orb.Point (longitude, latitude) to Point (Lat, Lon)
				ringPoints = append(ringPoints, Point{X: p[0], Y: p[1]})
			}
			r = append(r, ringPoints) // Add the converted ring to the result
		}
		result = append(result, r) // Add the converted ring to the result
	}

	return result
}

func getOp(op string) int {
	switch op {
	case "intersection":
		return Intersection
	case "union":
		return Union
	case "diff":
		return Difference
	case "xor":
		return XOR
	case "diff_ba":
		return -1
	default:
		return 0
	}
}

func load(file string) ([]TestFeature, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fc, err := geojson.UnmarshalFeatureCollection(content)
	if err != nil {
		return nil, err
	}
	var output []TestFeature
	for _, feature := range fc.Features {
		f := TestFeature{}
		geometry := feature.Geometry
		properties := feature.Properties
		f.Polygon, err = handleGeometry(geometry)
		if err != nil {
			return nil, err
		}
		f.Properties, err = handleProperties(properties)
		output = append(output, f)
	}
	return output, nil
}

func handleProperties(properties geojson.Properties) (map[string]string, error) {
	output := map[string]string{}
	for key, val := range properties {
		output[key] = val.(string)
	}
	return output, nil
}

func handleGeometry(geometry orb.Geometry) (orb.MultiPolygon, error) {
	switch geometry.GeoJSONType() {
	case "Polygon":
		p, ok := geometry.(orb.Polygon)
		if !ok {
			return nil, errors.New("failed to parse polygon")
		}
		return orb.MultiPolygon{p}, nil
	case "MultiPolygon":
		p, ok := geometry.(orb.MultiPolygon)
		if !ok {
			return nil, errors.New("failed to parse multipolygon")
		}
		return p, nil
	default:
		return nil, errors.New("invalid multipolygon")
	}
}

// func load2(file string) (*[]P, error) {
// 	content, err := os.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fc, err := geojson.UnmarshalFeatureCollection(content)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res []P
// 	var prop geojson.Properties
// 	var orbP orb.MultiPolygon
// 	for _, feature := range fc.Features {
// 		poly := P{}
// 		geometry := feature.Geometry
// 		// properties := feature.Properties
// 		switch geometry.GeoJSONType() {
// 		case "Polygon":
// 			p, ok := geometry.(orb.Polygon)
// 			if !ok {
// 				return nil, errors.New("failed to parse polygon")
// 			}
// 			orbP = orb.MultiPolygon{p}
// 		case "MultiPolygon":
// 			p, ok := geometry.(orb.MultiPolygon)
// 			if !ok {
// 				return nil, errors.New("failed to parse multipolygon")
// 			}
// 			orbP = p
// 		default:
// 			return nil, errors.New("invalid multipolygon")
// 		}
// 		if feature.Properties != nil {
// 			poly.Properties = feature.Properties
// 		}
// 		// poly.Polygons
// 	}
// 	// return &P{Polygons: res, op: prop}, nil
// }

// func checkType(i interface{}) (*geom.MultiPolygon, error) {
// 	switch i.(type) {
// 	case *geom.MultiPolygon:
// 		p := (i).(*geom.MultiPolygon)
// 		return p, nil
// 	case *geom.Polygon:
// 		p := geom.NewMultiPolygon(geom.XY)
// 		err := p.Push((i).(*geom.Polygon))
// 		if err != nil {
// 			return nil, err
// 		}
// 		return p, nil
// 	default:
// 		return nil, errors.New("unknown interface")
// 	}
// }
