package heatmap

import (
	"encoding/json"

	"github.com/jakecorrenti/fishfindr/types"
)

type Point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func PointsJSON(locations *[]types.Location) ([]byte, error) {
	points := make([]Point, len(*locations))

	for _, l := range *locations {
		points = append(points, Point{Latitude: l.Latitude, Longitude: l.Latitude})
	}

	return json.Marshal(points)
}
