package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// httpGoClient defines the minimal interface used (compatible with *http.Client).
type httpGoClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPClientConfig allows configuring the HTTP client's parameters.
type HTTPClientConfig struct {
	Timeout             time.Duration
	DisableKeepAlives   bool
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	TLSClientConfig     *tls.Config
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
	Ctx        context.Context
	ClientHTTP httpGoClient
	Headers    map[string]string
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
var defaultClient = NewClient()

// NewClient creates a new Client applying the provided options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		Ctx:        context.Background(),
		Headers:    map[string]string{"Content-Type": "application/json"},
		ClientHTTP: NewHTTPClientFromConfig(nil),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
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

// Get sends an HTTP GET request.
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodGet, nil)
}

// Post sends an HTTP POST request with a flexible body.
func (c *Client) Post(url string, body any) (*ClientResponse, error) {
	reader, err := parseBody(body)
	if err != nil {
		return nil, err
	}
	return c.createRequest(url, http.MethodPost, reader)
}

// Put sends an HTTP PUT request with a flexible body.
func (c *Client) Put(url string, body any) (*ClientResponse, error) {
	reader, err := parseBody(body)
	if err != nil {
		return nil, err
	}
	return c.createRequest(url, http.MethodPut, reader)
}

// Delete sends an HTTP DELETE request.
func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.createRequest(url, http.MethodDelete, nil)
}

// createRequest builds and executes the HTTP request.
func (c *Client) createRequest(url, method string, requestBody io.Reader) (*ClientResponse, error) {
	req, err := http.NewRequestWithContext(c.Ctx, method, url, requestBody)
	if err != nil {
		return nil, err
	}

	// Set the configured headers.
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.ClientHTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}, nil
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

	// If body is already an io.Reader, use it directly.
	if r, ok := body.(io.Reader); ok {
		return r, nil
	}

	// If body is a string, convert it to a reader.
	if s, ok := body.(string); ok {
		return strings.NewReader(s), nil
	}

	// Otherwise, attempt to marshal the body to JSON.
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
