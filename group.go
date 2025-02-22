package goquick

import (
	"net/http"

	"github.com/jeffotoni/goquick/internal/concat"
)

// Group represents a collection of routes that share a common prefix
type Group struct {
	prefix string
	routes []Route
	quick  *Quick
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
	pattern = concat.String(g.prefix, pattern)
	path, params, partternExist := extractParamsPattern(pattern)

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: extractParamsGet(g.quick, path, params, handlerFunc),
		Method:  http.MethodGet,
		Group:   g.prefix,
	}

	g.quick.appendRoute(&route)
	g.quick.mux.HandleFunc(path, route.handler)
}

// Post registers a new POST route within the group
// The result will Post(pattern string, handlerFunc HandleFunc)
func (g *Group) Post(pattern string, handlerFunc HandleFunc) {
	pattern = concat.String(g.prefix, pattern)
	_, params, partternExist := extractParamsPattern(pattern)

	pathPost := concat.String("post#", pattern)

	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: extractParamsPost(g.quick, handlerFunc),
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
	pattern = concat.String(g.prefix, pattern)
	_, params, partternExist := extractParamsPattern(pattern)

	pathPut := concat.String("put#", pattern)

	// Setting up the group
	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: extractParamsPut(g.quick, handlerFunc),
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
	pattern = concat.String(g.prefix, pattern)
	_, params, partternExist := extractParamsPattern(pattern)

	pathDelete := concat.String("delete#", pattern)

	// Setting up the group
	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		Params:  params,
		handler: extractParamsDelete(g.quick, handlerFunc),
		Method:  http.MethodDelete,
		Group:   g.prefix,
	}
	g.quick.appendRoute(&route)
	g.quick.mux.HandleFunc(pathDelete, route.handler)
}
