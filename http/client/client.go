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
	"strings"
	"time"
)

// httpGoClient defines the minimal interface (compatible with *http.Client).
type httpGoClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryConfig encapsulates the retry parameters, providing a clearer API.
type RetryConfig struct {
	MaxRetries int
	Delay      time.Duration
	UseBackoff bool
	Statuses   []int
	EnableLog  bool
}

// Client represents the custom HTTP client.
type Client struct {
	Ctx          context.Context
	ClientHTTP   httpGoClient
	Headers      map[string]string
	EnableLogger bool
	Logger       *slog.Logger // Logger instance
	MaxRetries   int          // Number of retry attempts
	RetryDelay   time.Duration
	UseBackoff   bool
	RetryStatus  []int
}

// RetryTransport implements the RoundTripper with retry logic.
type RetryTransport struct {
	Base         http.RoundTripper // Base transport (e.g., http.DefaultTransport)
	MaxRetries   int               // Maximum number of retries
	RetryDelay   time.Duration     // Delay between attempts
	UseBackoff   bool              // Enable exponential backoff
	RetryStatus  []int             // HTTP status codes that trigger a retry
	Logger       *slog.Logger      // Logger
	EnableLogger bool              // Flag to enable logging
}

// HTTPClientConfig allows configuring parameters for the HTTP client.
type HTTPClientConfig struct {
	Timeout             time.Duration
	DisableKeepAlives   bool
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	TLSClientConfig     *tls.Config
	MaxRetries          int           // Maximum number of retries (internal use)
	RetryDelay          time.Duration // Delay between retries (internal use)
	RetryStatus         []int         // HTTP status codes for retry (internal use)
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

// defaultClient is the default instance of Client.
var defaultClient = New()

// New creates a new Client applying the provided options.
func New(opts ...Option) *Client {
	c := &Client{
		Ctx:          context.Background(),
		Headers:      map[string]string{"Content-Type": "application/json"},
		ClientHTTP:   NewHTTPClientFromConfig(nil),
		Logger:       defaultLogger(),
		EnableLogger: false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithLogger enables or disables the custom logger.
func WithLogger(enableLogger bool) Option {
	return func(c *Client) {
		c.EnableLogger = enableLogger
	}
}

// defaultLogger returns the default logger.
func defaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// Global functions using the defaultClient.
func Get(url string) (*ClientResponse, error) {
	return defaultClient.Get(url)
}

func Post(url string, body any) (*ClientResponse, error) {
	return defaultClient.Post(url, body)
}

func Put(url string, body any) (*ClientResponse, error) {
	return defaultClient.Put(url, body)
}

func Delete(url string) (*ClientResponse, error) {
	return defaultClient.Delete(url)
}

func PostForm(url string, formData url.Values) (*ClientResponse, error) {
	return defaultClient.PostForm(url, formData)
}

// Get sends an HTTP GET request.
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodGet, nil)
}

// Post sends an HTTP POST request.
func (c *Client) Post(url string, body any) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodPost, body)
}

// Put sends an HTTP PUT request.
func (c *Client) Put(url string, body any) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodPut, body)
}

// Delete sends an HTTP DELETE request.
func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodDelete, nil)
}

// PostForm sends an HTTP POST request with form-encoded data.
// It automatically sets "Content-Type: application/x-www-form-urlencoded".
func (c *Client) PostForm(url string, formData url.Values) (*ClientResponse, error) {
	c.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	return c.createRequest(url, http.MethodPost, formData.Encode())
}

// createRequest builds and executes the HTTP request.
func (c *Client) createRequest(endpoint, httpMethod string, requestBody any) (*ClientResponse, error) {
	reader, err := parseBody(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := c.newHTTPRequest(endpoint, httpMethod, reader)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}, nil
}

func (c *Client) newHTTPRequest(endpoint, httpMethod string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(c.Ctx, httpMethod, endpoint, body)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

// readResponseBody reads the response body.
func readResponseBody(body io.ReadCloser) ([]byte, error) {
	return io.ReadAll(body)
}

// executeWithRetry attempts to send the HTTP request with retry logic.
// If ClientHTTP uses RetryTransport, the request is executed directly.
func (c *Client) executeWithRetry(req *http.Request) (*http.Response, error) {
	if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
		if _, ok := httpClient.Transport.(*RetryTransport); ok {
			return httpClient.Do(req)
		}
	}

	var resp *http.Response
	var err error
	var bodyData []byte
	if req.Body != nil {
		bodyData, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyData))
	}

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		resp, err = c.ClientHTTP.Do(req)
		if err == nil && !shouldRetry(resp, c.RetryStatus) {
			return resp, nil
		}

		if c.Logger != nil && c.EnableLogger {
			c.Logger.Warn("Retrying request",
				slog.String("url", req.URL.String()),
				slog.String("method", req.Method),
				slog.Int("attempt", attempt+1),
			)
		}

		if resp != nil {
			resp.Body.Close()
		}

		if len(bodyData) > 0 {
			req.Body = io.NopCloser(bytes.NewReader(bodyData))
		}

		waitTime := c.RetryDelay
		if c.UseBackoff {
			waitTime = time.Duration(math.Pow(2, float64(attempt))) * c.RetryDelay
		}
		time.Sleep(waitTime)
	}
	return resp, err
}

func shouldRetry(resp *http.Response, retryStatus []int) bool {
	for _, status := range retryStatus {
		if resp.StatusCode == status {
			return true
		}
	}
	return false
}

// parseBody converts the given value into an io.Reader.
// If body is nil, returns nil.
// If body is an io.Reader, it returns it directly.
// If body is a string, it creates a reader from the string.
// Otherwise, it marshals the body into JSON.
func parseBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	if r, ok := body.(io.Reader); ok {
		return r, nil
	}
	if s, ok := body.(string); ok {
		return strings.NewReader(s), nil
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

// WithRetry configures the retry behavior using RetryConfig.
func WithRetry(cfg RetryConfig) Option {
	return func(c *Client) {
		c.MaxRetries = cfg.MaxRetries
		c.RetryDelay = cfg.Delay
		c.UseBackoff = cfg.UseBackoff
		c.RetryStatus = cfg.Statuses
		c.EnableLogger = cfg.EnableLog
	}
}

// WithRetryRoundTripper configures the transport with retry using RetryConfig.
func WithRetryRoundTripper(cfg RetryConfig) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			logger := defaultLogger()
			if !cfg.EnableLog {
				logger = nil
			}
			httpClient.Transport = &RetryTransport{
				Base:         http.DefaultTransport,
				MaxRetries:   cfg.MaxRetries,
				RetryDelay:   cfg.Delay,
				UseBackoff:   cfg.UseBackoff,
				RetryStatus:  cfg.Statuses,
				Logger:       logger,
				EnableLogger: cfg.EnableLog,
			}
		}
	}
}

// WithTimeout sets the HTTP client's timeout and enables logging if desired.
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

// WithInsecureTLS allows insecure connections by setting InsecureSkipVerify.
func WithInsecureTLS(insecure bool) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
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

// WithTransport allows setting a custom HTTP transport.
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

// WithTransportConfig sets the HTTP transport for the client using a pre-configured *http.Transport.
func WithTransportConfig(tr *http.Transport) Option {
	return func(c *Client) {
		if httpClient, ok := c.ClientHTTP.(*http.Client); ok {
			httpClient.Transport = tr
		}
	}
}

// RoundTrip executes the HTTP request with retry logic in RetryTransport.
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return resp, err
		}
		req.Body.Close()
	}

	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = rt.Base.RoundTrip(req)
		if err == nil && !contains(rt.RetryStatus, resp.StatusCode) {
			return resp, nil
		}

		if rt.Logger != nil && rt.EnableLogger {
			rt.Logger.Warn("Retrying RoundTrip request",
				slog.String("url", req.URL.String()),
				slog.String("method", req.Method),
				slog.Int("attempt", attempt+1),
			)
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt == rt.MaxRetries {
			break
		}

		waitTime := rt.RetryDelay
		if rt.UseBackoff {
			waitTime = time.Duration(math.Pow(2, float64(attempt))) * rt.RetryDelay
		}
		time.Sleep(waitTime)
	}
	return resp, err
}

func contains(list []int, status int) bool {
	for _, s := range list {
		if s == status {
			return true
		}
	}
	return false
}
