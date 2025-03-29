# 🧩 HTML Template Engine for Quick ![Quick Logo](./quick.png)


### ⚡️ Quick Template Engine

**Quick** is a high-performance, lightweight web framework for Go, and includes a powerful HTML template engine built on top of Go’s standard `html/template` package.

This template engine supports layouts, partials, embedded files, aliasing, and custom functions — making it ideal for modern web apps and APIs.


Quick provides native support for html/template with layout and partial composition, plus custom functions and embed.FS support.
🔧 Features

 - Nested layouts using {{ .yield }}

 - Template aliasing (e.g., index, index.html, or views/index.html)

 - Custom functions via AddFunc

 - Support for html/template's {{ define }} and {{ template }}

 - Compatible with fs.FS for embedded templates (embed.FS)

 - Auto-loading and caching via Load()

### 📁 Template Structure Example


```bash
.
├── main.go
└── views/
    ├── index.html
    └── layouts/
        ├── main.html
        └── base.html
```

```

🧠 Important Note About define

When using the filesystem approach, your HTML files must use {{ define "name" }} to register correctly:


### 🚀 Basic Usage (Filesystem)

⚠️ Note: When using the filesystem loader (html.New(...)), each template file must start with {{ define "name" }} and end with {{ end }}.


Example:

views/index.html

```html

{{ define "index" }}
  <h2>{{ .Title }}</h2>
  <p>{{ .Message }}</p>
{{ end }}

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
    engine.AddFunc("upper", strings.ToUpper)
    engine.Load()

    app := quick.New(quick.Config{
        Views: engine,
    })

    app.Get("/", func(c *quick.Ctx) error {
        return c.HTML("index", map[string]any{
            "Title":   "Quick + Templates",
            "Message": "This is your index content",
        })
    })

    app.Listen(":8080")
}

```

# 📦 Using with embed.FS

When using embed, define is optional but still recommended for advanced template reuse (e.g., base layouts or partials).

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
    engine.Dir = "views" // Set base directory for trimming paths
    engine.AddFunc("upper", strings.ToUpper)
    engine.Load()

    app := quick.New(quick.Config{
        Views: engine,
    })

    app.Get("/", func(c *quick.Ctx) error {
        return c.HTML("index", map[string]any{
            "Title":   "Quick + Templates (embed)",
            "Message": "This is embedded content",
        })
    })

    app.Listen(":8080")
}

```

### 🧱 Layout Example Filesystem


views/layouts/main.html

```html
{{ define "layouts/main" }}
<!DOCTYPE html>
<html>
  <head><title>{{ .Title }}</title></head>
  <body>
    <header><h1>My Site</h1></header>
    <main>{{ .yield }}</main>
  </body>
</html>
{{ end }}
```

views/layouts/base.html
```html
{{ define "layouts/base" }}
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <div>
        <p>BASE HEADER</p>
        {{ .yield }}
        <p>BASE FOOTER</p>
    </div>
</body>
</html>
{{ end }}

```

views/index.html
```html
{{ define "index" }}
<h2>{{ .Title }}</h2>
<p>{{ .Message }}</p>
{{ end }}

```


### 🧱 Layout Example embed.FS


views/layouts/main.html

```html
<!DOCTYPE html>
<html>
  <head><title>{{ .Title }}</title></head>
  <body>
    <header><h1>My Site</h1></header>
    <main>{{ .yield }}</main>
  </body>
</html>
{{ end }}
```

views/layouts/base.html
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <div>
        <p>BASE HEADER</p>
        {{ .yield }}
        <p>BASE FOOTER</p>
    </div>
</body>
</html>
{{ end }}

```

views/index.html
```html
<h2>{{ .Title }}</h2>
<p>{{ .Message }}</p>
{{ end }}

```

🧩 Render with Layout
```go

app.Get("/layout", func(c *quick.Ctx) error {
    return c.HTML("index", map[string]any{
        "Title":   "Page with Layout",
        "Message": "This content is wrapped",
    }, "layouts/main")
})

```

🔧 Custom Functions

```go

engine.AddFunc("upper", strings.ToUpper)

```


Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

