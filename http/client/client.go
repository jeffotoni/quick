package client

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
)

type httpGoClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type Client struct {
	Ctx        context.Context
	ClientHttp httpGoClient
	Headers    map[string]string
}

var defaultClient = Client{
	Ctx:        context.Background(),
	ClientHttp: ClientSec,
	Headers: map[string]string{
		"Content-Type": "application/json",
	},
}

type ClientResponse struct {
	Body       []byte
	StatusCode int
}

var (
	ClientSec httpGoClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        10,
			MaxConnsPerHost:     10,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
)

// Global request Calls
func Get(url string) (*ClientResponse, error) {
	return defaultClient.Get(url)
}

func Post(url string, body io.Reader) (*ClientResponse, error) {
	return defaultClient.Post(url, body)
}

func Put(url string, body io.Reader) (*ClientResponse, error) {
	return defaultClient.Put(url, body)
}

func Delete(url string) (*ClientResponse, error) {
	return defaultClient.Delete(url)
}

// Client request Calls
func (c *Client) Get(url string) (*ClientResponse, error) {
	return c.createRequest(url, "GET", nil)
}

func (c *Client) Post(url string, body io.Reader) (*ClientResponse, error) {
	return c.createRequest(url, "POST", body)
}

func (c *Client) Put(url string, body io.Reader) (*ClientResponse, error) {
	return c.createRequest(url, "PUT", body)
}

func (c *Client) Delete(url string) (*ClientResponse, error) {
	return c.createRequest(url, "DELETE", nil)
}

func (c *Client) createRequest(url, method string, requestBody io.Reader) (*ClientResponse, error) {

	req, err := http.NewRequestWithContext(c.Ctx, method, url, requestBody)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.ClientHttp.Do(req)

	if err != nil {
		return nil, err
	}

	var errClose error
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			errClose = err
		}
	}()

	if errClose != nil {
		return nil, errClose
	}

	code := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}

	return &ClientResponse{Body: body, StatusCode: code}, nil
}
