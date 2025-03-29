// Package glog provides a lightweight and flexible logging library for Go.
// It supports multiple output formats (text, JSON, slog-like) and dynamic log patterns.
// The logger can be configured globally or used independently.
//
// Example usage:
//
//	package main
//
//	import (
//		"github.com/jeffotoni/quick/glog"
//	)
//
//	func main() {
//		Set(Config{
//			Format:     "text",
//			Pattern:    "[${time}] ${level} ${msg}",
//			TimeFormat: "2006-01-02 15:04:05",
//		})
//
//		Infof("Server started on port %d", 8080)
//		Error("Failed to connect to database", Fields{"retry": true})
//	}
//
// This will output colorized logs to stdout with the defined format and values.
package glog

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestGlogTextFormat verifies text output format and placeholder replacement.
func TestGlogTextFormat(t *testing.T) {
	var buf bytes.Buffer

	Set(Config{
		Format:     "text",
		Writer:     &buf,
		TimeFormat: time.RFC3339,
		Pattern:    "[${time}] ${level} ${msg}",
		Level:      DEBUG,
	})

	Info("Test log entry", Fields{"user": "jeff"})
	output := buf.String()

	if !strings.Contains(output, "Test log entry") {
		t.Errorf("Expected log message not found in output: %s", output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log level 'INFO' not found in output: %s", output)
	}
}

// TestGlogJsonFormat verifies JSON log output with dynamic fields.
func TestGlogJsonFormat(t *testing.T) {
	var buf bytes.Buffer

	Set(Config{
		Format:     "json",
		Writer:     &buf,
		TimeFormat: time.RFC3339,
		Level:      DEBUG,
	})

	Error("Something went wrong", Fields{"code": 500})
	output := buf.String()

	if !strings.Contains(output, "\"msg\":\"Something went wrong\"") {
		t.Errorf("Expected JSON 'msg' not found: %s", output)
	}
	if !strings.Contains(output, "\"level\":\"ERROR\"") {
		t.Errorf("Expected JSON 'level' not found: %s", output)
	}
	if !strings.Contains(output, "\"code\":500") {
		t.Errorf("Expected JSON field 'code' not found: %s", output)
	}
}

// TestGlogLevelFiltering ensures messages below the configured level are skipped.
func TestGlogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	Set(Config{
		Format:     "text",
		Writer:     &buf,
		TimeFormat: time.RFC3339,
		Pattern:    "[${time}] ${level} ${msg}",
		Level:      WARN, // Minimum level set to WARN
	})

	Info("This should not appear", nil)
	Warn("This should appear", nil)

	output := buf.String()

	if strings.Contains(output, "This should not appear") {
		t.Errorf("INFO log should have been filtered out, but it appeared")
	}
	if !strings.Contains(output, "This should appear") {
		t.Errorf("WARN log expected but not found")
	}
}

// TestGlogCustomFields validates inclusion of global custom fields in the output.
func TestGlogCustomFields(t *testing.T) {
	var buf bytes.Buffer

	Set(Config{
		Format:       "text",
		Writer:       &buf,
		TimeFormat:   time.RFC3339,
		Pattern:      "[${time}] ${level} ${msg} (${service})",
		Level:        DEBUG,
		CustomFields: map[string]string{"service": "quick-api"},
	})

	Debug("Debugging app", nil)
	output := buf.String()

	if !strings.Contains(output, "quick-api") {
		t.Errorf("Custom field 'service' not found in log output: %s", output)
	}
}

// TestGlogFormattedLogs checks formatted logging using Infof, Debugf, etc.
func TestGlogFormattedLogs(t *testing.T) {
	var buf bytes.Buffer

	Set(Config{
		Format:     "text",
		Writer:     &buf,
		TimeFormat: time.RFC3339,
		Pattern:    "[${time}] ${level} ${msg}",
		Level:      DEBUG,
	})

	Infof("Hello %s, your score is %d", "Arthur", 99)
	output := buf.String()

	if !strings.Contains(output, "Arthur") || !strings.Contains(output, "99") {
		t.Errorf("Formatted log output incorrect: %s", output)
	}
}

// TestGlogFullCoverage tests glog behaviors for all formats, patterns, dynamic fields, colors, and separators.
func TestGlogFullCoverage(t *testing.T) {
	var buf bytes.Buffer

	// TEXT format with proper pattern spacing
	Set(Config{
		Format:     "text",
		Writer:     &buf,
		TimeFormat: "2006-01-02 15:04:05",
		Pattern:    "${time} ${level} ${msg} ",
		Level:      DEBUG,
		CustomFields: map[string]string{
			"service": "example-api",
		},
	})

	Debug("Debugging", Fields{"user": "jeff"})
	Info("User login", Fields{"trace": "abc123"})
	Warn("Warning issued")
	Error("Error occurred", Fields{"retry": true})

	out := buf.String()

	if !strings.Contains(out, "Debugging") {
		t.Errorf("Expected 'Debugging': %s", out)
	}
	if !strings.Contains(out, "user jeff") {
		t.Errorf("Expected 'user jeff': %s", out)
	}
	if !strings.Contains(out, "service example-api") {
		t.Errorf("Expected 'service example-api': %s", out)
	}
	if !strings.Contains(out, "trace abc123") {
		t.Errorf("Expected 'trace abc123': %s", out)
	}
	if !strings.Contains(out, "retry true") {
		t.Errorf("Expected 'retry true': %s", out)
	}

	buf.Reset()

	// SLOG format with key=value and color output
	Set(Config{
		Format:     "slog",
		Writer:     &buf,
		Pattern:    "${time} ${level} ${msg} ",
		TimeFormat: time.RFC3339,
		Level:      INFO,
		CustomFields: map[string]string{
			"env": "dev",
		},
	})
	Info("Slog test", Fields{"xid": "xyz"})
	out = buf.String()

	if !strings.Contains(out, "xid=xyz") ||
		!strings.Contains(out, "env=dev") ||
		!strings.Contains(out, "msg=Slog test") {
		t.Errorf("Expected slog format fields missing: %s", out)
	}

	buf.Reset()

	// JSON format with serialized fields
	Set(Config{
		Format:     "json",
		Writer:     &buf,
		TimeFormat: time.RFC3339,
		Level:      DEBUG,
	})
	Debug("Json test", Fields{"id": 123, "ok": true})
	out = buf.String()

	if !strings.Contains(out, "\"id\":123") || !strings.Contains(out, "\"ok\":true") {
		t.Errorf("Expected fields in JSON output: %s", out)
	}
	if strings.Contains(out, "\033[") {
		t.Errorf("JSON output must not contain ANSI color codes: %s", out)
	}

	buf.Reset()

	// TEXT with custom separator " | " in pattern
	Set(Config{
		Format:     "text",
		Writer:     &buf,
		Pattern:    "${time} | ${level} | ${msg} | ",
		TimeFormat: time.RFC3339,
	})
	Info("Pattern separator test", Fields{
		"ip":     "127.0.0.1",
		"region": "us-east",
	})
	out = buf.String()

	if !strings.Contains(out, "ip 127.0.0.1") || !strings.Contains(out, "region us-east") {
		t.Errorf("Expected dynamic fields in output: %s", out)
	}
	if strings.Count(out, "|") < 3 {
		t.Errorf("Expected at least 3 pipes (|) in log: %s", out)
	}
}

func TestGlogFormattedMethods(t *testing.T) {
	var buf bytes.Buffer

	// Setup with TEXT format so we can match string content easily
	Set(Config{
		Format:     "text",
		Writer:     &buf,
		TimeFormat: "2006-01-02 15:04:05",
		Pattern:    "${time} ${level} ${msg} ",
		Level:      DEBUG,
	})

	buf.Reset()
	Debugf("Hello %s", "Jeff")
	if !strings.Contains(buf.String(), "Hello Jeff") {
		t.Errorf("Expected formatted debug message to appear")
	}

	buf.Reset()
	Warnf("Warning at step %d", 3)
	if !strings.Contains(buf.String(), "Warning at step 3") {
		t.Errorf("Expected formatted warning message to appear")
	}

	buf.Reset()
	Errorf("Failed with code %d", 500)
	if !strings.Contains(buf.String(), "Failed with code 500") {
		t.Errorf("Expected formatted error message to appear")
	}
}

func GetConfig() Config {
	std.mu.RLock()
	defer std.mu.RUnlock()
	return std.config
}

func TestGlogSet_Defaults(t *testing.T) {
	// Set config without Writer and TimeFormat
	cfg := Config{
		Format:       "text",
		Level:        DEBUG,
		Pattern:      "", // intentionally empty
		CustomFields: map[string]string{"app": "test"},
		// Writer and TimeFormat left empty on purpose
	}

	Set(cfg)

	stdConfig := GetConfig()

	if stdConfig.Writer == nil {
		t.Errorf("Expected Writer to be defaulted to os.Stdout")
	}

	if stdConfig.TimeFormat != time.RFC3339 {
		t.Errorf("Expected default TimeFormat to be RFC3339, got %s", stdConfig.TimeFormat)
	}

	if stdConfig.Pattern != "${time} ${level} ${msg} " {
		t.Errorf("Expected default pattern to be injected")
	}
}

func TestDetectSeparatorFallback(t *testing.T) {
	sep := detectSeparator("${time}${level}${msg}")
	if sep != " " {
		t.Errorf("Expected fallback separator ' ', got: %q", sep)
	}
}

func TestIncludeCallerAddsFileLine(t *testing.T) {
	var buf bytes.Buffer

	// Set logger with IncludeCaller true
	Set(Config{
		Format:        "text",
		Writer:        &buf,
		Pattern:       "${time} ${level} ${msg} ", // file not in pattern
		TimeFormat:    "2006-01-02 15:04:05",
		Level:         DEBUG,
		IncludeCaller: true,
	})

	buf.Reset()
	Info("Testing caller inclusion")

	logOutput := buf.String()

	// Expect that the auto-included 'file' field appears as dynamic field (since not in pattern)
	if !strings.Contains(logOutput, "file") || !strings.Contains(logOutput, ".go:") {
		t.Errorf("Expected 'file' field with .go:<line> to be included, got: %s", logOutput)
	}
}
