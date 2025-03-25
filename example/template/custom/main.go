// .
// ├── main.go
// ├── views/
// │   ├── index.html
// │   └── layouts/
// │       ├── main.html
// │       └── base.html

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
		return c.HTML("index", quick.M{
			"Title":   "Quick with Layout",
			"Message": "layout with main.html",
		}, "layouts/main")
	})

	app.Get("/layout-nested", func(c *quick.Ctx) error {
		return c.HTML("index", quick.M{
			"Title":   "Nested Layouts",
			"Message": "this is nested layout content",
		}, "layouts/main", "layouts/base")
	})

	app.Listen(":8080")
}
