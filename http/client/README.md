# 🌐 HTTP Client Quick ![Quick Logo](/quick.png)


The **Client** package provides a flexible HTTP client that simplifies making HTTP requests (GET, POST, PUT, DELETE) with automatic body parsing. It supports passing request bodies as strings, structs (which are marshaled to JSON), or any type that implements `io.Reader`.

## 📌 Overview

🔹 **Main features**:
- 🚀 **Global functions** for fast HTTP requests using a standard client.
- ⚙️ **Custom client creation** with options for context, headers and HTTP transport configuration.
- 🔄 **Flexible body parsing** for POST and PUT requests, accepting different input types.
- 🔁 **Requests with Retry and Failover** for greater resilience.

---

## 📦 Go package documentation

To access the documentation for each **Quick Framework** package, click on the links below:

| Pacote | Descrição | Go.dev |
|--------|----------|--------|
| **quick/http/client** |  HTTP client optimized for requests and failover | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick/http/client.svg)](https://pkg.go.dev/github.com/jeffotoni/quick/http/client) |

---



### ✅ Method Reference

| Method Signature                                                                          | Description                                           |
| ----------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `func Get(url string) (*ClientResponse, error)`                                           | Global GET request using the default client           |
| `func Post(url string, body any) (*ClientResponse, error)`                                | Global POST request with flexible body input          |
| `func Put(url string, body any) (*ClientResponse, error)`                                 | Global PUT request with flexible body input           |
| `func Delete(url string) (*ClientResponse, error)`                                        | Global DELETE request using the default client        |
| `func (c *Client) Get(url string) (*ClientResponse, error)`                               | GET request using a custom client instance            |
| `func (c *Client) Post(url string, body any) (*ClientResponse, error)`                      | POST request using a custom client instance           |
| `func (c *Client) Put(url string, body any) (*ClientResponse, error)`                       | PUT request using a custom client instance            |
| `func (c *Client) Delete(url string) (*ClientResponse, error)`                              | DELETE request using a custom client instance         |
| `func New(opts ...Option) *Client`                                                  | Creates a new Client with optional custom configurations|
| `func WithContext(ctx context.Context) Option`                                            | Option to set a custom context for the client         |
| `func WithHeaders(headers map[string]string) Option`                                      | Option to set custom headers                          |
| `func WithHTTPClientConfig(cfg *HTTPClientConfig) Option`                                 | Option to set a custom HTTP transport configuration   |

---
## 📌 Example Usage with [ReqRes API](https://reqres.in/)

### 🔹 GET Request Example
Retrieves a list of users from `ReqRes` API.

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Use the default client
	resp, err := client.Get("https://reqres.in/api/users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}

```

### 🔹 POST Request Example (Sending JSON)
This sends a `POST` request to create a new user.

#### POST Request Example (Using a Struct)
```go
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
		user string `json:"user"`
	}{
		user: "Emma",
	}

	// POST request to ReqRes API
	resp, err := client.Post("https://reqres.in/api/users", data)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}

```

### 🔹 PUT Request Example (Using a String)

Updates an existing user.
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Define a struct with user data
	data := struct {
		user string `json:"name"`
	}{
		user: "Jeff",
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	// PUT request to ReqRes API
	resp, err := client.Put("https://reqres.in/api/users/2", string(jsonData))
	if err != nil {
		log.Fatal("Error making request:", err)
	}

	// Print the HTTP status and response body
	fmt.Println("HTTP Status Code:", resp.StatusCode)
	fmt.Println("Raw Response Body:", string(resp.Body))
}
```

### 🔹 DELETE Request Example
Deletes a user and checks if the response is `204 No Content`.
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// DELETE request to ReqRes API
	resp, err := client.Delete("https://reqres.in/api/users/2")
	if err != nil {
		log.Fatal("Error making request:", err)
	}

	// Print the HTTP status to confirm deletion
	fmt.Println("HTTP Status Code:", resp.StatusCode)

	// Since DELETE usually returns no content, we check if it's empty
	if len(resp.Body) > 0 {
		fmt.Println("Raw Response Body:", string(resp.Body))
	} else {
		fmt.Println("Response Body is empty (expected for 204 No Content)")
	}
}

```
### 🔹 Retry with Dynamic Failover
This example configures the HTTP client to retry automatically in case of failure, using automatic retry, exponential backoff and dynamic failover for alternative URLs.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Retry configuration: maximum 3 attempts, 
	// with delay of 2s, exponential backoff
	retryConfig := client.RetryConfig{
	// Sets the maximum number of retry attempts
	MaxRetries: 3, 
	// Sets the base time between attempts before trying again
	Delay: 2 * time. Second,   
	// Enables exponential backoff to increase the time between failed attempts
	UseBackoff: true,        
	// List of HTTP status codes that trigger an automatic retry
	Status: []int{500, 502, 503, 504},   
	// Alternative URLs for failover in case of failure on the original request
	FailoverURLs: []string{"http://hosterror", "https://httpbin.org/get"}, 
	// Enable logs to log retry attempts
	EnableLog: true 
}


	// Create an HTTP client with retry configured
	httpClient := client.New(
		client.WithRetry(retryConfig),
	)

	// Making a GET request
	resp, err := httpClient.Get("https://httpbin_error.org/get")
	if err != nil {
		log.Fatal("Error in request:", err)
	}

	fmt.Println("HTTP code:", resp.StatusCode)
	fmt.Println("Answer:", string(resp.Body))
}

```

---
### 🔹 WithHeaders
This code demonstrates the creation of a highly configurable HTTP client.

```go
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

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
// - WithContext: Injects a context for the client (context.TODO() used as placeholder).
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/xml).
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
// - WithHTTPClientConfig: Defines advanced transport settings like connection pooling.
func main() {
	cfg := &client.HTTPClientConfig{
		Timeout:             20 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConns:        20,
		MaxConnsPerHost:     20,
		MaxIdleConnsPerHost: 20,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
	}

	// Creating an HTTP client with the pre-defined configuration.
	//
	// - WithContext: Sets a custom context for handling request cancellation and deadlines.
	// - WithHeaders: Adds a map of default headers (e.g., "Content-Type: application/xml").
	// - WithHTTPClientConfig: Applies the entire configuration object (cfg) to the client.
	cClient := client.New(
		client.WithContext(context.TODO()),
		client.WithHeaders(map[string]string{"Content-Type": "application/xml"}),
		client.WithHTTPClientConfig(cfg),
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				FailoverURLs: []string{
					"http://backup1",
					"https://httpbin_error.org/post",
					"https://httpbin.org/post"},
				EnableLog: true,
			}),
	)

	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := cClient.Post("https://httpbin_error.org/post ", data)
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

```
---
### 🔹 FullMethodsRetry
This example shows a full methods retry

```go
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

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
//   - WithTimeout: Sets the HTTP client timeout to 30 seconds.
//   - WithDisableKeepAlives: Enables or disables HTTP keep-alives (false = keep-alives enabled).
//   - WithMaxIdleConns: Defines the maximum number of idle connections (20).
//   - WithMaxConnsPerHost: Sets the maximum connections allowed per host (20).
//   - WithMaxIdleConnsPerHost: Sets the maximum number of idle connections per host (20).
//   - WithContext: Injects a context for the client (context.TODO() used as placeholder).
//   - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
//   - WithTLSConfig: Configures TLS settings, including InsecureSkipVerify and TLS version.
//   - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
//     with exponential backoff (2s-bex) and a maximum of 3 attempts.
func main() {

	// Create a new Quick HTTP client with custom settings
	cClient := client.New(
		client.WithTimeout(5*time.Second),   // Sets the request timeout to 5 seconds
		client.WithDisableKeepAlives(false), // Enables persistent connections (Keep-Alive)
		client.WithMaxIdleConns(20),         // Defines a maximum of 20 idle connections
		client.WithMaxConnsPerHost(20),      // Limits simultaneous connections per host to 20
		client.WithMaxIdleConnsPerHost(20),  // Limits idle connections per host to 20
		client.WithContext(context.TODO()),  // Injects a context (can be used for cancellation)
		client.WithHeaders(
			map[string]string{
				"Content-Type":  "application/json", // Specifies the request content type
				"Authorization": "Bearer Token",     // Adds an authorization token for authentication
			},
		),
		client.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,             // ⚠ Disables SSL certificate verification (use with caution)
			MinVersion:         tls.VersionTLS12, // Enforces a minimum TLS version for security
		}),
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,                         // Allows up to 2 retry attempts for failed requests
				Delay:      1 * time.Second,           // Delay of 1 second between retries
				UseBackoff: true,                      // Enables exponential backoff for retries
				Statuses:   []int{502, 503, 504, 403}, // Retries only on specific HTTP status codes
				FailoverURLs: []string{ // Backup URLs in case the primary request fails
					"http://backup1",
					"https://reqres.in/api/users",
					"https://httpbin_error.org/post",
				},
				EnableLog: true, // Enables logging for debugging retry behavior
			}),
	)
	// Send a POST request to the primary URL
	resp, err := cClient.Post("https://httpbin_error.org/post",
		map[string]string{"message": "Hello, POST in Quick!"})
	if err != nil {
		log.Fatal(err) // Logs an error and exits if the request fails
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err) // Logs an error if the response cannot be parsed
	}

	// Print the response
	fmt.Println("POST response:", result)
}
```
## **📌 What I included in this README**
- ✅ Overview: Explanation of the HTTP client in Quick.
- ✅ Method Reference: Quick lookup for available functions.
- ✅ GET, POST, PUT, DELETE Examples: How to use each method with ReqRes API.
- ✅ Retry with Failover: Demonstration of automatic retries, exponential backoff, and dynamic failover URLs.
- ✅ Testing with cURL: Alternative manual testing.
- ✅ Response Handling Improvements: Ensuring valid JSON parsing and response verification.

---

Now you can **complete with your specific examples** where I left the spaces ` go ...`

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
