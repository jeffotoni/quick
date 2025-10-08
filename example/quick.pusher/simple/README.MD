# HTTP/2 Server Push with Quick

This project demonstrates how to use HTTP/2 Server Push in Go with the Quick framework to improve web page loading performance.

### 📋 Overview
The server implements HTTP/2 Server Push to send static resources (CSS and JavaScript) to the client before they even request them, reducing total page load time.

### 🚀 Features

* HTTP/2 Server Push: Proactive sending of static resources

* TLS Server: HTTPS with SSL certificates

* Serve static files: ./public directory accessible via /public

* Flexible routing: Captures all routes with /*

### 🛠️ Code Structure
```go
package main

import (
	"time"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Main route with HTTP/2 Push
	q.Get("/*", func(c *quick.Ctx) error {
		pusher, ok := c.Pusher()
		if ok {
			pusher.Push("/public/style.css", nil)
			pusher.Push("/public/app.js", nil)
		}

		time.Sleep(2 * time.Second) // Simulates processing

		html := `<!DOCTYPE html>
		<html>
		<head><link rel="stylesheet" href="/public/style.css"></head>
		<body>
			<h1>HTTP/2 Push Test</h1>
			<script src="/public/app.js"></script>
		</body>
		</html>`

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Status(200).SendString(html)
	})

	q.Static("/public", "./public")
	q.ListenTLS(":8443", "cert.pem", "key.pem", true)
}
```

### 📁 File Structure

project/
├── main.go
├── cert.pem
├── key.pem
└── public/
    ├── style.css
    └── app.js

#### 🔧 Setup

1. Generate SSL Certificates
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
```

2. Create Directory and Static Files
```bash
mkdir public
```
public/style.css:

```css

body {
    font-family: Arial, sans-serif;
    margin: 40px;
    background-color: #f0f0f0;
}

h1 {
    color: #333;
    text-align: center;
}
```

public/app.js:

```javascript

console.log('HTTP/2 Push working!');
document.addEventListener('DOMContentLoaded', function() {
    console.log('Page loaded successfully');
});
```
3. Run the Server
```bahs
go run main.go
```
Access: **https://localhost:8443**

### 🔍 How It Works

HTTP/2 Server Push

1. Pusher Detection: Server checks if client supports HTTP/2 Push

2. Proactive Push: Sends style.css and app.js before HTML

3. Client Cache: Resources are immediately available

📊 HTTP/2 Push Benefits
* RTT Reduction: Eliminates additional round trips

* Better Performance: Critical resources load faster

* Optimization: Especially useful for 
high-latency connections

### 🎯 Usage Examples
**Example 1: Page with Multiple Resources**


```go
q.Get("/dashboard", func(c *quick.Ctx) error {
    pusher, ok := c.Pusher()
    if ok {
        pusher.Push("/public/css/bootstrap.css", nil)
        pusher.Push("/public/css/dashboard.css", nil)
        pusher.Push("/public/js/chart.js", nil)
        pusher.Push("/public/js/dashboard.js", nil)
    }
    
    // Dashboard HTML...
})
```
**Example 2: API with Specific Resources**

```go
q.Get("/api/page/:id", func(c *quick.Ctx) error {
    pageID := c.Param("id")
    pusher, ok := c.Pusher()
    
    if ok {
        // Conditional push based on page type
        if pageID == "admin" {
            pusher.Push("/public/admin.css", nil)
            pusher.Push("/public/admin.js", nil)
        } else {
            pusher.Push("/public/guest.css", nil)
            pusher.Push("/public/guest.js", nil)
        }
    }
    
    // Return specific HTML...
})
```

**Example 3: E-commerce Page**
```go
q.Get("/product/:id", func(c *quick.Ctx) error {
    pusher, ok := c.Pusher()
    if ok {
        pusher.Push("/public/css/product.css", nil)
        pusher.Push("/public/js/gallery.js", nil)
        pusher.Push("/public/js/cart.js", nil)
        pusher.Push("/public/images/loading-spinner.svg", nil)
    }
    
    // Product HTML...
})
```

⚠️ Important Considerations

1. Cache Awareness
```go
// Example with cache headers
pusher.Push("/public/style.css", &http.PushOptions{
    Header: http.Header{
        "Cache-Control": []string{"public, max-age=31536000"},
    },
})
```

2. Avoid Excessive Push
    * Only critical above-the-fold resources

    * Monitor real performance

    * Consider user connections

3. HTTP/1.1 Fallback
```go
pusher, ok := c.Pusher()
if ok {
    // HTTP/2 Push
    pusher.Push("/public/style.css", nil)
} else {
    // HTTP/1.1 - preload via header
    c.Set("Link", "</public/style.css>; rel=preload; as=style")

}

```

### 🧪 Testing
**Verify HTTP/2 Push**

1. Chrome DevTools:

    * Network tab → check "Initiator" column

    * Push resources appear as "Push"



2. CURL:

```bash
curl -v --http2 https://localhost:8443
```
3. Browser:

* Access **chrome://net-export/** for detailed logs

📈 Monitoring

```go
// Example: Push statistics logging
q.Get("/*", func(c *quick.Ctx) error {
    pusher, ok := c.Pusher()
    if ok {
        start := time.Now()
        pusher.Push("/public/style.css", nil)
        pusher.Push("/public/app.js", nil)
        elapsed := time.Since(start)
        fmt.Printf("HTTP/2 Push completed in %v\n", elapsed)
    }
    // ...
})
```
### 🔒 Security
* SSL certificates are mandatory for HTTP/2

* Use valid certificates in production

* Consider rate limiting for pushes

### 🚀 Production Deployment
1. **Real Certificates:** Use Let's Encrypt or trusted CA

2. **Optimize Resources:** Compress CSS/JS

3. **Monitoring:** Track performance metrics

4. **A/B Testing:** Compare with and without push

This example demonstrates how HTTP/2 Server Push can significantly improve user experience by reducing critical resource loading time.

