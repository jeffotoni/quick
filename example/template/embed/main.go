// project/
// │
// ├── main.go
// ├── views/
// │   ├── index.html
// │   └── layouts/
// │       ├── main.html
// │       └── base.html

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
	engine.Dir = "views"
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
