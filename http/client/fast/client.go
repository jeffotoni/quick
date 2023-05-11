package fast

import (
	"context"
	"crypto/tls"

	f "github.com/valyala/fasthttp"
)

type fastHttpClient interface {
	Do(request *f.Request, resp *f.Response) error
}

type hostClient interface {
	GetHostClient() *f.HostClient
}

type Client struct {
	Ctx            context.Context
	ClientFastHttp fastHttpClient
	Headers        map[string]string
	IsTLS          bool
	req            *f.Request
}

var defaultClient = Client{
	Ctx: context.Background(),
	Headers: map[string]string{
		"Content-Type": "application/json",
	},
	IsTLS: true,
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
		c.req = new(f.Request)
		c.req.SetRequestURI(url)

	}

	if requestBody != nil {
		c.req.SetBody(requestBody)
	}

	for k, v := range c.Headers {
		c.req.Header.Set(k, v)
	}

	var resp *f.Response

	err := c.ClientFastHttp.Do(c.req, resp)

	if err != nil {
		return nil, err
	}

	return &ClientResponse{Body: resp.Body(), StatusCode: resp.StatusCode()}, nil
}

func (c *Client) defaultConfig(URL string) {
	c.req = new(f.Request)
	c.req.SetRequestURI(URL)
	var cDefault fastHttpClient = &f.HostClient{
		Addr:                     f.AddMissingPort(string(c.req.URI().Host()), c.IsTLS),
		Name:                     "quick_fast",
		NoDefaultUserAgentHeader: false,
		DialDualStack:            false,
		IsTLS:                    false,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
		MaxConns:                      100,
		MaxConnDuration:               100,
		MaxIdleConnDuration:           100,
		MaxIdemponentCallAttempts:     100,
		ReadBufferSize:                100,
		WriteBufferSize:               100,
		ReadTimeout:                   10000,
		WriteTimeout:                  10000,
		MaxResponseBodySize:           100,
		DisableHeaderNamesNormalizing: false,
		DisablePathNormalizing:        false,
		SecureErrorLogMessage:         true,
		MaxConnWaitTimeout:            100,
		ConnPoolStrategy:              100,
		StreamResponseBody:            false,
	}

	c.ClientFastHttp = cDefault
}
