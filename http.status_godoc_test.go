package goquick

import (
	"fmt"
	"testing"
)

// This function is named ExampleStatusText()
// it with the Examples type.
func ExampleStatusText() {
	fmt.Println(StatusText(200))
	fmt.Println(StatusText(404))
	fmt.Println(StatusText(500))

	// Out put:
	// OK
	// Not Found
	// Internal Server Error
}

// go test -v -run ^TestStatusText
func TestStatusText(t *testing.T) {
	if StatusText(200) != "OK" {
		t.Errorf("Expected 'OK', but got '%s'", StatusText(200))
	}

	if StatusText(404) != "Not Found" {
		t.Errorf("Expected 'Not Found', but got '%s'", StatusText(404))
	}

	if StatusText(500) != "Internal Server Error" {
		t.Errorf("Expected 'Internal Server Error', but got '%s'", StatusText(500))
	}

	if StatusText(999) != "" {
		t.Errorf("Expected empty string for unknown status code, but got '%s'", StatusText(999))
	}
}
