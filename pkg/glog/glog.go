// Package glog provides a lightweight and flexible logging library for Go.
// It supports text, slog-style, and JSON formats with dynamic fields and custom patterns.
//
// Format behavior:
// - text:     colorized values, "field value" style, dynamic fields use pattern's first separator
// - slog:     colorized, "field=value" for all fields
// - json:     no color, serialized JSON
//
// Example:
//
//	glog.Set(glog.Config{
//		Format:  "text",
//		Pattern: "[${time}] ${level} ${msg} |",
//		Level:   glog.DEBUG,
//	})
//	glog.InfoT("User login", glog.Fields{"user": "jeff"})
//	glog.Info("User login").Str("user", "jeff").Send()
package glog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// Level defines the severity level of a log message.
type Level string

// Level defines the severity of a log entry.
// These constants represent the standard log levels used throughout glog.
// They control the verbosity of output and can be filtered via Config.Level.
const (
	DEBUG Level = "DEBUG" // Fine-grained information for debugging
	INFO  Level = "INFO"  // General operational messages (default)
	WARN  Level = "WARN"  // Indications of possible issues or important changes
	ERROR Level = "ERROR" // Critical errors that need attention
)

// Config holds global logger configuration.
type Config struct {
	Format        string            // "text", "json", "slog"
	Pattern       string            // Format pattern, e.g. "[${time}] ${level} ${msg}"
	Writer        io.Writer         // Output target (default: os.Stdout)
	TimeFormat    string            // Timestamp format (default: RFC3339)
	Level         Level             // Minimum level to log
	CustomFields  map[string]string // Global fields always included
	IncludeCaller bool              // Include file:line information
	Separator     string            // default: " ", applies only when Pattern is not defined in format text
}

// Fields is a shorthand for generic field map, kept for compatibility.
type Fields = map[string]any

// Fields represent contextual dynamic fields passed per log.
type Field struct {
	Key string
	Val any
}

// Entry is the fluent log entry builder.
type Entry struct {
	level  Level
	msg    string
	fields []Field
}

// Logger defines the internal logger structure.
type Logger struct {
	mu     sync.RWMutex
	config Config
}

// std is the default global logger instance used internally by glog.
// It holds the initial configuration and is shared across all logging calls
// unless overridden via Set(Config). This ensures thread-safe logging with
// sensible defaults: text format, standard output, timestamp, and DEBUG level.
var std = &Logger{
	config: Config{
		Format:     "text",                     // Default to human-readable text output
		Writer:     os.Stdout,                  // Logs are written to standard output
		TimeFormat: time.RFC3339,               // ISO 8601 timestamp format
		Pattern:    "${time} ${level} ${msg} ", // Basic log pattern
		Level:      DEBUG,                      // Lowest level to log by default
	},
}

// Set sets the global logger configuration.
func Set(cfg Config) {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339
	}
	if cfg.Format == "text" && cfg.Pattern == "" && cfg.Separator != "" {
		cfg.Pattern = ""
	} else if cfg.Pattern == "" {
		cfg.Pattern = "${time} ${level} ${msg} "
	}
	if cfg.Level == "" {
		cfg.Level = INFO
	}
	std.mu.Lock()
	defer std.mu.Unlock()
	std.config = cfg
}

// Info creates a fluent INFO-level log entry.
func Info(msg string) *Entry {
	return &Entry{level: INFO, msg: msg}
}

// Debug creates a fluent DEBUG-level log entry.
func Debug(msg string) *Entry {
	return &Entry{level: DEBUG, msg: msg}
}

// Warn creates a fluent WARN-level log entry.
func Warn(msg string) *Entry {
	return &Entry{level: WARN, msg: msg}
}

// Error creates a fluent ERROR-level log entry.
func Error(msg string) *Entry {
	return &Entry{level: ERROR, msg: msg}
}

// InfoT logs with a structured Fields map (compat mode).
func InfoT(msg string, fields ...map[string]any) {
	std.log(INFO, msg, flattenMap(fieldsSafe(fields...)))
}

// DebugT logs with a structured Fields map (compat mode).
func DebugT(msg string, fields ...map[string]any) {
	std.log(DEBUG, msg, flattenMap(fieldsSafe(fields...)))
}

// WarnT logs with a structured Fields map (compat mode).
func WarnT(msg string, fields ...map[string]any) {
	std.log(WARN, msg, flattenMap(fieldsSafe(fields...)))
}

// ErrorT logs with a structured Fields map (compat mode).
func ErrorT(msg string, fields ...map[string]any) {
	std.log(ERROR, msg, flattenMap(fieldsSafe(fields...)))
}

// Infof logs a formatted INFO message.
func Infof(format string, args ...any) { std.log(INFO, fmt.Sprintf(format, args...), nil) }

// Debugf logs a formatted DEBUG message.
func Debugf(format string, args ...any) { std.log(DEBUG, fmt.Sprintf(format, args...), nil) }

// Warnf logs a formatted WARN message.
func Warnf(format string, args ...any) { std.log(WARN, fmt.Sprintf(format, args...), nil) }

// Errorf logs a formatted ERROR message.
func Errorf(format string, args ...any) { std.log(ERROR, fmt.Sprintf(format, args...), nil) }

// Str adds a string field to the log entry.
func (e *Entry) Str(key, val string) *Entry {
	e.fields = append(e.fields, Field{key, val})
	return e
}

// Int adds an integer field to the log entry.
func (e *Entry) Int(key string, val int) *Entry {
	e.fields = append(e.fields, Field{key, val})
	return e
}

// Float64 adds a float64 field to the log entry.
func (e *Entry) Float64(key string, val float64) *Entry {
	e.fields = append(e.fields, Field{key, val})
	return e
}

// Bool adds a boolean field to the log entry.
func (e *Entry) Bool(key string, val bool) *Entry {
	e.fields = append(e.fields, Field{key, val})
	return e
}

// Duration adds a time.Duration field to the log entry.
func (e *Entry) Duration(key string, val time.Duration) *Entry {
	e.fields = append(e.fields, Field{key, val.String()})
	return e
}

// Time adds a time.Time field to the log entry.
func (e *Entry) Time(key string, val time.Time) *Entry {
	e.fields = append(e.fields, Field{key, val.Format(time.RFC3339)})
	return e
}

// Err adds an error field to the log entry.
func (e *Entry) Err(key string, err error) *Entry {
	if err != nil {
		e.fields = append(e.fields, Field{key, err.Error()})
	}
	return e
}

// Func adds a field by executing a function that returns any.
func (e *Entry) Func(key string, fn func() any) *Entry {
	if fn != nil {
		e.fields = append(e.fields, Field{key, fn()})
	}
	return e
}

// Any adds a generic value field to the log entry.
func (e *Entry) Any(key string, val any) *Entry {
	e.fields = append(e.fields, Field{key, val})
	return e
}

// Send finalizes and sends the log entry.
func (e *Entry) Send() {
	std.log(e.level, e.msg, e.fields)
}

// flattenMap converts a generic map[K]V into a slice of Field structs.
// This is useful for integrating arbitrary key-value pairs into the log entry system.
// Keys are converted to strings via fmt.Sprint to ensure compatibility.
func flattenMap[K comparable, V any](m map[K]V) []Field {
	result := make([]Field, 0, len(m))
	for k, v := range m {
		result = append(result, Field{Key: fmt.Sprint(k), Val: v})
	}
	return result
}

// fieldsSafe safely unwraps the first map from a variadic slice of maps.
// Returns nil if no map is provided or if the first map is nil.
// Commonly used to handle optional log fields passed to log functions.
func fieldsSafe(m ...map[string]any) map[string]any {
	if len(m) == 0 || m[0] == nil {
		return nil
	}
	return m[0]
}

// log renders and writes the log message.
func (l *Logger) log(level Level, msg string, fields []Field) {
	l.mu.RLock()
	cfg := l.config
	l.mu.RUnlock()

	if !shouldLog(cfg.Level, level) {
		return
	}

	ts := time.Now().Format(cfg.TimeFormat)
	merged := make(map[string]any)
	for k, v := range cfg.CustomFields {
		merged[k] = v
	}
	for _, f := range fields {
		merged[f.Key] = f.Val
	}

	merged["level"] = level
	merged["time"] = ts
	merged["msg"] = msg

	if cfg.IncludeCaller {
		if _, file, line, ok := runtime.Caller(3); ok {
			merged["file"] = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}

	switch cfg.Format {
	case "json":
		printJSON(cfg.Writer, merged)
	case "slog":
		printSlog(cfg.Writer, merged)
	default:
		printText(cfg.Writer, merged, cfg.Pattern)
	}
}

// printJSON writes log data as JSON.
func printJSON(w io.Writer, data map[string]any) {
	b, _ := json.Marshal(data)
	fmt.Fprintln(w, string(b))
}

// printSlog formats all fields as key=value with color.
func printSlog(w io.Writer, data map[string]any) {
	var parts []string
	keys := sortedKeys(data)
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, colorValue(k, data[k])))
	}
	fmt.Fprintln(w, strings.Join(parts, " "))
}

// printText formats fields based on pattern and appends extra fields with same separator.
func printText(w io.Writer, data map[string]any, pattern string) {
	line := pattern
	used := map[string]bool{}

	for k, v := range data {
		token := "${" + k + "}"
		if strings.Contains(line, token) {
			line = strings.ReplaceAll(line, token, fmt.Sprint(colorValue(k, v)))
			used[k] = true
		}
	}

	sep := detectSeparator(pattern)

	if pattern == "" {
		std.mu.RLock()
		sep = std.config.Separator
		if sep == "" {
			sep = " "
		}
		std.mu.RUnlock()
	}

	var extra []string
	keys := sortedKeys(data) // ensures consistent output order
	for _, k := range keys {
		if !used[k] {
			extra = append(extra, fmt.Sprintf("%s %v", k, colorValue(k, data[k])))
		}
	}

	// Create a list of extras (fields not used in the Pattern)
	if len(extra) > 0 {
		if pattern == "" {
			line = strings.Join(extra, sep)
		} else {
			if !strings.HasSuffix(line, sep) {
				line += sep
			}
			line += strings.Join(extra, sep)
		}
	}

	fmt.Fprintln(w, line)
}

// shouldLog determines whether a log message should be emitted
// based on the current message level and the minimum configured level.
// Returns true if the current level is equal to or more severe than the minimum.
func shouldLog(min Level, current Level) bool {
	priorities := map[Level]int{
		DEBUG: 1,
		INFO:  2,
		WARN:  3,
		ERROR: 4,
	}
	return priorities[current] >= priorities[min]
}

// detectSeparator infers the field separator from the pattern.
func detectSeparator(pattern string) string {
	r := regexp.MustCompile(`\}\s*([^$\{\s]+)\s*\$\{`)
	match := r.FindStringSubmatch(pattern)
	if len(match) > 1 && match[1] != "" {
		return match[1]
	}
	return " "
}

// sortedKeys returns keys sorted alphabetically.
func sortedKeys(m map[string]any) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// colorValue returns colorized string based on field name or level.
func colorValue(field string, val any) string {
	s := fmt.Sprint(val)
	if field == "level" {
		switch Level(strings.ToUpper(s)) {
		case DEBUG:
			return "\033[36m" + s + "\033[0m"
		case INFO:
			return "\033[32m" + s + "\033[0m"
		case WARN:
			return "\033[33m" + s + "\033[0m"
		case ERROR:
			return "\033[31m" + s + "\033[0m"
		}
	}
	return s
}
