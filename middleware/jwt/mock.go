package jwt

import (
	"log"
	"net/http"
	"strings"
)

type testJWT struct {
	Request     *http.Request
	HandlerFunc http.HandlerFunc
}

var (
	successDefaultJWTHeaders = map[string][]string{}

	successCustomJWTHeaders = map[string][]string{}
)

var testJWTSuccess = testJWT{
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
