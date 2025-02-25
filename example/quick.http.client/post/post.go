package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := client.Post("http://localhost:3000/v1/user", data)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result["message"])
}
