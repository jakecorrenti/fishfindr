package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Arafatk/glot"
	"github.com/gorilla/mux"
	"github.com/jakecorrenti/fishfindr/analysis"
	"github.com/jakecorrenti/fishfindr/db"
	"github.com/jakecorrenti/fishfindr/heatmap"
	"github.com/jakecorrenti/fishfindr/types"
	_ "github.com/mattn/go-sqlite3"
)

var (
	locationsRepository *db.SQLiteRepository
	lock                sync.Mutex
)

func graphHandler(w http.ResponseWriter, r *http.Request) {
	locations, err := locationsRepository.All()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	clusters, noise, err := analysis.GetClusters(&locations)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	// initialize graph
	dimensions := 2
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)
	style := "points"

	xVals := make([]float64, 0)
	yVals := make([]float64, 0)

	// for each cluster in clusters, C stands for the cluster #, and Points
	// is all of the individual points within the cluster
	for _, cluster := range clusters {
		// each cluster needs to be plotted on the graph with a different color to signify they are different
		for _, i := range cluster.Points {
			xVals = append(xVals, locations[i].Latitude)
			yVals = append(yVals, locations[i].Longitude)
		}
		points := [][]float64{xVals, yVals}
		plot.AddPointGroup("Cluster "+fmt.Sprint(cluster.C), style, points)
		xVals = make([]float64, 0)
		yVals = make([]float64, 0)
	}

	for _, i := range noise {
		// each noise point needs to be plotted on the graph with a distinct color, such as black
		// to signify that they are just noise and do not matter much to the overall density of the points
		xVals = append(xVals, locations[i].Latitude)
		yVals = append(yVals, locations[i].Longitude)
	}
	points := [][]float64{xVals, yVals}
	plot.AddPointGroup("noise", style, points)
	plot.SetTitle("Catches")
	plot.SetXLabel("Lat")
	plot.SetYLabel("Long")
	plot.SavePlot("plot.png")

	img, err := os.Open("plot.png")
	if err != nil {
		log.Fatal(err)
	}

	defer img.Close()
	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, img)
}

func createNewLocation(w http.ResponseWriter, r *http.Request) {
	// create a sense of thread safety so that when there are multiple requests at the same time,
	// there isn't a data race when trying to store the locations
	lock.Lock()
	defer lock.Unlock()

	// get the body of POST request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// unmarshal the body into a new Location struct
	var location types.Location
	json.Unmarshal(body, &location)

	fmt.Println(location)

	_, err = locationsRepository.Create(location)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	json.NewEncoder(w).Encode(location)
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getJSON(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	locations, err := locationsRepository.All()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	j, err := heatmap.PointsJSON(&locations)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(string(j))
}

func main() {
	database, err := sql.Open("sqlite3", "locations.db")
	if err != nil {
		log.Fatal(err)
	}

	locationsRepository = db.NewSQLiteRepository(database)

	if err := locationsRepository.Migrate(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.HandleFunc("/graph", graphHandler)
	router.HandleFunc("/api/v1/json", getJSON)
	go router.HandleFunc("/api/v1/location", createNewLocation).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	server.ListenAndServe()
}
