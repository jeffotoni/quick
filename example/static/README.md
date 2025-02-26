## 📦 Files Server in Quick ![Quick Logo](/quick.png)

#### 📂 STATIC FILES
🔹 How It Works
    
1. The server listens for HTTP requests targeting static file paths.
2. If a requested file exists in the configured directory, the server reads and returns the file as a response.
3. MIME types are automatically determined based on the file extension.

:zap: Key Features
- Efficient handling: Serves files directly without additional processing.
- MIME type detection: Automatically identifies file types for proper rendering.
- Caching support: Can be configured to improve performance via HTTP headers.
- Directory listing: (Optional) Allows browsing available static files.

:warning: Security Considerations
- Restrict access to sensitive files (.env, .git, etc.).
- Configure CORS policies when necessary.
- Use a Content Security Policy (CSP) to mitigate XSS risks.

#### Serving Static Files with Quick Framework

This example sets up a basic web server that serves static files, such as HTML, CSS, or JavaScript.

```go    
    // Create a new Quick server instance
    q := quick.New()

    // Static Files Setup
    // Serves files from the "./static" directory under the "/static" URL path.
    q.Static("/static", "./static")

    // Route Definition
    // Defines a route to serve the "index.html" file when accessing "/".
    q.Get("/", func(c *quick.Ctx) error {
        c.File("./static/index.html")
        return nil
    })

    // Starting the Server
    // Starts the server to listen for incoming requests on port 8080.
    q.Listen("0.0.0.0:8080")

```
#### 📁 EMBED
🔹 How Embedded Static Files Work
    
1. Static assets are compiled directly into the binary at build time (e.g., using Go’s embed package).
2. The application serves these files from memory instead of reading from disk.
3. This eliminates external dependencies, making deployment easier.

:zap:  Advantages of Embedded Files
- Portability: Single binary distribution without extra files.
- Performance: Faster access to static assets as they are stored in memory.
- Security: Reduces exposure to external file system attacks.

### Embedding Files

This example incorporates static files into the binary using the embed package and serves them using the Quick structure.

```go
//go:embed static/*
var staticFiles embed.FS

func main() {
	// Server Initialization
	// Creates a new instance of the Quick server
	q := quick.New()

	// Static Files Setup
	// Defines the directory for serving static files using the embedded files
	q.Static("/static", staticFiles)

	// Route Definition
	// Defines a route that serves the HTML index file
	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/index.html") // Renders the index.html file
		return nil
	})

	// Starting the Server
	// Starts the server on port 8080, listening on all addresses
	q.Listen("0.0.0.0:8080")
}

```
#### 📌 Testing with cURL

##### 🔹Testing the Static File Directory:
```bash
# Accessing the index.html file from the static directory
curl http://localhost:8080/static/index.html
```

##### 🔹Testing Embedded Static Files
```bash
# Accessing embedded static files
curl http://localhost:8080/
```

##### 🚀 You can now implement fast and efficient file uploads in Quick! 🔥

#### **📌 What I included in this README**
- ✅ Basic explanation of static files.
- ✅ Example of serving static files with Quick.
- ✅ Example of embedding static files with Quick.
- ✅ Example of testing with cURL.

Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥