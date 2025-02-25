## 🌐 HTTP Client - Quick Framework

The **Quick HTTP Client** provides an intuitive and flexible way to make HTTP requests, including **GET, POST, PUT, and DELETE** methods. It simplifies handling requests, responses, and custom configurations.

---

##### ✅ Key Features

| Feature                  | Benefit |
|--------------------------|---------|
| 🌍 **RESTful Requests**  | Supports GET, POST, PUT, DELETE, and more. |
| ⚡ **Easy JSON Handling** | Automatically marshals and unmarshals JSON data. |
| 🔧 **Custom Headers**    | Allows setting custom request headers. |
| ⏳ **Timeout Support**   | Configurable request timeouts for reliability. |
| 🔄 **TLS Configuration** | Enables custom TLS settings for security. |

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

	// Unmarshal JSON response
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result["message"])
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
- ✅ Simplified API: No need to manually create and configure HTTP clients.
- ✅ Flexible: Supports multiple request methods and custom configurations.
- ✅ Optimized for Performance: Efficient handling of HTTP requests and responses.

---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
