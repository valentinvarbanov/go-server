package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
)

type POIRecord struct {
	Name       string
	Route      string
	ExitNumber string
	Category   string
	Latitude   float64
	Longitude  float64
}

func createEntriesFromData(data [][]string) []POIRecord {
	var entries []POIRecord
	for i, line := range data {
		if i > 0 { // omit header line
			var rec POIRecord
			for j, field := range line {
				if j == 2 {
					rec.Route = field
				} else if j == 3 {
					rec.ExitNumber = field
				} else if j == 4 {
					rec.Category = field
				} else if j == 5 {
					rec.Name = field
				} else if j == 15 {
					rec.Latitude, _ = strconv.ParseFloat(field, 64)
				} else if j == 16 {
					rec.Longitude, _ = strconv.ParseFloat(field, 64)
				}
			}
			entries = append(entries, rec)
		}
	}
	return entries
}

func readDataFromCSV(filename string) ([]POIRecord, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := createEntriesFromData(data)

	return entries, nil
}

func serializePOIRecords(records []POIRecord) []byte {

	var pois []map[string]interface{} = make([]map[string]interface{}, 0)

	for _, record := range records {
		poi := make(map[string]interface{})

		poi["name"] = record.Name
		poi["route"] = record.Route
		poi["exitNumber"] = record.ExitNumber
		poi["category"] = record.Latitude
		poi["distance"] = 1000 // TODO:

		pois = append(pois, poi)
	}

	json, _ := json.MarshalIndent(pois, "", "    ")
	return json
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// https://stackoverflow.com/questions/27928/calculate-distance-between-two-latitude-longitude-points-haversine-formula/27943#27943
func getDistanceFromLatLon(lat1, lon1, lat2, lon2 float64) float64 {
	var R float64 = 6371 // Radius of the earth in km
	var dLat = deg2rad(lat2 - lat1)
	var dLon = deg2rad(lon2 - lon1)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c
	return d * 1000 // km to m
}

func main() {

	// entries, err := readDataFromCSV("data/mock.csv")
	entries, err := readDataFromCSV("data/data.csv")

	if err != nil {
		fmt.Printf("Failed to read input data, exiting.\n")
	}

	http.HandleFunc("/poi", func(w http.ResponseWriter, r *http.Request) {

		radius := 5000 // meters
		// TODO: get radius from query
		// category := r.URL.Query().Get("category")
		longitude, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		latitude, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)

		var filtered []POIRecord = make([]POIRecord, 0)
		for _, poi := range entries {

			distance := getDistanceFromLatLon(latitude, longitude, poi.Latitude, poi.Longitude)

			if distance < float64(radius) {
				filtered = append(filtered, poi)
			}
		}

		raw_data := serializePOIRecords(filtered)
		w.Write(raw_data)
	})

	const port = 8080
	fmt.Printf("Starting server at port %v\n", port)

	// sample requests
	// http://localhost:8080/poi?lon=-120.4&lat=37.3
	// http://localhost:8080/poi?lon=-120&lat=30

	var host = fmt.Sprintf(":%v", port)
	http.ListenAndServe(host, nil)
}
