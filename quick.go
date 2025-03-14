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
	"net/http/httptest"
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

const SO_REUSEPORT = 0x0F // Manual definition for Linux

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
	BodyLimit      int64 // Deprecated: Use MaxBodySize instead
	MaxBodySize    int64 // Maximum request body size allowed.
	MaxHeaderBytes int   // Maximum number of bytes allowed in the HTTP headers.

	GOMAXPROCS      int   // defines the maximum number of CPU cores
	GCHeapThreshold int64 // GCHeapThreshold sets the memory threshold (in bytes)
	BufferPoolSize  int   // BufferPoolSize determines the size (in bytes)

	RouteCapacity     int           // Initial capacity of the route slice.
	MoreRequests      int           // Value to set GCPercent. influences the garbage collector performance. 0-1000
	ReadTimeout       time.Duration // Maximum duration for reading the entire request.
	WriteTimeout      time.Duration // Maximum duration before timing out writes of the response.
	IdleTimeout       time.Duration // Maximum amount of time to wait for the next request when keep-alives are enabled.
	ReadHeaderTimeout time.Duration // Amount of time allowed to read request headers.
	GCPercent         int           // Renamed to be more descriptive (0-1000) - influences the garbage collector performance.
	TLSConfig         *tls.Config   // Integrated TLS configuration
	CorsConfig        *CorsConfig   // Specific type for CORS

	NoBanner bool // Flag to disable the Quick startup Display.
}

var defaultConfig = Config{
	BodyLimit:      2 * 1024 * 1024, // 2MB
	MaxBodySize:    2 * 1024 * 1024, // 2MB
	MaxHeaderBytes: 1 * 1024 * 1024, // 1MB

	GOMAXPROCS:      runtime.NumCPU(),
	GCHeapThreshold: 1 << 30, // 1GB
	BufferPoolSize:  32768,

	RouteCapacity: 1000,  // Initial capacity of 1000 routes.
	MoreRequests:  290,   // default GC value equilibrium value
	NoBanner:      false, // Display Quick banner by default.
}

type Zeroth int

const (
	Zero Zeroth = 0
)

type CorsConfig struct {
	Enabled  bool              // Enable cors
	Options  map[string]string // Add custom options
	AllowAll bool              // Enable all access
}

// Quick is the main structure of the framework, holding routes and configurations.
type Quick struct {
	config        Config         // Configuration settings.
	Cors          bool           // Indicates if CORS is enabled.
	groups        []Group        // List of route groups.
	handler       http.Handler   // The primary HTTP handler.
	mux           *http.ServeMux // Multiplexer for routing requests.
	routes        []*Route       // Registered routes.
	routeCapacity int            // The maximum number of routes allowed.
	mws2          []any          // List of registered middlewares.

	CorsSet     func(http.Handler) http.Handler // CORS middleware handler function.
	CorsOptions map[string]string               // CORS options map
	// corsConfig    *CorsConfig // Specific type for CORS // Removed unused field
	embedFS embed.FS     // File system for embedded static files.
	server  *http.Server // Http server

	bufferPool *sync.Pool
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
// Method Used Internally
// The result will Use(mw any)
func (q *Quick) Use(mw any) {
	switch mwc := mw.(type) {
	case func(http.Handler) http.Handler:
		// Automatically detects if it is CORS
		if isCorsMiddleware(mwc) {
			q.Cors = true
			q.CorsSet = mwc
			return
		}
	}
	q.mws2 = append(q.mws2, mw)
}

// Helper function to automatically detect whether the middleware is CORS
func isCorsMiddleware(mw func(http.Handler) http.Handler) bool {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	testRequest := httptest.NewRequest("OPTIONS", "/", nil)
	testResponse := httptest.NewRecorder()

	mw(testHandler).ServeHTTP(testResponse, testRequest)

	// If the middleware sets Access-Control-Allow-Origin, it's CORS
	return testResponse.Header().Get("Access-Control-Allow-Origin") != ""
}

// Responsible for clearing the path to be accepted in
// Servemux receives something like get#/v1/user/_id:[0-9]+_, without {}
// Method Used Internally
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
// Method Used Internally
// The result will registerRoute(method, pattern string, handlerFunc HandleFunc)
func (q *Quick) registerRoute(method, pattern string, handlerFunc HandleFunc) {
	path, params, patternExist := extractParamsPattern(pattern)
	formattedPath := concat.String(strings.ToLower(method), "#", clearRegex(pattern))

	for _, route := range q.routes {
		if route.Method == method && route.Path == path {
			fmt.Printf("Warning: Route '%s %s' is already registered, ignoring duplicate registration.\n", method, path)
			return // Ignore duplication instead of generating panic
		}
	}

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

// handleOptions processes HTTP OPTIONS requests for CORS preflight checks.
// This function is automatically called before routing when an OPTIONS request is received.
// It ensures that the appropriate CORS headers are included in the response.
//
// If CORS middleware is enabled, it applies the middleware before setting default headers.
//
// Headers added by this function:
// - Access-Control-Allow-Origin: Allows cross-origin requests (set dynamically).
// - Access-Control-Allow-Methods: Specifies allowed HTTP methods (GET, POST, PUT, DELETE, OPTIONS).
// - Access-Control-Allow-Headers: Defines which headers are allowed in the request.
//
// If no Origin header is provided in the request, a 204 No Content response is returned.
//
// Parameters:
// - w: http.ResponseWriter – The response writer to send headers and status.
// - r: *http.Request – The incoming HTTP request.
//
// Response:
// - 204 No Content (if the request is valid and processed successfully).
//
// Example Usage:
// This function is automatically triggered in `ServeHTTP()` when an OPTIONS request is received.
func (q *Quick) handleOptions(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		w.WriteHeader(StatusNoContent)
		return
	}

	// Apply CORS middleware before setting headers
	if q.Cors && q.CorsSet != nil {
		q.CorsSet(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
	}

	// Set default CORS headers
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Ajustável pelo middleware
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.WriteHeader(http.StatusNoContent) // Returns 204 No Content
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
// Method Used Internally
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
// Method Used Internally
// However, both methods often handle request parameters and body parsing in the same way.
func extractParamsPatch(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return extractParamsPut(q, handlerFunc)
}

// extractParamsOptions processes an HTTP OPTIONS request, setting appropriate
// headers to handle CORS preflight requests. It reuses a pooled Ctx instance
// for optimized memory usage and performance.
//
// If a handlerFunc is provided, it executes that handler with the pooled context.
// If no handlerFunc is given, it responds with HTTP 204 No Content.
//
// Parameters:
//   - q: The Quick instance providing configurations and routing context.
//   - method: The HTTP method being handled (typically "OPTIONS").
//   - path: The route path being handled.
//   - handlerFunc: An optional handler to execute for the OPTIONS request.
//
// Returns:
//   - http.HandlerFunc: A handler function optimized for handling OPTIONS requests.
func extractParamsOptions(q *Quick, method, path string, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Acquire a pooled context
		ctx := acquireCtx()
		defer releaseCtx(ctx) // Ensure context is returned to the pool after handling

		// Populate the pooled context
		ctx.Response = w
		ctx.Request = r
		ctx.MoreRequests = q.config.MoreRequests

		if q.Cors && q.CorsSet != nil {
			wrappedHandler := q.CorsSet(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Middleware CORS apply success
			}))
			wrappedHandler.ServeHTTP(w, r)
		}

		if ctx.Response.Header().Get("Access-Control-Allow-Origin") == "" {
			allowMethods := []string{MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch, MethodOptions}
			ctx.Set("Allow", strings.Join(allowMethods, ", "))
			ctx.Set("Access-Control-Allow-Origin", "*")
			ctx.Set("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
			ctx.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// Execute handler function if provided
		if handlerFunc != nil {
			if err := handlerFunc(ctx); err != nil {
				http.Error(w, err.Error(), StatusInternalServerError)
			}
		} else {
			w.WriteHeader(StatusNoContent) // 204 No Content if no handlerFunc
		}
	}
}

// extractHeaders extracts all headers from an HTTP request and returns them
// Method Used Internally
// The result will extractHeaders(req http.Request) map[string][]string
func extractHeaders(req http.Request) map[string][]string {
	headersMap := make(map[string][]string)
	for key, values := range req.Header {
		headersMap[key] = values
	}
	return headersMap
}

// extractParamsBind decodes request bodies for JSON/XML payloads using a pooled buffer
// to minimize memory allocations and garbage collection overhead.
//
// Parameters:
//   - c: The Quick context containing request information.
//   - v: The target structure to decode the JSON/XML payload.
//
// Returns:
//   - error: Any decoding errors encountered or unsupported content-type errors.
func extractParamsBind(c *Ctx, v interface{}) error {
	contentType := strings.ToLower(c.Request.Header.Get("Content-Type"))

	// Check supported Content-Type
	if !strings.HasPrefix(contentType, ContentTypeAppJSON) &&
		!strings.HasPrefix(contentType, ContentTypeAppXML) &&
		!strings.HasPrefix(contentType, ContentTypeTextXML) {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	switch {
	case strings.HasPrefix(contentType, ContentTypeAppJSON):

		// Acquire pooled buffer
		buf := acquireJSONBuffer()
		defer releaseJSONBuffer(buf)

		// Read body content into buffer
		if _, err := buf.ReadFrom(c.Request.Body); err != nil {
			return err
		}

		// Reset the Request.Body after reading, enabling re-reads if needed
		c.Request.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

		return json.Unmarshal(buf.Bytes(), v)
	case strings.HasPrefix(contentType, ContentTypeAppXML), strings.HasPrefix(contentType, ContentTypeTextXML):

		// Acquire pooled buffer
		buf := acquireXMLBuffer()
		defer releaseXMLBuffer(buf)

		// Read body content into buffer
		if _, err := buf.ReadFrom(c.Request.Body); err != nil {
			return err
		}

		// Reset the Request.Body after reading, enabling re-reads if needed
		c.Request.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

		return xml.Unmarshal(buf.Bytes(), v)
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// extractParamsPattern extracts the fixed path and dynamic parameters from a given route pattern
// Method Used Internally
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

// extractParamsGet processes an HTTP GET request for a dynamic route,
// extracting query parameters, headers, and handling the request using
// the provided handler function.
//
// This function ensures efficient processing by leveraging a pooled
// Ctx instance, which minimizes memory allocations and reduces garbage
// collection overhead.
//
// The request context (`myContextKey`) is retrieved to extract dynamic
// parameters mapped to the route.
//
// Parameters:
//   - q: The Quick instance that provides configurations and routing context.
//   - pathTmp: The template path used for dynamic route matching.
//   - paramsPath: The actual path used to extract route parameters.
//   - handlerFunc: The function that processes the HTTP request.
//
// Returns:
//   - http.HandlerFunc: A function that processes the request efficiently.
func extractParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Acquire a context from the pool
		ctx := acquireCtx()
		defer releaseCtx(ctx)

		// Retrieve the custom context from the request (myContextKey)
		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		cval := v.(ctxServeHttp)

		// Fill the pooled context with request-specific data
		ctx.Response = w
		ctx.Request = req
		ctx.Params = cval.ParamsMap

		// Initialize Query and Headers maps properly
		ctx.Query = make(map[string]string)
		for key, val := range req.URL.Query() {
			ctx.Query[key] = val[0]
		}

		ctx.Headers = extractHeaders(*req)
		ctx.MoreRequests = q.config.MoreRequests

		// Execute the handler function using the pooled context
		execHandleFunc(ctx, handlerFunc)
	}
}

// extractParamsPost processes an HTTP POST request, extracting the request body
// and headers and handling the request using the provided handler function.
//
// This function ensures that the request body is within the allowed size limit,
// extracts headers, and reuses a pooled Ctx instance to optimize memory usage.
//
// Parameters:
//   - q: The Quick instance that provides configurations and routing context.
//   - handlerFunc: The function that processes the HTTP request.
//
// Returns:
//   - http.HandlerFunc: A handler function that processes the request efficiently.
func extractParamsPost(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Validate body size before processing
		if req.ContentLength > q.config.MaxBodySize {
			http.Error(w, "Request body too large", StatusRequestEntityTooLarge)
			return
		}

		// Acquire a pooled context for request processing
		ctx := acquireCtx()
		defer releaseCtx(ctx) // Ensure the context is returned to the pool after execution

		// Retrieve the custom context from the request
		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req) // Return 404 if no context value is found
			return
		}

		// Extract headers into the pooled Ctx
		ctx.Headers = extractHeaders(*req)

		// Read the request body while minimizing allocations
		bodyBytes, bodyReader := extractBodyBytes(req.Body)

		// Populate the Ctx with relevant data
		ctx.Response = w
		ctx.Request = req
		ctx.bodyByte = bodyBytes
		ctx.MoreRequests = q.config.MoreRequests

		// Reset `Request.Body` with the new bodyReader to allow re-reading
		ctx.Request.Body = bodyReader

		// Execute the handler function using the pooled context
		execHandleFunc(ctx, handlerFunc)
	}
}

// extractParamsPut processes an HTTP PUT request, extracting the request body,
// headers, and route parameters while efficiently reusing a pooled Ctx instance.
//
// This function ensures that the request body does not exceed the configured
// size limit, extracts headers, and minimizes memory allocations by leveraging
// a preallocated Ctx from the sync.Pool.
//
// Parameters:
//   - q: The Quick instance that provides configurations and routing context.
//   - handlerFunc: The function that processes the HTTP request.
//
// Returns:
//   - http.HandlerFunc: A function that processes the request efficiently.
func extractParamsPut(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Validate body size before processing
		if req.ContentLength > q.config.MaxBodySize {
			http.Error(w, "Request body too large", StatusRequestEntityTooLarge)
			return
		}

		// Acquire a pooled context for request processing
		ctx := acquireCtx()
		defer releaseCtx(ctx) // Ensure the context is returned to the pool after execution

		// Retrieve the custom context from the request
		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req) // Return 404 if no context value is found
			return
		}

		cval := v.(ctxServeHttp)

		// Extract headers into the pooled Ctx
		ctx.Headers = extractHeaders(*req)

		// Read the request body while minimizing allocations
		bodyBytes, bodyReader := extractBodyBytes(req.Body)

		// Populate the Ctx with relevant data
		ctx.Response = w
		ctx.Request = req
		ctx.bodyByte = bodyBytes
		ctx.Params = cval.ParamsMap
		ctx.MoreRequests = q.config.MoreRequests

		// Reset `Request.Body` with the new bodyReader to allow re-reading
		ctx.Request.Body = bodyReader

		// Execute the handler function using the pooled context
		execHandleFunc(ctx, handlerFunc)
	}
}

// extractParamsDelete processes an HTTP DELETE request, extracting request parameters
// and headers before executing the provided handler function.
//
// This function optimizes memory usage by reusing a pooled Ctx instance,
// reducing unnecessary allocations and garbage collection overhead.
//
// Parameters:
//   - q: The Quick instance that provides configurations and routing context.
//   - handlerFunc: The function that processes the HTTP request.
//
// Returns:
//   - http.HandlerFunc: A function that processes the request efficiently.
func extractParamsDelete(q *Quick, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Acquire a pooled context for request processing
		ctx := acquireCtx()
		defer releaseCtx(ctx) // Ensure the context is returned to the pool after execution

		// Retrieve the custom context from the request
		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req) // Return 404 if no context value is found
			return
		}

		cval := v.(ctxServeHttp)

		// Extract headers into the pooled Ctx
		ctx.Headers = extractHeaders(*req)

		// Populate the Ctx with relevant data
		ctx.Response = w
		ctx.Request = req
		ctx.Params = cval.ParamsMap
		ctx.MoreRequests = q.config.MoreRequests

		// Execute the handler function using the pooled context
		execHandleFunc(ctx, handlerFunc)
	}
}

// execHandleFunc executes the provided handler function and handles errors if they occur
// Method Used Internally
// The result will execHandleFunc(c *Ctx, handleFunc HandleFunc)
func execHandleFunc(c *Ctx, handleFunc HandleFunc) {
	err := handleFunc(c)
	if err != nil {
		c.Set("Content-Type", "text/plain; charset=utf-8")
		// #nosec G104
		c.Status(500).SendString(err.Error())
	}
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
// Method Used Internally
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
// Method Used Internally
// The result will appendRoute(route *Route)
func (q *Quick) appendRoute(route *Route) {
	route.handler = q.mwWrapper(route.handler).ServeHTTP
	//q.routes = append(q.routes, *route)
	q.routes = append(q.routes, route)
}

func (rw *pooledResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// ServeHTTP processes incoming HTTP requests and matches registered routes.
// It uses a pooledResponseWriter to reduce memory allocations and improve performance.
func (q *Quick) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// call options
	if req.Method == http.MethodOptions {
		q.handleOptions(w, req)
		return
	}

	// Acquire a ResponseWriter from the pool for efficient request handling.
	rw := acquireResponseWriter(w)
	defer releaseResponseWriter(rw) // Ensure it returns to the pool.

	// Acquiring Ctx from the pool
	ctx := newCtx(rw, req) // <- creates a new, clean instance of the context
	defer releaseCtx(ctx)  // Returns it to the pool

	for i := range q.routes {
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
			continue // This route doesn't match, continue checking.
		}

		var c = ctxServeHttp{
			Path:      requestURI,
			ParamsMap: paramsMap,
			Method:    q.routes[i].Method,
		}
		req = req.WithContext(context.WithValue(req.Context(), myContextKey, c))

		// Pass the rw (pooledResponseWriter) to the handler
		q.routes[i].handler(rw, req)
		return
	}

	// If no route matches, send a 404 response.
	http.NotFound(rw, req)
}

// createParamsAndValid create params map and check if the request URI and pattern URI are valid
// Method Used Internally
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
// Method Used Internally
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

// ListenWithShutdown starts an HTTP server and returns both the server instance and a shutdown function.
//
// This method initializes performance tuning settings, creates a TCP listener, and starts the server in a background goroutine.
// The returned shutdown function allows for a graceful termination of the server.
//
// Parameters:
//   - addr: The address (host:port) where the server should listen.
//   - handler: Optional HTTP handlers that can be provided to the server.
//
// Returns:
//   - *http.Server: A reference to the initialized HTTP server.
//   - func(): A shutdown function to gracefully stop the server.
//   - error: Any error encountered during the server setup.
func (q *Quick) ListenWithShutdown(addr string, handler ...http.Handler) (*http.Server, func(), error) {
	q.setupPerformanceTuning()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	server := q.httpServer(listener.Addr().String(), handler...)
	q.server = server

	// Shutdown function to gracefully terminate the server.
	shutdownFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		listener.Close()
	}

	// Start the server in a background goroutine.
	go func() {
		server.Serve(listener)
	}()

	return server, shutdownFunc, nil
}

// setupPerformanceTuning configures performance settings for the Quick server.
//
// This method:
//   - Tunes garbage collection behavior dynamically if MoreRequests is configured.
//   - Adjusts the GOMAXPROCS value if specified in the configuration.
//   - Initializes a buffer pool to optimize memory allocation.
func (q *Quick) setupPerformanceTuning() {
	if q.config.MoreRequests > 0 {
		go q.adaptiveGCTuner()
	}

	if q.config.GOMAXPROCS > 0 {
		runtime.GOMAXPROCS(q.config.GOMAXPROCS)
	}

	q.bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, q.config.BufferPoolSize))
		},
	}
}

// adaptiveGCTuner periodically monitors memory usage and triggers garbage collection if necessary.
//
// This function runs in a background goroutine and:
//   - Checks heap memory usage every 15 seconds.
//   - If the heap usage exceeds a defined threshold, it triggers garbage collection and frees OS memory.
func (q *Quick) adaptiveGCTuner() {
	var threshold uint64 = uint64(q.config.GCHeapThreshold)
	if threshold == 0 {
		threshold = 1 << 30 // Default threshold: 1GB
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var m runtime.MemStats
	for range ticker.C {
		runtime.ReadMemStats(&m)

		if m.HeapInuse > threshold {
			debug.FreeOSMemory()
			runtime.GC()
		}
	}
}

// Listen calls ListenWithShutdown and blocks with select{}
// The result will Listen(addr string, handler ...http.Handler) error
func (q *Quick) Listen(addr string, handler ...http.Handler) error {
	_, shutdown, err := q.ListenWithShutdown(addr, handler...)
	if err != nil {
		return err
	}
	defer shutdown()

	q.Display("http", addr)
	// Locks indefinitely
	<-make(chan struct{}) // Bloqueio sem consumo de CPU
	return nil
}

// Shutdown gracefully shuts down the server without interrupting any active connections
// The result will (q *Quick) Shutdown() error
func (q *Quick) Shutdown() error {
	// Create a context with a timeout to control the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure the context is cancelled to free resources

	// Check if the server is initialized before attempting to shut it down
	if q.server != nil {
		q.server.SetKeepAlivesEnabled(false)
		err := q.server.Shutdown(ctx)
		q.releaseResources()
		return err // Attempt to shutdown the server gracefully
	}
	return nil // Return nil if there is no server to shutdown
}

func (q *Quick) releaseResources() {
	// System settings reset
	debug.SetGCPercent(100) // Return to default GC
	runtime.GOMAXPROCS(0)   // Reset to automatic thread configuration
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
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
			PreferServerCipherSuites: true,                              // Prioritize server ciphers
			SessionTicketsDisabled:   false,                             // Enable Session Resumption (minus TLS Handshakes)
			ClientSessionCache:       tls.NewLRUClientSessionCache(128), // Cache TLS sessions for reuse
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
					if runtime.GOOS == "linux" {
						if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, SO_REUSEPORT, 1); err != nil {
							log.Fatalf("Erro ao definir SO_REUSEPORT: %v", err)
						}
					}
					// } else {
					// 	if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
					// 		log.Fatalf("Erro ao definir SO_REUSEPORT: %v", err)
					// 	}
					// }
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
