package httpclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

type (
	HTTPClient interface {
		Get(URL string) *ClientResponse
	}

	httpGoClient interface {
		Do(request *http.Request) (*http.Response, error)
	}
)

type (
	HttpClient struct {
		Ctx        context.Context
		ClientHttp httpGoClient
		Headers    map[string]string
	}

	ClientGoConfig struct {
		Transport *http.Transport
		Timeout   time.Duration
	}
)

func New(hc ...HttpClient) HTTPClient {
	conf := HttpClient{
		Ctx:        context.Background(),
		ClientHttp: ClientSec,
	}

	if len(hc) > 0 {
		conf = hc[0]
	}

	return &conf
}

func NewGoClientConfig(cgc ...ClientGoConfig) httpGoClient {
	conf := ClientSec

	if len(cgc) > 0 {
		conf = &http.Client{
			Transport: cgc[0].Transport,
			Timeout:   cgc[0].Timeout,
		}
	}

	return conf
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
			MaxIdleConns:        1000,
			MaxConnsPerHost:     1000,
			MaxIdleConnsPerHost: 1000,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ClientSec httpGoClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        1000,
			MaxConnsPerHost:     1000,
			MaxIdleConnsPerHost: 1000,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
)
