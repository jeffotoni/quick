package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// httpGoClient defines the minimal interface used (compatible with *http.Client).
type httpGoClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryTransport implements http.RoundTripper with retry logic.
type RetryTransport struct {
	Base         http.RoundTripper // Base transport (e.g. http.DefaultTransport)
	MaxRetries   int               // Maximum number of retries
	RetryDelay   time.Duration     // Time between attempts
	UseBackoff   bool              // Enables exponential backoff
	RetryStatus  []int             // HTTP status for retry
	Logger       *slog.Logger      // New Logger field
	EnableLogger bool
}

// HTTPClientConfig allows configuring the HTTP client's parameters.
type HTTPClientConfig struct {
	Timeout             time.Duration
	DisableKeepAlives   bool
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	TLSClientConfig     *tls.Config
	MaxRetries          int           // Número máximo de tentativas
	RetryDelay          time.Duration // Tempo entre as tentativas
	RetryStatus         []int         // Lista de códigos de status HTTP que devem ser re-tentados
}

// NewHTTPClientFromConfig creates an HTTP client using the provided configuration.
// If cfg is nil, default values are used.
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

// Client represents the custom HTTP client.
type Client struct {
	Ctx          context.Context
	ClientHTTP   httpGoClient
	Headers      map[string]string
	EnableLogger bool
	Logger       *slog.Logger  // New Logger field
	MaxRetries   int           // Number of retry attempts
	RetryDelay   time.Duration // Delay between retries
	UseBackoff   bool          // Enable exponential backoff
	RetryStatus  []int         // HTTP status codes that trigger a retry

}

// ClientResponse represents the response obtained.
type ClientResponse struct {
	Body       []byte
	StatusCode int
}

// Option defines a function that modifies the Client.
type Option func(*Client)

// WithContext sets a custom context.
func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.Ctx = ctx
	}
}

// WithHeaders sets custom headers.
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.Headers = headers
	}
}

// WithHTTPClientConfig sets a custom configuration for the HTTP client.
func WithHTTPClientConfig(cfg *HTTPClientConfig) Option {
	return func(c *Client) {
		c.ClientHTTP = NewHTTPClientFromConfig(cfg)
	}
}

// defaultClient is the default client instance using standard values.
var defaultClient = New()

// New creates a new Client applying the provided options.
func New(opts ...Option) *Client {
	c := &Client{
		Ctx:        context.Background(),
		Headers:    map[string]string{"Content-Type": "application/json"},
		ClientHTTP: NewHTTPClientFromConfig(nil),
		Logger:     defaultLogger(), // Set default logger
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithLogger allows setting a custom logger.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.Logger = logger
	}
}

// Default logger (if not provided)
func defaultLogger() *slog.Logger {
	// slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// return slog.New(slog.NewTextHandler(os.Stderr, nil))
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// Global functions using the defaultClient.
func Get(url string) (*ClientResponse, error) {
	return defaultClient.Get(url)
}

// Post sends an HTTP POST request with content-type data using the default client.
func Post(url string, body any) (*ClientResponse, error) {
	return defaultClient.Post(url, body)
}

// Put sends an HTTP POST request with content-type data using the default client.
func Put(url string, body any) (*ClientResponse, error) {
	return defaultClient.Put(url, body)
}

// Delete HTTP DELETE request with content-type data using the default client.
func Delete(url string) (*ClientResponse, error) {
	return defaultClient.Delete(url)
}

// PostForm sends an HTTP POST request with form-encoded data using the default client.
func PostForm(url string, formData url.Values) (*ClientResponse, error) {
	return defaultClient.PostForm(url, formData)
}

// Get sends an HTTP GET request.
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodGet, nil)
}

// Post sends an HTTP POST request with a flexible body.
func (c *Client) Post(url string, body any) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodPost, body)
}

// Put sends an HTTP PUT request with a flexible body.
func (c *Client) Put(url string, body any) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodPut, body)
}

// Delete sends an HTTP DELETE request.
func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodDelete, nil)
}

// It automatically sets "Content-Type: application/x-www-form-urlencoded".
// PostForm sends an HTTP POST request with form-encoded data.
func (c *Client) PostForm(url string, formData url.Values) (*ClientResponse, error) {
	// Ensure the correct Content-Type header is set
	c.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	// Encode the form data and call createRequest
	return c.createRequest(url, http.MethodPost, formData.Encode())
}

// createRequest builds and executes the HTTP request.
func (c *Client) createRequest(url, method string, requestBody any) (*ClientResponse, error) {
	reader, err := parseBody(requestBody)
	if err != nil {
		return nil, err
	}

	// call NewRequestWithContext
	req, err := http.NewRequestWithContext(c.Ctx, method, url, reader)
	if err != nil {
		return nil, err
	}

	// Set headers
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// Execute request with retry logic
	resp, err := c.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}, nil
}

// executeWithRetry attempts to send the provided HTTP request multiple times,
// recreating its Body on each attempt to ensure a valid payload. If the response
// returns a non-retryable status or no error occurs, it returns immediately.
// Otherwise, it waits (optionally with exponential backoff) before retrying until
// the maximum number of retries is reached.
func (c *Client) executeWithRetry(req *http.Request) (*http.Response, error) {
	// Check if the ClientHTTP is an *http.Client and if it is using RetryTransport
	if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
		if _, ok := httpClient.Transport.(*RetryTransport); ok {
			// Log the retry attempt
			// if c.Logger != nil {
			// 	c.Logger.Warn("RetryTransport check RoundTrip",
			// 		slog.String("url", req.URL.String()),
			// 		slog.String("method", req.Method),
			// 		//slog.Int("attempt", attempt+1),
			// 		// slog.Any("error", err),
			// 	)
			// }

			// If RoundTripper is active, just execute the request without manual retry
			return httpClient.Do(req)
		}
	}

	var resp *http.Response
	var err error

	// If there is a Body in the request, we need to read it all and store it,
	// since Go consumes the Body on the first read.
	var bodyData []byte
	if req.Body != nil {
		bodyData, _ = io.ReadAll(req.Body)
		req.Body.Close() // Prevents resource leaks
		req.Body = io.NopCloser(bytes.NewReader(bodyData))
	}

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		resp, err = c.ClientHTTP.Do(req)
		if err == nil && !shouldRetry(resp, c.RetryStatus) {
			return resp, nil
		}

		if c.Logger != nil {
			// Log the retry attempt
			c.Logger.Warn("Retrying Quick request",
				slog.String("url", req.URL.String()),
				slog.String("method", req.Method),
				slog.Int("attempt", attempt+1),
				//slog.Any("error", err),
			)
		}

		// Close the response body to prevent leaking
		if resp != nil {
			resp.Body.Close()
		}

		// If there is a next attempt, we reset the Body of the request
		// to allow a new reading in the next loop.
		if len(bodyData) > 0 {
			req.Body = io.NopCloser(bytes.NewReader(bodyData))
		}

		// Wait before retrying (exponential if configured)
		waitTime := c.RetryDelay
		if c.UseBackoff {
			waitTime = time.Duration(math.Pow(2, float64(attempt))) * c.RetryDelay
		}
		time.Sleep(waitTime)
	}
	return resp, err
}

// parseBody converts the given value to an io.Reader.
// If the body is nil, returns nil.
// If body is an io.Reader, it is returned as-is.
// If body is a string, it creates a reader from the string.
// Otherwise, it marshals the body to JSON.
func parseBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	// If body is already an io.Reader, use it directly
	if r, ok := body.(io.Reader); ok {
		return r, nil
	}

	// If body is a string, convert it to a reader
	if s, ok := body.(string); ok {
		return strings.NewReader(s), nil
	}

	// Marshal struct or map to JSON
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// The result will WithRetry(maxRetries int, retryDelayStr string, retryStatusStr string) Option
func WithRetry(maxRetries int, retryDelayStr string, retryStatusStr string) Option {
	return func(c *Client) {
		c.MaxRetries = maxRetries
		c.RetryDelay, c.UseBackoff = parseRetryDelay(retryDelayStr)
		c.RetryStatus = parseRetryStatus(retryStatusStr)
	}
}

// Convert "2s", "2s-bex", "2mil-bex" to time.Duration and detect backoff
// Parses the retry delay string into a time.Duration and detects if backoff is enabled
// The result will parseRetryDelay(retryDelayStr string) (time.Duration, bool)
func parseRetryDelay(retryDelayStr string) (time.Duration, bool) {
	useBackoff := strings.Contains(retryDelayStr, "-bex")
	retryDelayStr = strings.Replace(retryDelayStr, "-bex", "", 1)

	var duration time.Duration
	switch {
	case strings.HasSuffix(retryDelayStr, "mil"):
		val, _ := strconv.Atoi(strings.TrimSuffix(retryDelayStr, "mil"))
		duration = time.Duration(val) * time.Millisecond
	case strings.HasSuffix(retryDelayStr, "s"):
		val, _ := strconv.Atoi(strings.TrimSuffix(retryDelayStr, "s"))
		duration = time.Duration(val) * time.Second
	case strings.HasSuffix(retryDelayStr, "min"):
		val, _ := strconv.Atoi(strings.TrimSuffix(retryDelayStr, "min"))
		duration = time.Duration(val) * time.Minute
	default:
		duration = 2 * time.Second // Default value
	}

	return duration, useBackoff
}

// Converts "500,502,503,504" to []int
// The result willparseRetryStatus(retryStatusStr string) []int
func parseRetryStatus(retryStatusStr string) []int {
	var statusList []int
	for _, s := range strings.Split(retryStatusStr, ",") {
		if code, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			statusList = append(statusList, code)
		}
	}
	return statusList
}

// shouldRetry(resp *http.Response, retryStatus []int) bool
func shouldRetry(resp *http.Response, retryStatus []int) bool {
	for _, status := range retryStatus {
		if resp.StatusCode == status {
			return true
		}
	}
	return false
}

// WithTimeout sets the timeout for the HTTP client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Timeout = d
		}
	}
}

// WithDisableKeepAlives enables or disables HTTP keep-alives.
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
func WithMaxIdleConns(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxIdleConns = max
			}
		}
	}
}

// WithMaxConnsPerHost sets the maximum number of connections per host.
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
func WithMaxIdleConnsPerHost(max int) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.MaxIdleConnsPerHost = max
			}
		}
	}
}

// WithTLSConfig sets the TLS configuration for the HTTP client.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				transport.TLSClientConfig = tlsConfig
			}
		}
	}
}

// WithTransport allows setting a custom HTTP transport for advanced configurations.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = transport
		}
	}
}

// WithCustomHTTPClient allows setting a fully custom *http.Client.
func WithCustomHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.ClientHTTP = client
	}
}

// WithRetryTransport enables Go's native retry mechanism in http.Transport.
func WithRetryTransport(
	maxIdleConns int,
	maxIdleConnsPerHost int,
	maxConnsPerHost int,
	disableKeepAlives bool,
	forceHTTP2 bool,
	proxy func(*http.Request) (*url.URL, error),
	tlsConfig *tls.Config,
) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = &http.Transport{
				Proxy:               proxy,               // Configure the proxy (can be nil)
				TLSClientConfig:     tlsConfig,           // TLS configuration (can be nil)
				ForceAttemptHTTP2:   forceHTTP2,          // Force HTTP/2 when available
				MaxIdleConns:        maxIdleConns,        // Idle connections in the pool
				MaxIdleConnsPerHost: maxIdleConnsPerHost, // Idle connections per host
				MaxConnsPerHost:     maxConnsPerHost,     // Maximum connections per host
				DisableKeepAlives:   disableKeepAlives,   // Disable Keep-Alive
			}
		}
	}
}

// WithRetryRoundTripper applies a custom RoundTripper for retries.
func WithRetryRoundTripper(maxRetries int, retryDelayStr string, useBackoff bool, retryStatusStr string, enableLogger bool) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			retryDelay, _ := parseRetryDelay(retryDelayStr)

			c.EnableLogger = enableLogger
			logger := defaultLogger() // Logger default

			// If the user has disabled logs, we do not initialize the logger.
			if !enableLogger {
				logger = nil
			}

			httpClient.Transport = &RetryTransport{
				Base:        http.DefaultTransport, // Uses the default transport as a base
				MaxRetries:  maxRetries,
				RetryDelay:  retryDelay,
				UseBackoff:  useBackoff,
				RetryStatus: parseRetryStatus(retryStatusStr),
				Logger:      logger,
			}
		}
	}
}

//retry RoundTrip
///

// RoundTrip executes an HTTP request with retry logic.
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		resp, err = rt.Base.RoundTrip(req) // Send the request

		// If there is no error and the status is not in the retry list, we return the response
		if err == nil && !contains(rt.RetryStatus, resp.StatusCode) {
			return resp, nil
		}

		if rt.Logger != nil && rt.EnableLogger {
			// Log the retry attempt
			rt.Logger.Warn("Retrying RoundTrip request",
				slog.String("url", req.URL.String()),
				slog.String("method", req.Method),
				slog.Int("attempt", attempt+1),
				//slog.Any("error", err),
			)
		}

		// Close the response body to avoid connection leaks
		if resp != nil {
			resp.Body.Close()
		}

		// If this is the last attempt, return the error
		if attempt == rt.MaxRetries {
			break
		}

		// Apply exponential backoff if enabled
		waitTime := rt.RetryDelay
		if rt.UseBackoff {
			waitTime = time.Duration(math.Pow(2, float64(attempt))) * rt.RetryDelay
		}
		time.Sleep(waitTime)
	}
	return resp, err
}

// contains checks if an HTTP status is in the retry list.
func contains(list []int, status int) bool {
	for _, s := range list {
		if s == status {
			return true
		}
	}
	return false
}
