package martinez_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestScenarios(t *testing.T) {
	baseDir := "test/featureTypes" // Adjust the base directory to match your file structure
	clipping, err := loadGeoJSON(filepath.Join(baseDir, "clippingPoly.geojson"))
	assert.Nil(t, err)
	outDir := filepath.Join(baseDir, "out")

	tests := []struct {
		testName    string
		subjectPoly string
	}{
		{"polyToClipping", "poly"},
		{"polyWithHoleToClipping", "polyWithHole"},
		{"multiPolyToClipping", "multiPoly"},
		{"multiPolyWithHoleToClipping", "multiPolyWithHole"},
	}

	ops := []struct {
		dirname string
		op      int
	}{
		{dirname: "intersection", op: Intersection},
		{dirname: "xor", op: XOR},
		{dirname: "difference", op: Difference},
		{dirname: "union", op: Union},
	}

	for _, ts := range tests {
		t.Run(ts.testName, func(t *testing.T) {
			subject, _ := loadGeoJSON(filepath.Join(baseDir, ts.subjectPoly+".geojson"))
			for _, op := range ops {
				t.Run(fmt.Sprintf("%s -> %s", ts.testName, op.dirname), func(t *testing.T) {
					expectedIntResult, err := loadGeoJSON(filepath.Join(outDir, op.dirname, ts.testName+".geojson"))
					assert.Nil(t, err)
					intResult := Compute(subject, clipping, op.op)
					assert.Equal(t, expectedIntResult, intResult)
				})
			}
		})
	}
}

type GeoJSON struct {
	Type     string   `json:"type"`
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Type        string            `json:"type"`
	Coordinates [][][]interface{} `json:"coordinates"`
}

func loadGeoJSON(path string) ([][][]Point, error) {
	var multipolygon [][][]Point
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var geojson GeoJSON
	err = json.Unmarshal(bytes, &geojson)
	if err != nil {
		return nil, err
	}
	multipolygon = [][][]Point{}
	if geojson.Geometry.Type == "Polygon" {
		poly := [][]Point{}

		for _, rings := range geojson.Geometry.Coordinates {
			ring := []Point{}
			for _, gring := range rings {
				p := Point{}
				for x, coordinate := range gring {
					switch v := coordinate.(type) {
					case interface{}:
						if vf, ok := v.(float64); ok {
							if x == 0 {
								p.X = vf
							} else if x == 1 {
								p.Y = vf
							} else {
								return nil, errors.New("invalid coordinate length")
							}
						} else {
							return nil, errors.New("unable to convert point coordinate to float")
						}
					default:
						return nil, errors.New("invalid value type")
					}
				}
				ring = append(ring, p)
			}
			poly = append(poly, ring)
		}
		multipolygon = append(multipolygon, poly)
	} else if geojson.Geometry.Type == "MultiPolygon" {
		for i, poly := range geojson.Geometry.Coordinates {
			multipolygon = append(multipolygon, [][]Point{})
			for z, rings := range poly {
				multipolygon[i] = append(multipolygon[i], []Point{})
				for _, ring := range rings {
					switch v := ring.(type) {
					case []interface{}:
						p := Point{}
						if vf, ok := v[0].(float64); ok {
							p.X = vf
						} else {
							return nil, errors.New("unable to convert point coordinate to float")
						}
						if vf, ok := v[1].(float64); ok {
							p.Y = vf
						} else {
							return nil, errors.New("unable to convert point coordinate to float")
						}
						multipolygon[i][z] = append(multipolygon[i][z], p)
					default:
						return nil, errors.New("invalid value type")
					}
				}
			}
		}
	}
	return multipolygon, nil
}

func TestMPolyLoading(t *testing.T) {
	path := "test/featureTypes/multiPoly.geojson"
	mpolygon, err := loadGeoJSON(path)
	assert.Nil(t, err)
	assert.NotNil(t, mpolygon)
	assert.Equal(t, 2, len(mpolygon))
	assert.Equal(t, 1, len(mpolygon[0]))
	assert.Equal(t, 1, len(mpolygon[1]))
	assert.Equal(t, 5, len(mpolygon[0][0]))
	assert.Equal(t, 5, len(mpolygon[1][0]))

}

func TestPolyLoading(t *testing.T) {
	path := "test/featureTypes/poly.geojson"
	mpolygon, err := loadGeoJSON(path)
	assert.Nil(t, err)
	assert.NotNil(t, mpolygon)
	assert.Equal(t, 1, len(mpolygon))
	assert.Equal(t, 5, len(mpolygon[0][0]))
}
