package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jeffotoni/goquick/http/client"
)

func main() {
	callLocally()
	// callLocally2()
	// callLocally3()
	// callLetsGoQuick()
	// callGoDev()
	// callWithCustomClient()
}

func callLocally() {
	resp, err := client.Post("http://localhost:8000/post", io.NopCloser(strings.NewReader(`{"data": "client quick!"}`)))
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callLocally2() {
	buffData := []byte(`{"data": "client quick!"}`)
	resp, err := client.Post("http://localhost:8000/post", bytes.NewBuffer(buffData))
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callLocally3() {
	urlData := url.Values{}
	urlData.Set("quick", "is awesome!")

	resp, err := client.Post("http://localhost:8000/post", strings.NewReader(urlData.Encode()))
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callLetsGoQuick() {
	resp, err := client.Post("https://letsgoquick.com", io.NopCloser(strings.NewReader(`{"data": "client quick!"}`)))
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}

func callGoDev() {
	resp, err := client.Post("https://go.dev", io.NopCloser(strings.NewReader(`{"data": "client quick!"}`)))
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

	resp, err := c.Post("http://localhost:8000/post", io.NopCloser(strings.NewReader(`{"data": "client quick!"}`)))
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("response: %s | statusCode: %d", resp.Body, resp.StatusCode)
}
