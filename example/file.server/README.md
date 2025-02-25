## 📦 Files Server in Quick ![Quick Logo](/quick.png)


#### 📌 What is Static Files?
**`Static files`** are files that are served as-is to the client without any server-side processing. These files include HTML, CSS, JavaScript, images, and other assets that are not generated dynamically by the server. Using static files allows you to build efficient and fast web applications.


### Serving Static Files with Quick Framework

This example sets up a basic web server that serves static files, such as HTML, CSS, or JavaScript, using the Quick framework.

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

### Embedding Files

This example incorporates static files into the binary using the embed package and serves them using the Quick framework. This allows you to include all static files directly into the executable, making deployment easier.

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