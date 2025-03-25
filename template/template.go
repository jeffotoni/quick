package template

import (
	"io"
)

// TemplateEngine is the interface for pluggable template engines.
type TemplateEngine interface {
	Render(w io.Writer, name string, data interface{}, layouts ...string) error
	AddFunc(name string, fn interface{})
}
