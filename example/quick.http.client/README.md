## 🌐 HTTP Client - Quick Framework

The **Quick HTTP Client** provides an intuitive and flexible way to make HTTP requests, including **GET, POST, PUT, and DELETE** methods. It simplifies handling requests, responses, and custom configurations.

---

#### 🚀 Running the Example
Before using the Quick HTTP Client, make sure to start the Quick server by running:

```bash
$ go run server.go
```
This will start the API server on http://localhost:3000, making it ready to receive HTTP requests.



#### ✅ Key Features

| Feature                  | Benefit |
|--------------------------|---------|
| 🌍 **RESTful Requests**  | Supports GET, POST, PUT, DELETE, and more. |
| ⚡ **Easy JSON Handling** | Automatically marshals and unmarshals JSON data. |
| 📝 **Form Data Support**  | Easily send application/x-www-form-urlencoded requests with PostForm. |
| 🔧 **Custom Headers**    | Allows setting custom request headers. |
| ⏳ **Configurable Timeouts**  | Ensures timely responses, enhancing reliability. |
| 🔄 **TLS Configuration** | Enables custom TLS settings for security. |
| 🔀 **Failover Mechanism** | Automatically switch to backup URLs if the primary server fails. |
| 🔐 **Secure TLS Support** | Customizable TLS settings for enhanced security. |
| 🏎 **High Performance**  | Optimized HTTP client with keep-alive and connection pooling. |


### ✅ Method Reference

| Method Signature                                                                          | Description                                           |
| ----------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `func Get(url string) (*ClientResponse, error)`                                           | Global GET request using the default client           |
| `func Post(url string, body any) (*ClientResponse, error)`                                | Global POST request with flexible body input          |
| `func PostForm(url string, formData url.Values) (*ClientResponse, error)`                 | Global POST request sending form-data (URL-encoded)  |
| `func Put(url string, body any) (*ClientResponse, error)`                                 | Global PUT request with flexible body input           |
| `func Delete(url string) (*ClientResponse, error)`                                        | Global DELETE request using the default client        |
| `func (c *Client) Get(url string) (*ClientResponse, error)`                               | GET request using a custom client instance            |
| `func (c *Client) Post(url string, body any) (*ClientResponse, error)`                    | POST request using a custom client instance           |
| `func (c *Client) PostForm(url string, formData url.Values) (*ClientResponse, error)`     | POST request sending form-data (URL-encoded) with a custom client |
| `func (c *Client) Put(url string, body any) (*ClientResponse, error)`                     | PUT request using a custom client instance            |
| `func (c *Client) Delete(url string) (*ClientResponse, error)`                            | DELETE request using a custom client instance         |
| `func New(opts ...Option) *Client`                                                  | Creates a new Client with optional custom configurations|
| `func WithContext(ctx context.Context) Option`                                            | Option to set a custom context for the client         |
| `func WithHeaders(headers map[string]string) Option`                                      | Option to set custom headers                          |
| `func WithHTTPClientConfig(cfg *HTTPClientConfig) Option`                                 | Option to set a custom HTTP transport configuration   |
| `func WithRetry(cfg RetryConfig) Option`                                                 | Option to enable **automatic retries** for requests.   |
| `func WithTLSConfig(cfg *tls.Config) Option`                                             | Option to configure **TLS settings** for secure requests. |
| `func WithTimeout(timeout time.Duration) Option`                                         | Option to set a **timeout** for all requests.          |
| `func WithFailover(urls []string) Option`                                                | Option to enable **failover mechanism** with backup URLs. |
| `func WithCustomHTTPClient(client *http.Client) Option`                                  | Option to use a **custom HTTP client**.                |


---

##### 🚀 How It Works

The Quick HTTP Client allows **quick API calls** while providing **flexibility** through custom configurations. You can use **default clients** or **custom clients** with timeout, headers, and TLS settings.

---

##### 📌 Example Usage

##### 🔹 **GET Request**
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Use the default client
	resp, err := client.Get("http://localhost:3000/v1/user/1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}
```

##### 🔹 **POST Request (JSON Struct)**
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
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := client.Post("http://localhost:3000/v1/user", data)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result["message"])
}
```

##### 🔹 **PUT Request**
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, PUT Quick!",
	}

	resp, err := client.Put("http://localhost:3000/v1/user/1234", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PUT response:", string(resp.Body))
}
```

##### 🔹 **DELETE Request**
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	resp, err := client.Delete("http://localhost:3000/v1/user/1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DELETE response:", string(resp.Body))
}

```

##### 🔹 **Custom HTTP Client Configuration**
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

func main() {
	// Custom configuration
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

	// Create a new custom client
	cClient := client.New(
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

	// Unmarshal JSON response
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result["message"])
}

```

##### 🔹 **HTTP Request with Retry Support**
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)
cClient := client.New(
    client.WithRetry(
    	3,                 // Maximum number of retries
		"2s",              // Delay between attempts
		true,              // Use exponential backoff
		"500,502,503,504", // HTTP status for retry
		true,              // show Logger
    ),
)

resp, err := cClient.Get("http://localhost:3000/v1/user/1234")
if err != nil {
    log.Fatal(err)
}
fmt.Println("GET response:", string(resp.Body))
```
---
##### 🔹 **Custom HTTP Client with Advanced Configuration & Retries**
```go

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Creating a CookieJar to manage cookies automatically.
	jar, _ := cookiejar.New(nil)

	// Creating a fully custom *http.Client.
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second, // Sets a global timeout of 10 seconds.
		Jar:     jar,              // Uses a CookieJar to store cookies.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allows up to 3 redirects.
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Allows insecure TLS (not recommended for production).
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        50,    // Maximum idle connections.
			MaxConnsPerHost:     30,    // Max simultaneous connections per host.
			MaxIdleConnsPerHost: 10,    // Max idle connections per host.
			DisableKeepAlives:   false, // Enables keep-alive.
		},
	}

	// Creating a quick client using the custom *http.Client.
	cClient := client.New(
		client.WithCustomHTTPClient(customHTTPClient), // Uses the pre-configured HTTP client.
		client.WithContext(context.Background()),      // Sets a request context.
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		// Enables retry for specific HTTP status codes using the new RetryConfig.
		client.WithRetry(client.RetryConfig{
			MaxRetries: 3,                         // Maximum number of retries.
			Delay:      1 * time.Second,           // Delay between attempts.
			UseBackoff: true,                      // Use exponential backoff.
			Statuses:   []int{500, 502, 503, 504}, // HTTP statuses for retry.
			EnableLog:  true,                      // Enable logger.
		}),
	)

	// Performing a GET request.
	resp, err := cClient.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}
	fmt.Println("GET Response:", string(resp.Body))

	// Performing a POST request.
	data := map[string]string{"name": "QuickFramework", "version": "1.0"}
	resp, err = cClient.Post("https://httpbin.org/post", data)
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Response:", string(resp.Body))
}
```
---
##### 🔹 HTTP Client Configuration for GET Requests
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Define a GET endpoint that forwards requests to an external API.
	q.Get("/api/users", func(c *quick.Ctx) error {
		// Create an HTTP client with specific configurations.
		cClient := client.New(
			// Set the timeout for the HTTP client to 10 seconds.
			client.WithTimeout(10*time.Second),
			client.WithMaxConnsPerHost(20),
			client.WithDisableKeepAlives(false),
			// Add custom headers, including content
			//  type and authorization token.
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
			// Use a background context for the HTTP client.
			//  This context cannot be cancelled
			// and does not carry any deadline. It is 
			// suitable for operations that run
			// indefinitely or until the application is shut down.
			client.WithContext(context.Background()),
		)

		// Perform a GET request to the external API.
		resp, err := cClient.Get("https://reqres.in/api/users/2")
		if err != nil {
			// Log the error and return a server error response if the GET request fails.
			log.Println("GET Error:", err)
			return c.Status(500).SendString("Failed to connect to external API")
		}

		// Log and return the response body from the external API.
		log.Println("GET Response:", string(resp.Body))
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start listening on port 3000 for incoming HTTP requests.
	q.Listen(":3000")
}

```
---
##### 🔹 HTTP Client Configuration for POST Requests
```go
package main

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Initialize the Quick framework.
	q := quick.New()

	// Define a POST endpoint to process incoming requests.
	q.Post("/api/users", func(c *quick.Ctx) error {
		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.Status(400).SendString("Error reading request body: " + err.Error())
		}

		// Check if the request body is empty
		if len(body) == 0 {
			return c.Status(400).SendString("Error: Request body is empty")
		}

		// Validate that the request body is valid JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			return c.Status(400).SendString("Error: Invalid JSON")
		}

		// Create a modular HTTP client with customizable options.
		cClient := client.New(
			// Sets the HTTP timeout to 5 seconds.
			client.WithTimeout(5*time.Second),

			// Enables or disables HTTP Keep-Alive 
			// connections (false = keep-alives enabled).
			client.WithDisableKeepAlives(false),

			// Adds custom headers to the request, including 
			// Content-Type and Authorization.
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
		)

		// Forward the request to the external API
		resp, err := cClient.Post("https://reqres.in/api/users", json.RawMessage(body))
		if err != nil {
			log.Println("Error making request to external API:", err)
			return c.Status(500).SendString("Error connecting to external API")
		}

		// Log response from external API for debugging
		log.Println("External API response:", string(resp.Body))

		// Return the response from the external API to the client
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start the server on port 3000
	q.Listen(":3000")
}

```
---
##### 🔹 HTTP Client Configuration for PUT Requests
```go
package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Define a PUT endpoint to update user data.
	q.Put("/api/users/2", func(c *quick.Ctx) error {
		// Read the request body from the client
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Read Error:", err)
			return c.Status(500).SendString("Failed to read request body")
		}

		// Create an HTTP client with specific configurations.
		cClient := client.New(
			// Set the timeout for the HTTP client to 10 seconds.
			client.WithTimeout(10*time.Second),
			// Add custom headers, including content type and authorization token.
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
			// Use a background context for the HTTP client. 
			// This context cannot be cancelled
			// and does not carry any deadline. It is suitable 
			// for operations that run
			// indefinitely or until the application is shut down.
			client.WithContext(context.Background()),
		)

		// Perform a PUT request to the external API with the data received from the client.
		resp, err := cClient.Put("https://reqres.in/api/users/2", requestBody)
		if err != nil {
			log.Println("PUT Error:", err)
			return c.Status(500).SendString("Failed to connect to external API")
		}

		// Log and return the response body from the external API.
		log.Println("PUT Response:", string(resp.Body))
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start listening on port 3000 for incoming HTTP requests.
	q.Listen(":3000")
}

```
---
##### 🔹 HTTP Client Configuration for DELETE Requests
```go
package main

import (
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Define a DELETE endpoint to delete user data.
	q.Delete("/api/users/2", func(c *quick.Ctx) error {
		// Create an HTTP client with specific configurations.
		cClient := client.New(
			client.WithTimeout(2*time.Second),
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
		)

		// Perform a DELETE request to the external API.
		resp, err := cClient.Delete("https://reqres.in/api/users/2")
		if err != nil {
			log.Println("DELETE Error:", err)
			return c.Status(500).SendString("Failed to connect to external API")
		}

		// Log and return the response body from the external API.
		log.Println("DELETE Response:", string(resp.Body))
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start listening on port 3000 for incoming HTTP requests.
	q.Listen(":3000")
}
```


### 📌 Testing HTTP Client Settings with Curl

#### 🔹 GET Request
```bash
$ curl --location 'http://localhost:3000/api/users' \
--header 'Authorization: Bearer EXAMPLE_TOKEN'
```

#### 🔹 POST Request
```bash
$ curl -X POST http://localhost:3000/api/users \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer EXAMPLE_TOKEN" \
    -d '{"name": "John Doe", "job": "Software Engineer"}'
```

#### 🔹 PUT Request
```bash
$ curl -X PUT https://reqres.in/api/users/2 \
    -H "Content-Type: application/json" \
    -d '{"name": "Morpheus", "job": "zion resident"}'
```

#### 🔹 DELETE Request
```bash
$ curl -X DELETE https://reqres.in/api/users/2 \
    -H "Authorization: Bearer EXAMPLE_TOKEN"
```


---

##### 🔹 **Full HTTP Client with Retries, Failover & Secure Requests**
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

func main() {

	cClient := client.New(
		client.WithTimeout(5*time.Second),
		client.WithDisableKeepAlives(false),
		client.WithMaxIdleConns(20),
		client.WithMaxConnsPerHost(20),
		client.WithMaxIdleConnsPerHost(20),
		client.WithContext(context.TODO()),
		client.WithHeaders(
			map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer Token"},
		),
		client.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}),

		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{502, 503, 504, 403},
				FailoverURLs: []string{
					"http://backup1",
					"https://httpbin_error.org/post",
					"https://httpbin.org/post"},
				EnableLog: true,
			}),
	)

	resp, err := cClient.Post("http://api.quick/v1/user",
		map[string]string{"message": "Hello, POST in Quick!"})
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
##### 🔹 **Advanced HTTP Retry Configuration**
```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// when I declare the 2 retrys, WithRetry RoundTripper and WithRetry ,
	// the With Retry RoundTripper overrides it which is executed.
	cClient := client.New(
		client.WithTimeout(5*time.Second),
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Advanced HTTP transport configuration
		client.WithTransportConfig(&http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     20,
			DisableKeepAlives:   false,
		}),

		// WithRetry
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				EnableLog:  true,
			}),

		// Retry quick
		// client.WithRetry(5, "2s-bex", "500,502,503,504"),
	)

	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Form Response:", string(resp.Body))
}
```
---
##### 🔹 **Custom HTTP Client with Retry and TLS**
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

func main() {

	// Configuring the HTTP client using a structured approach.
	//
	// The following settings are applied to the HTTP client:
	// - Timeout: Sets the maximum duration for requests (20 seconds).
	// - DisableKeepAlives: Controls whether keep-alive connections are disabled (false = keep-alives enabled).
	// - MaxIdleConns: Defines the maximum number of idle connections across all hosts (20).
	// - MaxConnsPerHost: Sets the maximum number of simultaneous connections to a single host (20).
	// - MaxIdleConnsPerHost: Defines the maximum number of idle connections per host (20).
	// - TLSClientConfig: Configures TLS settings, including:
	//     * InsecureSkipVerify: false (enables strict TLS verification).
	//     * MinVersion: TLS 1.2 (ensures a minimum TLS version for security).
	//
	// Using WithHTTPClientConfig(cfg), all the configurations are applied at once.
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
				EnableLog:  true,
			}),
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
```
---
##### 🔹 **Quick Server with Form Handling and HTTP Client**
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Define a route to process POST form-data
	q.Post("/postform", func(c *quick.Ctx) error {
		form := c.FormValues()
		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Quick server running at http://localhost:3000")
		if err := q.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Quick server: %v", err)
		}
	}()

	// Criando um cliente HTTP antes de chamar PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Define um timeout de 5s
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // Tipo correto para forms
		}),
	)

	// Check if the HTTP client was initialized correctly
	if cClient == nil {
		log.Fatal("Erro: cliente HTTP não foi inicializado corretamente")
	}

	// Declare Values
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Send a POST request
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Check if the response is valid
	if resp == nil || resp.Body == nil {
		log.Fatal("Erro: resposta vazia ou inválida")
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}
```

---
##### 🔹 **Retry Mechanism for GET Requests**
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	cClient := client.New(
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				EnableLog:  true,
			}),
	)

	resp, err := cClient.Get("http://localhost:3000/v1/user/1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}

```

---
##### 🔹 **POST Request with Retry and Failover URLs**
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cC := client.New(
		client.WithRetry(client.RetryConfig{
			MaxRetries:   3,
			Delay:        1 * time.Second,
			UseBackoff:   false,
			Statuses:     []int{502, 503, 504, 403},
			FailoverURLs: []string{"http://backup1", "https://httpbin_error.org/post", "http://hosterror"},
			EnableLog:    true,
		}),
		client.WithHeaders(map[string]string{
			"Authorization": "Bearer token",
		}),
	)

	// Perform the POST request
	resp, err := cC.Post("http://localhost:3000/v1/user", map[string]string{
		"name":  "Jefferson",
		"email": "jeff@example.com",
	})
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}

	// Print the response body and status code
	fmt.Println("POST Response Status:", resp.StatusCode)
	fmt.Println("POST Response Body:", string(resp.Body))

}
```

---
##### 🔹 **POST Request with Retries and Backoff**
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	cClient := client.New(
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				EnableLog:  true,
			}),
	)

	// Perform the POST request
	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{
		"name":  "Jefferson",
		"email": "jeff@example.com",
	})
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}

	// Print the response body and status code
	fmt.Println("POST Response Status:", resp.StatusCode)
	fmt.Println("POST Response Body:", string(resp.Body))
}

```
---
##### 🔹 **Quick Server with Form Submission**
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Define a route to process POST form-data
	q.Post("/postform", func(c *quick.Ctx) error {
		form := c.FormValues()
		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Quick server running at http://localhost:3000")
		if err := q.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Quick server: %v", err)
		}
	}()

	// Creating an HTTP client before calling PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Define um timeout de 5s
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // Correct type for forms
		}),
	)

	// Check if the HTTP client was initialized correctly
	if cClient == nil {
		log.Fatal("Erro: cliente HTTP não foi inicializado corretamente")
	}

	// Declare Values
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Send a POST request
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Check if the response is valid
	if resp == nil || resp.Body == nil {
		log.Fatal("Erro: resposta vazia ou inválida")
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}

```
---
### 📌 Testing with cURL

### 🔹 GET Request
```bash
$ curl -X GET http://localhost:3000/v1/user/1234
```

### 🔹 POST Request (Sending JSON)
```bash
$ curl -X POST http://localhost:3000/v1/user \
   -H "Content-Type: application/json" \
   -d '{"message": "Hello, POST!"}'
```

### 🔹 PUT Request
```bash
$ curl -X PUT http://localhost:3000/v1/user/1234 \
   -H "Content-Type: application/json" \
   -d '{"message": "Hello, PUT!"}'
```

### 🔹  DELETE Request
```bash
$ curl -X DELETE http://localhost:3000/v1/user/1234
```
---
#### 🔹 Advanced HTTP client configuration with failover mechanism

This code example showcases the setup of an HTTP client capable of handling network interruptions and server failures gracefully. It features custom transport configurations, including enhanced security settings, connection management, and a robust failover mechanism. Such a setup ensures that the application remains resilient and responsive under various network conditions. 

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	
	customTransport := &http.Transport{
		// Uses system proxy settings if available.
		Proxy: http.ProxyFromEnvironment, 
		TLSClientConfig: &tls.Config{
			// Allows insecure TLS connections (not recommended for production).
			InsecureSkipVerify: true,    
			// Enforces a minimum TLS version for security.         
			MinVersion:         tls.VersionTLS12, 
		},
		 // Maximum number of idle connections across all hosts.
		MaxIdleConns:        50,  
		// Maximum simultaneous connections per host. 
		MaxConnsPerHost:     30, 
		// Maximum number of idle connections per host.   
		MaxIdleConnsPerHost: 10, 
		// Enables persistent connections (Keep-Alive).   
		DisableKeepAlives:   false, 
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		// Sets a global timeout for all requests.
		Timeout: 5 * time.Second,
	}

	// Creating a client using both the custom transport and other configurations.
	cClient := client.New(
		// Applying the custom HTTP client.
		client.WithCustomHTTPClient(customHTTPClient), 
		// Custom context for request cancellation and deadlines.
		client.WithContext(context.Background()),      
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		// Applying the custom transport.
		client.WithTransport(customTransport), 
		// Setting a timeout for requests.
		client.WithTimeout(5*time.Second),     
		// Retry on specific status codes.
		client.WithRetry(
			client.RetryConfig{
				MaxRetries:   2,
				Delay:        1 * time.Second,
				UseBackoff:   true,
				Statuses:     []int{500},
				FailoverURLs: []string{"https://httpbin_error.org/post", "https://httpbin.org/post"},
				EnableLog:    true,
			}),
	)

	// call client to POST
	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"message": "Hello Post!!"})
	if err != nil {
		log.Fatal(err)
	}

	// show resp
	fmt.Println("POST response:\n", string(resp.Body))
}

```
---
#### 🔹 HTTP Client with Advanced Transport and Failover Capabilities

Explore how to set up an HTTP client that not only adheres to security best practices with TLS configurations but also ensures your application remains operational through network issues. This example includes detailed setups for handling HTTP client retries and switching to failover URLs when typical requests fail. Ideal for systems requiring high reliability and fault tolerance. 

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// Creating a custom HTTP transport with advanced settings.
	customTransport := &http.Transport{
		// Uses system proxy settings if available.
		Proxy: http.ProxyFromEnvironment, 
		TLSClientConfig: &tls.Config{
			 // Allows insecure TLS connections (not recommended for production).
			InsecureSkipVerify: true,      
			// Enforces a minimum TLS version for security.      
			MinVersion:         tls.VersionTLS12, 
		},
		// Maximum number of idle connections across all hosts.
		MaxIdleConns:        50, 
		// Maximum simultaneous connections per host.   
		MaxConnsPerHost:     30,   
		 // Maximum number of idle connections per host. 
		MaxIdleConnsPerHost: 10,  
		// Enables persistent connections (Keep-Alive). 
		DisableKeepAlives:   false, 
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		 // Sets a global timeout for all requests.
		Timeout:   5 * time.Second,
		// Uses the custom transport.
		Transport: customTransport, 
	}

	// Creating a client using both the custom transport and other configurations.
	cClient := client.New(
		// Applying the custom HTTP client.
		client.WithCustomHTTPClient(customHTTPClient), 
		 // Custom context for request cancellation and deadlines.
		client.WithContext(context.Background()),     
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		client.WithTimeout(5*time.Second), // Setting a timeout for requests.
		// Retry on specific status codes.
		client.WithRetry(
			client.RetryConfig{
				MaxRetries:   2,
				Delay:        1 * time.Second,
				UseBackoff:   true,
				Statuses:     []int{500},
				FailoverURLs: []string{"http://hosterror", "https://httpbin_error.org/post"},
				EnableLog:    true,
			}),
	)

	resp, err := cClient.Post("https://httpbin.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	// show resp
	fmt.Println("POST response:", string(resp.Body))
}
```
---
### 🔹Configuring HTTP Client with Retry and Failover Mechanisms

Discover how to build an HTTP client capable of dealing with network instabilities and server failures. This setup includes detailed retry configurations and introduces failover URLs to ensure that your application can maintain communication under adverse conditions. The example demonstrates using exponential backoff for retries and provides multiple endpoints to guarantee the availability of services. 

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create a new HTTP client with specific configurations.
	cClient := client.New(
		// Set a timeout for all requests made by this client to 10 seconds.
		// This helps prevent the client from hanging indefinitely on requests.
		client.WithTimeout(10*time.Second),

		// Set default headers for all requests made by this client.
		// Here, 'Content-Type' is set to 'application/json'
		//  which is typical for API calls.
		client.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),

		// Enable automatic retry mechanism with specific configurations.
		// This is useful for handling intermittent errors and ensuring robustness.
		client.WithRetry(
			client.RetryConfig{
				 // Retry failed requests up to two times.
				MaxRetries: 2,      
				// Wait for 1 second before retrying.             
				Delay:      1 * time.Second, 
				 // Use exponential backoff strategy for retries.     
				UseBackoff: true,         
				// HTTP status codes that trigger a retry.       
				Statuses:   []int{500, 502, 503}, 
				// Alternate URLs to try if the main request fails.
				FailoverURLs: []string{ 
					"http://hosterror",
					"https://httpbin_error.org/post",
				},
				// Enable logging for retry operations.
				EnableLog: true, 
			}),
	)

	// Perform a POST request using the configured HTTP client.
	// Includes a JSON payload with a "name" key.
	resp, err := cClient.Post("https://httpbin.org/post", map[string]string{
		"name": "jeffotoni in action with Quick!!!",
	})

	// Check if there was an error with the POST request.
	if err != nil {
		// If an error occurs, log the error and terminate the program.
		log.Fatalf("POST request failed: %v", err)
	}

	// Print the response from the server to the console.
	fmt.Println("POST Form Response:", string(resp.Body))
}

```
---
### 🔹Advanced HTTP Client Configuration with Transport and Retry Settings

Explore the configuration of an HTTP client designed for high reliability and security in network communications. This example includes sophisticated transport settings, featuring TLS configurations for enhanced security, and a robust retry mechanism to handle request failures gracefully. These settings are essential for applications requiring reliable data exchange with external APIs, especially in environments where network stability might be a concern. 

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create an HTTP client with custom configurations using the Quick framework.
	cClient := client.New(
		// Set a global timeout for all requests made by this client to 10 seconds.
		// This helps prevent the client from hanging indefinitely on requests.
		client.WithTimeout(10*time.Second),

		// Set default headers for all requests made by this client.
		// Here, we specify that we expect to send and receive JSON data.
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Configure the underlying transport for the HTTP client.
		client.WithTransportConfig(&http.Transport{
			// Use the system environment settings for proxy configuration.
			Proxy: http.ProxyFromEnvironment,

			// Configure TLS settings to skip verification of the server's 
			// certificate chain and hostname.
			// Warning: Setting InsecureSkipVerify to true is not recommended for
			//  production as it is insecure.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},

			// Enable HTTP/2 for supported servers.
			ForceAttemptHTTP2: true,

			// Set the maximum number of idle connections in the connection pool for all hosts.
			MaxIdleConns: 20,

			// Set the maximum number of idle connections in the connection pool per host.
			MaxIdleConnsPerHost: 10,

			// Set the maximum number of simultaneous connections per host.
			MaxConnsPerHost: 20,

			// Keep connections alive between requests. This can help improve performance.
			DisableKeepAlives: false,
		}),
	)

	// Perform a POST request with a JSON payload.
	// The payload includes a single field "name" with a value.
	resp, err := cClient.Post("https://httpbin.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		// Log the error and stop the program if the POST request fails.
		log.Fatalf("POST request failed: %v", err)
	}

	// Output the response from the POST request.
	fmt.Println("POST Form Response:", string(resp.Body))
}

```
---
## 🔎📝 Regex

### 🔹Accepts only numbers in id
```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Route that accepts only numeric IDs (using regex [0-9]+)
	q.Get("/users/{id:[0-9]+}", func(c *quick.Ctx) error {
		id := c.Param("id")
		return c.JSON(map[string]string{
			"message": "User found",
			"user_id": id,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}
```

### 🔹Accepts only lowercase letters in the slug 

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Route that accepts only lowercase slugs (words with lowercase letters)
		q.Get("/profile/{slug:[a-z]+}", func(c *quick.Ctx) error {
		slug := c.Param("slug")
		return c.JSON(map[string]string{
			"message": "Profile found",
			"profile": slug,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

```

### 🔹Supports API version and numeric Id

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Route that accepts an API version (v1, v2, etc.) and a numeric user ID
	q.Get("/api/{version:v[0-9]+}/users/{id:[0-9]+}", func(c *quick.Ctx) error {
		version := c.Param("version")
		id := c.Param("id")
		return c.JSON(map[string]string{
			"message": "API Versioned User",
			"version": version,
			"user_id": id,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}
```
---

### 🎯 Why Use Quick HTTP Client?
- ✅ **Simplified API** – Eliminates the need to manually create and configure HTTP clients.
- ✅ **Flexible** – Supports multiple request methods (**GET, POST, PUT, DELETE**) and customizable configurations.
- ✅ **Optimized Performance** – Efficient connection handling with keep-alive, connection pooling, and reduced latency.
- ✅ **Automatic Retries** – Configurable retry logic with exponential backoff for handling transient failures.
- ✅ **Failover Mechanism** – Automatically switches to backup URLs if the primary server fails.
- ✅ **Secure Requests** – Customizable TLS settings for enhanced security and encrypted communication.
- ✅ **Timeout Control** – Prevents hanging requests by setting timeouts at the client level.
- ✅ **Custom Headers & Context** – Allows setting headers dynamically and supports request cancellation via context.Context.
- ✅ **Middleware Friendly** – Easily integrates with logging, authentication, and other middleware solutions.


---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
