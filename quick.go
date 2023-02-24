package quick

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Ctx struct {
	Response http.ResponseWriter
	Request  *http.Request
	Headers  map[string][]string
	Params   map[string]string
	Query    map[string]string
	JSON     map[string]interface{}
	BodyByte []byte
	JsonStr  string
}

type Route struct {
	//Pattern *regexp.Regexp
	Pattern string
	Path    string
	Params  string
	handler http.HandlerFunc
	Method  string
}

type ctxServeHttp struct {
	Path      string
	Params    string
	ParamsMap map[string]string
	Method    string
}

type Config struct {
	MaxBodySize       int64
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxHeaderBytes    int64
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

var defaultConfig = Config{
	MaxBodySize:    3 * 1024 * 1024,
	MaxHeaderBytes: 1 * 1024 * 1024,
	//ReadTimeout:  10 * time.Second,
	//WriteTimeout: 10 * time.Second,
	//IdleTimeout:       1 * time.Second,
	//ReadHeaderTimeout: 3 * time.Second,
}

type Group struct {
	prefix string
	routes []Route
	quick  *Quick
}

type Quick struct {
	routes  []Route
	groups  []Group
	mws     []func(http.Handler) http.Handler
	mux     *http.ServeMux
	handler http.Handler
	config  Config
	//groupp  string
}

func New(c ...Config) *Quick {
	var config Config
	if len(c) > 0 {
		config = c[0]
	} else {
		config = defaultConfig
	}

	return &Quick{
		mux:     http.NewServeMux(),
		handler: http.NewServeMux(),
		config:  config,
	}
}

func (q *Quick) Use(mw func(http.Handler) http.Handler) {
	q.mws = append(q.mws, mw)
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

func (q *Quick) Post(pattern string, handlerFunc func(*Ctx)) {
	pathPost := ConcatStr("post#", pattern)
	route := Route{
		Pattern: "",
		Path:    pattern,
		handler: extractParamsPost(q, pattern, handlerFunc),
		Method:  http.MethodPost,
	}

	q.routes = append(q.routes, route)
	q.registerHandler(pathPost, route.handler)
}

func extractHeaders(req http.Request) map[string][]string {
	headersMap := make(map[string][]string)
	for key, values := range req.Header {
		headersMap[key] = values
	}
	return headersMap
}

func extractBodyByte(req http.Request) ([]byte, error) {
	var bodyByte []byte
	var err error

	if req.Header.Get("Content-Type") == "application/json" {
		bodyByte, err = io.ReadAll(req.Body)
	}
	return bodyByte, err
}

func extractParamsPattern(pattern string) (path, params, partternExist string) {
	path = pattern
	index := strings.Index(pattern, ":")
	if index > 0 {
		path = pattern[:index]
		path = strings.TrimSuffix(path, "/")
		params = strings.TrimPrefix(pattern, path)
		partternExist = pattern
	}

	return
}

func extractParamsPost(q *Quick, pathTmp string, handlerFunc func(*Ctx)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(0)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		if req.ContentLength > q.config.MaxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		headersMap := extractHeaders(*req)

		bodyByte, err := extractBodyByte(*req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c := &Ctx{
			Response: w,
			Request:  req,
			Headers:  headersMap,
			BodyByte: bodyByte,
		}
		handlerFunc(c)
	}
}

func extractParamsPut(q *Quick, pathTmp string, handlerFunc func(*Ctx)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(0)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		if req.ContentLength > q.config.MaxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		headersMap := extractHeaders(*req)

		cval := v.(ctxServeHttp)

		bodyByte, err := extractBodyByte(*req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c := &Ctx{
			Response: w,
			Request:  req,
			Headers:  headersMap,
			BodyByte: bodyByte,
			Params:   cval.ParamsMap,
		}
		handlerFunc(c)
	}
}

func (c *Ctx) Param(key string) string {
	val, ok := c.Params[key]
	if ok {
		return val
	}
	return ""
}

func (q *Quick) registerHandler(pattern string, handler http.Handler) {
	for i := range q.mws {
		handler = q.mws[i](handler)
	}

	q.mux.Handle(pattern, handler)
}

func (c *Ctx) Body(v interface{}) (err error) {
	if c.Request.Header.Get("Content-Type") == "application/json" {
		if len(c.BodyByte) > 0 {
			err = json.Unmarshal(c.BodyByte, v)
			if err != nil {
				return
			}
			c.JsonStr = string(c.BodyByte)
		}
	}
	return nil
}

func (c *Ctx) BodyString() string {
	return c.JsonStr
}

func (g *Group) Get(pattern string, handlerFunc func(*Ctx)) {
	pattern = ConcatStr(g.prefix, pattern)
	path, params, partternExist := extractParamsPattern(pattern)

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: extractParamsGet(path, params, handlerFunc),
		Method:  http.MethodGet,
	}

	g.quick.routes = append(g.quick.routes, route)
	g.quick.registerHandler(path, route.handler)
}

func (g *Group) Post(pattern string, handlerFunc func(*Ctx)) {
	pattern = ConcatStr(g.prefix, pattern)
	_, params, partternExist := extractParamsPattern(pattern)
	pathPost := ConcatStr("post#", pattern)
	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: extractParamsPost(g.quick, pattern, handlerFunc),
		Method:  http.MethodPost,
		Params:  params,
	}

	g.quick.routes = append(g.quick.routes, route)
	g.quick.registerHandler(pathPost, route.handler)
}

func (q *Quick) Get(pattern string, handlerFunc func(*Ctx)) {
	path, params, partternExist := extractParamsPattern(pattern)

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: extractParamsGet(path, params, handlerFunc),
		Method:  http.MethodGet,
	}

	q.routes = append(q.routes, route)
	q.registerHandler(path, route.handler)
}

func (q *Quick) Put(pattern string, handlerFunc func(*Ctx)) {
	_, params, partternExist := extractParamsPattern(pattern)
	pathPut := ConcatStr("put#", pattern)
	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: extractParamsPost(q, pattern, handlerFunc),
		Method:  http.MethodPut,
		Params:  params,
	}

	q.routes = append(q.routes, route)
	q.registerHandler(pathPut, route.handler)
}

func extractParamsGet(pathTmp, paramsPath string, handlerFunc func(*Ctx)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(0)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		cval := v.(ctxServeHttp)
		querys := make(map[string]string)
		queryParams := req.URL.Query()
		for key, values := range queryParams {
			querys[key] = values[0]
		}
		headersMap := extractHeaders(*req)

		c := &Ctx{
			Response: w,
			Request:  req,
			Params:   cval.ParamsMap,
			Query:    querys,
			Headers:  headersMap,
		}
		handlerFunc(c)
	}
}

func (q *Quick) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range q.routes {
		var pathTmp string = req.URL.Path
		var paramsMap = make(map[string]string)
		var paramUriIndex = make([]int, 0)

		if len(route.Pattern) > 0 && route.Method == req.Method {
			var routeParams string
			ppurl := strings.Split(route.Pattern, "/")
			paramsCount := 0

			for i := 0; i < len(ppurl); i++ {
				if strings.Contains(ppurl[i], ":") {
					routeParams = routeParams + ppurl[i]
					paramUriIndex = append(paramUriIndex, i)
					paramsCount++
				}
			}

			reqParams := strings.Split(pathTmp, "/")
			pathParams := strings.Split(routeParams, ":")[1:]
			prefix := ConcatStr(route.Path, "/")
			if strings.HasPrefix(pathTmp, prefix) {
				newPath := pathTmp[len(prefix):]
				if len(newPath) > 0 {
					pathTmp = route.Path
				}
			}

			if route.Path == pathTmp && route.Method == req.Method {
				if len(pathParams) != paramsCount {
					continue
				}

				for i, p := range pathParams {
					paramsMap[p] = reqParams[paramUriIndex[i]]
				}

				var c = ctxServeHttp{Path: pathTmp, ParamsMap: paramsMap, Method: route.Method}
				req = req.WithContext(context.WithValue(req.Context(), 0, c))
				route.handler(w, req)
				return
			}
		}

		if route.Path == pathTmp && route.Method == req.Method {
			var c = ctxServeHttp{Path: pathTmp, Method: route.Method}
			req = req.WithContext(context.WithValue(req.Context(), 0, c))
			route.handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func (c *Ctx) Json(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", "application/json")
	_, err = c.Response.Write(b)
	return err
}

func (c *Ctx) Byte(b []byte) (err error) {
	_, err = c.Response.Write(b)
	return err
}

func (c *Ctx) SendString(s string) error {
	_, err := c.Response.Write([]byte(s))
	return err
}

func (c *Ctx) SendFile(file []byte) error {
	_, err := c.Response.Write(file)
	return err
}

func (c *Ctx) Set(key, value string) {
	c.Response.Header().Set(key, value)
}

func (c *Ctx) Accepts(acceptType string) *Ctx {
	c.Response.Header().Set("Accept", acceptType)
	return c
}

func (c *Ctx) Status(status int) *Ctx {
	c.Response.WriteHeader(status)
	return c
}

func (q *Quick) GetRoute() []Route {
	return q.routes
}

func (q *Quick) Static(staticFolder string) {
	// generate route get with a pattern like this: /static/:file
	pattern := ConcatStr(staticFolder, ":file")
	q.Get(pattern, func(c *Ctx) {
		path, _, _ := extractParamsPattern(pattern)
		file := c.Params["file"]
		filePath := ConcatStr(".", path, "/", file)

		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			c.Status(http.StatusNotFound).SendString("File Not Found")
		}

		c.Status(http.StatusOK).SendFile(fileBytes)
	})
}

func (q *Quick) Listen(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: q,
		// ReadTimeout:
		// WriteTimeout:
		// MaxHeaderBytes:
		// IdleTimeout:
		// ReadHeaderTimeout:
	}

	Print("\033[0;33mRun Server Quick:", addr, "\033[0m")
	return server.ListenAndServe()
}
