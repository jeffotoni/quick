package logger

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type testLogger struct {
	Request     *http.Request
	HandlerFunc http.HandlerFunc
}

var (
	testLoggerSuccess = testLogger{
		Request: &http.Request{
			Header:     http.Header{},
			Host:       "localhost:3000",
			RemoteAddr: "127.0.0.1:3000",
			URL: &url.URL{
				Scheme: "http",
				Host:   "letsgoquick.comm",
			},
		},
		HandlerFunc: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Printf("data req: %v", req)
		}),
	}

	testLoggerSuccessBody = testLogger{
		Request: &http.Request{
			Header:     http.Header{},
			Host:       "localhost:3000",
			RemoteAddr: "127.0.0.1:3000",
			URL: &url.URL{
				Scheme: "http",
				Host:   "letsgoquick.com",
			},
			Body: io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)),
		},
		HandlerFunc: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Printf("data req: %v", req)
		}),
	}
)

var (
	testLoggerErrorBody = testLogger{
		Request: &http.Request{
			Header:     http.Header{},
			Host:       "localhost:3000",
			RemoteAddr: "127.0.0.1:3000",
			URL: &url.URL{
				Scheme: "http",
				Host:   "letsgoquick.com",
			},
			ContentLength: 0,
			Body:          io.NopCloser(strings.NewReader(`<=>`)),
		},
		HandlerFunc: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	}
)
