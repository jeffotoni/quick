package quick

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gojeffotoni/quick/internal/concat"
	"github.com/gojeffotoni/quick/internal/print"
)

type Ctx struct {
	Response  http.ResponseWriter
	Request   *http.Request
	resStatus int
	bodyByte  []byte
	JsonStr   string
	Headers   map[string][]string
	Params    map[string]string
	Query     map[string]string
}

type Route struct {
	//Pattern *regexp.Regexp
	Group   string
	Pattern string
	Path    string
	Params  string
	Method  string
	handler http.HandlerFunc
}

type ctxServeHttp struct {
	Path      string
	Params    string
	Method    string
	ParamsMap map[string]string
}

type Config struct {
	MaxBodySize       int64
	MaxHeaderBytes    int64
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

var defaultConfig = Config{
	MaxBodySize:    3 * 1024 * 1024,
	MaxHeaderBytes: 1 * 1024 * 1024,
	//ReadTimeout:  10 * time.Second,
	//WriteTimeout: 10 * time.Second,
	//IdleTimeout:       1 * time.Second,
	ReadHeaderTimeout: time.Duration(3) * time.Second,
}

type Group struct {
	prefix string
	quick  *Quick
}

type Quick struct {
	config  Config
	cors    bool
	group   *Group
	handler http.Handler
	mux     *http.ServeMux
	routes  []Route
	mws2    []any
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

func (q *Quick) Use(mw any, nf ...string) {
	if len(nf) > 0 {
		if strings.ToLower(nf[0]) == "cors" {
			q.cors = true
		}
	}
	q.mws2 = append(q.mws2, mw)
}

func (q *Quick) Group(prefix string) {
	g := &Group{
		prefix: prefix,
		quick:  q,
	}
	q.group = g
}

func (q *Quick) Get(pattern string, handlerFunc func(*Ctx)) {
	path, params, partternExist := extractParamsPattern(pattern)

	var gr string
	// Setting up the group
	if q.group != nil {
		path = concat.String(q.group.prefix, path)
		gr = q.group.prefix
	}

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: extractParamsGet(path, params, handlerFunc),
		Method:  http.MethodGet,
		Group:   gr,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(path, route.handler)
}

func (q *Quick) Post(pattern string, handlerFunc func(*Ctx)) {

	var (
		gr       string
		pathPost string
	)

	pathPost = concat.String("post#", pattern)

	// Setting up the group
	if q.group != nil {
		pathPost = concat.String("post#", q.group.prefix, pattern)
		pattern = concat.String(q.group.prefix, pattern)
		gr = q.group.prefix
	}

	route := Route{
		Pattern: "",
		Path:    pattern,
		handler: extractParamsPost(q, pattern, handlerFunc),
		Method:  http.MethodPost,
		Group:   gr,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(pathPost, route.handler)
}

func (q *Quick) Put(pattern string, handlerFunc func(*Ctx)) {
	_, params, partternExist := extractParamsPattern(pattern)

	var (
		gr      string
		pathPut string
	)

	pathPut = concat.String("put#", pattern)

	// Setting up the group
	if q.group != nil {
		pathPut = concat.String("put#", q.group.prefix, pathPut)
		pattern = concat.String(q.group.prefix, pattern)
		gr = q.group.prefix
	}

	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: extractParamsPut(q, pattern, handlerFunc),
		Method:  http.MethodPut,
		Params:  params,
		Group:   gr,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(pathPut, route.handler)
}

func extractHeaders(req http.Request) map[string][]string {
	headersMap := make(map[string][]string)
	for key, values := range req.Header {
		headersMap[key] = values
	}
	return headersMap
}

func extractBind(c *Ctx, v interface{}) (obj interface{}, err error) {
	var req http.Request = *c.Request
	if strings.ToLower(req.Header.Get("Content-Type")) == "application/json" ||
		strings.ToLower(req.Header.Get("Content-Type")) == "application/json; charset=utf-8" ||
		strings.ToLower(req.Header.Get("Content-Type")) == "application/json;charset=utf-8" {
		err = json.NewDecoder(bytes.NewReader(c.bodyByte)).Decode(v)
		obj = v
	}
	return obj, err
}

func extractParamsPattern(pattern string) (path, params, partternExist string) {
	path = pattern
	index := strings.Index(pattern, ":")

	if index > 0 {
		path = pattern[:index]
		path = strings.TrimSuffix(path, "/")
		if index == 1 {
			path = "/"
		}
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

		c := &Ctx{
			Response: w,
			Request:  req,
			bodyByte: extractBodyBytes(req.Body),
			Headers:  headersMap,
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

		c := &Ctx{
			Response: w,
			Request:  req,
			Headers:  headersMap,
			bodyByte: extractBodyBytes(req.Body),
			Params:   cval.ParamsMap,
		}

		handlerFunc(c)
	}
}

func extractBodyBytes(r io.ReadCloser) []byte {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	return b
}

func (c *Ctx) Param(key string) string {
	val, ok := c.Params[key]
	if ok {
		return val
	}
	return ""
}

func (q *Quick) mwWrapper(handler http.Handler) http.Handler {
	for i := range q.mws2 {
		switch mw := q.mws2[i].(type) {
		case func(http.Handler) http.Handler:
			handler = mw(handler)
		case func(http.ResponseWriter, *http.Request, http.Handler):
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw(w, r, handler)
			})
		}
	}
	return handler
}

func (q *Quick) appendRoute(route *Route) {
	route.handler = q.mwWrapper(route.handler).ServeHTTP
	q.routes = append(q.routes, *route)
}

func (c *Ctx) Bind(v interface{}) (err error) {
	if c.Request.Header.Get("Content-Type") == "application/json" {
		obj, err := extractBind(c, v)
		if err != nil {
			return err
		}
		v = obj
	}

	return nil
}

func (c *Ctx) BodyParse(v interface{}) (err error) {
	if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
		err = json.Unmarshal(c.bodyByte, v)
		if err != nil {
			return err
		}
	}

	if strings.Contains(c.Request.Header.Get("Content-Type"), "text/xml") ||
		strings.Contains(c.Request.Header.Get("Content-Type"), "application/xml") {
		err = xml.Unmarshal(c.bodyByte, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Ctx) Body() []byte {
	return c.bodyByte
}

func (c *Ctx) BodyString() string {
	return string(c.bodyByte)
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
			bodyByte: extractBodyBytes(req.Body),
			Headers:  headersMap,
		}
		handlerFunc(c)
	}
}

func (q *Quick) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range q.routes {
		var requestURI = req.URL.Path
		var paramsMap = make(map[string]string)
		var patternUri = route.Pattern

		if route.Method != req.Method {
			continue
		}

		if len(patternUri) == 0 {
			patternUri = route.Path
		}

		if len(route.Group) != 0 && strings.Contains(route.Pattern, ":") {
			patternUri = concat.String(route.Group, route.Pattern)
		}

		paramsMap, isValid := createParamsAndValid(requestURI, patternUri)

		if !isValid {
			continue
		}

		var c = ctxServeHttp{Path: requestURI, ParamsMap: paramsMap, Method: route.Method}
		req = req.WithContext(context.WithValue(req.Context(), 0, c))
		route.handler(w, req)
		return
	}
	http.NotFound(w, req)
}

func createParamsAndValid(reqURI, patternURI string) (map[string]string, bool) {
	params := make(map[string]string)
	var tmpPath string

	reqURISplt := strings.Split(reqURI, "/")
	patternURISplt := strings.Split(patternURI, "/")

	if len(reqURISplt) != len(patternURISplt) {
		return nil, false
	}

	for pttrn := 0; pttrn < len(patternURISplt); pttrn++ {
		if strings.Contains(patternURISplt[pttrn], ":") {
			params[patternURISplt[pttrn]] = reqURISplt[pttrn]
			tmpPath = concat.String(tmpPath, "/", reqURISplt[pttrn])
		} else {
			tmpPath = concat.String(tmpPath, "/", patternURISplt[pttrn])
		}
	}

	// This tmpPath is to check if request's uri is the same that our pattern
	if tmpPath[1:] != reqURI {
		return nil, false
	}

	return params, true
}

func (c *Ctx) JSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", "application/json")
	return c.writeResponse(b)
}

func (c *Ctx) XML(v interface{}) error {
	b, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", "text/xml")
	return c.writeResponse(b)
}

func (c *Ctx) writeResponse(b []byte) error {
	if c.resStatus != 0 {
		c.Response.WriteHeader(c.resStatus)
	}
	_, err := c.Response.Write(b)
	return err
}

func (c *Ctx) Byte(b []byte) (err error) {
	return c.writeResponse(b)
}

func (c *Ctx) SendString(s string) error {
	return c.writeResponse([]byte(s))

}

func (c *Ctx) String(s string) error {
	return c.writeResponse([]byte(s))
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
	c.resStatus = status
	return c
}

func (q *Quick) GetRoute() []Route {
	return q.routes
}

func (q *Quick) Static(staticFolder string) {
	// generate route get with a pattern like this: /static/:file
	// pattern := concat.String(staticFolder, ":file")
	// q.Get(pattern, func(c *Ctx) {
	// 	path, _, _ := extractParamsPattern(pattern)
	// 	file := c.Params["file"]
	// 	filePath := concat.String(".", path, "/", file)
	// 	err := qos.FileExist(filePath)
	// 	if err != nil {
	// 		c.Status(http.StatusForbidden).SendString(err.Error())
	// 		return
	// 	}

	// 	fileBytes, err := os.ReadFile(filePath)
	// 	if err != nil {
	// 		c.Status(http.StatusInternalServerError).SendString(err.Error())
	// 		return
	// 	}
	// 	c.Status(http.StatusOK).SendFile(fileBytes)
	// })
}

func (q *Quick) Listen(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: q,
		// ReadTimeout:
		// WriteTimeout:
		// MaxHeaderBytes:
		// IdleTimeout:
		ReadHeaderTimeout: q.config.ReadHeaderTimeout,
	}

	print.Stdout("\033[0;33mRun Server Quick:", addr, "\033[0m")
	return server.ListenAndServe()
}
