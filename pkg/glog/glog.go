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
//	glog.Info("User login", glog.Fields{"user": "jeff"})
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

const (
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
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
}

// Fields represent contextual dynamic fields passed per log.
type Fields map[string]any

// Logger defines the internal logger structure.
type Logger struct {
	mu     sync.RWMutex
	config Config
}

var std = &Logger{
	config: Config{
		Format:     "text",
		Writer:     os.Stdout,
		TimeFormat: time.RFC3339,
		Pattern:    "${time} ${level} ${msg} ",
		Level:      DEBUG,
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
	if cfg.Pattern == "" {
		cfg.Pattern = "${time} ${level} ${msg} "
	}
	if cfg.Level == "" {
		cfg.Level = INFO
	}
	std.mu.Lock()
	defer std.mu.Unlock()
	std.config = cfg
}

// Debug logs a DEBUG level message.
func Debug(msg string, fields ...Fields) { std.log(DEBUG, msg, getFields(fields...)) }

// Info logs an INFO level message.
func Info(msg string, fields ...Fields) { std.log(INFO, msg, getFields(fields...)) }

// Warn logs a WARN level message.
func Warn(msg string, fields ...Fields) { std.log(WARN, msg, getFields(fields...)) }

// Error logs an ERROR level message.
func Error(msg string, fields ...Fields) { std.log(ERROR, msg, getFields(fields...)) }

// Debugf logs a DEBUG formatted message.
func Debugf(format string, args ...any) { std.log(DEBUG, fmt.Sprintf(format, args...), nil) }

// Infof logs an INFO formatted message.
func Infof(format string, args ...any) { std.log(INFO, fmt.Sprintf(format, args...), nil) }

// Warnf logs a WARN formatted message.
func Warnf(format string, args ...any) { std.log(WARN, fmt.Sprintf(format, args...), nil) }

// Errorf logs an ERROR formatted message.
func Errorf(format string, args ...any) { std.log(ERROR, fmt.Sprintf(format, args...), nil) }

// log renders and writes the log message.
func (l *Logger) log(level Level, msg string, fields Fields) {
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
	for k, v := range fields {
		merged[k] = v
	}
	merged["level"] = level
	merged["time"] = ts
	merged["msg"] = msg

	if cfg.IncludeCaller {
		if _, file, line, ok := runtime.Caller(2); ok {
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

	// Replaces existing tokens in the Pattern
	for k, v := range data {
		token := "${" + k + "}"
		if strings.Contains(line, token) {
			line = strings.ReplaceAll(line, token, fmt.Sprint(colorValue(k, v)))
			used[k] = true
		}
	}

	// Detect separator (ex: " " or " | ")
	sep := detectSeparator(pattern)

	// Create a list of extras (fields not used in the Pattern)
	var extra []string
	for k, v := range data {
		if !used[k] {
			extra = append(extra, fmt.Sprintf("%s %v", k, colorValue(k, v)))
		}
	}

	// Concatenate extras, if any
	if len(extra) > 0 {
		if !strings.HasSuffix(line, sep) {
			line += sep
		}
		line += strings.Join(extra, sep)
	}

	fmt.Fprintln(w, line)
}

// shouldLog checks if the current level is loggable.
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

// getFields ensures variadic Fields is safely unwrapped.
func getFields(fields ...Fields) Fields {
	if len(fields) > 0 && fields[0] != nil {
		return fields[0]
	}
	return Fields{}
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

	// Only color the "level" field
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
