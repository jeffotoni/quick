package cors

import (
	"log"
	"net/http"
	"strings"
)

type testCors struct {
	Request     *http.Request
	HandlerFunc http.HandlerFunc
}

var (
	successDefaultCorsHeaders = map[string][]string{
		"Access-Control-Allow-Headers": {"Origin", "Content-Type"},
		"Access-Control-Allow-Methods": {"POST", "GET", "PUT", "DELETE", "PATH", "HEAD", "OPTIONS"},
		"Access-Control-Allow-Origin":  {"*"},
		"X-Cors":                       {"true"},
	}

	successCustomCorsHeaders = map[string][]string{
		"Access-Control-Allow-Headers": {"Origin", "Content-Type"},
		"Access-Control-Allow-Methods": {"GET", "POST"},
		"Access-Control-Allow-Origin":  {"*"},
		"X-Cors":                       {"true"},
	}
)

var testCorsSuccess = testCors{
	Request: &http.Request{
		Header: http.Header{},
	},
	HandlerFunc: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("wr uuid -> %s", rw.Header())
		log.Printf("req uuid -> %s", req.Header)
	},
	),
}

func isHeaderEqual(come, want []string) bool {
	var isAlltrue = false

	for i := 0; i < len(come); i++ {
		for j := 0; j < len(want); j++ {
			dh := strings.Split(come[i], ",")
			if len(dh) != len(want) {
				return false
			}
			if dh[j] == want[j] {
				isAlltrue = true
			}
		}
	}

	return isAlltrue
}

func isHeaderEqualDefault(come, want []string) bool {
	var isAlltrue = false

	for i := 0; i < len(come); i++ {
		for j := 0; j < len(want); j++ {
			if come[i] == want[j] {
				isAlltrue = true
			}
		}
	}

	return isAlltrue
}
