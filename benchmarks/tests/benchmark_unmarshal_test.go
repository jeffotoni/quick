package benchmark

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"
)

// Example struct
type UnPerson struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// Benchmark using io.ReadAll + json.Unmarshal
func BenchmarkReadAllUnmarshal(b *testing.B) {
	p := UnPerson{Name: "John Doe", Age: 30, Email: "john@example.com"}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		b.Fatalf("Error encoding JSON: %v", err)
	}
	jsonData := buf.Bytes()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(jsonData)
		data, err := io.ReadAll(reader)
		if err != nil {
			b.Fatalf("Error reading JSON: %v", err)
		}

		var p2 UnPerson
		if err := json.Unmarshal(data, &p2); err != nil {
			b.Fatalf("Error decoding JSON: %v", err)
		}
	}
}

// Benchmark using json.NewDecoder().Decode(&v)
func BenchmarkReadAlJSONDecoder(b *testing.B) {
	p := UnPerson{Name: "John Doe", Age: 30, Email: "john@example.com"}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		b.Fatalf("Error encoding JSON: %v", err)
	}
	jsonData := buf.Bytes()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(jsonData)
		var p2 UnPerson
		if err := json.NewDecoder(reader).Decode(&p2); err != nil {
			b.Fatalf("Error decoding JSON: %v", err)
		}
	}
}

// Creating pools for reusable buffers and structs
var unbufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Benchmark using io.ReadAll + json.Unmarshal with sync.Pool
func BenchmarkReadAllUnmarshalPool(b *testing.B) {
	p := UnPerson{Name: "John Doe", Age: 30, Email: "john@example.com"}
	buf := unbufPool.Get().(*bytes.Buffer)
	buf.Reset()
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		b.Fatalf("Error encoding JSON: %v", err)
	}
	jsonData := buf.Bytes()
	unbufPool.Put(buf)

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(jsonData)
		data, err := io.ReadAll(reader)
		if err != nil {
			b.Fatalf("Error reading JSON: %v", err)
		}

		p2 := unbufPool.Get().(*UnPerson)
		if err := json.Unmarshal(data, p2); err != nil {
			b.Fatalf("Error decoding JSON: %v", err)
		}
		unbufPool.Put(p2)
	}
}

// Benchmark using json.NewDecoder().Decode(&v) with sync.Pool
func BenchmarkReadAlJSONDecoderPool(b *testing.B) {
	p := UnPerson{Name: "John Doe", Age: 30, Email: "john@example.com"}
	buf := unbufPool.Get().(*bytes.Buffer)
	buf.Reset()
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		b.Fatalf("Error encoding JSON: %v", err)
	}
	jsonData := buf.Bytes()
	unbufPool.Put(buf)

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(jsonData)
		p2 := unbufPool.Get().(*UnPerson)

		if err := json.NewDecoder(reader).Decode(p2); err != nil {
			b.Fatalf("Error decoding JSON: %v", err)
		}
		unbufPool.Put(p2)
	}
}
