// Package quick provides a high-performance, lightweight web framework for building
// modern HTTP applications in Go. It is designed for speed, efficiency, and simplicity.
//
// Features:
// - Middleware support for request/response processing.
// - Optimized routing with low overhead.
// - Built-in support for JSON, XML, and form parsing.
// - Efficient request handling using sync.Pool for memory optimization.
// - Customizable response handling with structured output.
//
// Quick is ideal for building RESTful APIs, microservices, and high-performance web applications.
package quick

import (
	"bytes"
	"net/http"
	"sync"
)

// Ctx Pool
var ctxPool = sync.Pool{
	New: func() interface{} {
		// Initialize a new Ctx with empty maps to avoid nil checks in usage.
		return &Ctx{
			Params:  make(map[string]string),
			Query:   make(map[string]string),
			Headers: make(map[string][]string),
		}
	},
}

// acquireCtx retrieves a Ctx instance from the sync.Pool.
func acquireCtx() *Ctx {
	return ctxPool.Get().(*Ctx)
}

// releaseCtx resets the Ctx fields and returns it to the sync.Pool for reuse.
func releaseCtx(ctx *Ctx) {
	// clear maps without reallocating
	for k := range ctx.Params {
		delete(ctx.Params, k)
	}
	for k := range ctx.Query {
		delete(ctx.Query, k)
	}
	for k := range ctx.Headers {
		delete(ctx.Headers, k)
	}

	ctx.Response = nil
	ctx.Request = nil
	ctx.bodyByte = nil
	ctx.JsonStr = ""
	ctx.resStatus = 0
	ctx.MoreRequests = 0
	ctx.App = nil

	ctxPool.Put(ctx)
}

// pooledResponseWriter wraps http.ResponseWriter and provides a buffer for potential response optimizations.
type pooledResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// responseWriterPool is a sync.Pool for pooledResponseWriter instances to reduce allocations.
var responseWriterPool = sync.Pool{
	New: func() interface{} {
		return &pooledResponseWriter{
			buf: bytes.NewBuffer(make([]byte, 0, 4096)), // initial 4KB buffer
		}
	},
}

// acquireResponseWriter retrieves a pooledResponseWriter instance from the pool.
func acquireResponseWriter(w http.ResponseWriter) *pooledResponseWriter {
	rw := responseWriterPool.Get().(*pooledResponseWriter)
	rw.ResponseWriter = w
	return rw
}

// releaseResponseWriter resets and returns the pooledResponseWriter to the pool for reuse.
func releaseResponseWriter(rw *pooledResponseWriter) {
	rw.buf.Reset()
	rw.ResponseWriter = nil
	responseWriterPool.Put(rw)
}

var jsonBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096)) // 4KB buffer
	},
}

// acquireJSONBuffer retrieves a buffer from the pool.
func acquireJSONBuffer() *bytes.Buffer {
	return jsonBufferPool.Get().(*bytes.Buffer)
}

// releaseJSONBuffer resets and returns the buffer to the pool.
func releaseJSONBuffer(buf *bytes.Buffer) {
	buf.Reset()
	jsonBufferPool.Put(buf)
}

// xmlBufferPool is a sync.Pool for optimizing XML serialization by reusing buffers.
var xmlBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096)) // 4KB buffer
	},
}

// acquireXMLBuffer retrieves a buffer from the pool.
func acquireXMLBuffer() *bytes.Buffer {
	return xmlBufferPool.Get().(*bytes.Buffer)
}

// releaseXMLBuffer resets and returns the buffer to the pool.
func releaseXMLBuffer(buf *bytes.Buffer) {
	buf.Reset()
	xmlBufferPool.Put(buf)
}

var bufferPool = sync.Pool{
	// Create new buffers with an initial capacity of 4KB.
	// Adjust this size based on expected request body sizes.
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

// acquireBuffer retrieves a *bytes.Buffer from the sync.Pool.
func acquireBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// releaseBuffer resets the buffer and places it back in the sync.Pool for reuse.
func releaseBuffer(buf *bytes.Buffer) {
	buf.Reset() // Clear any existing data
	bufferPool.Put(buf)
}

// newCtx returns a new, clean instance of Ctx
func newCtx(w http.ResponseWriter, r *http.Request, q *Quick) *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.Reset(w, r)
	ctx.App = q // Set the App reference
	return ctx
}

// Reset clears Ctx data for safe reuse
func (c *Ctx) Reset(w http.ResponseWriter, r *http.Request) {
	c.Response = w
	c.Request = r
	c.resStatus = 0

	// Clear existing maps for reuse
	for k := range c.Params {
		delete(c.Params, k)
	}
	for k := range c.Query {
		delete(c.Query, k)
	}
	for k := range c.Headers {
		delete(c.Headers, k)
	}
}
