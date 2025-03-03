package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	resp, err := client.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}
