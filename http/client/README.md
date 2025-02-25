## http client Quick ![Quick Logo](/quick.png)

The **Client** package provides a flexible HTTP client that simplifies making HTTP requests (GET, POST, PUT, DELETE) with automatic body parsing. It supports passing request bodies as strings, structs (which are marshaled to JSON), or any type that implements `io.Reader`.

### Overview

This package offers:
- **Global convenience functions** for quick HTTP requests using a default client.
- **Custom client creation** using options to set context, headers, and HTTP transport configurations.
- **Flexible body parsing** for POST and PUT requests that accepts various input types.

### Method Signature

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
| `func NewClient(opts ...Option) *Client`                                                  | Creates a new Client with optional custom configurations|
| `func WithContext(ctx context.Context) Option`                                            | Option to set a custom context for the client         |
| `func WithHeaders(headers map[string]string) Option`                                      | Option to set custom headers                          |
| `func WithHTTPClientConfig(cfg *HTTPClientConfig) Option`                                 | Option to set a custom HTTP transport configuration   |

### Examples

#### GET Request Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Use the default client
	resp, err := client.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}

```

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
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := client.Post("https://example.com", data)
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

#### PUT Request Example (Using a String)
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Use a simple string as the PUT body
	resp, err := client.Put("https://example.com", "Hello, PUT!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PUT response:", string(resp.Body))
}

```

#### DELETE Request Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	resp, err := client.Delete("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DELETE response:", string(resp.Body))
}

```