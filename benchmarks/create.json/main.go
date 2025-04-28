package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Structure for data generation
type CreateJSON struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Year     int                    `json:"year"`
	Price    float64                `json:"price"`
	Big      bool                   `json:"big"`
	Car      bool                   `json:"car"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Options  []Option               `json:"options"`
	Extra    interface{}            `json:"extra"`
	Dynamic  map[string]interface{} `json:"dynamic"`
}

// Auxiliary structure for key/value inside Options
type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// List of random names
var names = []string{"Jefferson", "Arthur", "Mariana", "Carlos",
	"Fernanda", "João", "Beatriz", "Gabriel", "Maria", "Jaque",
	"Gael", "Vanessa", "Felipe", "Tatiane", "Lucas", "Larissa", "Pedro"}

// Generate a random name 
func randomName() string {
	return names[rand.Intn(len(names))]
}

// Function to generate a random instance
func generateRandomCreateJSON() CreateJSON {
	rand.Seed(time.Now().UnixNano())
	
	return CreateJSON{
		ID:    fmt.Sprintf("%d", rand.Intn(999999)),
		Name:  randomName(),
		Year:  rand.Intn(2050-2000) + 2000,
		Price: rand.Float64() * 1000,
		Big:   rand.Intn(2) == 1,
		Car:   rand.Intn(2) == 1,
		Tags:  []string{"golang", "api", "performance"},
		Metadata: map[string]interface{}{
			"created_at": time.Now().Format(time.RFC3339),
			"author":     randomName(),
		},
		Options: []Option{
			{Key: "theme", Value: "dark"},
			{Key: "language", Value: "pt-BR"},
		},
		Extra: nil,
		Dynamic: map[string]interface{}{
			"type":        "admin",
			"permissions": []string{"read", "write"},
		},
	}
}

// Function to generate and save JSON in list `[]` or single `{}` based on mode
func generateAndSaveJSON(numRecords int, mode string, filename string) error {
	var jsonData []byte
	var err error

	if mode == "single" {
		// Generate a JSON with a single object `{}` ("single" mode)
		CreateJSON := generateRandomCreateJSON()
		jsonData, err = json.MarshalIndent(CreateJSON, "", " ")
	} else {
		// Generate a list `[]` with multiple objects ("list" mode, default)
		data := make([]CreateJSON, numRecords)
		for i := 0; i < numRecords; i++ {
			data[i] = generateRandomCreateJSON()
		}
		jsonData, err = json.MarshalIndent(data, "", " ")
	}

	if err != nil {
		return err
	}

	// Save to file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	// Display the size of the generated file
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}

	fmt.Printf("Generated file: %s | Records: %d | Size: %.2f KB | Mode: %s\n",
		filename, numRecords, float64(fileInfo.Size())/1024, mode)

	return nil
}

func main() {
	// Check if the argument was past
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <number of records> [list|single]")
		os.Exit(1)
	}

	// Convert argument to integer
	numRecords, err := strconv.Atoi(os.Args[1])
	if err != nil || numRecords <= 0 {
		fmt.Println("Error: The number of records must be a positive integer.")
		os.Exit(1)
	}

	// Set the generation mode (default: "list")
	mode := "list"
	if len(os.Args) > 2 && (os.Args[2] == "single" || os.Args[2] == "list") {
		mode = os.Args[2]
	}

	// Set the file name dynamically
	filename := fmt.Sprintf("data_%dk_%s.json", numRecords/1000, mode)

	// Generate and save the JSON in the chosen mode
	err = generateAndSaveJSON(numRecords, mode, filename)
	if err != nil {
		fmt.Println("Error generating JSON:", err)
	}
}
