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
| ⏳ **Timeout Support**   | Configurable request timeouts for reliability. |
| 🔄 **TLS Configuration** | Enables custom TLS settings for security. |
| 🔀 **Failover Mechanism** | Automatically switch to backup URLs if the primary server fails. |
| 🔐 **Secure TLS Support** | Customizable TLS settings for enhanced security. |
| ⏳ **Timeout Support**   | Prevents hanging requests by setting timeouts. |
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
##### 🔹 **Full HTTP Client with Custom Transport & Retries**
```go
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Initialize the quick framework.
	q := quick.New()

	// Define routes using quick.
	q.Get("/get", func(c *quick.Ctx) error {
		return c.Status(200).SendString("GET OK")
	})
	q.Post("/post", func(c *quick.Ctx) error {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(201).SendString("POST: " + string(body))
	})
	q.Put("/put", func(c *quick.Ctx) error {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(200).SendString("PUT: " + string(body))
	})
	q.Delete("/delete", func(c *quick.Ctx) error {
		return c.Status(200).SendString("DELETE OK")
	})
	q.Post("/postform", func(c *quick.Ctx) error {
		// Assume FormValues returns map[string][]string.
		form := c.FormValues()
		vals := url.Values(form)
		return c.Status(200).SendString("POSTFORM: " + vals.Encode())
	})

	// Create a test server using the quick handler.
	ts := httptest.NewServer(q)
	defer ts.Close()

	// Creating a custom HTTP transport with advanced settings.
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // Uses system proxy settings if available.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,             // Allows insecure TLS connections (not recommended for production).
			MinVersion:         tls.VersionTLS12, // Enforces a minimum TLS version for security.
		},
		MaxIdleConns:        50,    // Maximum number of idle connections across all hosts.
		MaxConnsPerHost:     30,    // Maximum simultaneous connections per host.
		MaxIdleConnsPerHost: 10,    // Maximum number of idle connections per host.
		DisableKeepAlives:   false, // Enables persistent connections (Keep-Alive).
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		Timeout:   15 * time.Second, // Global timeout for all requests.
		Transport: customTransport,  // Uses the custom transport.
	}

	// Create a client with extended options.
	cClient := client.New(
		client.WithTimeout(5*time.Second),
		client.WithDisableKeepAlives(false),
		client.WithMaxIdleConns(20),
		client.WithMaxConnsPerHost(20),
		client.WithMaxIdleConnsPerHost(20),
		client.WithContext(context.Background()),
		client.WithCustomHTTPClient(customHTTPClient),
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer EXAMPLE_TOKEN",
		}),
	
	// Also configure client retry settings (for manual retry logic).
		client.WithRetry(client.RetryConfig{
			MaxRetries: 2,
			Delay:      1 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500},
			EnableLog:  false,
		}),
	)

	// GET request.
	resp, err := cClient.Get(ts.URL + "/get")
	if err != nil {
		fmt.Println("GET Error:", err)
		return
	}
	fmt.Println("GET:", string(resp.Body))

	// POST request with a string body.
	resp, err = cClient.Post(ts.URL+"/post", "Hello, extended POST!")
	if err != nil {
		fmt.Println("POST Error:", err)
		return
	}
	fmt.Println("POST:", string(resp.Body))

	// PUT request with a struct body (marshaled to JSON).
	data := struct {
		Data string `json:"data"`
	}{
		Data: "Hello, extended PUT!",
	}
	resp, err = cClient.Put(ts.URL+"/put", data)
	if err != nil {
		fmt.Println("PUT Error:", err)
		return
	}
	// To display the JSON response as a string, unmarshal and marshal it back.
	var putResult map[string]string
	_ = json.Unmarshal(resp.Body, &putResult)
	putJSON, _ := json.Marshal(putResult)
	fmt.Println("PUT:", string(putJSON))

	// DELETE request.
	resp, err = cClient.Delete(ts.URL + "/delete")
	if err != nil {
		fmt.Println("DELETE Error:", err)
		return
	}
	fmt.Println("DELETE:", string(resp.Body))

	// POSTFORM request.
	formData := url.Values{}
	formData.Set("key", "value")
	resp, err = cClient.PostForm(ts.URL+"/postform", formData)
	if err != nil {
		fmt.Println("POSTFORM Error:", err)
		return
	}
	fmt.Println("POSTFORM:", string(resp.Body))	
}

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
					"https://reqres.in/api/users",
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
			FailoverURLs: []string{"http://backup1", "https://reqres.in/api/users", "https://httpbin.org/post"},
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
