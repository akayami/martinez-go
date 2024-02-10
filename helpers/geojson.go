package helpers

import (
	"github.com/paulmach/orb/geojson"
	"os"
)

func LoadGeoJson(path string) (*geojson.FeatureCollection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}
	return fc, nil
}
