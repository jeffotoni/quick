package uuid

import (
	"bytes"
	"fmt"
	"strings"
)

// This function is named ExampleNew()
// it with the Examples type.
func ExampleNew() {
	u := New()
	fmt.Println(len(u), u.Version())
	// Output: 16 VERSION_4
}

// This function is named ExampleNewString()
// it with the Examples type.
func ExampleNewString() {
	s := NewString()
	fmt.Println(len(s), strings.Count(s, "-"))
	// Output: 36 4
}

// This function is named ExampleNewRandom()
// it with the Examples type.
func ExampleNewRandom() {
	u, err := NewRandom()
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(len(u), u.Version())
	// Output: 16 VERSION_4
}

// This function is named ExampleNewRandomFromReader()
// it with the Examples type.
func ExampleNewRandomFromReader() {
	reader := bytes.NewReader([]byte("abcdefghijklmnopABCDEFGHIJKLMNOP"))
	u, err := NewRandomFromReader(reader)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(u.String())
}
