package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type GeoJSONFeature struct {
	Type       string       `json:"type"`
	Properties GeoJSONProps `json:"properties"`
	Geometry   GeoJSONGeom  `json:"geometry"`
}

type GeoJSONProps struct {
	Timestamp string `json:"Timestamp"`
	IP        string `json:"IP"`
	City      string `json:"City"`
	Region    string `json:"Region"`
	Country   string `json:"Country"`
}

type GeoJSONGeom struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run csv2geojson.go <csvFilePath>")
		return
	}

	csvFilePath := os.Args[1]

	// Open the CSV file
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer csvFile.Close()

	// Parse the CSV file
	reader := csv.NewReader(csvFile)

	// Read and discard the header line
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV header:", err)
		return
	}

	// Read the rest of the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Extract the base name of the CSV file without extension
	baseName := strings.TrimSuffix(filepath.Base(csvFilePath), filepath.Ext(csvFilePath))

	// Create GeoJSON features
	var features []GeoJSONFeature
	for _, record := range records {
		// Parse latitude and longitude as float64
		latitude, _ := strconv.ParseFloat(record[5], 64)
		longitude, _ := strconv.ParseFloat(record[6], 64)

		// Create GeoJSON feature
		feature := GeoJSONFeature{
			Type: "Feature",
			Properties: GeoJSONProps{
				Timestamp: record[0],
				IP:        record[1],
				City:      record[2],
				Region:    record[3],
				Country:   record[4],
			},
			Geometry: GeoJSONGeom{
				Type:        "Point",
				Coordinates: []float64{longitude, latitude},
			},
		}

		// Append feature to features slice
		features = append(features, feature)
	}

	// Create GeoJSON structure
	geoJSON := struct {
		Type     string           `json:"type"`
		Features []GeoJSONFeature `json:"features"`
	}{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Generate output GeoJSON file name
	outputFileName := baseName + ".geojson"

	// Convert GeoJSON to JSON
	jsonData, err := json.MarshalIndent(geoJSON, "", "    ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Write JSON to file
	jsonFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating GeoJSON file:", err)
		return
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)

	fmt.Printf("GeoJSON file '%s' created successfully.\n", outputFileName)
}

