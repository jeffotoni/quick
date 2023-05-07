package httpclient

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
	ClientHttp: ClientInsec,
	Headers: map[string]string{
		"Content-Type": "application/json",
	},
}

type ClientResponse struct {
	Body       []byte
	StatusCode int
	Error      error
}

var (
	ClientInsec httpGoClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        10,
			MaxConnsPerHost:     10,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ClientSec httpGoClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        10,
			MaxConnsPerHost:     10,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
)

func (c *Client) createRequest(url, method string, requestBody io.Reader) *ClientResponse {

	req, err := http.NewRequestWithContext(c.Ctx, method, url, requestBody)

	if err != nil {
		return &ClientResponse{Error: err}
	}

	if err != nil {
		return &ClientResponse{Error: err}
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.ClientHttp.Do(req)

	if err != nil {
		return &ClientResponse{Error: err}
	}

	defer resp.Body.Close()
	code := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ClientResponse{StatusCode: code, Error: err}

	}

	return &ClientResponse{Body: body, StatusCode: code, Error: err}
}
