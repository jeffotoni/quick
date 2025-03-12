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
	"strings"
	"sync"
	"time"
)

// httpGoClient defines the minimal interface for HTTP clients
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

// GetDefaultClient returns a thread-safe singleton instance of Client
// The result will GetDefaultClient() *Client
func GetDefaultClient() *Client {
	once.Do(func() {
		defaultClient = New()
	})
	return defaultClient
}

// NewHTTPClientFromConfig creates a configured HTTP client
// If cfg is nil, sensible defaults are used
// The result will NewHTTPClientFromConfig(cfg *HTTPClientConfig) *httpGoClient
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

// WithHTTPClientConfig sets a custom configuration for the HTTP client
// The result will WithHTTPClientConfig(cfg *HTTPClientConfig) Option
func WithHTTPClientConfig(cfg *HTTPClientConfig) Option {
	return func(c *Client) {
		c.ClientHTTP = NewHTTPClientFromConfig(cfg)
	}
}

// New creates a new Client with optional configurations
// The result will New(opts ...Option) *Client
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

// cloneHeaders creates a thread-safe copy of the headers.
// This ensures that concurrent access does not modify the original headers.
// Method Used Internally
// The result will cloneHeaders() map[string]string
func (c *Client) cloneHeaders() map[string]string {
	c.headersLock.RLock()
	defer c.headersLock.RUnlock()

	headers := make(map[string]string, len(c.Headers))
	for k, v := range c.Headers {
		headers[k] = v
	}
	return headers
}

// log writes a log message if logging is enabled and a logger is set.
// Method Used Internally
// The result will log(msg string, args ...interface{})
func (c *Client) log(msg string, args ...interface{}) {
	if c.EnableLogger && c.Logger != nil {
		// Format and log the message at the INFO level.
		c.Logger.Info(msg, args...)
	}
}

// WithLogger enables or disables logging
// The result will WithLogger(enable bool) Option
func WithLogger(enable bool) Option {
	return func(c *Client) {
		c.EnableLogger = enable
	}
}

// defaultLogger creates a default JSON logger
// Method Used Internally
// The result will defaultLogger() *slog.Logger
func defaultLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Get performs a GET request using the default client
// The result will Get(url string) (*ClientResponse, error)
func Get(url string) (*ClientResponse, error) {
	return defaultClient.Get(url)
}

// Post performs a POST request using the default client
// The result will Post(url string, body any) (*ClientResponse, error)
func Post(url string, body any) (*ClientResponse, error) {
	return defaultClient.Post(url, body)
}

// Put performs a PUT request using the default client
// The result will Put(url string, body any) (*ClientResponse, error)
func Put(url string, body any) (*ClientResponse, error) {
	return defaultClient.Put(url, body)
}

// Delete performs a DELETE request using the default client
// The result will Delete(url string) (*ClientResponse, error)
func Delete(url string) (*ClientResponse, error) {
	return defaultClient.Delete(url)
}

// PostForm performs a form POST request using the default client
// The result will PostForm(url string, formData url.Values) (*ClientResponse, error)
func PostForm(url string, formData url.Values) (*ClientResponse, error) {
	return defaultClient.PostForm(url, formData)
}

// Get performs a GET request
// The result will Get(url string) (*ClientResponse, error)
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodGet, nil)
}

// Post performs a POST request
// The result will Post(url string, body any) (*ClientResponse, error)
func (c *Client) Post(url string, body any) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodPost, body)
}

// Put sends an HTTP PUT request with a request body.
func (c *Client) Put(url string, body any) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodPut, body)
}

// Delete sends an HTTP DELETE request.
// The result will Delete(url string) (*ClientResponse, error)
func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.doRequest(url, http.MethodDelete, nil)
}

// PostForm performs a form POST request with URL-encoded data
// The result will PostForm(url string, formData url.Values) (*ClientResponse, error)
func (c *Client) PostForm(url string, formData url.Values) (*ClientResponse, error) {
	headers := c.cloneHeaders()
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return c.doRequestWithHeaders(url, http.MethodPost, formData.Encode(), headers)
}

// doRequest executes an HTTP request with default headers
// Method Used Internally
// The result will doRequest(url, method string, body any) (*ClientResponse, error)
func (c *Client) doRequest(url, method string, body any) (*ClientResponse, error) {
	return c.doRequestWithHeaders(url, method, body, c.cloneHeaders())
}

// doRequestWithHeaders executes an HTTP request with custom headers
// Method Used Internally
// The result will doRequestWithHeaders(endpoint, method string, body any, headers map[string]string) (*ClientResponse, error)
func (c *Client) doRequestWithHeaders(endpoint, method string, body any, headers map[string]string) (*ClientResponse, error) {
	reader, err := parseBody(body)
	if err != nil {
		return nil, err
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
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{
		Body:       responseBody,
		StatusCode: resp.StatusCode,
	}, nil
}

// parseBody converts various types into an io.Reader
// Method Used Internally
// The result will parseBody(body any) (io.Reader, error)
func parseBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	switch v := body.(type) {
	case io.Reader:
		return v, nil
	case string:
		return strings.NewReader(v), nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(data), nil
	}
}

// WithRetry configures retry behavior using RetryConfig
// The result will WithRetry(cfg RetryConfig) Option
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

// WithTimeout sets the HTTP client's timeout
// The result will WithTimeout(d time.Duration) Option
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Timeout = d
		}
	}
}

// WithDisableKeepAlives enables or disables HTTP keep-alive connections.
// The result will WithDisableKeepAlives(disable bool) Option
func WithDisableKeepAlives(disable bool) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.DisableKeepAlives = disable
			}
		}
	}
}

// WithMaxIdleConns sets the maximum number of idle connections for the HTTP client.
// The result will WithMaxIdleConns(max int) Option
func WithMaxIdleConns(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxIdleConns = max
			}
		}
	}
}

// WithMaxConnsPerHost sets the maximum number of concurrent connections per host.
// The result will WithMaxConnsPerHost(max int) Option
func WithMaxConnsPerHost(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxConnsPerHost = max
			}
		}
	}
}

// WithMaxIdleConnsPerHost sets the maximum number of idle connections per host.
// The result will WithMaxIdleConnsPerHost(max int) Option
func WithMaxIdleConnsPerHost(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxIdleConnsPerHost = max
			}
		}
	}
}

// WithTLSConfig sets the TLS configuration for the HTTP client
// The result will WithTLSConfig(tlsConfig *tls.Config) Option
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.TLSClientConfig = tlsConfig
			}
		}
	}
}

// WithInsecureTLS enables or disables certificate verification (not recommended for production).
// The result will WithInsecureTLS(insecure bool) Option
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
// The result will WithTransport(transport http.RoundTripper) Option
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = transport
		}
	}
}

// WithCustomHTTPClient replaces the default HTTP client with a fully custom one.
// The result will WithCustomHTTPClient(client *http.Client) Option
func WithCustomHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client.Transport == nil {
			client.Transport = http.DefaultTransport
		}
		c.ClientHTTP = client
	}
}

// WithTransportConfig applies a pre-configured HTTP transport.
// The result will WithTransportConfig(tr *http.Transport) Option
func WithTransportConfig(tr *http.Transport) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = tr
		}
	}
}

// WithContext sets a custom context for the client
// The result will WithContext(ctx context.Context) Option
func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.Ctx = ctx
	}
}

// WithHeaders sets custom headers for the client
// The result will WithHeaders(headers map[string]string) Option
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.headersLock.Lock()
		defer c.headersLock.Unlock()

		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// The strategy used is Fallback Dinâmico
// RoundTrip executes the HTTP request with retry and failover logic
// The result will RoundTrip(req *http.Request) (*http.Response, error)
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var body []byte

	urls := append([]string{req.URL.String()}, rt.FailoverURLs...) // Include original URL

	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		req.Body.Close()
	}

	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		index := attempt % len(urls) // Alterna entre URLs conforme o número da tentativa
		u := urls[index]

		parsedURL, err := url.Parse(u)
		if err != nil {
			continue
		}

		req.URL = parsedURL
		if body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		resp, err = rt.Base.RoundTrip(req)
		if rt.shouldRetry(resp, err) {
			if rt.EnableLogger && rt.Logger != nil {
				rt.Logger.Warn("Retrying request",
					slog.String("url", u),
					slog.String("method", req.Method),
					slog.Int("attempt", attempt+1),
					slog.Int("failover", index+1),
				)
			}
			if resp != nil {
				resp.Body.Close()
			}
			rt.sleep(attempt)
			continue
		}
		return resp, err
	}
	return resp, fmt.Errorf("max retries exceeded: %w", err)
}

// shouldRetry determines if a request should be retried based on response and error
// Method Used Internally
// The result will shouldRetry(resp *http.Response, err error) bool
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

// sleep implements the backoff strategy for retries
// Method Used Internally
// The result will sleep(attempt int)
func (rt *RetryTransport) sleep(attempt int) {
	delay := rt.RetryDelay
	if rt.UseBackoff {
		delay = time.Duration(math.Pow(2, float64(attempt))) * rt.RetryDelay
	}
	time.Sleep(delay)
}
