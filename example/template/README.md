# 🖼️ Template Engine for Quick ![Quick Logo](/quick.png)


This package provides a flexible and extensible **template rendering engine** for the [Quick](https://github.com/jeffotoni/quick) web framework.  
It allows you to build dynamic HTML views using Go's standard `html/template` package, enriched with features like layout support, custom functions, and file system abstraction.

---

## 📌 What is a Template?

In web development, a **template** is a file that defines the structure of the output (usually HTML) with dynamic placeholders.  
You can inject data into these placeholders at runtime to render personalized content for each request.

### Example `index.html`:

```html
<h1>{{ .Title }}</h1>
<p>{{ .Message }}</p>
```

Given the data:
```go
map[string]interface{}{
  "Title": "Welcome",
  "Message": "This is your homepage.",
}
```
The rendered output becomes:
```html
<h1>Welcome</h1>
<p>This is your homepage.</p>
```

## 🧩 Rendering Templates with and without Layouts

The example below shows how to render templates in Quick using:

- A **basic template** (`/`)
- A template wrapped with a **single layout** (`/layout`)
- A template wrapped with **nested layouts** (`/layout-nested`)

It also demonstrates how to register custom template functions (e.g., `upper`) and how to configure the `html.Engine` to load `.html` files from the `views/` directory.

### 📁 Project Structure

```text
.
├── main.go
└── views/
    ├── index.html
    └── layouts/
        ├── main.html
        └── base.html
```

```go

package main

import (
	"strings"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/template/html"
)

func main() {
	engine := html.New("./views", ".html")

	// Example of adding a custom function
	engine.AddFunc("upper", strings.ToUpper)
	engine.Load()

	app := quick.New(quick.Config{
		Views: engine,
	})

	app.Get("/", func(c *quick.Ctx) error {
		return c.HTML("index", map[string]interface{}{
			"Title":   "Quick + Templates",
			"Message": "this is your index content in views",
		})
	})

	app.Get("/layout", func(c *quick.Ctx) error {
		return c.HTML("index", map[string]interface{}{
			"Title":   "Quick with Layout",
			"Message": "layout with main.html",
		}, "layouts/main")
	})

	app.Get("/layout-nested", func(c *quick.Ctx) error {
		return c.HTML("index", map[string]interface{}{
			"Title":   "Nested Layouts",
			"Message": "this is nested layout content",
		}, "layouts/main", "layouts/base")
	})

	app.Listen(":8080")
}
```

## 📦 Rendering Templates with embed.FS (Go 1.16+)

This example demonstrates how to embed templates into your Go binary using the `embed` package.  
This is useful for distributing a single executable without external template files.

---

### 📁 Embedded Project Structure

```text
project/
├── main.go
└── views/
    ├── index.html
    └── layouts/
        ├── main.html
        └── base.html
```
### 🧩 Example Using embed.FS
```go
package main

import (
	"embed"
	"strings"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/template/html"
)

//go:embed views/*.html views/layouts/*.html
var viewsFS embed.FS

func main() {
	engine := html.NewFileSystem(viewsFS, ".html")
	engine.Dir = "views" // required for path normalization
	engine.AddFunc("upper", strings.ToUpper)
	engine.Load()

	app := quick.New(quick.Config{
		Views: engine,
	})

	app.Get("/", func(c *quick.Ctx) error {
		return c.HTML("index", map[string]interface{}{
			"Title":   "Quick + Templates (embed)",
			"Message": "this is your index content in views (embedded)",
		})
	})

	app.Get("/layout", func(c *quick.Ctx) error {
		return c.HTML("index", map[string]interface{}{
			"Title":   "Quick with Layout",
			"Message": "layout with main.html",
		}, "layouts/main")
	})

	app.Get("/layout-nested", func(c *quick.Ctx) error {
		return c.HTML("index.html", map[string]interface{}{
			"Title":   "Nested Layouts",
			"Message": "this is nested layout content",
		}, "layouts/main", "layouts/base")
	})

	app.Listen(":8080")
}
```
### 🧪 Test with cURL
```bash
curl -i http://localhost:8080/
curl -i http://localhost:8080/layout
curl -i http://localhost:8080/layout-nested
```