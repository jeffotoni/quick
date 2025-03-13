package benchmark

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

// Struct example
type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// Pool for bytes.Buffer
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Benchmark using json.Marshal
func BenchmarkJSONMarshal(b *testing.B) {
	p := Person{Name: "John Doe", Age: 30, Email: "john@example.com"}
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(p)
		if err != nil {
			b.Fatalf("Error marshaling JSON: %v", err)
		}
	}
}

// Benchmark using json.NewEncoder(buf).Encode(p)
func BenchmarkJSONEncoder(b *testing.B) {
	p := Person{Name: "John Doe", Age: 30, Email: "john@example.com"}
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(p); err != nil {
			b.Fatalf("Error encoding JSON: %v", err)
		}
		if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] == '\n' {
			buf.Truncate(buf.Len() - 1)
		}
	}
}

// Benchmark using json.Marshal with sync.Pool
func BenchmarkJSONMarshalPool(b *testing.B) {
	p := Person{Name: "John Doe", Age: 30, Email: "john@example.com"}
	for i := 0; i < b.N; i++ {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		data, err := json.Marshal(p)
		if err != nil {
			b.Fatalf("Error marshaling JSON: %v", err)
		}

		buf.Write(data)
		bufPool.Put(buf)
	}
}

// Benchmark using json.NewEncoder(buf).Encode(p) with sync.Pool
func BenchmarkJSONEncoderPool(b *testing.B) {
	p := Person{Name: "John Doe", Age: 30, Email: "john@example.com"}
	for i := 0; i < b.N; i++ {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if err := json.NewEncoder(buf).Encode(p); err != nil {
			b.Fatalf("Error encoding JSON: %v", err)
		}

		if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] == '\n' {
			buf.Truncate(buf.Len() - 1)
		}

		bufPool.Put(buf)
	}
}
