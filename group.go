package quick

import (
	"net/http"
	"strings"

	"github.com/jeffotoni/quick/internal/concat"
)

// Group represents a collection of routes that share a common prefix
type Group struct {
	prefix      string
	routes      []Route
	middlewares []func(http.Handler) http.Handler
	quick       *Quick
}

// Use add middlewares to the group
// The result will Use(mw func(http.Handler) http.Handler)
func (g *Group) Use(mw func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, mw)
}

// Group creates a new route group with a shared prefix
// The result will Group(prefix string) *Group
func (q *Quick) Group(prefix string) *Group {
	g := &Group{
		prefix: prefix,
		routes: []Route{},
		quick:  q,
	}
	q.groups = append(q.groups, *g)
	return g
}

// Get registers a new GET route within the group
// The result will Get(pattern string, handlerFunc HandleFunc)
func (g *Group) Get(pattern string, handlerFunc HandleFunc) {
	pattern = concat.String(strings.TrimRight(g.prefix, "/"), "/", strings.TrimLeft(pattern, "/"))
	path, params, partternExist := extractParamsPattern(pattern)

	// Create the original handler
	handler := http.HandlerFunc(extractParamsGet(g.quick, path, params, handlerFunc))

	// Apply group middlewares (if any)
	for _, mw := range g.middlewares {
		handler = http.HandlerFunc(mw(handler).ServeHTTP)
	}

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: handler,
		Method:  http.MethodGet,
		Group:   g.prefix,
	}

	g.quick.appendRoute(&route)
	g.quick.mux.HandleFunc(path, route.handler)
}

// Post registers a new POST route within the group
// The result will Post(pattern string, handlerFunc HandleFunc)
func (g *Group) Post(pattern string, handlerFunc HandleFunc) {
	pattern = concat.String(strings.TrimRight(g.prefix, "/"), "/", strings.TrimLeft(pattern, "/"))
	_, params, partternExist := extractParamsPattern(pattern)

	pathPost := concat.String("post#", pattern)

	// Create the original handler
	handler := http.HandlerFunc(extractParamsPost(g.quick, handlerFunc))

	// Apply group middlewares (if any)
	for _, mw := range g.middlewares {
		handler = http.HandlerFunc(mw(handler).ServeHTTP)
	}

	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: handler,
		Method:  http.MethodPost,
		Params:  params,
		Group:   g.prefix,
	}

	g.quick.appendRoute(&route)
	g.quick.mux.HandleFunc(pathPost, route.handler)
}

// Put registers a new PUT route within the group
// The result will  Put(pattern string, handlerFunc HandleFunc)
func (g *Group) Put(pattern string, handlerFunc HandleFunc) {
	pattern = concat.String(strings.TrimRight(g.prefix, "/"), "/", strings.TrimLeft(pattern, "/"))
	_, params, partternExist := extractParamsPattern(pattern)

	pathPut := concat.String("put#", pattern)

	// Create the original handler
	handler := http.HandlerFunc(extractParamsPut(g.quick, handlerFunc))

	// Apply group middlewares (if any)
	for _, mw := range g.middlewares {
		handler = http.HandlerFunc(mw(handler).ServeHTTP)
	}

	// Setting up the group
	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: handler,
		Method:  http.MethodPut,
		Params:  params,
		Group:   g.prefix,
	}

	g.quick.appendRoute(&route)
	g.quick.mux.HandleFunc(pathPut, route.handler)
}

// Delete registers a new DELETE route within the group.
// The result will Delete(pattern string, handlerFunc HandleFunc)
func (g *Group) Delete(pattern string, handlerFunc HandleFunc) {
	pattern = concat.String(strings.TrimRight(g.prefix, "/"), "/", strings.TrimLeft(pattern, "/"))
	_, params, partternExist := extractParamsPattern(pattern)

	pathDelete := concat.String("delete#", pattern)

	// Create the original handler
	handler := http.HandlerFunc(extractParamsDelete(g.quick, handlerFunc))

	// Apply group middlewares (if any)
	for _, mw := range g.middlewares {
		handler = http.HandlerFunc(mw(handler).ServeHTTP)
	}

	// Setting up the group
	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		Params:  params,
		handler: handler,
		Method:  http.MethodDelete,
		Group:   g.prefix,
	}
	g.quick.appendRoute(&route)
	g.quick.mux.HandleFunc(pathDelete, route.handler)
}
