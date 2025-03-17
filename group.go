// Package group provides functionality for managing grouped routes within the Quick framework.
//
// This package allows for the creation and organization of route groups, enabling
// the reuse of middleware and the application of common prefixes to related routes.
// It simplifies route management by structuring them into logical collections.
//
// Features:
// - Grouping of related routes under a shared prefix.
// - Support for middleware application at the group level.
// - Simplified registration of HTTP methods (GET, POST, PUT, DELETE, etc.).
// - Automatic handling of parameter extraction in routes.

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

// Group represents a collection of routes that share a common prefix.
//
// Fields:
//   - prefix: The URL prefix shared by all routes in the group.
//   - routes: A list of registered routes within the group.
//   - middlewares: A list of middleware functions applied to the group.
//   - quick: A reference to the Quick router.
type Group struct {
	prefix      string
	routes      []Route
	middlewares []func(http.Handler) http.Handler
	quick       *Quick
}

// Use adds middlewares to the group.
//
// Parameters:
//   - mw: A middleware function that modifies the HTTP handler.
//
// Example:
//
//	g.Use(loggingMiddleware)
func (g *Group) Use(mw func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, mw)
}

// Group creates a new route group with a shared prefix.
//
// Parameters:
//   - prefix: The common prefix for all routes in this group.
//
// Returns:
//   - *Group: A new Group instance.
//
// Example:
//
//	api := q.Group("/api")
func (q *Quick) Group(prefix string) *Group {
	g := &Group{
		prefix: prefix,
		routes: []Route{},
		quick:  q,
	}
	q.groups = append(q.groups, *g)
	return g
}

// normalizePattern constructs the full path with the group prefix.
//
// Parameters:
//   - prefix: The group's base URL path.
//   - pattern: The specific route pattern.
//
// Returns:
//   - string: The normalized URL path.
//
// Example:
//
//	fullPath := normalizePattern("/api", "/users") // "/api/users"
func normalizePattern(prefix, pattern string) string {
	return concat.String(strings.TrimRight(prefix, "/"), "/", strings.TrimLeft(pattern, "/"))
}

// resolveParamExtractor ensures the correct function signature for paramExtractor.
//
// Parameters:
//   - q: The Quick router instance.
//   - handlerFunc: The function handling the route.
//   - paramExtractor: A function for extracting parameters.
//   - path: The normalized route path.
//   - params: URL parameters.
//
// Returns:
//   - http.HandlerFunc: A wrapped handler with parameter extraction.
//
// Example:
//
//	handler := resolveParamExtractor(q, handlerFunc, extractParams, "/users/:id", "id")
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

// applyMiddlewares applies all middlewares to a handler.
//
// Parameters:
//   - handler: The original HTTP handler function.
//   - middlewares: A list of middleware functions.
//
// Returns:
//   - http.HandlerFunc: The wrapped handler with applied middlewares.
func applyMiddlewares(handler http.HandlerFunc, middlewares []func(http.Handler) http.Handler) http.HandlerFunc {
	for _, mw := range middlewares {
		handler = mw(handler).(http.HandlerFunc) // CORREÇÃO: Garante conversão correta
	}
	return handler
}

// createAndRegisterRoute creates a new route and registers it in the Quick router.
//
// Parameters:
//   - g: The group to which the route belongs.
//   - method: The HTTP method (GET, POST, etc.).
//   - pattern: The route pattern.
//   - compiledPattern: The compiled route pattern with parameters.
//   - params: URL parameters for the route.
//   - handler: The HTTP handler function.
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

// Handle registers a new route dynamically.
//
// Parameters:
//   - method: The HTTP method (GET, POST, etc.).
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//   - paramExtractor: The function to extract parameters.
//
// Example:
//
//	g.Handle("GET", "/users/:id", userHandler, extractParamsGet)
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

// Get registers a new GET route.
//
// Parameters:
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//
// Example:
//
//	g.Get("/users", listUsersHandler)
func (g *Group) Get(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodGet, pattern, handlerFunc, extractParamsGet)
}

// Post registers a new POST route.
//
// Parameters:
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//
// Example:
//
//	g.Post("/users", createUserHandler)
func (g *Group) Post(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodPost, pattern, handlerFunc, extractParamsPost)
}

// Put registers a new PUT route.
//
// Parameters:
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//
// Example:
//
//	g.Put("/users/:id", updateUserHandler)
func (g *Group) Put(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodPut, pattern, handlerFunc, extractParamsPut)
}

// Delete registers a new DELETE route.
//
// Parameters:
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//
// Example:
//
//	g.Delete("/users/:id", deleteUserHandler)
func (g *Group) Delete(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodDelete, pattern, handlerFunc, extractParamsDelete)
}

// Patch registers a new PATCH route.
//
// Parameters:
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//
// Example:
//
//	g.Patch("/users/:id", partialUpdateHandler)
func (g *Group) Patch(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodPatch, pattern, handlerFunc, extractParamsPatch)
}

// Options registers a new OPTIONS route.
//
// Parameters:
//   - pattern: The route pattern.
//   - handlerFunc: The function handling the request.
//
// Example:
//
//	g.Options("/users", optionsHandler)
func (g *Group) Options(pattern string, handlerFunc HandleFunc) {
	g.Handle(http.MethodOptions, pattern, handlerFunc, extractParamsOptions)
}
