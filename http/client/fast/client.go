package fast

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	f "github.com/valyala/fasthttp"
)

type fastHttpClient interface {
	Do(*f.Request, *f.Response) error
	DoTimeout(*f.Request, *f.Response, time.Duration) error
}

type hostClient interface {
	GetHostClient() *f.HostClient
}

type Client struct {
	Ctx            context.Context
	ClientFastHttp fastHttpClient
	Headers        map[string]string
	IsTLS          bool
	Name           string
	Timeout        time.Duration
	req            *f.Request
}

var defaultClient = Client{
	Ctx: context.Background(),
	Headers: map[string]string{
		"Content-Type": "application/json",
	},
	Name:    "quick_fast",
	IsTLS:   true,
	Timeout: 100000,
}

type ClientResponse struct {
	Body       []byte
	StatusCode int
}

// Global request Calls
func Get(url string) (*ClientResponse, error) {
	defaultClient.defaultConfig(url)
	return defaultClient.Get(url)
}

func Post(url string, body []byte) (*ClientResponse, error) {
	defaultClient.defaultConfig(url)
	return defaultClient.Post(url, body)
}

func Put(url string, body []byte) (*ClientResponse, error) {
	defaultClient.defaultConfig(url)
	return defaultClient.Put(url, body)
}

func Delete(url string) (*ClientResponse, error) {
	defaultClient.defaultConfig(url)
	return defaultClient.Delete(url)
}

// Client request Calls
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.createRequest(url, "GET", nil)
}

func (c *Client) Post(url string, body []byte) (*ClientResponse, error) {
	return c.createRequest(url, "POST", body)
}

func (c *Client) Put(url string, body []byte) (*ClientResponse, error) {
	return c.createRequest(url, "PUT", body)
}

func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.createRequest(url, "DELETE", nil)
}

func (c *Client) createRequest(url, method string, requestBody []byte) (*ClientResponse, error) {

	if c.req == nil {
		c.req = f.AcquireRequest()
		c.req.SetRequestURI(url)
		c.req.SetTimeout(c.Timeout)
	}

	if requestBody != nil {
		c.req.SetBodyRaw(requestBody)
	}

	for k, v := range c.Headers {
		c.req.Header.Set(k, v)
	}

	var resp = f.AcquireResponse()

	err := c.ClientFastHttp.DoTimeout(c.req, resp, c.Timeout)

	// defer f.ReleaseRequest(c.req)
	// defer f.ReleaseResponse(resp)

	if err != nil {
		return nil, err
	}

	return &ClientResponse{Body: resp.Body(), StatusCode: resp.StatusCode()}, nil
}

func (c *Client) defaultConfig(URL string) {
	c.req = f.AcquireRequest()
	c.req.SetRequestURI(URL)
	c.req.SetTimeout(c.Timeout)

	addr := f.AddMissingPort(string(c.req.URI().Host()), c.IsTLS)
	var cDefault fastHttpClient = &f.HostClient{
		Addr:                     addr,
		Name:                     c.Name,
		NoDefaultUserAgentHeader: true,
		DialDualStack:            true,
		IsTLS:                    c.IsTLS,
		Dial: func(addr string) (net.Conn, error) {
			return (&f.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: 5 * time.Minute,
			}).DialTimeout(addr, time.Hour)
		},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
		MaxConns:                      100,
		MaxConnDuration:               100,
		MaxIdleConnDuration:           100,
		MaxIdemponentCallAttempts:     100,
		ReadBufferSize:                5 * 1024 * 1024,
		WriteBufferSize:               5 * 1024 * 1024,
		ReadTimeout:                   c.Timeout,
		WriteTimeout:                  c.Timeout,
		MaxResponseBodySize:           100,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		SecureErrorLogMessage:         true,
		MaxConnWaitTimeout:            10000,
		StreamResponseBody:            false,
	}

	c.ClientFastHttp = cDefault
}
