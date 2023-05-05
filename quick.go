package quick

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jeffotoni/quick/internal/concat"
	p "github.com/jeffotoni/quick/internal/print"
)

const (
	ContentTypeAppJSON = `application/json`
	ContentTypeAppXML  = `application/xml`
	ContentTypeTextXML = `text/xml`
)

type contextKey int

const myContextKey contextKey = 0

type HandleFunc func(*Ctx) error

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
	BodyLimit         int64
	MaxBodySize       int64
	MaxHeaderBytes    int64
	RouteCapacity     int
	MoreRequests      int // 0 a 1000
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

var defaultConfig = Config{
	BodyLimit:      2 * 1024 * 1024,
	MaxBodySize:    2 * 1024 * 1024,
	MaxHeaderBytes: 1 * 1024 * 1024,
	RouteCapacity:  1000,
	MoreRequests:   290, // valor de equilibrio
	//ReadTimeout:  10 * time.Second,
	//WriteTimeout: 10 * time.Second,
	//IdleTimeout:       1 * time.Second,
	// ReadHeaderTimeout: time.Duration(3) * time.Second,
}

type Zeroth int

const (
	Zero Zeroth = 0
)

type Quick struct {
	config        Config
	Cors          bool
	groups        []Group
	handler       http.Handler
	mux           *http.ServeMux
	routes        []*Route
	routeCapacity int
	mws2          []any
	CorsSet       func(http.Handler) http.Handler
	CorsOptions   map[string]string
}

func New(c ...Config) *Quick {
	var config Config
	if len(c) > 0 {
		config = c[0]
	} else {
		config = defaultConfig
	}
	if config.RouteCapacity == 0 {
		config.RouteCapacity = 1000
	}

	return &Quick{
		routes:        make([]*Route, 0, config.RouteCapacity),
		routeCapacity: config.RouteCapacity,
		mux:           http.NewServeMux(),
		handler:       http.NewServeMux(),
		config:        config,
	}
}

func (q *Quick) Use(mw any, nf ...string) {
	if len(nf) > 0 {
		if strings.ToLower(nf[0]) == "cors" {
			switch mwc := mw.(type) {
			case func(http.Handler) http.Handler:
				if strings.ToLower(nf[0]) == "cors" {
					q.Cors = true
					q.CorsSet = mwc
				}
			}
		}
	}
	q.mws2 = append(q.mws2, mw)
}

func (q *Quick) Get(pattern string, handlerFunc HandleFunc) {
	path, params, partternExist := extractParamsPattern(pattern)

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: extractParamsGet(q, path, params, handlerFunc),
		Method:  MethodGet,
	}
	q.appendRoute(&route)
	q.mux.HandleFunc(path, route.handler)
}

func (q *Quick) Post(pattern string, handlerFunc HandleFunc) {
	_, params, partternExist := extractParamsPattern(pattern)
	pathPost := concat.String("post#", pattern)

	route := Route{
		Pattern: partternExist,
		Params:  params,
		Path:    pattern,
		handler: extractParamsPost(q, handlerFunc),
		Method:  MethodPost,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(pathPost, route.handler)
}

func (q *Quick) Put(pattern string, handlerFunc HandleFunc) {
	_, params, partternExist := extractParamsPattern(pattern)

	pathPut := concat.String("put#", pattern)

	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		handler: extractParamsPut(q, handlerFunc),
		Method:  MethodPut,
		Params:  params,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(pathPut, route.handler)
}

func (q *Quick) Delete(pattern string, handlerFunc HandleFunc) {
	_, params, partternExist := extractParamsPattern(pattern)
	pathDelete := concat.String("delete#", pattern)

	route := Route{
		Pattern: partternExist,
		Path:    pattern,
		Params:  params,
		handler: extractParamsDelete(q, handlerFunc),
		Method:  MethodDelete,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(pathDelete, route.handler)
}

func extractHeaders(req http.Request) map[string][]string {
	headersMap := make(map[string][]string)
	for key, values := range req.Header {
		headersMap[key] = values
	}
	return headersMap
}

func extractBind(c *Ctx, v interface{}) (err error) {
	var req http.Request = *c.Request
	if strings.ToLower(req.Header.Get("Content-Type")) == ContentTypeAppJSON ||
		strings.ToLower(req.Header.Get("Content-Type")) == "application/json; charset=utf-8" ||
		strings.ToLower(req.Header.Get("Content-Type")) == "application/json;charset=utf-8" {
		err = json.NewDecoder(bytes.NewReader(c.bodyByte)).Decode(v)
	} else if strings.ToLower(req.Header.Get("Content-Type")) == ContentTypeTextXML ||
		strings.ToLower(req.Header.Get("Content-Type")) == ContentTypeAppXML {
		err = xml.NewDecoder(bytes.NewReader(c.bodyByte)).Decode(v)
	}
	return err
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

func extractParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(myContextKey)
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
			//bodyByte: extractBodyBytes(req.Body),
			//bodyByte: extractBodyBytes(req.Body),
			Headers:      headersMap,
			MoreRequests: q.config.MoreRequests,
		}
		execHandleFunc(c, handlerFunc)
	}
}

func extractParamsPost(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(myContextKey)
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
			Response:     w,
			Request:      req,
			bodyByte:     extractBodyBytes(req.Body),
			Headers:      headersMap,
			MoreRequests: q.config.MoreRequests,
		}
		execHandleFunc(c, handlerFunc)
	}
}

func extractParamsPut(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(myContextKey)
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
			Response:     w,
			Request:      req,
			Headers:      headersMap,
			bodyByte:     extractBodyBytes(req.Body),
			Params:       cval.ParamsMap,
			MoreRequests: q.config.MoreRequests,
		}

		execHandleFunc(c, handlerFunc)
	}
}

func extractParamsDelete(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		headersMap := extractHeaders(*req)

		cval := v.(ctxServeHttp)

		c := &Ctx{
			Response:     w,
			Request:      req,
			Headers:      headersMap,
			Params:       cval.ParamsMap,
			MoreRequests: q.config.MoreRequests,
		}

		execHandleFunc(c, handlerFunc)
	}
}

func execHandleFunc(c *Ctx, handleFunc HandleFunc) {
	err := handleFunc(c)
	if err != nil {
		c.Set("Content-Type", "text/plain; charset=utf-8")
		// #nosec G104
		c.Status(500).SendString(err.Error())
	}
}

func extractBodyBytes(r io.ReadCloser) []byte {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	return b
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
	//q.routes = append(q.routes, *route)
	q.routes = append(q.routes, route)
}

func (q *Quick) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for i := 0; i < len(q.routes); i++ {
		var requestURI = req.URL.Path
		var patternUri = q.routes[i].Pattern

		if q.routes[i].Method != req.Method {
			continue
		}

		if len(patternUri) == 0 {
			patternUri = q.routes[i].Path
		}

		paramsMap, isValid := createParamsAndValid(requestURI, patternUri)

		if !isValid {
			continue
		}

		var c = ctxServeHttp{Path: requestURI, ParamsMap: paramsMap, Method: q.routes[i].Method}
		req = req.WithContext(context.WithValue(req.Context(), myContextKey, c))
		q.routes[i].handler(w, req)
		return
	}
	http.NotFound(w, req)
}

// createParamsAndValid: Create params map and check if the request URI and pattern URI are valid
func createParamsAndValid(reqURI, patternURI string) (map[string]string, bool) {
	params := make(map[string]string)
	var tmpPath string

	reqURISplt := strings.Split(reqURI, "/")
	patternURISplt := strings.Split(patternURI, "/")

	if len(reqURISplt) != len(patternURISplt) {
		return nil, false
	}

	for pttrn := 0; pttrn < len(patternURISplt); pttrn++ { // collecting params
		if strings.Contains(patternURISplt[pttrn], ":") {
			params[patternURISplt[pttrn][1:]] = reqURISplt[pttrn]
			tmpPath = concat.String(tmpPath, "/", reqURISplt[pttrn])
		} else if strings.Contains(patternURISplt[pttrn], "{") { // regex support
			regexPattern := patternURISplt[pttrn][1:]
			regexPattern = regexPattern[:len(regexPattern)-1]
			rgx := regexp.MustCompile(regexPattern)
			params[patternURISplt[pttrn]] = rgx.FindString(reqURISplt[pttrn])
			tmpPath = concat.String(tmpPath, "/", rgx.FindString(reqURISplt[pttrn]))
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

func (q *Quick) GetRoute() []*Route {
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

func (q *Quick) execHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (q *Quick) corsHandler() http.Handler {
	return q.CorsSet(q)
}

func (q *Quick) httpServer(addr string, handler ...http.Handler) *http.Server {
	var server *http.Server

	if len(handler) > 0 {
		//assume o nosso mux
		server = &http.Server{
			Addr:    addr,
			Handler: q.execHandler(handler[0]),
			// ReadTimeout:
			// WriteTimeout:
			// MaxHeaderBytes:
			// IdleTimeout:
			ReadHeaderTimeout: q.config.ReadHeaderTimeout,
		}

	} else if q.Cors {
		server = &http.Server{
			Addr:    addr,
			Handler: q.corsHandler(),
			// ReadTimeout:
			// WriteTimeout:
			// MaxHeaderBytes:
			// IdleTimeout:
			ReadHeaderTimeout: q.config.ReadHeaderTimeout,
		}
	} else {
		server = &http.Server{
			Addr:    addr,
			Handler: q,
			// ReadTimeout:
			// WriteTimeout:
			// MaxHeaderBytes:
			// IdleTimeout:
			ReadHeaderTimeout: q.config.ReadHeaderTimeout,
		}
	}
	return server
}

func (q *Quick) Listen(addr string, handler ...http.Handler) error {

	if q.config.MoreRequests > 0 {
		debug.SetGCPercent(q.config.MoreRequests)
	}

	server := q.httpServer(addr, handler...)
	p.Stdout("\033[0;33mRun Server Quick:", addr, "\033[0m\n")
	return server.ListenAndServe()
}
