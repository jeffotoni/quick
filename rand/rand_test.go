package rand

import (
	"strconv"
	"testing"
)

// TestRandomIntRange verifies that RandomInt returns values within the specified range [10, 20).
// It also ensures that no error occurs during generation.
//
// To run:
//
//	go test -v -run ^TestRandomIntRange$
func TestRandomIntRange(t *testing.T) {
	for i := 0; i < 10; i++ {
		n, err := RandomInt(10, 20)
		if err != nil {
			t.Fatal(err)
		}
		if n < 10 || n >= 20 {
			t.Errorf("expected value in range [10, 20), got %d", n)
		}
	}
}

// TestTraceIDLength ensures that the generated trace ID has a fixed length of 16 characters.
//
// To run:
//
//	go test -v -run ^TestTraceIDLength$
func TestTraceIDLength(t *testing.T) {
	id := TraceID()
	if len(id) != 16 {
		t.Errorf("expected ID length 16, got %d", len(id))
	}
}

// TestAlgoDefaultRange checks whether AlgoDefault generates a non-empty string
// and the resulting integer falls within the expected range [Start, Start+End).
//
// To run:
//
//	go test -v -run ^TestAlgoDefaultRange$
func TestAlgoDefaultRange(t *testing.T) {
	Start := 1000
	End := 9999

	for i := 0; i < 10; i++ {
		id := AlgoDefault(Start, End)
		if id == "" {
			t.Errorf("expected non-empty string, got empty")
		}
		n, err := strconv.Atoi(id)
		if err != nil {
			t.Fatal("invalid integer:", err)
		}
		// Acceptable range is [Start, Start+End)
		if n < Start || n >= Start+End {
			t.Errorf("expected value in range [%d, %d), got %d", Start, Start+End, n)
		}
	}
}
