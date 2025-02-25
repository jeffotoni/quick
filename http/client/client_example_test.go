package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// ExampleClient_Get demonstrates using the Client's Get method.
// The result will ExampleClient_Get()
func ExampleClient_Get() {
	// Create a test server that returns "GET OK" for GET requests.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET OK"))
	}))
	defer ts.Close()

	// Create a default client
	c := NewClient()

	// Send a GET request.
	resp, err := c.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// GET OK
}

// ExampleClient_Post demonstrates using the Client's Post method with a flexible body.
// The result will ExampleClient_Post()
func ExampleClient_Post() {
	// Create a test server that echoes the request body.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	}))
	defer ts.Close()

	// Create a default client.
	c := NewClient()

	// Example 1: Using a string as the POST body.
	resp, err := c.Post(ts.URL, "Hello, POST!")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String body:", string(resp.Body))

	// Example 2: Using a struct as the POST body, which will be marshaled to JSON.
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, JSON POST!",
	}
	resp, err = c.Post(ts.URL, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Println("Struct body:", result["message"])

	// Example 3: Using an io.Reader as the POST body.
	reader := strings.NewReader("Reader POST")
	resp, err = c.Post(ts.URL, reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("io.Reader body:", string(resp.Body))

	// Out put:
	// String body: Hello, POST!
	// Struct body: Hello, JSON POST!
	// io.Reader body: Reader POST
}

// ExampleClient_Put demonstrates using the Client's Put method with a flexible body.
// The result will ExampleClient_Put()
func ExampleClient_Put() {
	// Create a test server that echoes the request body.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer ts.Close()

	// Create a default client.
	c := NewClient()

	// Example 1: Using a string as the PUT body.
	resp, err := c.Put(ts.URL, "Hello, PUT!")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String body:", string(resp.Body))

	// Example 2: Using a struct as the PUT body, which will be marshaled to JSON.
	data := struct {
		Value int `json:"value"`
	}{Value: 42}
	resp, err = c.Put(ts.URL, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var result map[string]int
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Println("Struct body:", result["value"])

	// Out put:
	// String body: Hello, PUT!
	// Struct body: 42
}

// ExampleClient_Delete demonstrates using the Client's Delete method.
// The result will ExampleClient_Delete()
func ExampleClient_Delete() {
	// Create a test server that returns "DELETE OK" for DELETE requests.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DELETE OK"))
	}))
	defer ts.Close()

	// Create a default client.
	c := NewClient()

	// Send a DELETE request.
	resp, err := c.Delete(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// DELETE OK
}
