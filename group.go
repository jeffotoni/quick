package quick

import (
	"net/http"
	"strings"

	"github.com/jeffotoni/quick/internal/concat"
)

// Constants for route processing
const (
	methodSeparator     = "#"
	errInvalidExtractor = "Invalid function signature for paramExtractor"
)

// Group represents a collection of routes that share a common prefix
type Group struct {
	prefix      string
	routes      []Route
	middlewares []func(http.Handler) http.Handler
	quick       *Quick
}

// Use adds middlewares to the group
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

// normalizePattern constructs the full path with the group prefix
// The result will normalizePattern(prefix, pattern string) string
func normalizePattern(prefix, pattern string) string {
	return concat.String(strings.TrimRight(prefix, "/"), "/", strings.TrimLeft(pattern, "/"))
}

// resolveParamExtractor ensures the correct function signature for paramExtractor
// The result will resolveParamExtractor(q *Quick, handlerFunc HandleFunc, paramExtractor interface{}, path, params string) http.HandlerFunc
func resolveParamExtractor(q *Quick, handlerFunc HandleFunc, paramExtractor interface{}, path, params string) http.HandlerFunc {
	switch fn := paramExtractor.(type) {
	case func(*Quick, HandleFunc) http.HandlerFunc:
		return fn(q, handlerFunc)
	case func(*Quick, string, string, HandleFunc) http.HandlerFunc:
		return fn(q, path, params, handlerFunc)
	default:
		panic(errInvalidExtractor)
	}
}

// applyMiddlewares applies all middlewares to a handler
// The result will applyMiddlewares(handler http.HandlerFunc, middlewares []func(http.Handler) http.Handler) http.HandlerFunc
func applyMiddlewares(handler http.HandlerFunc, middlewares []func(http.Handler) http.Handler) http.HandlerFunc {
	for _, mw := range middlewares {
		handler = mw(handler).(http.HandlerFunc) // CORREÇÃO: Garante conversão correta
	}
	return handler
}

// createAndRegisterRoute creates a new route and registers it in the Quick router
// The result will createAndRegisterRoute(g *Group, method, pattern, compiledPattern, params string, handler http.HandlerFunc)
func createAndRegisterRoute(g *Group, method, pattern, compiledPattern, params string, handler http.HandlerFunc) {
	route := Route{
		Pattern: compiledPattern,
		Path:    pattern,
		Params:  params,
		handler: handler,
		Method:  method,
		Group:   g.prefix,
	}
	g.quick.appendRoute(&route)

	// FIX: Adjust path in mux to maintain compatibility with tests
	if method == http.MethodGet {
		g.quick.mux.HandleFunc(pattern, handler)
	} else {
		g.quick.mux.HandleFunc(concat.String(strings.ToLower(method), methodSeparator, pattern), handler)
	}
}

// Handle registers a new route dynamically
// The result will Handle(method, pattern string, handlerFunc HandleFunc, paramExtractor interface{})
func (g *Group) Handle(method, pattern string, handlerFunc HandleFunc, paramExtractor any) {
	// Normalize pattern and extract parameters
	pattern = normalizePattern(g.prefix, pattern)
	path, params, compiledPattern := extractParamsPattern(pattern)

	// Resolve parameter extractor and apply middlewares
	handler := resolveParamExtractor(g.quick, handlerFunc, paramExtractor, path, params)
	handler = applyMiddlewares(handler, g.middlewares)

	// Register route
	createAndRegisterRoute(g, method, pattern, compiledPattern, params, handler)
}

// Get registers a new GET route
// The result will Get(pattern string, handlerFunc HandleFunc)
func (g *Group) Get(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodGet, pattern, handlerFunc, extractParamsGet)
}

// Post registers a new POST route
// The result will Post(pattern string, handlerFunc HandleFunc)
func (g *Group) Post(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodPost, pattern, handlerFunc, extractParamsPost)
}

// Put registers a new PUT route
// The result will Put(pattern string, handlerFunc HandleFunc)
func (g *Group) Put(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodPut, pattern, handlerFunc, extractParamsPut)
}

// Delete registers a new DELETE route
// The result will Delete(pattern string, handlerFunc HandleFunc)
func (g *Group) Delete(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodDelete, pattern, handlerFunc, extractParamsDelete)
}

// Patch registers a new PATCH route
// The result will Patch(pattern string, handlerFunc HandleFunc)
func (g *Group) Patch(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodPatch, pattern, handlerFunc, extractParamsPatch)
}

// Options registers a new OPTIONS route
// The result will Options(pattern string, handlerFunc HandleFunc)
func (g *Group) Options(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodOptions, pattern, handlerFunc, extractParamsOptions)
}
