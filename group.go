package quick

import (
	"net/http"

	"github.com/jeffotoni/quick/internal/concat"
)

type Group struct {
	prefix string
	routes []Route
	quick  *Quick
}

func (q *Quick) Group(prefix string) *Group {
	g := &Group{
		prefix: prefix,
		routes: []Route{},
		quick:  q,
	}
	q.groups = append(q.groups, *g)
	return g
}

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
