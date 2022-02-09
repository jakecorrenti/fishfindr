package analysis

import (
	"github.com/jakecorrenti/fishfindr/types"
	cluster "github.com/smira/go-point-clustering"
)

func GetClusters(locations *[]types.Location) ([]cluster.Cluster, []int, error) {

	var points []cluster.Point
	for _, l := range *locations {
		points = append(points, cluster.Point{l.Latitude, l.Longitude})
	}

	var pointsList cluster.PointList
	for _, point := range points {
		pointsList = append(pointsList, point)
	}

	// epsilon is 0.5 miles in km, 3 points minimum in epsilon-neighborhood
	clusters, noise := cluster.DBScan(pointsList, 0.804672, 3)

	return clusters, noise, nil
}
