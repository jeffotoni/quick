package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// config set
	cfg := &client.HTTPClientConfig{
		Timeout:             20 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConns:        20,
		MaxConnsPerHost:     20,
		MaxIdleConnsPerHost: 20,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
	}

	// Create a client with the custom configuration.
	cClient := client.NewClient(
		client.WithContext(context.TODO()),
		client.WithHeaders(map[string]string{"Content-Type": "application/xml"}),
		client.WithHTTPClientConfig(cfg),
	)

	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := cClient.Post("http://localhost:3000/v1/user", data)
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
