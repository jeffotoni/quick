package msgid

import (
	"log"
	"net/http"
)

type testMsgID struct {
	Request     *http.Request
	HandlerFunc http.HandlerFunc
}

var testMsgIDSuccess = testMsgID{
	Request: &http.Request{
		Header: http.Header{},
	},
	HandlerFunc: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("wr uuid -> %s", rw.Header().Get(KeyMsgID))
		log.Printf("req uuid -> %s", req.Header.Get(KeyMsgID))
	},
	),
}
