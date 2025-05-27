package glog_test

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jeffotoni/quick/glog"
)

// TestNew_DefaultConfig verifies that calling glog.New() with no config
// falls back to default values like os.Stdout, default layout, and INFO level.
//
// To run:
//
//	go test -v -run ^TestNew_DefaultConfig$
func TestNew_DefaultConfig(t *testing.T) {
	logger := glog.New() // no config passed, should fallback to defaults

	if logger == nil {
		t.Fatal("Expected logger instance, got nil")
	}

	// Validate internal config fallback values
	v := reflect.ValueOf(logger).Elem()
	cfg := v.FieldByName("config")

	if cfg.FieldByName("Writer").IsNil() {
		t.Error("Expected default Writer (os.Stdout), got nil")
	}

	if cfg.FieldByName("TimeFormat").String() != glog.LayoutDefault {
		t.Errorf("Expected default TimeFormat, got: %s", cfg.FieldByName("TimeFormat").String())
	}

	if glog.Level(cfg.FieldByName("Level").String()) != glog.INFO {
		t.Errorf("Expected default Level=INFO, got: %s", cfg.FieldByName("Level").String())
	}
}

// TestNew_WithCustomConfig checks if a logger initialized with a full custom config
// writes structured logs (JSON) as expected.
//
// To run:
//
//	go test -v -run ^TestNew_WithCustomConfig$
func TestNew_WithCustomConfig(t *testing.T) {
	var buf bytes.Buffer
	cfg := glog.Config{
		Format:     "json",
		Writer:     &buf,
		TimeFormat: glog.LayoutDateTime,
		Level:      glog.DEBUG,
		Separator:  " :: ",
	}

	logger := glog.New(cfg)

	if logger == nil {
		t.Fatal("Expected logger instance, got nil")
	}

	logger.Debug().Str("event", "test").Msg("ok").Send()

	if !strings.Contains(buf.String(), `"event":"test"`) {
		t.Errorf("Expected field in JSON output, got: %s", buf.String())
	}
}

// TestCallerIncluded checks that the caller field is present when .Caller() is enabled.
//
// To run:
//
//	go test -v -run ^TestCallerIncluded$
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

// TestTextFormatMinimal checks if log in text format avoids key=value by default.
//
// To run:
//
//	go test -v -run ^TestTextFormatMinimal$
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

// TestTextFormatFull validates inclusion of time, level, and msg in text format.
//
// To run:
//
//	go test -v -run ^TestTextFormatFull$
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

// TestSlogFormatOrder ensures field order is preserved in slog format output.
//
// To run:
//
//	go test -v -run ^TestSlogFormatOrder$
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

// TestSlogWithTimeAndLevel ensures that time and level fields appear when enabled in slog format.
//
// To run:
//
//	go test -v -run ^TestSlogWithTimeAndLevel$
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

// TestJSONFormatFull ensures that all fields are present and correctly serialized in JSON output.
//
// To run:
//
//	go test -v -run ^TestJSONFormatFull$
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

// TestLevelFiltering verifies that log entries below the configured level are filtered out.
//
// To run:
//
//	go test -v -run ^TestLevelFiltering$
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

// TestDefaultFallbacks ensures that a logger with empty config still works without panic.
//
// To run:
//
//	go test -v -run ^TestDefaultFallbacks$
func TestDefaultFallbacks(t *testing.T) {
	logger := glog.Set(glog.Config{})
	logger.Info().Str("fallback", "true").Send()
	// Just checking that it doesn't panic and writes to os.Stdout.
}

// TestWarnAndErrorLevels ensures WARN and ERROR level logs are correctly handled.
//
// To run:
//
//	go test -v -run ^TestWarnAndErrorLevels$
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

// TestIntAndBoolFields checks that integer and boolean fields are serialized correctly.
//
// To run:
//
//	go test -v -run ^TestIntAndBoolFields$
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

// TestTextMsgField ensures messages are included in output when using text format.
//
// To run:
//
//	go test -v -run ^TestTextMsgField$
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

// TestSlogMsgField ensures that "msg" is present in slog output.
//
// To run:
//
//	go test -v -run ^TestSlogMsgField$
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

// TestLevelPriorityFallback verifies fallback to priority 0 for unknown levels.
//
// To run:
//
//	go test -v -run ^TestLevelPriorityFallback$
func TestLevelPriorityFallback(t *testing.T) {
	if prio := glog.DEBUG; glog.TestLevelPriority("UNKNOWN") != 0 {
		t.Errorf("Expected fallback priority 0 for unknown level, got %d prio %v", glog.TestLevelPriority("UNKNOWN"), prio)
	}
}

// TestTextFormatWithUnsupportedType ensures custom/unsupported types are rendered as "null".
//
// To run:
//
//	go test -v -run ^TestTextFormatWithUnsupportedType$
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

// TestJSONFormatWithUnsupportedType ensures unsupported types are serialized as null in JSON.
//
// To run:
//
//	go test -v -run ^TestJSONFormatWithUnsupportedType$
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

// TestSmallIntOptimization checks optimized rendering for small integers.
//
// To run:
//
//	go test -v -run ^TestSmallIntOptimization$
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

// TestJSONStringEscaping ensures special characters are escaped in JSON.
//
// To run:
//
//	go test -v -run ^TestJSONStringEscaping$
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

// TestBytesToStringNoAlloc verifies byte-to-string conversion is zero-allocation.
//
// To run:
//
//	go test -v -run ^TestBytesToStringNoAlloc$
func TestBytesToStringNoAlloc(t *testing.T) {
	original := []byte("convert")
	converted := glog.TestBytesToString(original)
	if converted != "convert" {
		t.Errorf("Expected 'convert', got '%s'", converted)
	}
}

// TestNilHandlingInJSON ensures nil values are encoded as null in JSON output.
//
// To run:
//
//	go test -v -run ^TestNilHandlingInJSON$
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

// TestBoolHandlingText checks if boolean fields appear correctly in text format.
//
// To run:
//
//	go test -v -run ^TestBoolHandlingText$
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

// TestEntryFloat64Field verifies correct rendering of float64 values in text format.
//
// To run:
//
//	go test -v -run ^TestEntryFloat64Field$
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

// TestEntryDurationField ensures duration fields are rendered as time strings.
//
// To run:
//
//	go test -v -run ^TestEntryDurationField$
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

// TestEntryTimeFieldWithDefaultFormat checks that time fields use RFC3339 by default.
//
// To run:
//
//	go test -v -run ^TestEntryTimeFieldWithDefaultFormat$
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

// TestEntryTimeFieldWithCustomLayout ensures time fields respect a custom layout.
//
// To run:
//
//	go test -v -run ^TestEntryTimeFieldWithCustomLayout$
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

// TestEntryErrField ensures errors are correctly rendered in log output.
//
// To run:
//
//	go test -v -run ^TestEntryErrField$
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
