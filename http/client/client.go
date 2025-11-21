// Package client provides an advanced HTTP client for making HTTP requests.
//
// It supports:
// - Customizable HTTP clients with configurable timeouts, connection pooling, and TLS settings.
// - Built-in retry logic with exponential backoff and failover support.
// - Thread-safe header management to prevent data races in concurrent requests.
// - Structured logging using `slog.Logger` for debugging and observability.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// httpGoClient defines the minimal interface for an HTTP client.
type httpGoClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryConfig defines configuration parameters for retry logic
type RetryConfig struct {
	MaxRetries   int           // Maximum number of retry attempts
	Delay        time.Duration // Base delay between retries
	UseBackoff   bool          // Enable exponential backoff
	Statuses     []int         // HTTP status codes that trigger retry
	FailoverURLs []string      // Alternative URLs for failover
	EnableLog    bool          // Enable logging for retry attempts
}

// Client represents a configurable HTTP client with advanced features
type Client struct {
	Ctx          context.Context   // Context for request cancellation
	ClientHTTP   httpGoClient      // Underlying HTTP client implementation
	Headers      map[string]string // Default headers for requests
	EnableLogger bool              // Flag to enable/disable logging
	Logger       *slog.Logger      // Logger instance
	headersLock  sync.RWMutex      // Mutex for thread-safe header access
}

// RetryTransport implements http.RoundTripper with retry and failover logic
type RetryTransport struct {
	Base          http.RoundTripper // Base transport implementation
	MaxRetries    int               // Maximum retry attempts
	RetryDelay    time.Duration     // Delay between retries
	UseBackoff    bool              // Use exponential backoff
	RetryStatuses []int             // Status codes triggering retry
	Logger        *slog.Logger      // Logger instance
	EnableLogger  bool              // Enable logging
	FailoverURLs  []string          // Alternative URLs for failover
}

// HTTPClientConfig defines configuration for the underlying HTTP client
type HTTPClientConfig struct {
	Timeout             time.Duration // Request timeout
	DisableKeepAlives   bool          // Disable HTTP keep-alive
	MaxIdleConns        int           // Maximum idle connections
	MaxConnsPerHost     int           // Maximum connections per host
	MaxIdleConnsPerHost int           // Maximum idle connections per host
	TLSClientConfig     *tls.Config   // TLS configuration
}

// ClientResponse represents the response from an HTTP request
type ClientResponse struct {
	Body       []byte // Response body
	StatusCode int    // HTTP status code
}

// Option defines a functional option for configuring the Client
type Option func(*Client)

var (
	defaultClient *Client
	once          sync.Once
)

// GetDefaultClient returns a singleton instance of the default HTTP client.
//
// This ensures that the same client instance is reused across the application,
// reducing the overhead of creating multiple clients.
//
// /Return:
//   - *Client: Pointer to the default Client instance.
func GetDefaultClient() *Client {
	once.Do(func() {
		defaultClient = New()
	})
	return defaultClient
}

// NewHTTPClientFromConfig creates and configures an HTTP client based on the provided settings.
//
// If no configuration is provided (cfg is nil), default settings are applied with sensible values.
//
// /Parameters:
//   - cfg: Pointer to `HTTPClientConfig` containing custom configurations. If nil, defaults are used.
//
// /Return:
//   - httpGoClient: Configured HTTP client with timeout, connection settings, and TLS options.
func NewHTTPClientFromConfig(cfg *HTTPClientConfig) httpGoClient {
	if cfg == nil {
		cfg = &HTTPClientConfig{
			Timeout:             30 * time.Second,
			DisableKeepAlives:   true,
			MaxIdleConns:        10,
			MaxConnsPerHost:     10,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				},
			},
		}
	}

	return &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives:   cfg.DisableKeepAlives,
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			TLSClientConfig:     cfg.TLSClientConfig,
		},
	}
}

// WithHTTPClientConfig applies a custom configuration to the HTTP client.
//
// This function allows customization of the HTTP client's settings, such as
// timeout duration, connection pooling, keep-alive behavior, and TLS configuration.
//
// /Parameters:
//   - cfg (*HTTPClientConfig): Pointer to a struct containing the desired configuration settings.
//
// /Return:
//   - Option: A function that applies the custom HTTP configuration to the client.
//
// Example Usage:
//
//	client := New(WithHTTPClientConfig(&HTTPClientConfig{
//	    Timeout: 10 * time.Second,
//	    MaxIdleConns: 20,
//	}))
func WithHTTPClientConfig(cfg *HTTPClientConfig) Option {
	return func(c *Client) {
		c.ClientHTTP = NewHTTPClientFromConfig(cfg)
	}
}

// New creates a new Client instance with optional configurations.
//
// This function initializes a client with default values, but allows additional
// configurations to be passed as options. These options are applied sequentially
// to modify the client settings.
//
// /Parameters:
//   - opts: Variadic slice of `Option` functions to customize the client.
//
// /Return:
//   - *Client: Pointer to the configured HTTP client.
func New(opts ...Option) *Client {
	c := &Client{
		Ctx:        context.Background(),
		Headers:    make(map[string]string),
		ClientHTTP: NewHTTPClientFromConfig(nil),
		Logger:     defaultLogger(),
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// cloneHeaders creates a copy of the client's headers while ensuring thread safety.
//
// This function prevents concurrent modifications to the headers map,
// ensuring data integrity in multi-threaded applications.
//
// /Return:
//   - map[string]string: A thread-safe copy of the headers.
func (c *Client) cloneHeaders() map[string]string {
	c.headersLock.RLock()
	defer c.headersLock.RUnlock()

	headers := make(map[string]string, len(c.Headers))
	for k, v := range c.Headers {
		headers[k] = v
	}
	return headers
}

// log writes log messages if logging is enabled for the client.
//
// This function is used internally to log messages at the INFO level,
// helping with debugging and monitoring request activity.
//
// /Parameters:
//   - msg: The log message format string.
//   - args: Additional arguments to format within the log message.
//
// /Return:
//   - None (void function).
func (c *Client) Log(msg string, args ...interface{}) {
	if c.EnableLogger && c.Logger != nil {
		// Format and log the message at the INFO level.
		c.Logger.Info(msg, args...)
	}
}

// WithLogger enables or disables logging for the HTTP client.
//
// When enabled, the client will log request details, response statuses, and errors
// using the structured logging system provided by `slog.Logger`.
//
// /Parameters:
//   - enable (bool): If true, enables logging. If false, disables it.
//
// /Return:
//   - Option: A function that enables or disables logging in the client.
//
// Example Usage:
//
//	client := New(WithLogger(true))
func WithLogger(enable bool) Option {
	return func(c *Client) {
		c.EnableLogger = enable
	}
}

// defaultLogger initializes and returns a default JSON logger.
//
// This logger writes structured logs in JSON format to standard output
// and is useful for debugging HTTP requests and responses.
//
// /Return:
//   - *slog.Logger: Pointer to the default JSON logger instance.
func defaultLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Get performs a GET request using the default client.
//
// /Parameters:
//   - url: The URL to send the GET request to.
//
// /Return:
//   - *ClientResponse: Pointer to the response containing the status code and body.
//   - error: Error object if the request fails, otherwise nil.
func Get(url string) (*ClientResponse, error) {
	return defaultClient.Get(url)
}

// Post performs a POST request using the default client.
//
// /Parameters:
//   - url: The target URL for the request.
//   - body: The request payload, which can be a string, struct, or JSON object.
//
// /Return:
//   - *ClientResponse: Pointer to the response containing the status code and body.
//   - error: Error object if the request fails, otherwise nil.
func Post(url string, body any) (*ClientResponse, error) {
	return defaultClient.Post(url, body)
}

// Put performs a PUT request using the default client.
//
// /Parameters:
//   - url: The target URL for the request.
//   - body: The request payload, which can be a string, struct, or JSON object.
//
// /Return:
//   - *ClientResponse: Pointer to the response containing the status code and body.
//   - error: Error object if the request fails, otherwise nil.
func Put(url string, body any) (*ClientResponse, error) {
	return defaultClient.Put(url, body)
}

// Delete performs a DELETE request using the default client.
//
// /Parameters:
//   - url: The target URL for the request.
//
// /Return:
//   - *ClientResponse: Pointer to the response containing the status code and body.
//   - error: Error object if the request fails, otherwise nil.
func Delete(url string) (*ClientResponse, error) {
	return defaultClient.Delete(url)
}

// PostForm performs a form-encoded POST request.
//
// Sends a `application/x-www-form-urlencoded` POST request with form data.
//
// /Parameters:
//   - url (string): The target URL.
//   - formData (url.Values): Form fields as key-value pairs.
//
// /Return:
//   - *ClientResponse: Response with status code and body.
//   - error: Error if the request fails.
//
// Example:
//
//	PostForm("/login", formData)
func PostForm(url string, formData url.Values) (*ClientResponse, error) {
	return defaultClient.PostForm(url, formData)
}

// Get performs an HTTP GET request.
//
// Sends a GET request to the specified URL.
//
// /Parameters:
//   - url (string): The target URL.
//
// /Return:
//   - *ClientResponse: Response containing status and body.
//   - error: Error if the request fails.
//
// Example:
//
//	client.Get("/users/2")
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodGet, nil)
}

// Post performs an HTTP POST request.
//
// Sends data to the specified URL using a POST request.
//
// /Parameters:
//   - url (string): The target URL.
//   - body (any): Request body, auto-serialized if needed.
//
// /Return:
//   - *ClientResponse: Response containing status and body.
//   - error: Error if the request fails.
//
// Example:
//
//	client.Post("/users", userData)
func (c *Client) Post(url string, body any) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodPost, body)
}

// Put performs an HTTP PUT request.
//
// Sends a PUT request to update a resource at the specified URL.
//
// /Parameters:
//   - url (string): The target URL.
//   - body (any): Request body, auto-serialized.
//
// /Return:
//   - *ClientResponse: Response containing status and body.
//   - error: Error if the request fails.
//
// Example:
//
//	client.Put("/users/2", userData)
func (c *Client) Put(url string, body any) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodPut, body)
}

// Delete performs an HTTP DELETE request.
//
// Sends a DELETE request to remove a resource.
//
// /Parameters:
//   - url (string): The target URL.
//
// /Return:
//   - *ClientResponse: Response containing status and body.
//   - error: Error if the request fails.
//
// Example:
//
//	client.Delete("/users/2")
func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodDelete, nil)
}

// PostForm performs an HTTP POST request with form data.
//
// Sends form-encoded data (`application/x-www-form-urlencoded`).
//
// /Parameters:
//   - url (string): The target URL.
//   - formData (url.Values): Form fields as key-value pairs.
//
// /Return:
//   - *ClientResponse: Response containing status and body.
//   - error: Error if the request fails.
//
// Example:
//
//	client.PostForm("/login", formData)
func (c *Client) PostForm(url string, formData url.Values) (*ClientResponse, error) {
	headers := c.cloneHeaders()
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return c.doRequestWithHeaders(url, http.MethodPost, formData.Encode(), headers)
}

// doRequest executes an HTTP request using default headers.
//
// This function constructs and sends an HTTP request with the specified
// method and body, applying default headers stored in the client.
//
// /Parameters:
//   - url (string): The target URL.
//   - method (string): The HTTP method (GET, POST, PUT, DELETE, etc.).
//   - body (any): The request body, which is automatically serialized if needed.
//
// /Return:
//   - *ClientResponse: Contains the response body and status code.
//   - error: If an error occurs during the request, it is returned.
func (c *Client) doRequest(url, method string, body any) (*ClientResponse, error) {
	return c.doRequestWithHeaders(url, method, body, c.cloneHeaders())
}

// doRequestWithHeaders executes an HTTP request with custom headers.
//
// This function sends an HTTP request to the specified endpoint using the given method
// and request body. It also allows setting custom headers for the request.
//
// /Parameters:
//   - endpoint (string): The target URL for the request.
//   - method (string): The HTTP method (GET, POST, PUT, DELETE, etc.).
//   - body (any): The request payload, which will be automatically serialized if needed.
//   - headers (map[string]string): A map of custom headers to include in the request.
//
// /Return:
//   - *ClientResponse: A struct containing the HTTP response body and status code.
//   - error: An error object if the request fails, otherwise nil.
//
// Example Usage:
//
//	headers := map[string]string{"Authorization": "Bearer token"}
//	resp, err := client.Get("https://reqres.in/api/users")
func (c *Client) doRequestWithHeaders(endpoint, method string, body any, headers map[string]string) (*ClientResponse, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	reader, err := parseBody(body)
	if err != nil {
		return nil, err
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}
	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = "*/*"
	}

	req, err := http.NewRequestWithContext(c.Ctx, method, endpoint, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.ClientHTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{
		Body:       responseBody,
		StatusCode: resp.StatusCode,
	}, nil
}

// parseBody converts various data types into an io.Reader.
//
// This function is used internally to handle different types of input
// and convert them into a format that can be sent in an HTTP request.
//
// /Parameters:
//   - body (any): The request payload, which can be a string, reader, or struct.
//
// /Return:
//   - io.Reader: A reader containing the serialized data.
//   - error: If serialization fails, an error is returned.
func parseBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	switch v := body.(type) {
	case io.Reader:
		return v, nil
	case string:
		return bytes.NewReader([]byte(v)), nil
	case []byte:
		return bytes.NewReader(v), nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(data), nil
	}
}

// WithRetry configures automatic request retries.
//
// Adds retry behavior to the HTTP client, allowing it to attempt
// a request multiple times if it fails due to certain conditions.
//
// /Parameters:
//   - cfg (RetryConfig): Configuration specifying retry conditions.
//
// /Return:
//   - Option: A functional option that can be applied to a Client.
func WithRetry(cfg RetryConfig) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {

			var logger *slog.Logger
			if c.Logger != nil && cfg.EnableLog {
				logger = c.Logger
			}

			if httpClient.Transport == nil {
				httpClient.Transport = http.DefaultTransport
			}

			httpClient.Transport = &RetryTransport{
				Base:          httpClient.Transport,
				MaxRetries:    cfg.MaxRetries,
				RetryDelay:    cfg.Delay,
				UseBackoff:    cfg.UseBackoff,
				RetryStatuses: cfg.Statuses,
				Logger:        logger,
				EnableLogger:  cfg.EnableLog,
				FailoverURLs:  cfg.FailoverURLs,
			}
		}
	}
}

// WithTimeout sets the timeout for HTTP requests.
//
// /Parameters:
//   - d (time.Duration): The maximum duration before a request times out.
//
// /Return:
//   - Option: A functional option for setting the timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Timeout = d
		}
	}
}

// WithDisableKeepAlives enables or disables HTTP keep-alive connections.
//
// This function modifies the client's transport settings to control keep-alive behavior.
// Disabling keep-alives can be useful in cases where short-lived connections are preferred.
//
// /Parameters:
//   - disable (bool): If true, keep-alives are disabled; if false, they remain enabled.
//
// /Return:
//   - Option: A functional option that configures the client's transport.
//
// Example Usage:
//
//	client := New(WithDisableKeepAlives(true))
func WithDisableKeepAlives(disable bool) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.DisableKeepAlives = disable
			}
		}
	}
}

// WithMaxIdleConns sets the maximum number of idle connections.
//
// This function configures the maximum number of idle (persistent) connections
// the HTTP client can maintain. Increasing this value can improve performance
// for frequent requests to the same server.
//
// /Parameters:
//   - max (int): The maximum number of idle connections.
//
// /Return:
//   - Option: A functional option that configures the client's transport.
//
// Example Usage:
//
//	client := New(WithMaxIdleConns(50))
func WithMaxIdleConns(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxIdleConns = max
			}
		}
	}
}

// WithMaxConnsPerHost sets the maximum concurrent connections per host.
//
// This function limits the number of concurrent TCP connections to a single host.
// It helps control resource allocation when making multiple requests.
//
// /Parameters:
//   - max (int): The maximum number of concurrent connections per host.
//
// /Return:
//   - Option: A functional option that configures the client's transport.
//
// Example Usage:
//
//	client := New(WithMaxConnsPerHost(10))
func WithMaxConnsPerHost(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxConnsPerHost = max
			}
		}
	}
}

// WithMaxIdleConnsPerHost sets the maximum idle connections per host.
//
// This function defines the maximum number of idle (reusable) connections
// that the client can maintain per host.
//
// /Parameters:
//   - max (int): The maximum number of idle connections per host.
//
// /Return:
//   - Option: A functional option that configures the client's transport.
//
// Example Usage:
//
//	client := New(WithMaxIdleConnsPerHost(5))
func WithMaxIdleConnsPerHost(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxIdleConnsPerHost = max
			}
		}
	}
}

// WithTLSConfig applies a custom TLS configuration.
//
// /Parameters:
//   - tlsConfig (*tls.Config): Custom TLS settings.
//
// /Return:
//   - Option: A functional option for setting the TLS configuration.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.TLSClientConfig = tlsConfig
			}
		}
	}
}

// WithInsecureTLS enables or disables TLS certificate verification.
//
// This function modifies the TLS configuration of the HTTP client, allowing insecure connections
// by skipping certificate verification. This is useful for testing environments but should not
// be used in production due to security risks.
//
// /Parameters:
//   - insecure (bool): If true, certificate verification is disabled.
//
// /Return:
//   - Option: A functional option that configures the client's TLS settings.
//
// Example Usage:
//
//	client := New(WithInsecureTLS(true))
func WithInsecureTLS(insecure bool) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				// Ensure TLSClientConfig exists before modifying it
				if transport.TLSClientConfig == nil {
					transport.TLSClientConfig = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
				transport.TLSClientConfig.InsecureSkipVerify = insecure
			}
		}
	}
}

// WithTransport sets a custom HTTP transport.
//
// This function allows users to set a custom `http.RoundTripper` transport
// for the client, which is useful for scenarios requiring advanced transport
// configurations, such as proxies or custom connection handling.
//
// /Parameters:
//   - transport (http.RoundTripper): The custom transport to use.
//
// /Return:
//   - Option: A functional option that configures the client’s transport.
//
// Example Usage:
//
//	client := New(WithTransport(http.DefaultTransport))
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = transport
		}
	}
}

// WithIdleConnTimeout sets the maximum amount of time an idle connection will remain idle before closing.
//
// This function configures how long an idle (keep-alive) connection remains open
// before being closed. Setting this helps manage connection lifecycle and resource usage.
//
// /Parameters:
//   - timeout (time.Duration): The maximum idle time for connections.
//
// /Return:
//   - Option: A functional option that configures the client's transport.
//
// Example Usage:
//
//	client := New(WithIdleConnTimeout(90 * time.Second))
func WithIdleConnTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.IdleConnTimeout = timeout
			}
		}
	}
}

// WithCustomHTTPClient replaces the default HTTP client with a fully custom one.
//
// This function allows users to replace the internal HTTP client with a custom
// instance, providing complete control over transport, timeouts, and other
// configurations.
//
// /Parameters:
//   - client (*http.Client): A fully configured HTTP client.
//
// /Return:
//   - Option: A functional option that sets the custom HTTP client.
//
// Example Usage:
//
//	client := New(WithCustomHTTPClient(&http.Client{Timeout: 10 * time.Second}))
func WithCustomHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client.Transport == nil {
			client.Transport = http.DefaultTransport
		}
		c.ClientHTTP = client
	}
}

// WithTransportConfig applies a pre-configured HTTP transport.
//
// This function sets a pre-configured `http.Transport` for the client,
// allowing customization of low-level transport options such as connection pooling
// and timeouts.
//
// /Parameters:
//   - tr (*http.Transport): The pre-configured HTTP transport.
//
// /Return:
//   - Option: A functional option that applies the transport settings.
//
// Example Usage:
//
//	transport := &http.Transport{MaxIdleConns: 100}
//	client := New(WithTransportConfig(transport))
func WithTransportConfig(tr *http.Transport) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = tr
		}
	}
}

// WithContext sets a custom context for the client.
//
// This function assigns a custom `context.Context` to the client,
// enabling request cancellation and deadline management.
//
// /Parameters:
//   - ctx (context.Context): The custom context to use.
//
// /Return:
//   - Option: A functional option that sets the client's context.
//
// Example Usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	client := New(WithContext(ctx))
func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.Ctx = ctx
	}
}

// WithHeaders sets custom headers for the client.
//
// This function defines default headers that will be included in every request
// made by the client. Existing headers with the same key will be overwritten.
//
// /Parameters:
//   - headers (map[string]string): A map containing key-value pairs of headers.
//
// /Return:
//   - Option: A functional option that sets the default headers.
//
// Example Usage:
//
//	client := New(WithHeaders(map[string]string{"Authorization": "Bearer token"}))
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.headersLock.Lock()
		defer c.headersLock.Unlock()

		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// RoundTrip executes an HTTP request with retry and failover support.
//
// This function implements the http.RoundTripper interface and retries
// requests based on configured conditions.
//
// /Parameters:
//   - req (*http.Request): The HTTP request to send.
//
// /Return:
//   - *http.Response: The HTTP response from the server.
//   - error: If all retry attempts fail, an error is returned.
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var body []byte

	urls := append([]string{req.URL.String()}, rt.FailoverURLs...)

	if req.Body != nil && req.GetBody == nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		_ = req.Body.Close()
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
		req.ContentLength = int64(len(body))
	}

	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		req2 := req.Clone(req.Context())
		index := attempt % len(urls)
		u := urls[index]

		parsedURL, err := url.Parse(u)
		if err != nil {
			continue
		}
		req2.URL = parsedURL

		if req2.GetBody != nil {
			rc, _ := req2.GetBody()
			req2.Body = rc
		} else {
			req2.Body = http.NoBody
		}

		if req2.Header.Get("Accept") == "" {
			req2.Header.Set("Accept", "*/*")
		}
		if req2.Header.Get("Content-Type") == "" {
			req2.Header.Set("Content-Type", "application/json")
		}

		resp, err = rt.Base.RoundTrip(req2)
		if rt.shouldRetry(resp, err) {
			if rt.EnableLogger && rt.Logger != nil {
				rt.Logger.Warn("Retrying request",
					slog.String("url", u),
					slog.String("method", req2.Method),
					slog.Int("attempt", attempt+1),
					slog.Int("failover", index+1),
				)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
			rt.sleep(attempt)
			continue
		}
		return resp, err
	}
	return resp, fmt.Errorf("max retries exceeded: %d", rt.MaxRetries)
}

// shouldRetry determines if a request should be retried.
//
// Checks the response status code and any request errors to decide
// if a retry should be attempted.
//
// /Parameters:
//   - resp (*http.Response): The HTTP response (may be nil).
//   - err (error): An error that occurred during the request.
//
// /Return:
//   - bool: True if the request should be retried, false otherwise.
func (rt *RetryTransport) shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	for _, status := range rt.RetryStatuses {
		if resp.StatusCode == status {
			return true
		}
	}
	return false
}

// sleep waits before retrying a request.
//
// Implements an exponential backoff strategy to avoid overloading the server.
//
// /Parameters:
//   - attempt (int): The current retry attempt number.
func (rt *RetryTransport) sleep(attempt int) {
	delay := rt.RetryDelay
	if rt.UseBackoff {
		delay = time.Duration(math.Pow(2, float64(attempt))) * rt.RetryDelay
	}
	time.Sleep(delay)
}
