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
| 🔧 **Custom Headers**    | Allows setting custom request headers. |
| ⏳ **Timeout Support**   | Configurable request timeouts for reliability. |
| 🔄 **TLS Configuration** | Enables custom TLS settings for security. |


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
- ✅ Simplified API: No need to manually create and configure HTTP clients.
- ✅ Flexible: Supports multiple request methods and custom configurations.
- ✅ Optimized for Performance: Efficient handling of HTTP requests and responses.

---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
