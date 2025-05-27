package gcolor

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestStyleSprint verifies that the Sprint method correctly returns
// an ANSI-styled string containing the original content, ANSI prefix,
// and reset suffix.
//
// To run:
//
//	go test -v -run ^TestStyleSprint$
func TestStyleSprint(t *testing.T) {
	styled := New().Fg("red").Bg("black").Bold().Underline().Sprint("test")
	if !strings.Contains(styled, "test") {
		t.Error("Sprint should include original text")
	}
	if !strings.HasPrefix(styled, "\033[") {
		t.Error("Sprint should contain ANSI prefix")
	}
	if !strings.HasSuffix(styled, ansiReset) {
		t.Error("Sprint should end with reset code")
	}
}

// TestStyleSprintf ensures that Sprintf formats and styles text correctly,
// preserving printf-style formatting with applied ANSI styles.
//
// To run:
//
//	go test -v -run ^TestStyleSprintf$
func TestStyleSprintf(t *testing.T) {
	msg := New().Fg("green").Sprintf("hello %s", "world")
	if !strings.Contains(msg, "hello world") {
		t.Error("Sprintf should return formatted and styled string")
	}
}

// TestFluentAPI checks whether chained style methods produce correct styled output.
//
// To run:
//
//	go test -v -run ^TestFluentAPI$
func TestFluentAPI(t *testing.T) {
	style := New().Fg("blue").Bg("white").Bold().Underline()
	styled := style.Sprint("fluent")
	if !strings.Contains(styled, "fluent") {
		t.Error("Fluent API should build correct styled string")
	}
}

// TestEmptyStyle verifies that using no styles still returns a string
// with the original text and includes a reset code at the end.
//
// To run:
//
//	go test -v -run ^TestEmptyStyle$
func TestEmptyStyle(t *testing.T) {
	styled := New().Sprint("plain")
	if !strings.HasSuffix(styled, ansiReset) {
		t.Error("Even empty style should end with reset")
	}
	if !strings.Contains(styled, "plain") {
		t.Error("Original text should remain intact")
	}
}

// TestStylePrint ensures that Print writes styled output to stdout.
//
// To run:
//
//	go test -v -run ^TestStylePrint$
func TestStylePrint(t *testing.T) {
	r, w, _ := os.Pipe()
	saved := os.Stdout
	os.Stdout = w

	New().Fg("cyan").Print("print test")

	w.Close()
	os.Stdout = saved

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "print test") {
		t.Error("Print should write styled text to stdout")
	}
}

// TestStylePrintln ensures that Println writes styled text followed by a newline to stdout.
//
// To run:
//
//	go test -v -run ^TestStylePrintln$
func TestStylePrintln(t *testing.T) {
	r, w, _ := os.Pipe()
	saved := os.Stdout
	os.Stdout = w

	New().Fg("purple").Println("print line")

	w.Close()
	os.Stdout = saved

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "print line") {
		t.Error("Println should write styled text with newline to stdout")
	}
	if !strings.HasSuffix(output, "\n") {
		t.Error("Println output should end with newline")
	}
}
