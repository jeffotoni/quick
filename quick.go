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
	"path"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jeffotoni/quick/internal/concat"
	"github.com/jeffotoni/quick/template"
)

// SO_REUSEPORT is a constant manually defined for Linux systems
const SO_REUSEPORT = 0x0F

// Content-Type constants used for response headers
const (
	ContentTypeAppJSON = `application/json`
	ContentTypeAppXML  = `application/xml`
	ContentTypeTextXML = `text/xml`
)

// contextKey is a custom type used for storing values in context
type contextKey int

// myContextKey is a predefined key used for context storage
const myContextKey contextKey = 0

// HandleFunc represents a function signature for route handlers in Quick.
//
// This function type is used for defining request handlers within Quick's
// routing system. It receives a pointer to `Ctx`, which encapsulates
// request and response data.
//
// Example Usage:
//
//	q.Get("/example", func(c *quick.Ctx) error {
//	    return c.Status(quick.StatusOK).SendString("Hello, Quick!")
//	})
type HandleFunc func(*Ctx) error

// nextFunc is an internal type used to control next handler execution.
type nextFunc func(*Ctx) error

// HandlerFunc defines the function signature for request handlers in Quick.
//
// This type provides a way to implement request handlers as standalone
// functions while still conforming to the `Handler` interface. It allows
// functions of type `HandlerFunc` to be passed as middleware or endpoint handlers.
//
// Example Usage:
//
//	func myHandler(c *quick.Ctx) error {
//	    return c.Status(quick.StatusOK).SendString("HandlerFunc example")
//	}
//
//	q.Use(quick.HandlerFunc(myHandler))
type HandlerFunc func(c *Ctx) error

// M is a shortcut for map[string]interface{}, allowing `c.M{}`
type M map[string]interface{}

// Handler defines an interface that wraps the ServeQuick method.
//
// Any type implementing `ServeQuick(*Ctx) error` can be used as a request
// handler in Quick. This abstraction allows for more flexible middleware
// and handler implementations, including struct-based handlers.
//
// Example Usage:
//
//	type MyHandler struct{}
//
//	func (h MyHandler) ServeQuick(c *quick.Ctx) error {
//	    return c.Status(quick.StatusOK).SendString("Struct-based handler")
//	}
//
//	q.Use(MyHandler{})
type Handler interface {
	// ServeQuick processes an HTTP request in the Quick framework.
	//
	// Parameters:
	//   - c *Ctx: The request context containing request and response details.
	//
	// Returns:
	//   - error: Any error encountered while processing the request.
	ServeQuick(*Ctx) error
}

// ServeQuick allows a HandlerFunc to satisfy the Handler interface.
//
// This method enables `HandlerFunc` to be used wherever a `Handler`
// is required by implementing the `ServeQuick` method.
//
// Example Usage:
//
//	q.Use(quick.HandlerFunc(func(c *quick.Ctx) error {
//	    return c.Status(quick.StatusOK).SendString("Hello from HandlerFunc!")
//	}))
func (h HandlerFunc) ServeQuick(c *Ctx) error {
	return h(c)
}

// allMethods lists all supported HTTP methods used by the Any method.
var allMethods = []string{
	MethodGet,
	MethodPost,
	MethodPut,
	MethodPatch,
	MethodDelete,
	MethodOptions,
	MethodHead,
}

// Any registers the same handlerFunc for all standard HTTP methods (GET, POST, PUT, etc.).
//
// This is useful when you want to attach a single handler to a path regardless of the HTTP method,
// and handle method-based logic inside the handler itself (e.g., returning 405 if not GET).
//
// Example:
//
//	app := quick.New()
//	app.Any("/health", func(c *quick.Ctx) error {
//		if c.Method() != quick.MethodGet {
//			return c.Status(quick.StatusMethodNotAllowed).SendString("Method Not Allowed")
//		}
//		return c.Status(quick.StatusOK).SendString("OK")
//	})
//
// Note: The handlerFunc will be registered individually for each method listed in allMethods.
func (q *Quick) Any(path string, handlerFunc HandleFunc) {
	for _, method := range allMethods {
		q.registerRoute(method, path, handlerFunc)
	}
}

// MaxBytesReader is a thin wrapper around http.MaxBytesReader to limit the
// size of the request body in Quick applications.
//
// It returns an io.ReadCloser that reads from r but stops with an error
// after n bytes.  The sink just sees an io.EOF.
//
// This is useful to protect against large request bodies.
//
// Example usage:
//
//	c.Request.Body = quick.MaxBytesReader(c.Response, c.Request.Body, 10_000) // 10KB
func MaxBytesReader(w http.ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser {
	// Internally, just call the standard library function.
	return http.MaxBytesReader(w, r, n)
}

// Route represents a registered HTTP route in the Quick framework
type Route struct {
	Group   string           // Route group for organization
	Pattern string           // URL pattern associated with the route
	Path    string           // The registered path for the route
	Params  string           // Parameters extracted from the URL
	Method  string           // HTTP method associated with the route (GET, POST, etc.)
	handler http.HandlerFunc // Handler function for processing the request
}

// ctxServeHttp represents the structure for handling HTTP requests
type ctxServeHttp struct {
	Path      string            // Requested URL path
	Params    string            // Query parameters from the request
	Method    string            // HTTP method of the request
	ParamsMap map[string]string // Parsed parameters mapped as key-value pairs
}

// Config defines various configuration options for the Quick server
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
	Views             template.TemplateEngine

	NoBanner bool // Flag to disable the Quick startup Display.
}

// defaultConfig defines the default values for the Quick server configuration
var defaultConfig = Config{
	BodyLimit:      2 * 1024 * 1024, // Deprecated: Use MaxBodySize instead
	MaxBodySize:    2 * 1024 * 1024, // 2MB max request body size
	MaxHeaderBytes: 1 * 1024 * 1024, // 1MB max header size

	GOMAXPROCS:      runtime.NumCPU(), // Use all available CPU cores
	GCHeapThreshold: 1 << 30,          // 1GB memory threshold for GC
	BufferPoolSize:  32768,            // Buffer pool size

	RouteCapacity: 1000,  // Default initial route capacity
	MoreRequests:  290,   // Default GC value
	NoBanner:      false, // Show Quick banner by default
}

// Zeroth is a custom type for zero-value constants
type Zeroth int

// Zero is a predefined constant of type Zeroth
const (
	Zero Zeroth = 0
)

// CorsConfig defines the CORS settings for Quick
type CorsConfig struct {
	Enabled  bool              // If true, enables CORS support
	Options  map[string]string // Custom CORS options
	AllowAll bool              // If true, allows all origins
}

// Quick is the main structure of the framework, holding routes and configurations.
type Quick struct {
	config        Config                          // Configuration settings.
	Cors          bool                            // Indicates if CORS is enabled.
	groups        []Group                         // List of route groups.
	handler       http.Handler                    // The primary HTTP handler.
	mux           *http.ServeMux                  // Multiplexer for routing requests.
	routes        []*Route                        // Registered routes.
	routeCapacity int                             // The maximum number of routes allowed.
	mws2          []any                           // List of registered middlewares.
	CorsSet       func(http.Handler) http.Handler // CORS middleware handler function.
	CorsOptions   map[string]string               // CORS options map
	// corsConfig    *CorsConfig // Specific type for CORS // Removed unused field
	embedFS    embed.FS     // File system for embedded static files.
	server     *http.Server // Http server
	bufferPool *sync.Pool   // Reusable buffer pool to reduce allocations and improve performance
}

// indeed to Quick
type App = Quick

//	(q *Quick) Config() Config
//
// example: c.App.GetConfig().Views
func (q *Quick) GetConfig() Config {
	return q.config
}

// HandlerFunc adapts a quick.HandlerFunc to a standard http.HandlerFunc.
// It creates a new Quick context (Ctx) for each HTTP request,
// allowing Quick handlers to access request and response objects seamlessly.
//
// Usage Example:
//
//	http.HandleFunc(\"/\", app.HandlerFunc(func(c *quick.Ctx) error {
//		return c.Status(200).JSON(map[string]string{\"message\": \"Hello, Quick!\"})
//	}))
func (q *Quick) HandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := &Ctx{
			Response: w,
			Request:  req,
			App:      q,
		}

		if err := h(c); err != nil {
			http.Error(w, err.Error(), StatusInternalServerError)
		}
	}
}

// Handler returns the main HTTP handler for Quick, allowing integration with standard http.Server and testing frameworks.
func (q *Quick) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q.ServeHTTP(w, r)
	})
}

// MiddlewareFunc defines the signature for middleware functions in Quick.
// A middleware function receives the next HandlerFunc in the chain and returns a new HandlerFunc.
// Middleware can perform actions before and/or after calling the next handler.
//
// Example:
//
//	func LoggingMiddleware() quick.MiddlewareFunc {
//		return func(next quick.HandlerFunc) quick.HandlerFunc {
//			return func(c *quick.Ctx) error {
//				// Before handler logic (e.g., logging request details)
//				log.Printf("Request received: %s %s", c.Request.Method, c.Request.URL)
//
//				err := next(c) // Call the next handler
//
//				// After handler logic (e.g., logging response status)
//				log.Printf("Response sent with status: %d", c.ResponseWriter.Status())
//
//				return err
//			}
//		}
//	}
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// GetDefaultConfig returns the default configuration pre-defined for the system.
//
// This function provides a standardized configuration setup, ensuring that
// new instances use a consistent and optimized set of defaults.
//
// Returns:
//   - Config: A struct containing the default system configuration.
//
// Example Usage:
//
//	// This function is typically used when initializing a new Quick instance
//	// to ensure it starts with the default settings if no custom config is provided.
func GetDefaultConfig() Config {
	return defaultConfig
}

// New creates a new instance of the Quick structure to manage HTTP routes and handlers.
//
// This function initializes a Quick instance with optional configurations provided
// through the `Config` parameter. If no configuration is provided, it uses the `defaultConfig`.
//
// Parameters:
//   - c ...Config: (Optional) Configuration settings for customizing the Quick instance.
//
// Returns:
//   - *Quick: A pointer to the initialized Quick instance.
//
// Example Usage:
//
//	// Basic usage - Create a default Quick instance
//	q := quick.New()
//
//	// Custom usage - Create a Quick instance with specific configurations
//	q := quick.New(quick.Config{
//		RouteCapacity: 500,
//	})
//
//	q.Get("/", func(c quick.Ctx) error {
//		return c.SendString("Hello, Quick!")
//	})
func New(c ...Config) *Quick {
	var config Config
	// Check if a custom configuration is provided
	if len(c) > 0 {
		config = c[0] // Use the provided configuration
	} else {
		config = defaultConfig // Use the default configuration
	}

	// Ensure a minimum route capacity if not set
	if config.RouteCapacity == 0 {
		config.RouteCapacity = 1000
	}

	// Initialize and return the Quick instance
	return &Quick{
		routes:        make([]*Route, 0, config.RouteCapacity),
		routeCapacity: config.RouteCapacity,
		mux:           http.NewServeMux(),
		handler:       http.NewServeMux(),
		config:        config,
	}
}

// Use function adds middleware to the Quick server, with special treatment for CORS.
//
// This method allows adding custom middleware functions to process requests before they
// reach the final handler. If a CORS middleware is detected, it is automatically applied.
//
// Parameters:
//   - mw any: Middleware function to be added. It must be of type `func(http.Handler) http.Handler`.
//
// Example Usage:
//
//	q := quick.New()
//
//	q.Use(maxbody.New(50000))
//
//	q.Post("/v1/user/maxbody/any", func(c *quick.Ctx) error {
//	    c.Set("Content-Type", "application/json")//
//	    return c.Status(200).Send(c.Body())
//	})
func (q *Quick) Use(mw any) {
	switch mwc := mw.(type) {
	case func(http.Handler) http.Handler:
		// Detect if the middleware is related to CORS and apply it separately
		if isCorsMiddleware(mwc) {
			q.Cors = true
			q.CorsSet = mwc
			return
		}

	case func(HandleFunc) HandleFunc:
		q.mws2 = append(q.mws2, mwc)

	}

	// Append middleware to the list of registered middlewares
	q.mws2 = append(q.mws2, mw)
}

// isCorsMiddleware checks whether the provided middleware function is a CORS handler.
//
// This function detects if a middleware is handling CORS by sending an
// HTTP OPTIONS request and checking if it sets the `Access-Control-Allow-Origin` header.
//
// Parameters:
//   - mw func(http.Handler) http.Handler: The middleware function to be tested.
//
// Returns:
//   - bool: `true` if the middleware is identified as CORS, `false` otherwise.
//
// Example Usage:
// This function is automatically executed when a middleware is added to detect if it's a CORS handler.
func isCorsMiddleware(mw func(http.Handler) http.Handler) bool {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	testRequest := httptest.NewRequest("OPTIONS", "/", nil)
	testResponse := httptest.NewRecorder()

	mw(testHandler).ServeHTTP(testResponse, testRequest)

	// If the middleware sets Access-Control-Allow-Origin, it's CORS
	return testResponse.Header().Get("Access-Control-Allow-Origin") != ""
}

// clearRegex processes a route pattern, removing dynamic path parameters
// and replacing them with a simplified placeholder.
//
// This function is used internally to standardize dynamic routes in
// ServeMux, converting patterns like `/v1/user/{id:[0-9]+}` into
// `/v1/user/_id_`, making them easier to process.
//
// Parameters:
//   - route string: The route pattern containing dynamic parameters.
//
// Returns:
//   - string: A cleaned-up version of the route with placeholders instead of regex patterns.
//
// Example Usage:
// This function is automatically triggered internally to normalize route patterns.
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

// registerRoute is a helper function that centralizes the logic for registering routes.
//
// This function processes and registers an HTTP route, ensuring no duplicate routes
// are added. It extracts route parameters, formats the route, and associates the
// appropriate handler function.
//
// Parameters:
//   - method string: The HTTP method (e.g., "GET", "POST").
//   - pattern string: The route pattern, which may include dynamic parameters.
//   - handlerFunc HandleFunc: The function that will handle the route.
//
// Example Usage:
// This function is automatically triggered internally when a new route is added.
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

// Get registers an HTTP route with the GET method on the Quick server.
//
// This function associates a GET request with a specific route pattern and handler function.
// It ensures that the request is properly processed when received.
//
// Parameters:
//   - pattern string: The route pattern (e.g., "/users/:id").
//   - handlerFunc HandleFunc: The function that will handle the GET request.
//
// Example Usage:
//
//	// This function is automatically triggered when defining a GET route in Quick.
func (q *Quick) Get(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodGet, pattern, handlerFunc)
}

// Post registers an HTTP route with the POST method on the Quick server.
//
// This function associates a POST request with a specific route pattern and handler function.
// It is typically used for handling form submissions, JSON payloads, or data creation.
//
// Parameters:
//   - pattern string: The route pattern (e.g., "/users").
//   - handlerFunc HandleFunc: The function that will handle the POST request.
//
// Example Usage:
//
//	// This function is automatically triggered when defining a POST route in Quick.
func (q *Quick) Post(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodPost, pattern, handlerFunc)
}

// Put registers an HTTP route with the PUT method on the Quick server.
//
// This function associates a PUT request with a specific route pattern and handler function.
// It is typically used for updating existing resources.
//
// Parameters:
//   - pattern string: The route pattern (e.g., "/users/:id").
//   - handlerFunc HandleFunc: The function that will handle the PUT request.
//
// Example Usage:
//
//	// This function is automatically triggered when defining a PUT route in Quick.
func (q *Quick) Put(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodPut, pattern, handlerFunc)
}

// Delete registers an HTTP route with the DELETE method on the Quick server.
//
// This function associates a DELETE request with a specific route pattern and handler function.
// It is typically used for deleting existing resources.
//
// Parameters:
//   - pattern string: The route pattern (e.g., "/users/:id").
//   - handlerFunc HandleFunc: The function that will handle the DELETE request.
//
// Example Usage:
//
//	// This function is automatically triggered when defining a DELETE route in Quick.
func (q *Quick) Delete(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodDelete, pattern, handlerFunc)
}

// Patch registers an HTTP route with the PATCH method on the Quick server.
//
// This function associates a PATCH request with a specific route pattern and handler function.
// It is typically used for applying partial updates to an existing resource.
//
// Parameters:
//   - pattern string: The route pattern (e.g., "/users/:id").
//   - handlerFunc HandleFunc: The function that will handle the PATCH request.
//
// Example Usage:
//
//	// This function is automatically triggered when defining a PATCH route in Quick.
func (q *Quick) Patch(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodPatch, pattern, handlerFunc)
}

// Options registers an HTTP route with the OPTIONS method on the Quick server.
//
// This function associates an OPTIONS request with a specific route pattern and handler function.
// OPTIONS requests are typically used to determine the allowed HTTP methods for a resource.
//
// Parameters:
//   - pattern string: The route pattern (e.g., "/users").
//   - handlerFunc HandleFunc: The function that will handle the OPTIONS request.
//
// Example Usage:
//
//	// This function is automatically triggered when defining an OPTIONS route in Quick.
func (q *Quick) Options(pattern string, handlerFunc HandleFunc) {
	q.registerRoute(MethodOptions, pattern, handlerFunc)
}

// extractHandler selects the appropriate handler function for different HTTP methods.
//
// This function is responsible for determining which internal request processing function
// should handle a given HTTP method. It maps the method to the corresponding request parser.
//
// Parameters:
//   - q *Quick: The Quick instance managing the route and request context.
//   - method string: The HTTP method (e.g., "GET", "POST").
//   - path string: The route path associated with the request.
//   - params string: Route parameters extracted from the request URL.
//   - handlerFunc HandleFunc: The function that will handle the request.
//
// Returns:
//   - http.HandlerFunc: The appropriate handler function based on the HTTP method.
//
// Example Usage:
//
//	// This function is automatically executed internally when processing an HTTP request.
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

// extractParamsPatch processes an HTTP PATCH request by reusing the logic of the PUT method.
//
// The PATCH method is typically used for partial updates, while PUT replaces an entire resource.
// However, both methods often handle request parameters and body parsing in the same way,
// so this function delegates the processing to `extractParamsPut`.
//
// Parameters:
//   - q *Quick: The Quick instance managing the request context.
//   - handlerFunc HandleFunc: The function that will handle the PATCH request.
//
// Returns:
//   - http.HandlerFunc: A handler function that processes PATCH requests.
//
// Example Usage:
//
//	// This function is automatically executed internally when a PATCH request is received.
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
		ctx.App = q

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

// extractHeaders extracts all headers from an HTTP request and returns them as a map.
//
// This function iterates over all headers in the request and organizes them into a
// map structure, where each header key is mapped to its corresponding values.
//
// Parameters:
//   - req http.Request: The HTTP request from which headers will be extracted.
//
// Returns:
//   - map[string][]string: A map containing all request headers.
//
// Example Usage:
//
//	// This function is automatically executed internally when extracting headers from a request.
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
// This function checks the request's `Content-Type` and processes JSON or XML payloads accordingly.
// It ensures efficient memory usage by leveraging buffer pools for reading request bodies.
//
// Parameters:
//   - c *Ctx: The Quick context containing request information.
//   - v interface{}: The target structure where the decoded JSON/XML data will be stored.
//
// Returns:
//   - error: Returns any decoding errors encountered or an error for unsupported content types.
//
// Example Usage:
//
//	// This function is automatically executed internally when binding request data to a struct.
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

// extractParamsPattern extracts the fixed path and dynamic parameters from a given route pattern.
//
// This function is responsible for identifying and separating static paths from dynamic parameters
// in a route pattern. It ensures proper extraction of URL path segments and dynamic query parameters.
//
// Parameters:
//   - pattern string: The route pattern that may contain dynamic parameters.
//
// Returns:
//   - path string: The fixed portion of the route without dynamic parameters.
//   - params string: The extracted dynamic parameters (if any).
//   - patternExist string: The original pattern before extraction.
//
// Example Usage:
//
//	// This function is automatically executed internally when registering a dynamic route.
//	path, params, patternExist := extractParamsPattern("/users/:id")
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
			NotFound(w, req)
			return
		}

		cval := v.(ctxServeHttp)

		// Fill the pooled context with request-specific data
		ctx.Response = w
		ctx.Request = req
		ctx.Params = cval.ParamsMap
		ctx.App = q

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
		// req.Body = http.MaxBytesReader(w, req.Body, q.config.MaxBodySize)
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
			NotFound(w, req) // Return 404 if no context value is found
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

		ctx.App = q

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
			NotFound(w, req) // Return 404 if no context value is found
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

		ctx.App = q
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
			NotFound(w, req) // Return 404 if no context value is found
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

		ctx.App = q

		// Execute the handler function using the pooled context
		execHandleFunc(ctx, handlerFunc)
	}
}

// execHandleFunc executes the provided handler function and handles errors if they occur.
//
// This function ensures that the HTTP response is properly handled, including setting the
// appropriate content type and returning an error message if the handler function fails.
//
// Parameters:
//   - c *Ctx: The Quick context instance containing request and response data.
//   - handleFunc HandleFunc: The function that processes the HTTP request.
//
// Example Usage:
//
//	// This function is automatically executed internally after processing a request.
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
//
// Additionally, this function returns a new `io.ReadCloser` wrapping the same data,
// allowing it to be re-read if needed.
//
// Parameters:
//   - r io.ReadCloser: The original request body stream.
//
// Returns:
//   - []byte: A byte slice containing the full request body data.
//   - io.ReadCloser: A new ReadCloser that allows the body to be re-read.
//
// Example Usage:
//
//	// Read the request body into a byte slice and obtain a new ReadCloser.
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

// mwWrapper applies all registered middlewares to an HTTP handler.
//
// This function iterates through the registered middleware stack in reverse order
// (last added middleware is executed first) and wraps the final HTTP handler
// with each middleware layer.
//
// Parameters:
//   - handler http.Handler: The final HTTP handler to be wrapped with middlewares.
//
// Returns:
//   - http.Handler: The HTTP handler wrapped with all registered middlewares.
//
// Example Usage:
//
//	// This function is automatically executed internally before processing requests.
//
// mwWrapper applies all registered middlewares to an HTTP handler.
//
// This function iterates through the middleware stack in reverse order
// (last added middleware is executed first) and wraps the final HTTP handler
// with each middleware layer.
//
// It supports multiple middleware function signatures:
//   - `func(http.Handler) http.Handler`: Standard net/http middleware.
//   - `func(http.ResponseWriter, *http.Request, http.Handler)`: Middleware that
//     directly manipulates the response and request.
//   - `func(HandlerFunc) HandlerFunc`: Quick-specific middleware format.
//   - `func(Handler) Handler`: Another Quick middleware format.
//
// Parameters:
//   - handler http.Handler: The final HTTP handler to be wrapped.
//
// Returns:
//   - http.Handler: The HTTP handler wrapped with all registered middlewares.
func (q *Quick) mwWrapper(handler http.Handler) http.Handler {
	for i := len(q.mws2) - 1; i >= 0; i-- {
		switch mw := q.mws2[i].(type) {

		case func(http.Handler) http.Handler:
			// Apply standard net/http middleware
			handler = mw(handler)

		case func(http.ResponseWriter, *http.Request, http.Handler):
			// Apply middleware that takes ResponseWriter, Request, and the next handler
			originalHandler := handler // Avoid infinite reassignment
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw(w, r, originalHandler)
			})

		case func(HandlerFunc) HandlerFunc:
			// Convert net/http.Handler to Quick.HandlerFunc
			quickHandler := convertToQuickHandler(handler)
			// Apply Quick middleware
			quickHandler = mw(quickHandler)

			// Convert back to http.Handler
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c := &Ctx{
					Response: w,
					Request:  r,
					App:      q,
				}
				quickHandler(c)
			})

		case func(Handler) Handler:
			// Convert net/http.Handler to Quick.Handler
			qh := convertHttpToQuickHandler(handler)
			// Apply Quick middleware
			qh = mw(qh)
			// Convert back to http.Handler
			handler = convertQuickToHttpHandler(q, qh)
		}
	}
	return handler
}

// convertHttpToQuickHandler adapts a net/http.Handler to a Quick.Handler.
//
// This function allows standard HTTP handlers to be wrapped within Quick's middleware
// system by transforming them into the Quick.Handler interface.
//
// Parameters:
//   - h http.Handler: The standard HTTP handler to convert.
//
// Returns:
//   - Handler: The Quick-compatible handler.
func convertHttpToQuickHandler(h http.Handler) Handler {
	return HandlerFunc(func(c *Ctx) error {
		h.ServeHTTP(c.Response, c.Request)
		return nil
	})
}

// convertQuickToHttpHandler adapts a Quick.Handler to a net/http.Handler.
//
// This function wraps Quick handlers into a standard HTTP handler so they can
// be used within net/http's ecosystem.
//
// Parameters:
//   - h Handler: The Quick handler to convert.
//
// Returns:
//   - http.Handler: The net/http-compatible handler.
func convertQuickToHttpHandler(q *Quick, h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &Ctx{Response: w, Request: r, App: q}
		if err := h.ServeQuick(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// convertToQuickHandler adapts a net/http.Handler to a Quick.HandlerFunc.
//
// This function allows standard HTTP handlers to be used within Quick's middleware
// system by transforming them into Quick.HandlerFunc.
//
// Parameters:
//   - h http.Handler: The standard HTTP handler to convert.
//
// Returns:
//   - HandlerFunc: The Quick-compatible handler function.
func convertToQuickHandler(h http.Handler) HandlerFunc {
	return func(c *Ctx) error {
		h.ServeHTTP(c.Response, c.Request)
		return nil
	}
}

// func convertToQuickMiddleware(mw func(http.Handler) http.Handler) func(Handler) Handler {
// 	return func(next Handler) Handler {
// 		return HandlerFunc(func(c *Ctx) error {
// 			adaptedHandler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				c.Response = w
// 				c.Request = r
// 				next(c)
// 			}))

// 			adaptedHandler.ServeHTTP(c.Response, c.Request)
// 			return nil
// 		})
// 	}
// }

// appendRoute registers a new route in the Quick router and applies middlewares.
//
// This function ensures that the given route's handler is wrapped with all registered
// middlewares before being stored in the router. It optimizes performance by applying
// middleware only once during route registration instead of at runtime.
//
// Parameters:
//   - route *Route: The route to be registered in the Quick router.
//
// Example Usage:
//
//	// This function is automatically called when registering a new route.
func (q *Quick) appendRoute(route *Route) {
	route.handler = q.mwWrapper(route.handler).ServeHTTP
	//q.routes = append(q.routes, *route)
	q.routes = append(q.routes, route)
}

// Header retrieves the HTTP headers from the response writer.
//
// This method allows middleware or handlers to modify response headers
// before sending them to the client.
//
// Returns:
//   - http.Header: The set of response headers.
//
// Example Usage:
//
//	// Retrieve headers within a middleware or handler function
func (rw *pooledResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// ServeHTTP processes incoming HTTP requests and matches them to registered routes.
//
// This function efficiently routes HTTP requests to the appropriate handler while
// leveraging a **pooled response writer** and **context pooling** to minimize memory
// allocations and improve performance.
//
// If the request method is `OPTIONS`, it is handled separately via `handleOptions`.
// If no matching route is found, the function responds with `404 Not Found`.
//
// Example Usage:
// This function is automatically invoked by the `http.Server` when a request reaches
// the Quick router.
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
	ctx := newCtx(rw, req, q) // creates a new, clean instance of the context
	defer releaseCtx(ctx)     // Returns it to the pool

	for i := range q.routes {
		var requestURI = path.Clean(req.URL.Path)
		var patternUri = q.routes[i].Pattern

		if q.routes[i].Method != req.Method {
			continue
		}

		if len(patternUri) == 0 {
			patternUri = path.Clean(q.routes[i].Path)
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
	NotFound(rw, req)
}

// createParamsAndValid extracts dynamic parameters from a request URI and validates the pattern.
//
// This function compares the request URI with the registered pattern and extracts
// route parameters such as `:id` or `{id:[0-9]+}` dynamically.
//
// Example Usage:
// This function is internally used by `ServeHTTP()` to verify if a request matches a
// registered route pattern. If it does, it extracts the dynamic parameters and returns
// them as a map.
func createParamsAndValid(reqURI, patternURI string) (map[string]string, bool) {
	params := make(map[string]string)
	var builder strings.Builder

	reqURI = strings.TrimPrefix(reqURI, "/")
	patternURI = strings.TrimPrefix(patternURI, "/")

	// Ex: /xxx* or /x/y*
	if strings.HasSuffix(patternURI, "*") {
		base := strings.TrimSuffix(patternURI, "*")
		if strings.HasPrefix(reqURI, base) {
			return params, true
		}
	}

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
	return params, true
}

// GetRoute retrieves all registered routes in the Quick framework.
//
// This function returns a slice containing all the routes that have been
// registered in the Quick instance. It is useful for debugging, logging,
// or dynamically inspecting available routes.
//
// Example Usage:
//
//	routes := q.GetRoute()
//	for _, route := range routes {
//	    fmt.Println("Method:", route.Method, "Path:", route.Path)
//	}
//
// Returns:
//   - []*Route: A slice of pointers to the registered Route instances.
func (q *Quick) GetRoute() []*Route {
	return q.routes
}

// Static serves static files (HTML, CSS, JS, images, etc.) from a directory or embedded filesystem.
//
// This function allows you to register a static file server in the Quick framework, either using
// a local directory (`string`) or an embedded filesystem (`embed.FS`). By embedding files,
// they become part of the binary at compile time, eliminating the need for external file access.
//
// Example Usage:
// This function is useful for serving front-end assets or static resources directly from
// the Go application. It supports both local directories and embedded files.
//
// Parameters:
//   - route: The base path where static files will be served (e.g., "/static").
//   - dirOrFS: The source of the static files. It accepts either:
//   - `string`: A local directory path (e.g., `"./static"`).
//   - `embed.FS`: An embedded file system for compiled-in assets.
//
// Returns:
//   - None (void function).
//
// Notes:
//   - The function automatically trims trailing slashes from `route`.
//   - If an invalid parameter is provided, the function panics.
//   - When using an embedded filesystem, files are served directly from memory.
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

// execHandler wraps an HTTP handler with additional processing.
//
// This function ensures that the provided `http.Handler` is executed properly,
// allowing for additional middleware wrapping or request pre-processing.
//
// Example Usage:
// This function is automatically applied within Quick's internal request processing pipeline.
//
// Parameters:
//   - next: The next HTTP handler to be executed.
//
// Returns:
//   - http.Handler: A wrapped HTTP handler that ensures correct execution.
func (q *Quick) execHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// corsHandler applies the configured CORS middleware to the Quick router.
//
// This function checks if CORS is enabled in the Quick instance (`q.Cors`).
// If enabled, it wraps the request processing with the configured CORS middleware.
//
// Example Usage:
// Automatically used when CORS support is detected within `Use()`.
//
// Returns:
//   - http.Handler: The HTTP handler wrapped with CORS processing.
func (q *Quick) corsHandler() http.Handler {
	return q.CorsSet(q)
}

// httpServerTLS initializes and returns an HTTP server instance configured for TLS/HTTPS.
//
// This function sets up a secure server with `TLSConfig` and allows optional
// custom handlers to be provided. If no custom handler is specified, the default
// Quick router is used.
//
// Example Usage:
// Used internally by Quick when setting up an HTTPS server via `ListenTLS()`.
//
// Parameters:
//   - addr:      The network address the server should listen on (e.g., ":443").
//   - tlsConfig: A `*tls.Config` instance containing certificate and security settings.
//   - handler:   (Optional) One or more custom HTTP handlers.
//
// Returns:
//   - *http.Server: A configured HTTPS server instance.
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

// httpServer initializes and returns an HTTP server instance configured for plain HTTP.
//
// This function sets up a standard HTTP server, allowing optional custom handlers
// to be specified. If no handler is provided, the default Quick router is used.
//
// Example Usage:
// Used internally when starting an HTTP server via `Listen()`.
//
// Parameters:
//   - addr:    The network address the server should listen on (e.g., ":8080").
//   - handler: (Optional) One or more custom HTTP handlers.
//
// Returns:
//   - *http.Server: A configured HTTP server instance.
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
// This method optimizes performance by:
//   - Dynamically tuning garbage collection behavior if `MoreRequests` is configured.
//   - Adjusting the `GOMAXPROCS` value based on the configuration.
//   - Initializing a buffer pool to optimize memory allocation.
//
// Example Usage:
// Called automatically when initializing the Quick server.
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
//
// Example Usage:
// Automatically invoked by `setupPerformanceTuning()` when `MoreRequests` is enabled.
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

// Listen starts the Quick server and blocks indefinitely.
//
// This function initializes the HTTP server and prevents the application from exiting.
//
// Example Usage:
//
//	q.Listen(":8080")
//
// Parameters:
//   - addr: The address on which the server should listen (e.g., ":8080").
//   - handler: (Optional) Custom HTTP handlers.
//
// Returns:
//   - error: Any errors encountered while starting the server.
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

// Shutdown gracefully shuts down the Quick server without interrupting active connections.
//
// This function ensures that all ongoing requests are completed before shutting down,
// preventing abrupt connection termination.
//
// Example Usage:
//
//		q := quick.New()
//
//	    q.Get("/", func(c *quick.Ctx) error {
//	        return c.SendString("Server is running!")
//	    })
//
//		q.Shutdown()
//
// Returns:
//   - error: Any errors encountered during shutdown.
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

// releaseResources resets system-level performance settings after server shutdown.
//
// This function ensures that system configurations return to their default states
// after the Quick server is gracefully shut down, preventing excessive resource usage.
//
// Example Usage:
//
//	q.releaseResources()
//
// Actions Performed:
//   - Resets the garbage collection behavior to the default percentage (100).
//   - Restores `GOMAXPROCS` to automatic CPU thread allocation.
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

// NotFound sends a 404 Not Found response with optional custom body.
// It wraps http.NotFound and provides a Quick-style naming.
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
