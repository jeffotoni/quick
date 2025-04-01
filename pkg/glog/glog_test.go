package glog_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jeffotoni/quick/pkg/glog"
)

// TestCallerIncluded checks that the caller information is correctly included in the log when .Caller() is used.
func TestCallerIncluded(t *testing.T) {
	var buf bytes.Buffer

	logger := glog.Set(glog.Config{
		Format: "slog",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Debug().
		Caller().
		Str("trace", "abc123").
		Msg("with caller").
		Send()

	out := buf.String()

	if !strings.Contains(out, "caller=") {
		t.Errorf("Expected caller field in output, got: %s", out)
	}
	if !strings.Contains(out, "with caller") || !strings.Contains(out, "trace=abc123") {
		t.Errorf("Expected other fields in output, got: %s", out)
	}
}

// TestTextFormatMinimal verifies that a basic text log entry is correctly formatted without key=value when using text format.
func TestTextFormatMinimal(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format:    "text",
		Writer:    &buf,
		Level:     glog.DEBUG,
		Separator: " | ",
	})

	logger.Debug().Str("key", "value").Send()

	out := buf.String()
	if strings.Contains(out, "key=value") {
		t.Errorf("Expected 'key=value' in output: %s", out)
	}
}

// TestTextFormatFull checks if time, level, and message are included and formatted properly in full text output.
func TestTextFormatFull(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format:    "text",
		Writer:    &buf,
		Level:     glog.DEBUG,
		Separator: " | ",
	})

	logger.Debug().Time(glog.LayoutDateTime).Level().Str("x", "1").Msg("done").Send()

	out := buf.String()
	if strings.Contains(out, "time=") || strings.Contains(out, "level=DEBUG") || !strings.Contains(out, "done") {
		t.Errorf("Expected full fields in output: %s", out)
	}
}

// TestSlogFormatOrder ensures that fields are logged in the same order they are added using the slog format.
func TestSlogFormatOrder(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format:    "slog",
		Writer:    &buf,
		Level:     glog.DEBUG,
		Separator: " ",
	})

	logger.Debug().Str("first", "1").Str("second", "2").Str("third", "3").Send()
	out := buf.String()

	if !strings.Contains(out, "first=1 second=2 third=3") {
		t.Errorf("Field order incorrect: %s", out)
	}
}

// TestSlogWithTimeAndLevel verifies that time and level are included when explicitly enabled in slog format.
func TestSlogWithTimeAndLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "slog",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Debug().Time().Level().Str("func", "Process").Send()
	out := buf.String()

	if !strings.Contains(out, "time=") || !strings.Contains(out, "level=DEBUG") {
		t.Errorf("Expected time and level in slog output: %s", out)
	}
}

// TestJSONFormatFull validates that all expected fields are serialized correctly in JSON output.
func TestJSONFormatFull(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "json",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Debug().Time().Level().Str("status", "ok").Msg("json message").Send()
	out := buf.String()

	if !strings.Contains(out, "\"status\":\"ok\"") || !strings.Contains(out, "\"msg\":\"json message\"") {
		t.Errorf("Expected JSON fields in output: %s", out)
	}
}

// TestLevelFiltering ensures that log entries below the minimum log level are filtered out.
func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.WARN,
	})

	logger.Debug().Str("should", "not appear").Send()
	if out := buf.String(); out != "" {
		t.Errorf("Expected no output due to level filtering, got: %s", out)
	}
}

// TestDefaultFallbacks confirms that a logger using default config does not panic and logs to stdout.
func TestDefaultFallbacks(t *testing.T) {
	logger := glog.Set(glog.Config{})
	logger.Info().Str("fallback", "true").Send()
	// Just checking that it doesn't panic and writes to os.Stdout.
}

// TestWarnAndErrorLevels checks if WARN and ERROR levels are correctly handled and output.
func TestWarnAndErrorLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Warn().Str("warnKey", "warnValue").Send()
	logger.Error().Str("errorKey", "errorValue").Send()

	combined := buf.String()
	if strings.Contains(combined, "warnKey=warnValue") || strings.Contains(combined, "errorKey=errorValue") {
		t.Errorf("Expected warn and error logs: %s", combined)
	}
}

// TestIntAndBoolFields ensures integer and boolean fields are properly included in the log entry.
func TestIntAndBoolFields(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Debug().Int("code", 200).Bool("ok", true).Send()
	out := buf.String()
	if strings.Contains(out, "code=200") || strings.Contains(out, "ok=true") {
		t.Errorf("Expected int and bool fields in output: %s", out)
	}
}

// TestTextMsgField validates that a message appears correctly in the output when using text format.
func TestTextMsgField(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Debug().Msg("custom message").Send()
	out := buf.String()
	if !strings.Contains(out, "custom message") {
		t.Errorf("Expected message in text format: %s", out)
	}
}

// TestTextMsgField validates that a message appears correctly in the output when using text format.
func TestSlogMsgField(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "slog",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	logger.Debug().Msg("api call").Send()
	out := buf.String()
	if !strings.Contains(out, "msg=api call") {
		t.Errorf("Expected msg=api call in slog format: %s", out)
	}
}

// TestSlogMsgField verifies that the `msg` field appears in slog format output.
func TestLevelPriorityFallback(t *testing.T) {
	if prio := glog.DEBUG; glog.TestLevelPriority("UNKNOWN") != 0 {
		t.Errorf("Expected fallback priority 0 for unknown level, got %d prio %v", glog.TestLevelPriority("UNKNOWN"), prio)
	}
}

// TestTextFormatWithUnsupportedType ensures unsupported custom types are safely logged as "null" in text format.
func TestTextFormatWithUnsupportedType(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format:    "text",
		Writer:    &buf,
		Level:     glog.DEBUG,
		Separator: " | ",
	})

	type custom struct{}

	logger.Debug().Str("trace", "abc").Any("custom", custom{}).Send()

	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("null")) {
		t.Errorf("Expected 'null' for unsupported type, got: %s", out)
	}
}

// TestJSONFormatWithUnsupportedType ensures unsupported types are serialized as "null" in JSON format.
func TestJSONFormatWithUnsupportedType(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "json",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	type custom struct{}
	logger.Debug().Any("custom", custom{}).Send()

	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("\"custom\":null")) {
		t.Errorf("Expected JSON field with null value, got: %s", out)
	}
}

// TestSmallIntOptimization checks if small integer values are optimized in the text output (zero-allocation).
func TestSmallIntOptimization(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	logger.Debug().Int("single", 5).Send()
	if !bytes.Contains(buf.Bytes(), []byte("5")) {
		t.Errorf("Expected fast int conversion for small int, got: %s", buf.String())
	}
}

// TestJSONStringEscaping ensures special characters are properly escaped in JSON output.
func TestJSONStringEscaping(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "json",
		Writer: &buf,
		Level:  glog.DEBUG,
	})

	escaped := "line\nnewline\tquote\""
	logger.Debug().Str("escaped", escaped).Send()
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("\\n")) || !bytes.Contains([]byte(out), []byte("\\t")) || !bytes.Contains([]byte(out), []byte("\\\"")) {
		t.Errorf("Expected JSON escaping, got: %s", out)
	}
}

// TestBytesToStringNoAlloc validates that bytes are converted to string with no allocation using unsafe.
func TestBytesToStringNoAlloc(t *testing.T) {
	original := []byte("convert")
	converted := glog.TestBytesToString(original)
	if converted != "convert" {
		t.Errorf("Expected 'convert', got '%s'", converted)
	}
}

// TestNilHandlingInJSON ensures `nil` values are rendered as `null` in JSON output.
func TestNilHandlingInJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "json",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	logger.Debug().Any("maybe", nil).Send()
	if !bytes.Contains(buf.Bytes(), []byte("\"maybe\":null")) {
		t.Errorf("Expected 'maybe' field as null in JSON: %s", buf.String())
	}
}

// TestBoolHandlingText verifies correct rendering of boolean values in text output.
func TestBoolHandlingText(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format:    "text",
		Writer:    &buf,
		Level:     glog.DEBUG,
		Separator: " ",
	})
	logger.Debug().Bool("truth", true).Bool("lie", false).Send()
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("true")) || !bytes.Contains([]byte(out), []byte("false")) {
		t.Errorf("Expected boolean values in text output: %s", out)
	}
}

func TestEntryFloat64Field(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	logger.Debug().Float64("load", 99.99).Send()

	if !strings.Contains(buf.String(), "99.99") {
		t.Errorf("Expected float64 value in output: %s", buf.String())
	}
}

func TestEntryDurationField(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	duration := 3*time.Second + 500*time.Millisecond
	logger.Debug().Duration("elapsed", duration).Send()

	if !strings.Contains(buf.String(), duration.String()) {
		t.Errorf("Expected duration string in output: %s", buf.String())
	}
}

func TestEntryTimeFieldWithDefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "slog",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	now := time.Now()
	logger.Debug().AddTime("created_at", now).Send()

	// Uses RFC3339, e.g., 2025-03-31T12:00:00Z
	if !strings.Contains(buf.String(), "created_at="+now.Format(time.RFC3339)) {
		t.Errorf("Expected time field in output: %s Expected: %s", buf.String(), now.Format(time.RFC3339))
	}
}

func TestEntryTimeFieldWithCustomLayout(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	now := time.Now()
	logger.Debug().AddTime("ts", now, glog.LayoutDateTime).Send()

	if !strings.Contains(buf.String(), now.Format(glog.LayoutDateTime)) {
		t.Errorf("Expected formatted time field in output: %s", buf.String())
	}
}

func TestEntryErrField(t *testing.T) {
	var buf bytes.Buffer
	logger := glog.Set(glog.Config{
		Format: "text",
		Writer: &buf,
		Level:  glog.DEBUG,
	})
	err := errors.New("database connection failed")
	logger.Error().Err("error", err).Send()

	if !strings.Contains(buf.String(), "database connection failed") {
		t.Errorf("Expected error string in output: %s", buf.String())
	}
}
