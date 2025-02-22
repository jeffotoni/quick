package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// callLocally()
	// callLetsGoQuick()
	// callGoDev()
	callWithCustomClient()
}

func callLocally() {
	resp, err := client.Delete("http://localhost:8000/delete")
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callLetsGoQuick() {
	resp, err := client.Delete("https://letsgoquick.com")
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callGoDev() {
	resp, err := client.Delete("https://go.dev")
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callWithCustomClient() {
	c := client.Client{
		Ctx: context.Background(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		ClientHttp: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				MaxConnsPerHost: 10,
			},
			Timeout: 0,
		},
	}

	resp, err := c.Delete("http://localhost:8000/delete")
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}
