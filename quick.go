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
	"crypto/tls"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
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
	MaxHeaderBytes    int
	RouteCapacity     int
	MoreRequests      int // 0 a 1000
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	GCPercent         int         // Renamed to be more descriptive (0-1000)
	TLSConfig         *tls.Config // Integrated TLS configuration
	CorsConfig        *CorsConfig // Specific type for CORS
}

var defaultConfig = Config{
	BodyLimit:      2 * 1024 * 1024,
	MaxBodySize:    2 * 1024 * 1024,
	MaxHeaderBytes: 1 * 1024 * 1024,
	RouteCapacity:  1000,
	MoreRequests:   290, // equilibrium value
	// TLSConfig:      &tls.Config{},
	//ReadTimeout:  10 * time.Second,
	//WriteTimeout: 10 * time.Second,
	//IdleTimeout:       1 * time.Second,
	// ReadHeaderTimeout: time.Duration(3) * time.Second,
}

type Zeroth int

const (
	Zero Zeroth = 0
)

type CorsConfig struct {
	Enabled  bool
	Options  map[string]string
	AllowAll bool
}

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
	// corsConfig    *CorsConfig // Specific type for CORS
	embedFS embed.FS
	server  *http.Server
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

// Responsible for clearing the path to be accepted in
// Servemux receives something like get#/v1/user/_id:[0-9]+_, without {}
// The result will clearRegex(route string) string
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
// The result will registerRoute(method, pattern string, handlerFunc HandleFunc)
func (q *Quick) registerRoute(method, pattern string, handlerFunc HandleFunc) {
	path, params, patternExist := extractParamsPattern(pattern)
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

// Ctx Pool
var ctxPool = sync.Pool{
	New: func() interface{} {
		// Initialize a new Ctx with empty maps to avoid nil checks in usage.
		return &Ctx{
			Params:  make(map[string]string),
			Query:   make(map[string]string),
			Headers: make(map[string][]string),
		}
	},
}

// acquireCtx retrieves a Ctx instance from the sync.Pool.
func acquireCtx() *Ctx {
	return ctxPool.Get().(*Ctx)
}

// releaseCtx resets the Ctx fields and returns it to the sync.Pool for reuse.
func releaseCtx(ctx *Ctx) {
	// Clear or nil out the Ctx fields so we avoid data leaking between requests.
	ctx.Params = nil
	ctx.Query = nil
	ctx.Headers = nil
	ctx.Response = nil
	ctx.Request = nil
	ctx.bodyByte = nil
	ctx.MoreRequests = 0
	ctxPool.Put(ctx)
}

// extractParamsGet processes an HTTP request for a dynamic GET route,
// extracting query parameters, headers, and handling the request using
// the provided handler function.
func extractParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Acquire a pooled context at the beginning of the request.
		ctx := acquireCtx()
		// Ensure the context is released back to the pool once this handler finishes.
		defer releaseCtx(ctx)

		// Retrieve our custom context value from the request's context.
		v := req.Context().Value(myContextKey)
		if v == nil {
			// If there's no context value, respond with 404.
			http.NotFound(w, req)
			return
		}

		// Cast the interface{} to our internal type.
		cval := v.(ctxServeHttp)

		// Build a map of query parameters.
		// You could also reuse a pooled map if desired, but here we create a new one each time.
		queryParams := req.URL.Query()
		querys := make(map[string]string, len(queryParams))
		for key, values := range queryParams {
			// We only store the first value; adapt as needed for multi-value params.
			querys[key] = values[0]
		}

		// Extract headers into a map[string][]string
		headersMap := extractHeaders(*req)

		// Populate the Ctx object with current request data.
		ctx.Response = w
		ctx.Request = req
		ctx.Params = cval.ParamsMap
		ctx.Query = querys
		ctx.Headers = headersMap
		ctx.MoreRequests = q.config.MoreRequests

		// Finally, execute the provided handler function with this context.
		execHandleFunc(ctx, handlerFunc)
	}
}

// // extractParamsGet processes an HTTP request for a dynamic GET route, extracting query parameters, headers, and handling the request using the provided handler function
// // The result will extractParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc
// func extractParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		v := req.Context().Value(myContextKey)
// 		if v == nil {
// 			http.NotFound(w, req)
// 			return
// 		}

// 		cval := v.(ctxServeHttp)
// 		querys := make(map[string]string)
// 		queryParams := req.URL.Query()
// 		for key, values := range queryParams {
// 			querys[key] = values[0]
// 		}
// 		headersMap := extractHeaders(*req)

// 		c := &Ctx{
// 			Response:     w,
// 			Request:      req,
// 			Params:       cval.ParamsMap,
// 			Query:        querys,
// 			Headers:      headersMap,
// 			MoreRequests: q.config.MoreRequests,
// 		}
// 		execHandleFunc(c, handlerFunc)
// 	}
// }

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

var bufferPool = sync.Pool{
	// Create new buffers with an initial capacity of 4KB.
	// Adjust this size based on expected request body sizes.
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

// acquireBuffer retrieves a *bytes.Buffer from the sync.Pool.
func acquireBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// releaseBuffer resets the buffer and places it back in the sync.Pool for reuse.
func releaseBuffer(buf *bytes.Buffer) {
	buf.Reset() // Clear any existing data
	bufferPool.Put(buf)
}

// extractBodyBytes reads the entire request body into a pooled buffer, then
// copies the data to a new byte slice before returning it. This ensures that
// once the buffer is returned to the pool, the returned data remains valid.
// The function also returns a new io.ReadCloser wrapping that same data,
// allowing it to be re-read if needed.
//
// Note: If the request body is very large, the buffer will grow automatically
// and remain larger when placed back in the pool. If extremely large bodies
// are expected infrequently, you may want additional logic to discard overly
// large buffers rather than returning them to the pool.
func extractBodyBytes(r io.ReadCloser) ([]byte, io.ReadCloser) {
	// Acquire a reusable buffer from the pool
	buf := acquireBuffer()
	defer releaseBuffer(buf)

	// Read all data from the request body into the buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		// If there's an error, return an empty NopCloser
		// so downstream logic can handle gracefully.
		return nil, io.NopCloser(bytes.NewBuffer(nil))
	}

	// Copy the data from the buffer into a separate byte slice.
	// This step is crucial because once the buffer is released
	// back to the pool, its underlying memory can be reused.
	data := make([]byte, buf.Len())
	copy(data, buf.Bytes())

	// Return both the raw byte slice and a new ReadCloser
	// wrapping the same data, which allows for re-reading.
	return data, io.NopCloser(bytes.NewReader(data))
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

// ///// pool connection
type pooledResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

var responseWriterPool = sync.Pool{
	New: func() interface{} {
		return &pooledResponseWriter{
			buf: bytes.NewBuffer(nil),
		}
	},
}

func acquireResponseWriter(w http.ResponseWriter) *pooledResponseWriter {
	rw := responseWriterPool.Get().(*pooledResponseWriter)
	rw.ResponseWriter = w
	return rw
}

func releaseResponseWriter(rw *pooledResponseWriter) {
	rw.buf.Reset()
	rw.ResponseWriter = nil
	responseWriterPool.Put(rw)
}

//////// pool connection

// ServeHTTP is the main HTTP request dispatcher for the Quick router
// The result will ServeHTTP(w http.ResponseWriter, req *http.Request)
func (q *Quick) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//  rw := acquireResponseWriter(w)
	// defer releaseResponseWriter(rw)

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

// corsHandler returns an HTTP handler that applies the configured CORS settings.
// Internally, it uses q.CorsSet(q) to wrap the Quick router with CORS middleware
// if the feature is enabled.
func (q *Quick) corsHandler() http.Handler {
	return q.CorsSet(q)
}

// httpServerTLS creates and returns an HTTP server instance configured with Quick
// for TLS/HTTPS usage. This function accepts a tlsConfig for secure connections.
//
// Parameters:
//   - addr:      The network address the server should listen on (e.g., ":443").
//   - tlsConfig: A *tls.Config instance containing certificate and security settings.
//   - handler:   Optionally, one or more custom HTTP handlers.
//
// If no custom handler is provided, the default Quick router is used by default.
// If q.Cors is enabled, the returned handler includes CORS middleware.
func (q *Quick) httpServerTLS(addr string, tlsConfig *tls.Config, handler ...http.Handler) *http.Server {
	// Determine the handler to use based on optional arguments and CORS configuration.
	var h http.Handler = q
	if len(handler) > 0 {
		h = q.execHandler(handler[0])
	} else if q.Cors {
		h = q.corsHandler()
	}

	// Return a fully configured http.Server, including TLS settings.
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		TLSConfig:         tlsConfig,
		ReadTimeout:       q.config.ReadTimeout,
		WriteTimeout:      q.config.WriteTimeout,
		IdleTimeout:       q.config.IdleTimeout,
		ReadHeaderTimeout: q.config.ReadHeaderTimeout,
		MaxHeaderBytes:    q.config.MaxHeaderBytes,
	}
}

// httpServer creates and returns an HTTP server instance configured with Quick
// for plain HTTP (non-TLS) usage.
//
// Parameters:
//   - addr:    The network address the server should listen on (e.g., ":8080").
//   - handler: Optionally, one or more custom HTTP handlers.
//
// If no custom handler is provided, the default Quick router is used by default.
// If q.Cors is enabled, the returned handler includes CORS middleware.
func (q *Quick) httpServer(addr string, handler ...http.Handler) *http.Server {
	// Determine the handler to use based on optional arguments and CORS configuration.
	var h http.Handler = q
	if len(handler) > 0 {
		h = q.execHandler(handler[0])
	} else if q.Cors {
		h = q.corsHandler()
	}

	// Return a fully configured http.Server for plain HTTP usage.
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       q.config.ReadTimeout,
		WriteTimeout:      q.config.WriteTimeout,
		IdleTimeout:       q.config.IdleTimeout,
		ReadHeaderTimeout: q.config.ReadHeaderTimeout,
		MaxHeaderBytes:    q.config.MaxHeaderBytes,
	}
}

// ListenWithShutdown starts an HTTP server on the given address and returns both
// the *http.Server instance and a shutdown function. The server begins listening
// in a background goroutine, while the caller retains control over the server's
// lifecycle through the returned function.
//
// If q.config.MoreRequests > 0, the function sets the Go garbage collector's
// percentage via debug.SetGCPercent(q.config.MoreRequests). This can tune how
// aggressively the garbage collector operates.
//
// The returned shutdown function, when called, attempts a graceful shutdown with
// a 5-second timeout. This allows existing connections to finish processing
// before the server is closed. It also closes the underlying TCP listener.
//
// Parameters:
//   - addr:    The TCP network address to listen on (e.g. ":8080").
//   - handler: Optional HTTP handlers; if none is provided, the default handler is used.
//
// Returns:
//   - *http.Server: The configured HTTP server instance.
//   - func():        A function that triggers graceful shutdown of the server.
//   - error:         An error if the listener cannot be created.
//
// Usage Example:
//
//	server, shutdown, err := q.ListenWithShutdown(":8080")
//	if err != nil {
//	    log.Fatalf("failed to start server: %v", err)
//	}
//	// The server is now running in the background.
//
//	// At a later point, we can shut down gracefully:
//	shutdown()
func (q *Quick) ListenWithShutdown(addr string, handler ...http.Handler) (*http.Server, func(), error) {
	if q.config.MoreRequests > 0 {
		debug.SetGCPercent(q.config.MoreRequests)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	server := q.httpServer(listener.Addr().String(), handler...)

	// This function initiates a graceful shutdown of the server.
	shutdownFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		listener.Close()
	}

	// Start the server in the background goroutine.
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

// ListenTLS starts an HTTPS server on the specified address using the provided
// certificate and key files. It allows enabling or disabling HTTP/2 support.
// It also configures basic modern TLS settings, sets up a listener with
// SO_REUSEPORT (when possible), and applies a graceful shutdown procedure.
//
// Parameters:
//   - addr: the TCP network address to listen on (e.g., ":443")
//   - certFile: the path to the SSL certificate file
//   - keyFile: the path to the SSL private key file
//   - useHTTP2: whether or not to enable HTTP/2
//   - handler: optional HTTP handlers. If none is provided, the default handler is used.
//
// Returns:
//   - error: an error if something goes wrong creating the listener or starting the server.
func (q *Quick) ListenTLS(addr, certFile, keyFile string, useHTTP2 bool, handler ...http.Handler) error {
	// If the user has specified a custom GC percentage (> 0),
	// set it here to help control garbage collection aggressiveness.
	if q.config.GCPercent > 0 {
		debug.SetGCPercent(q.config.GCPercent)
	}

	// Extract or create a TLS configuration.
	// If q.config.TLSConfig is nil, set up a default TLS config with modern protocols
	// and ciphers. This includes TLS 1.3 and secure cipher suites.
	var tlsConfig = q.config.TLSConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			MinVersion:       tls.VersionTLS13, // Sets TLS 1.3 as the minimum version
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
			},
			NextProtos: []string{"h2", "http/1.1"}, // Default to supporting HTTP/2
		}
	}

	// Enable or disable HTTP/2 support based on the useHTTP2 parameter.
	if useHTTP2 {
		// HTTP/2 + HTTP/1.1
		tlsConfig.NextProtos = []string{"h2", "http/1.1"}
	} else {
		// Only HTTP/1.1
		tlsConfig.NextProtos = []string{"http/1.1"}
	}

	// Create a net.ListenConfig that attempts to set SO_REUSEPORT on supported platforms.
	// This feature can improve load balancing by letting multiple processes
	// bind to the same address.
	cfg := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// Avoid setting SO_REUSEPORT on macOS to prevent errors.
				if runtime.GOOS != "darwin" {
					syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
				}
			})
		},
	}

	// Listen on the specified TCP address using our custom ListenConfig.
	listener, err := cfg.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	// Create the HTTP server configured for TLS using the provided or default tlsConfig.
	// The address is taken from the listener to ensure correctness in case the actual
	// bound port differs (for example, if you used ":0" for a random port).
	q.server = q.httpServerTLS(listener.Addr().String(), tlsConfig, handler...)

	// Start the server and perform a graceful shutdown when a termination signal is received.
	return q.startServerWithGracefulShutdown(listener, certFile, keyFile)
}

// startServerWithGracefulShutdown starts the HTTPS server (using the provided TLS certificate
// and private key) on the given listener and blocks until the server either encounters
// an unrecoverable error or receives a termination signal.
//
// The server runs in a goroutine so that this function can simultaneously listen for
// interrupt signals (SIGINT, SIGTERM, SIGHUP). Once such a signal is detected, the function
// will gracefully shut down the server, allowing any ongoing requests to finish or timing
// out after 15 seconds.
//
// Parameters:
//   - listener: A net.Listener that the server will use to accept connections.
//   - certFile: Path to the TLS certificate file.
//   - keyFile:  Path to the TLS private key file.
//
// Returns:
//   - error: An error if the server fails to start, or if a forced shutdown occurs.
//     Returns nil on normal shutdown.
func (q *Quick) startServerWithGracefulShutdown(listener net.Listener, certFile, keyFile string) error {
	serverErr := make(chan error, 1)

	// Run ServeTLS in a goroutine. Any unrecoverable error that isn't http.ErrServerClosed
	// is sent to the channel for handling in the main select block.
	go func() {
		if err := q.server.ServeTLS(listener, certFile, keyFile); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
		close(serverErr)
	}()

	// Create a context that listens for SIGINT, SIGTERM, and SIGHUP signals.
	// When one of these signals occurs, the context is canceled automatically.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	select {
	case <-ctx.Done():
		// We've received a termination signal, so attempt a graceful shutdown.
		log.Println("Received shutdown signal. Stopping server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// If the server cannot gracefully shut down within 15 seconds,
		// it will exit with an error.
		if err := q.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("forced shutdown: %w", err)
		}
		return nil

	case err := <-serverErr:
		// If an unrecoverable error occurred in ServeTLS, return it here.
		return err
	}
}

func (q *Quick) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if q.server != nil {
		return q.server.Shutdown(ctx)
	}
	return nil
}
