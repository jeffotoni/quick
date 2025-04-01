package glog

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

// This function is named ExampleNew()
// it with the Examples type.
func ExampleNew() {
	var buf bytes.Buffer

	logger := New(Config{
		Format:    "slog", // Use slog para key=value
		Level:     DEBUG,
		Writer:    &buf,
		Separator: " | ",
	})

	logger.Info().
		Level().
		Str("service", "api").
		Str("status", "started").
		Msg("initialized").
		Send()

	fmt.Print(buf.String())

	// Output:
	//
	// level=INFO | service=api | status=started | msg=initialized
}

// This function is named ExampleLogger_Debug()
// it with the Examples type.
func ExampleLogger_Debug() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "slog",
		Level:  DEBUG,
		Writer: &buf,
	})

	logger.Debug().
		Level().
		Str("event", "debug_test").
		Msg("debug message").
		Send()

	fmt.Print(buf.String())

	// Output:
	// level=DEBUG event=debug_test msg=debug message
}

// This function is named ExampleLogger_Info()
// it with the Examples type.
func ExampleLogger_Info() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "slog",
		Level:  INFO,
		Writer: &buf,
	})

	logger.Info().
		Level().
		Str("event", "info_test").
		Msg("info message").
		Send()

	fmt.Print(buf.String())

	// Output:
	// level=INFO event=info_test msg=info message
}

// This function is named ExampleLogger_Warn()
// it with the Examples type.
func ExampleLogger_Warn() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "slog",
		Level:  WARN,
		Writer: &buf,
	})

	logger.Warn().
		Level().
		Str("event", "warn_test").
		Msg("warn message").
		Send()

	fmt.Print(buf.String())

	// Output:
	// level=WARN event=warn_test msg=warn message
}

// This function is named ExampleLogger_Error()
// it with the Examples type.
func ExampleLogger_Error() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "slog",
		Level:  ERROR,
		Writer: &buf,
	})

	logger.Error().
		Level().
		Str("event", "error_test").
		Msg("error message").
		Send()

	fmt.Print(buf.String())

	// Output:
	// level=ERROR event=error_test msg=error message
}

// This function is named ExampleEntry_Caller()
// it with the Examples type.
func ExampleEntry_Caller() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "json",
		Level:  DEBUG,
		Writer: &buf,
	})

	logger.Debug().
		Level().
		Str("test", "caller").
		Caller().
		Msg("with caller").
		Send()

	output := buf.String()
	fmt.Println(strings.Contains(output, `"caller"`)) // Just validate that "caller" field exists

	// Output:
	// true
}

// This function is named ExampleEntry_Any()
// it with the Examples type.
func ExampleEntry_Any() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "json",
		Level:  DEBUG,
		Writer: &buf,
	})

	data := map[string]int{"x": 1, "y": 2}

	logger.Info().
		Level().
		Any("payload", data).
		Msg("custom any field").
		Send()

	fmt.Print(buf.String())

	// Output:
	// {"payload":null,"level":"INFO","msg":"custom any field"}
}

// This function is named ExampleEntry_Str()
// it with the Examples type.
func ExampleEntry_Str() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})

	logger.Info().Str("status", "ok").Level().Msg("check").Send()
	fmt.Print(buf.String())
	// Output:
	// {"status":"ok","level":"INFO","msg":"check"}
}

// This function is named ExampleEntry_Int()
// it with the Examples type.
func ExampleEntry_Int() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})

	logger.Debug().Int("count", 3).Level().Msg("counter").Send()
	fmt.Print(buf.String())
	// Output:
	// {"count":3,"level":"DEBUG","msg":"counter"}
}

// This function is named ExampleEntry_Float64()
// it with the Examples type.
func ExampleEntry_Float64() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})

	logger.Info().Float64("pi", 3.14).Level().Msg("value").Send()
	fmt.Print(buf.String())
	// Output:
	// {"pi":3.14,"level":"INFO","msg":"value"}
}

// This function is named ExampleEntry_Duration()
// it with the Examples type.
func ExampleEntry_Duration() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})

	logger.Warn().Duration("elapsed", 2*time.Second).Level().Msg("timing").Send()
	fmt.Print(buf.String())
	// Output:
	// {"elapsed":"2s","level":"WARN","msg":"timing"}
}

// This function is named ExampleEntry_TimeField()
// it with the Examples type.
func ExampleEntry_TimeField() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})
	t := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
	logger.Debug().TimeField("fixed", t).Level().Msg("set").Send()
	fmt.Print(buf.String())
	// Output:
	// {"fixed":"2025-04-01T12:00:00Z","level":"DEBUG","msg":"set"}
}

// This function is named ExampleEntry_Err()
// it with the Examples type.
func ExampleEntry_Err() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})

	err := errors.New("file not found")
	logger.Error().Err("err", err).Level().Msg("problem").Send()
	fmt.Print(buf.String())
	// Output:
	// {"err":"file not found","level":"ERROR","msg":"problem"}
}

// This function is named ExampleEntry_AddTime()
// it with the Examples type.
func ExampleEntry_AddTime() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})
	t := time.Date(2025, 4, 1, 15, 0, 0, 0, time.UTC)
	logger.Info().AddTime("ts", t).Level().Msg("time captured").Send()
	fmt.Print(buf.String())
	// Output:
	// {"ts":"2025-04-01T15:00:00Z","level":"INFO","msg":"time captured"}
}

// This function is named ExampleEntry_Bool()
// it with the Examples type.
func ExampleEntry_Bool() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})
	logger.Debug().Bool("ok", true).Level().Msg("check").Send()
	fmt.Print(buf.String())
	// Output:
	// {"ok":true,"level":"DEBUG","msg":"check"}
}

// This function is named ExampleEntry_Msg()
// it with the Examples type.
func ExampleEntry_Msg() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})
	logger.Info().Str("op", "store").Msg("saved").Send()
	fmt.Print(buf.String())
	// Output:
	// {"op":"store","msg":"saved"}
}

// This function is named ExampleEntry_Level()
// it with the Examples type.
func ExampleEntry_Level() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})
	logger.Debug().Str("step", "start").Level().Msg("init").Send()
	fmt.Print(buf.String())
	// Output:
	// {"step":"start","level":"DEBUG","msg":"init"}
}

// This function is named ExampleEntry_Time()
// it with the Examples type.
func ExampleEntry_Time() {
	var buf bytes.Buffer
	logger := New(Config{Format: "json", Level: DEBUG, Writer: &buf})
	logger.Info().Str("check", "done").Time().Msg("completed").Send()
	out := buf.String()
	fmt.Println(strings.Contains(out, "time"))
	// Output:
	// true
}

// This function is named Example_logger_usage()
// it with the Examples type.
func Example_logger_use() {
	var buf bytes.Buffer

	logger := New(Config{
		Format: "json",
		Level:  DEBUG,
		Writer: &buf, // Redirect log output to buffer
	})

	traceID := "abc123"
	userID := "u456"
	taskID := "T789"

	logger.Debug().
		Str("trace_id", traceID).
		Str("user_id", userID).
		Str("task_id", taskID).
		Level().
		Msg("Starting task processing").
		Send()

	logger.Info().
		Str("trace_id", traceID).
		Str("elapsed", "100ms").
		Level().
		Msg("Task processed successfully").
		Send()

	logger.Warn().
		Str("trace_id", traceID).
		Bool("cache_used", false).
		Level().
		Msg("Cache disabled, using fallback").
		Send()

	logger.Error().
		Str("trace_id", traceID).
		Str("error", "simulated failure").
		Level().
		Msg("Failed to execute critical operation").
		Send()

	fmt.Print(buf.String())

	// Output:
	// {"trace_id":"abc123","user_id":"u456","task_id":"T789","level":"DEBUG","msg":"Starting task processing"}
	// {"trace_id":"abc123","elapsed":"100ms","level":"INFO","msg":"Task processed successfully"}
	// {"trace_id":"abc123","cache_used":false,"level":"WARN","msg":"Cache disabled, using fallback"}
	// {"trace_id":"abc123","error":"simulated failure","level":"ERROR","msg":"Failed to execute critical operation"}
}

// This function is named Example_basicJSON()
// it with the Examples type.
func Example_basicJSON() {
	var buf bytes.Buffer
	logger := New(Config{
		Format: "json",
		Level:  DEBUG,
		Writer: &buf,
	})

	logger.Debug().
		Str("trace_id", "abc123").
		Str("user", "jeff").
		Level().
		Msg("User login").
		Send()

	fmt.Print(buf.String())

	// Output:
	// {"trace_id":"abc123","user":"jeff","level":"DEBUG","msg":"User login"}
}

// This function is named Example_textSeparator()
// it with the Examples type.
func Example_textSeparator() {
	var buf bytes.Buffer
	logger := New(Config{
		Format:    "slog", // <- aqui trocamos para slog
		Level:     INFO,
		Writer:    &buf,
		Separator: " | ",
	})

	logger.Info().
		Level().
		Str("action", "save").
		Str("status", "ok").
		Msg("operation done").
		Send()

	fmt.Print(buf.String())

	// Output:
	// level=INFO | action=save | status=ok | msg=operation done
}

// This function is named ExampleEntry_Send_json()
// it with the Examples type.
func ExampleEntry_Send_json() {
	var buf bytes.Buffer
	logger := New(Config{
		Format: "json",
		Level:  DEBUG,
		Writer: &buf,
	})

	logger.Debug().Str("event", "json_test").Level().Msg("test").Send()
	fmt.Print(buf.String())
	// Output:
	// {"event":"json_test","level":"DEBUG","msg":"test"}
}

// This function is named ExampleEntry_Send_text()
// it with the Examples type.
func ExampleEntry_Send_text() {
	var buf bytes.Buffer
	logger := New(Config{
		Format:    "text",
		Level:     DEBUG,
		Writer:    &buf,
		Separator: " | ",
	})

	logger.Info().Str("status", "ok").Level().Msg("text mode").Send()
	fmt.Println(strings.Contains(buf.String(), "INFO"))
	// Output:
	// true
}
