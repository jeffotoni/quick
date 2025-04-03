package uuid

import (
	"crypto/rand"
	"fmt"
)

// This function is named ExampleNewV7()
// it with the Examples type.
func ExampleNewV7() {
	_, err := NewV7()
	//0195fd85-d68f-750a-b55b-6a0a82f17715
	if err != nil {
		fmt.Printf("error generating UUIDv7: %v\n", err)
		return
	}
	fmt.Printf("Generated UUIDv7")
	// Output: Generated UUIDv7
}

// This function is named ExampleNewV7FromReader()
// it with the Examples type.
func ExampleNewV7FromReader() {
	_, err := NewV7FromReader(rand.Reader)
	//0195fd88-13c5-7868-b042-f3aa076e28a7
	if err != nil {
		fmt.Printf("error generating UUIDv7 from reader: %v\n", err)
		return
	}
	fmt.Printf("Generated UUIDv7 from reader")
	// Output: Generated UUIDv7 from reader
}
