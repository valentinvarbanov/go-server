package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type POIRecord struct {
    Name string
    Route string
    ExitNumber string
    Category string
    Latitude string
    Longitude string
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
                    rec.Latitude = field
                } else if j == 16 {
                    rec.Longitude = field
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

func serializePOIRecord(record POIRecord) []byte {

	poi := make(map[string]string)

	poi["name"] = record.Name
	poi["route"] = record.Route
	poi["exitNumber"] = record.ExitNumber
	poi["category"] = record.Latitude
	poi["distance"] = "1000" // TODO: 
	
	json, _ := json.Marshal(poi)
	return json
}

func main() {

    entries, err := readDataFromCSV("data/mock.csv")
	// entries, err := readDataFromCSV("data/data.csv")

	if err != nil {
		fmt.Printf("Failed to read input data, exiting.\n")
	}

	first_poi_raw := serializePOIRecord(entries[0])

	http.HandleFunc("/poi", func(w http.ResponseWriter, r *http.Request) {
		w.Write(first_poi_raw)
// 		fmt.Fprintf(w, `{ 
//     "status": "success", 
//     "message": "Welcome to my website!" 
// }`)
	})

	const port = 8080
	fmt.Printf("Starting server at port %v\n", port)

	var host = fmt.Sprintf(":%v", port)
	http.ListenAndServe(host, nil)
}
