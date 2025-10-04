package quick

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"unsafe"
)

var byteBufferPoolTest = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 256) // Pre-alloc 256 bytes
		return &b
	},
}

func stringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

type mockWriter struct {
	data []byte
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)
	return len(p), nil
}

func (m *mockWriter) Header() http.Header        { return nil }
func (m *mockWriter) WriteHeader(statusCode int) {}

var message = strings.Repeat("Hello SSE Event Data", 50) // ~1KB

func BenchmarkFmtFprint(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		fmt.Fprint(w, message)
	}
}

func BenchmarkFmtFprintf(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		fmt.Fprintf(w, "data: %s\n\n", message)
	}
}

func BenchmarkIoWriteString(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		io.WriteString(w, message)
	}
}

func BenchmarkWriteBytes(b *testing.B) {
	w := &mockWriter{}
	msg := []byte(message)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		w.Write(msg)
	}
}

func BenchmarkMultipleWrites(b *testing.B) {
	w := &mockWriter{}
	msg := []byte(message)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		w.Write([]byte("data: "))
		w.Write(msg)
		w.Write([]byte("\n\n"))
	}
}

func BenchmarkStringsBuilder(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		var sb strings.Builder
		sb.WriteString("data: ")
		sb.WriteString(message)
		sb.WriteString("\n\n")
		w.Write([]byte(sb.String()))
	}
}

func BenchmarkOptimized(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		fullMessage := "event: error\ndata: " + message + "\n\n"
		w.Write([]byte(fullMessage))
	}
}

func BenchmarkUnsafe(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		w.Write(stringToBytes("event: error\ndata: "))
		w.Write(stringToBytes(message))
		w.Write(stringToBytes("\n\n"))
	}
}

func BenchmarkPooled(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		bufPtr := byteBufferPoolTest.Get().(*[]byte)
		buf := *bufPtr
		buf = buf[:0]
		buf = append(buf, "event: error\ndata: "...)
		buf = append(buf, message...)
		buf = append(buf, "\n\n"...)
		w.Write(buf)
		*bufPtr = buf
		byteBufferPoolTest.Put(bufPtr)
	}
}

// Benchmark with big menssage (10KB)
var largeMessage = strings.Repeat("Error: Something went wrong. ", 300) // ~10KB

func BenchmarkWriteBytesLarge(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		w.Write([]byte("event: error\ndata: " + largeMessage + "\n\n"))
	}
}

func BenchmarkPooledLarge(b *testing.B) {
	w := &mockWriter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.data = w.data[:0]
		bufPtr := byteBufferPoolTest.Get().(*[]byte)
		buf := *bufPtr
		buf = buf[:0]
		buf = append(buf, "event: error\ndata: "...)
		buf = append(buf, largeMessage...)
		buf = append(buf, "\n\n"...)
		w.Write(buf)
		*bufPtr = buf
		byteBufferPoolTest.Put(bufPtr)
	}
}
