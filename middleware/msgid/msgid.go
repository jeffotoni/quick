package msgid

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const MSGID_NAME string = "Msgid"

func New() http.Handler {
	var h http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(MSGID_NAME, RandAlgo1())
		w.Header().Set(MSGID_NAME, RandAlgo1())
		h.ServeHTTP(w, r)
	})
}

func msgID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(MSGID_NAME, RandAlgo1())
		w.Header().Set(MSGID_NAME, RandAlgo1())
		h.ServeHTTP(w, r)
	})
}

func RandAlgo1() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(int(rand.Intn(100000)))
}
