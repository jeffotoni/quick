package maxbody

import (
	"log"
	"net/http"
)

type testMaxBody struct {
	Request     *http.Request
	HandlerFunc http.HandlerFunc
}

var (
	testMaxBodySuccess = testMaxBody{
		Request: &http.Request{
			Header: http.Header{},
		},
		HandlerFunc: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Printf("req -> %v", req)
		},
		),
	}

	testMaxBodyFail = testMaxBody{
		Request: &http.Request{
			Header:        http.Header{},
			ContentLength: DefaultMaxBytes + 1,
		},
		HandlerFunc: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		},
		),
	}
)
