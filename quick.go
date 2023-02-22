package quick

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

type Quick struct {
	routes  []Route
	mws     []func(http.Handler) http.Handler
	mux     *http.ServeMux
	handler http.Handler
}

func New() *Quick {
	return &Quick{mux: http.NewServeMux(), handler: http.NewServeMux()}
}

func (q *Quick) Use(mw func(http.Handler) http.Handler) {
	q.mws = append(q.mws, mw)
}

func (r *Quick) Post(pattern string, handlerFunc func(*Ctx)) {
	pathPost := ConcatStr("post#", pattern)
	route := Route{
		Pattern: pattern,
		Path:    pattern,
		handler: extractParamsPost(pattern, handlerFunc),
		Method:  http.MethodPost,
	}
	r.routes = append(r.routes, route)
	r.mux.HandleFunc(pathPost, route.handler)
}

func extractParamsPost(pathTmp string, handlerFunc func(*Ctx)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		headersMap := make(map[string][]string)
		for key, values := range req.Header {
			headersMap[key] = values
		}

		var bodyByte []byte
		var err error
		if req.Header.Get("Content-Type") == "application/json" {
			bodyByte, err = io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
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

func (c *Ctx) Param(key string) string {
	val, ok := c.Params[key]
	if ok {
		return val
	}
	return ""
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

func (r *Quick) Get(pattern string, handlerFunc func(*Ctx)) {
	var path string = pattern
	var params string

	index := strings.Index(pattern, ":")
	if index > 0 {
		path = pattern[:index]
		path = strings.TrimSuffix(path, "/")
		params = strings.TrimPrefix(pattern, path)
	}

	route := Route{
		Pattern: pattern,
		Path:    path,
		Params:  params,
		handler: extractParamsGet(path, params, handlerFunc),
		Method:  http.MethodGet,
	}

	r.routes = append(r.routes, route)
	r.mux.HandleFunc(path, route.handler)
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
		headersMap := make(map[string][]string)
		for key, values := range req.Header {
			headersMap[key] = values
		}
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

	// for _, middleware := range q.middlewares {
	// 	middleware(w, req)
	// }

	for _, route := range q.routes {
		if route.Method != strings.ToUpper(req.Method) {
			continue
		}

		var pathTmp string = req.URL.Path
		routeParams := route.Params
		var paramsMap = make(map[string]string)

		if route.Method == "GET" {
			if !strings.HasPrefix(pathTmp, route.Path) {
				continue
			}

			if len(routeParams) > 0 {
				routeParams = strings.Replace(routeParams, "/", "", -1)
				tmppath := strings.Replace(req.URL.Path, route.Path, "", -1)
				ppurl := strings.Split(tmppath, "/")[1:]
				pathParams := strings.Split(routeParams, ":")[1:]
				if len(pathParams) != len(ppurl) {
					continue
				}
				pathTmp = route.Path
				for i, p := range pathParams {
					paramsMap[p] = ppurl[i]
				}
			}

			var c = ctxServeHttp{Path: pathTmp, Params: route.Params, ParamsMap: paramsMap, Method: route.Method}
			req = req.WithContext(context.WithValue(req.Context(), 0, c))
		}
		if route.Path != pathTmp {
			continue
		}

		route.handler(w, req)
		return
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

func (r *Quick) GetRoute() []Route {
	return r.routes
}

func (q *Quick) Listen(addr string) error {
	var handler http.Handler = q
	for i := len(q.mws) - 1; i >= 0; i-- {
		handler = q.mws[i](handler)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
		// ReadTimeout:
		// WriteTimeout:
		// MaxHeaderBytes:
		// IdleTimeout:
		// ReadHeaderTimeout:
	}
	println("\033[0;33mRun Server Quick:", addr, "\033[0m")
	return server.ListenAndServe()
}
