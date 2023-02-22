package logger

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func Logger(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	ip, port, _ := net.SplitHostPort(req.RemoteAddr)
	lw := &LogWriter{w, http.StatusOK}

	var bodySize int64
	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		bodySize = int64(len(body))
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	log.Printf("[%s]:%s %d - %s %s %v %d\n", ip, port, lw.status, req.Method, req.URL.Path, time.Since(start), bodySize)
}

type LogWriter struct {
	http.ResponseWriter
	status int
}

func (lw *LogWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}
