/*
🚀 Quick is a flexible and extensible route manager for the Go language.
It aims to be fast and performant, and 100% net/http compatible.
Quick is a project under constant development and is open for collaboration,
everyone is welcome to contribute. 😍
*/
package quick

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"encoding/xml"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jeffotoni/quick/internal/concat"
)

// show in console
// Run Server Quick:0.0.0.0:<PORT>
var PRINT_SERVER = os.Getenv("PRINT_SERVER")

const (
	ContentTypeAppJSON = `application/json`
	ContentTypeAppXML  = `application/xml`
	ContentTypeTextXML = `text/xml`
)

type contextKey int

const myContextKey contextKey = 0

type HandleFunc func(*Ctx) error

// Route represents a registered HTTP route in the Quick framework
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
	embedFS       embed.FS
}

// GetDefaultConfig Function is responsible for returning a default configuration that is pre-defined for the system
// The result will be GetDefaultConfig() Config
func GetDefaultConfig() Config {
	return defaultConfig
}

// New function creates a new instance of the Quick structure to manage HTTP routes and handlers
// The result will New(c ...Config) *Quick
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

// Use function adds middleware to the Quick server, with special treatment for CORS
// The result will Use(mw any, nf ...string)
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

func clearRegex(route string) string {
	// Here you transform "/v1/user/{id:[0-9]+}"
	// into something simple, like "/v1/user/_id_"
	// You can get more sophisticated if you want
	var re = regexp.MustCompile(`\{[^/]+\}`)
	return re.ReplaceAllStringFunc(route, func(s string) string {
		// s is "{id:[0-9]+}"
		// Let's just replace it with "_id_"
		// or any string that doesn't contain ":" or "{ }"
		return "_" + strings.Trim(s, "{}") + "_"
		//return "_" + strings.ReplaceAll(strings.ReplaceAll(strings.Trim(s, "{}"), ":", "_"), "[", "_") + "_"
	})
}

// registerRoute is a helper function to centralize route registration logic.
func (q *Quick) registerRoute(method, pattern string, handlerFunc HandleFunc) {
	path, params, patternExist := extractParamsPattern(pattern)
	// formattedPath := concat.String(strings.ToLower(method)+"#", pattern)
	formattedPath := strings.ToLower(method) + "#" + clearRegex(pattern)
	route := Route{
		Pattern: patternExist,
		Path:    path,
		Params:  params,
		handler: extractHandler(q, method, path, params, handlerFunc),
		Method:  method,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(formattedPath, route.handler)
}

// Get function is an HTTP route with the GET method on the Quick server
// The result will Get(pattern string, handlerFunc HandleFunc)
func (q *Quick) Get(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodGet, pattern, handlerFunc)
}

// Post function registers an HTTP route with the POST method on the Quick server
// The result will Post(pattern string, handlerFunc HandleFunc)
func (q *Quick) Post(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodPost, pattern, handlerFunc)
}

// Put function registers an HTTP route with the PUT method on the Quick server.
// The result will Put(pattern string, handlerFunc HandleFunc)
func (q *Quick) Put(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodPut, pattern, handlerFunc)
}

// Delete function registers an HTTP route with the DELETE method on the Quick server.
// The result will Delete(pattern string, handlerFunc HandleFunc)
func (q *Quick) Delete(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodDelete, pattern, handlerFunc)
}

// Path function registers an HTTP route with the PATH method on the Quick server.
// The result will Path(pattern string, handlerFunc HandleFunc)
func (q *Quick) Patch(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodPatch, pattern, handlerFunc)
}

// Options function registers an HTTP route with the Options method on the Quick server.
// The result will Options(pattern string, handlerFunc HandleFunc)
func (q *Quick) Options(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodOptions, pattern, handlerFunc)
}

// Generic handler extractor to minimize repeated logic across HTTP methods
// The result will extractHandler(q *Quick, method, path, params string, handlerFunc HandleFunc) http.HandlerFunc
func extractHandler(q *Quick, method, path, params string, handlerFunc HandleFunc) http.HandlerFunc {
	switch method {
	case MethodGet:
		return extractParamsGet(q, path, params, handlerFunc)
	case MethodPost:
		return extractParamsPost(q, handlerFunc)
	case MethodPut:
		return extractParamsPut(q, handlerFunc)
	case MethodDelete:
		return extractParamsDelete(q, handlerFunc)
	case MethodPatch:
		return extractParamsPatch(q, handlerFunc) // same as PUT
	case MethodOptions:
		return extractParamsOptions(q, method, path, handlerFunc)
	}
	return nil
}

// PATCH is generally used for partial updates, while PUT replaces the entire resource.
// However, both methods often handle request parameters and body parsing in the same way.
func extractParamsPatch(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return extractParamsPut(q, handlerFunc)
}

// extractParamsOptions processes an HTTP request for a dynamic route, extracting query parameters, headers, and handling the request using the provided handler function
// The result will extractParamsOptions(q *Quick, method, path string, handlerFunc HandleFunc) http.HandlerFunc
func extractParamsOptions(q *Quick, method, path string, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define allowed HTTP methods
		allowMethods := []string{MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch, MethodOptions}
		w.Header().Set("Allow", strings.Join(allowMethods, ", "))

		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If a handlerFunc exists, execute it
		if handlerFunc != nil {
			err := handlerFunc(&Ctx{Response: w, Request: r})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusNoContent) // Use 204 if no body is needed
		}
	}
}

// extractHeaders extracts all headers from an HTTP request and returns them
// The result will extractHeaders(req http.Request) map[string][]string
func extractHeaders(req http.Request) map[string][]string {
	headersMap := make(map[string][]string)
	for key, values := range req.Header {
		headersMap[key] = values
	}
	return headersMap
}

// extractBind decodes the request body into the provided interface based on the Content-Type header
// The result will extractBind(c *Ctx, v interface{}) (err error)
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

// extractParamsPattern extracts the fixed path and dynamic parameters from a given route pattern
// The result will extractParamsPattern(pattern string) (path, params, partternExist string)
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

// extractParamsGet processes an HTTP request for a dynamic GET route, extracting query parameters, headers, and handling the request using the provided handler function
// The result will extractParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc
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
			Response:     w,
			Request:      req,
			Params:       cval.ParamsMap,
			Query:        querys,
			Headers:      headersMap,
			MoreRequests: q.config.MoreRequests,
		}
		execHandleFunc(c, handlerFunc)
	}
}

// extractParamsPost processes an HTTP POST request, extracting the request body and headers and handling the request using the provided handler function
// The result will extractParamsPost(q *Quick, handlerFunc HandleFunc) http.HandlerFunc
func extractParamsPost(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Check if body size exceeds limit before further validations
		if req.ContentLength > q.config.MaxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		headersMap := extractHeaders(*req)
		bodyBytes, bodyReader := extractBodyBytes(req.Body)

		c := &Ctx{
			Response:     w,
			Request:      req,
			bodyByte:     bodyBytes,
			Headers:      headersMap,
			MoreRequests: q.config.MoreRequests,
		}

		// reset `Request.Body` with `bodyReader`
		c.Request.Body = bodyReader
		execHandleFunc(c, handlerFunc)
	}
}

// extractParamsPut processes an HTTP PUT request, extracting request parameters, headers and request body before executing the provided handler function
// The result will extractParamsPut(q *Quick, handlerFunc HandleFunc) http.HandlerFunc
func extractParamsPut(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength > q.config.MaxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		headersMap := extractHeaders(*req)
		cval := v.(ctxServeHttp)
		bodyBytes, bodyReader := extractBodyBytes(req.Body)

		c := &Ctx{
			Response:     w,
			Request:      req,
			Headers:      headersMap,
			bodyByte:     bodyBytes,
			Params:       cval.ParamsMap,
			MoreRequests: q.config.MoreRequests,
		}

		// reset `Request.Body` with `bodyReader`
		c.Request.Body = bodyReader
		execHandleFunc(c, handlerFunc)
	}
}

// extractParamsDelete processes an HTTP DELETE request, extracting request parameters and headers before executing the provided handler function
// The result will extractParamsDelete(q *Quick, handlerFunc HandleFunc) http.HandlerFunc
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

// execHandleFunc executes the provided handler function and handles errors if they occur
// The result will execHandleFunc(c *Ctx, handleFunc HandleFunc)
func execHandleFunc(c *Ctx, handleFunc HandleFunc) {
	err := handleFunc(c)
	if err != nil {
		c.Set("Content-Type", "text/plain; charset=utf-8")
		// #nosec G104
		c.Status(500).SendString(err.Error())
	}
}

// extractBodyBytes reads the request body and returns it as a byte slice
// The result will extractBodyBytes(r io.ReadCloser) []byte
func extractBodyBytes(r io.ReadCloser) ([]byte, io.ReadCloser) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, io.NopCloser(bytes.NewBuffer(nil))
	}
	return b, io.NopCloser(bytes.NewReader(b))
}

// mwWrapper applies all registered middlewares to an HTTP handler
// The result will mwWrapper(handler http.Handler) http.Handler
func (q *Quick) mwWrapper(handler http.Handler) http.Handler {
	for i := len(q.mws2) - 1; i >= 0; i-- {
		switch mw := q.mws2[i].(type) {
		case func(http.Handler) http.Handler:
			handler = mw(handler)
		case func(http.ResponseWriter, *http.Request, http.Handler):
			originalHandler := handler // Avoid infinite reassignment
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw(w, r, originalHandler)
			})
		}
	}
	return handler
}

// appendRoute registers a new route in the Quick router and applies middlewares
// The result will appendRoute(route *Route)
func (q *Quick) appendRoute(route *Route) {
	route.handler = q.mwWrapper(route.handler).ServeHTTP
	//q.routes = append(q.routes, *route)
	q.routes = append(q.routes, route)
}

// ServeHTTP is the main HTTP request dispatcher for the Quick router
// The result will ServeHTTP(w http.ResponseWriter, req *http.Request)
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

// createParamsAndValid create params map and check if the request URI and pattern URI are valid
// The result will createParamsAndValid(reqURI, patternURI string) (map[string]string, bool)
func createParamsAndValid(reqURI, patternURI string) (map[string]string, bool) {
	params := make(map[string]string)
	var builder strings.Builder

	reqURI = strings.TrimPrefix(reqURI, "/")
	patternURI = strings.TrimPrefix(patternURI, "/")

	reqSplit := strings.Split(reqURI, "/")
	patternSplit := strings.Split(patternURI, "/")
	if len(reqSplit) != len(patternSplit) {
		return nil, false
	}

	for i, seg := range patternSplit {
		reqSeg := reqSplit[i]

		switch {
		// Ex: :id => paramName = "id"
		case strings.HasPrefix(seg, ":"):
			paramName := seg[1:]
			if paramName == "" {
				return nil, false
			}
			params[paramName] = reqSeg
			builder.WriteString("/")
			builder.WriteString(reqSeg)

		// Ex: {id:[0-9]+}
		case strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}"):
			content := seg[1 : len(seg)-1]
			parts := strings.SplitN(content, ":", 2)
			// Check for name and regex
			if len(parts) != 2 || parts[0] == "" {
				return nil, false
			}
			paramName, regexPattern := parts[0], parts[1]

			rgx, err := regexp.Compile("^" + regexPattern + "$")
			if err != nil || !rgx.MatchString(reqSeg) {
				return nil, false
			}
			params[paramName] = reqSeg
			builder.WriteString("/")
			builder.WriteString(reqSeg)

		default:
			if seg != reqSeg {
				return nil, false
			}
			builder.WriteString("/")
			builder.WriteString(seg)
		}
	}

	//if "/"+reqURI != builder.String() {
	//	return nil, false
	//}

	return params, true
}

// GetRoute returns all registered routes in the Quick framework
// The result will GetRoute() []*Route
func (q *Quick) GetRoute() []*Route {
	return q.routes
}

// Static server files html, css, js etc
// Embed.FS allows you to include files directly into
// the binary during compilation, eliminating the need to load files
// from the file system at runtime. This means that
// static files (HTML, CSS, JS, images, etc.)
// are embedded into the executable.
// The result will Static(route string, dirOrFS any)
func (q *Quick) Static(route string, dirOrFS any) {
	route = strings.TrimSuffix(route, "/")

	var fileServer http.Handler

	// check of dirOrFS is a embed.FS
	switch v := dirOrFS.(type) {
	case string:
		fileServer = http.FileServer(http.Dir(v))
	case embed.FS:
		q.embedFS = v
		fileServer = http.FileServer(http.FS(v))
	default:
		panic("Static: invalid parameter, must be string or embed.FS")
	}

	q.mux.Handle(concat.String(route, "/"), http.StripPrefix(route, fileServer))
}

// execHandler wraps an HTTP handler with additional processing
// The result will execHandler(next http.Handler) http.Handler
func (q *Quick) execHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// corsHandler returns an HTTP handler that applies CORS settings
// The result will corsHandler() http.Handler
func (q *Quick) corsHandler() http.Handler {
	return q.CorsSet(q)
}

// httpServer creates and returns an HTTP server instance configured with Quick.
// The result will httpServer(addr string, handler ...http.Handler) *http.Server
func (q *Quick) httpServer(addr string, handler ...http.Handler) *http.Server {
	// Set the default handler
	var h http.Handler = q
	if len(handler) > 0 {
		h = q.execHandler(handler[0])
	} else if q.Cors {
		h = q.corsHandler()
	}

	// Returns a single http.Server struct without code repetition
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       q.config.ReadTimeout,
		WriteTimeout:      q.config.WriteTimeout,
		IdleTimeout:       q.config.IdleTimeout,
		ReadHeaderTimeout: q.config.ReadHeaderTimeout,
	}
}

// ListenWithShutdown starts the HTTP server and returns a shutdown function.
// ListenWithShutdown retorna server e shutdown, mas NÃO inicia o server
func (q *Quick) ListenWithShutdown(addr string, handler ...http.Handler) (*http.Server, func(), error) {
	if q.config.MoreRequests > 0 {
		debug.SetGCPercent(q.config.MoreRequests)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	server := q.httpServer(listener.Addr().String(), handler...)
	shutdownFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		listener.Close()
	}

	// Servidor inicia em background
	go func() {
		server.Serve(listener)
	}()

	return server, shutdownFunc, nil
}

// Listen chama ListenWithShutdown e bloqueia com `select{}`
func (q *Quick) Listen(addr string, handler ...http.Handler) error {
	_, shutdown, err := q.ListenWithShutdown(addr, handler...)
	if err != nil {
		return err
	}
	defer shutdown()

	// Bloqueia indefinidamente
	select {}
}
