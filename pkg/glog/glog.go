// Package glog provides a fast, zero-allocation, flexible and fluent logger for structured logging in Go.
// It supports multiple formats including:
// - text: human-readable format for development and terminal output
// - json: structured output ideal for machines and logging systems
// - slog: simplified structured log format with key=value pairs
// The logger supports fluent-style API with optional fields like time, level, and caller, which can be configured per entry.
// It is optimized for high performance with minimal allocations using internal buffer pools and unsafe conversions.

// Basic usage:
//
//	package main
//	import (
//		"github.com/jeffotoni/quick/pkg/glog"
//	)
//	func main() {
//		logger := glog.Set(glog.Config{
//			Format:    "text", // or "json", "slog"
//			Level:     glog.DEBUG,
//			Separator: " | ",
//		})
//		logger.Debug().
//			Time().Level().
//			Str("trace", "abc123").
//			Str("user", "jeff").
//			Msg("user login").
//			Send()
//	}
//
// This will produce output like:
//
//	time=2025-03-31T23:15:22-03:00 | level=DEBUG | trace=abc123 | user=jeff | msg=user login
//
// You can control the order and content of fields by chaining fluent methods, and configure output format, minimum level, and separator globally via glog.Config.
package glog

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// Level defines the severity level of a log entry.
type Level string

const (
	DEBUG Level = "DEBUG" // Fine-grained debugging information
	INFO  Level = "INFO"  // General operational entries
	WARN  Level = "WARN"  // Indications of potential issues
	ERROR Level = "ERROR" // Errors that should be investigated
)

// Predefined time layouts for log formatting.
const (
	LayoutDefault     = time.RFC3339      // Default layout: 2006-01-02T15:04:05Z07:00
	LayoutCompact     = "20060102T150405" // Compact layout: 20250101T150405
	LayoutDateTime    = "2006-01-02 15:04:05"
	LayoutDateOnly    = "2006-01-02"
	LayoutTimeOnly    = "15:04:05"
	LayoutISO8601Nano = time.RFC3339Nano // High-precision layout
)

// Constants for JSON formatting to avoid allocations.
const (
	openBrace  = '{'
	closeBrace = '}'
	comma      = ','
	quote      = '"'
	colon      = ':'
	newline    = '\n'
	nullStr    = "null"
	trueStr    = "true"
	falseStr   = "false"
	timeKey    = "time"
	levelKey   = "level"
	msgKey     = "msg"
)

// Config defines the configuration for the logger instance.
// Config defines the global configuration used to initialize a Logger.
// It can be passed to glog.Set or glog.New.
type Config struct {
	Format     string    // Output format: "text", "json", or "slog"
	Writer     io.Writer // Destination writer (e.g., os.Stdout, file, buffer)
	TimeFormat string    // Time layout used by Time() or AddTime()
	Level      Level     // Minimum level to output (DEBUG, INFO, WARN, ERROR)
	Separator  string    // Used for text and slog formats to separate fields
}

// Field represents a key-value log field with zero allocations.
type Field struct {
	key string
	val interface{} // interface{} instead of any for Go < 1.18 compatibility
}

// Entry represents a single log entry with fluent field composition.
type Entry struct {
	level     Level
	msg       string
	fields    []Field
	addTime   bool
	timeFmt   string
	addLevel  bool
	logger    *Logger
	addCaller bool
}

// Logger holds the configuration and synchronization for structured logging.
type Logger struct {
	mu     sync.RWMutex
	config Config
}

var (
	// bufferPool reuses byte buffers to minimize allocations during logging.
	bufferPool = sync.Pool{
		New: func() interface{} { return &bytes.Buffer{} },
	}

	// entryPool reuses log entries to minimize allocations per log call.
	entryPool = sync.Pool{
		New: func() interface{} {
			return &Entry{
				fields: make([]Field, 0, 8),
			}
		},
	}
)

// getBuffer retrieves a clean buffer from the pool.
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// putBuffer returns a buffer to the pool after use.
func putBuffer(buf *bytes.Buffer) {
	bufferPool.Put(buf)
}

// getEntry retrieves a reusable Entry instance from the pool.
func getEntry() *Entry {
	e := entryPool.Get().(*Entry)
	e.fields = e.fields[:0]
	e.msg = ""
	e.addTime = false
	e.addLevel = false
	// Ensure caller information is disabled by default.
	// It will only be included if the user explicitly calls .Caller().
	e.addCaller = false
	return e
}

// putEntry resets and returns the Entry to the pool to be reused.
func putEntry(e *Entry) {
	// Clear all references to help the garbage collector.
	for i := range e.fields {
		e.fields[i].key = ""
		e.fields[i].val = nil
	}
	e.logger = nil
	entryPool.Put(e)
}

// Set initializes and returns a new logger with the given configuration.
// If Writer is nil, it defaults to os.Stdout.
// If TimeFormat is not specified, it uses LayoutDefault.
// If Level is not specified, it defaults to INFO.
func Set(cfg Config) *Logger {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = LayoutDefault
	}
	if cfg.Level == "" {
		cfg.Level = INFO
	}
	return &Logger{config: cfg}
}

// New creates a new logger with optional configuration.
// If no config is passed, it uses default settings.
func New(cfgs ...Config) *Logger {
	var cfg Config
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = LayoutDefault
	}
	if cfg.Level == "" {
		cfg.Level = INFO
	}
	if cfg.Format == "" {
		cfg.Format = "text"
	}
	return &Logger{config: cfg}
}

// Debug starts a new log entry with the DEBUG level.
// It retrieves a reusable Entry instance and binds it to the logger.
func (l *Logger) Debug() *Entry {
	e := getEntry()
	e.level = DEBUG
	e.logger = l
	return e
}

// Info starts a new log entry with the INFO level.
func (l *Logger) Info() *Entry {
	e := getEntry()
	e.level = INFO
	e.logger = l
	return e
}

// Warn starts a new log entry with the WARN level.
func (l *Logger) Warn() *Entry {
	e := getEntry()
	e.level = WARN
	e.logger = l
	return e
}

// Error starts a new log entry with the ERROR level.
func (l *Logger) Error() *Entry {
	e := getEntry()
	e.level = ERROR
	e.logger = l
	return e
}

// Caller adds file:line information to the log entry
func (e *Entry) Caller() *Entry {
	e.addCaller = true
	return e
}

// Any adds a generic key-value field to the log entry.
func (e *Entry) Any(k string, v interface{}) *Entry {
	e.fields = append(e.fields, Field{key: k, val: v})
	return e
}

// Str adds a string field to the log entry.
func (e *Entry) Str(k, v string) *Entry {
	e.fields = append(e.fields, Field{key: k, val: v})
	return e
}

// Int adds an integer field to the log entry.
func (e *Entry) Int(k string, v int) *Entry {
	e.fields = append(e.fields, Field{key: k, val: v})
	return e
}

// Int adds an Float64 field to the log entry.
func (e *Entry) Float64(k string, v float64) *Entry {
	e.fields = append(e.fields, Field{key: k, val: v})
	return e
}

// Int adds an Duration field to the log entry.
func (e *Entry) Duration(k string, v time.Duration) *Entry {
	e.fields = append(e.fields, Field{key: k, val: v.String()})
	return e
}

// TimeField adds a time.Time field with custom format.
func (e *Entry) TimeField(k string, v time.Time, layout ...string) *Entry {
	format := time.RFC3339
	if len(layout) > 0 {
		format = layout[0]
	}
	e.fields = append(e.fields, Field{key: k, val: v.Format(format)})
	return e
}

// Err adds an error field to the log entry.
func (e *Entry) Err(k string, err error) *Entry {
	if err != nil {
		e.fields = append(e.fields, Field{key: k, val: err.Error()})
	}
	return e
}

// AddTime adds a time.Time field with optional format.
// If no format is provided, it uses RFC3339.
func (e *Entry) AddTime(k string, v time.Time, layout ...string) *Entry {
	format := time.RFC3339
	if len(layout) > 0 {
		format = layout[0]
	}
	e.fields = append(e.fields, Field{key: k, val: v.Format(format)})
	return e
}

// Bool adds a boolean field to the log entry.
func (e *Entry) Bool(k string, v bool) *Entry {
	e.fields = append(e.fields, Field{key: k, val: v})
	return e
}

// Msg sets the message content of the log entry.
func (e *Entry) Msg(m string) *Entry {
	e.msg = m
	return e
}

// Level includes the log level in the output.
func (e *Entry) Level() *Entry {
	e.addLevel = true
	return e
}

// Time includes the timestamp in the output.
// Accepts an optional layout format string.
func (e *Entry) Time(layout ...string) *Entry {
	e.addTime = true
	if len(layout) > 0 {
		e.timeFmt = layout[0]
	} else {
		e.timeFmt = e.logger.config.TimeFormat
	}
	return e
}

// Send finalizes the log entry and writes it to the configured writer.
func (e *Entry) Send() {
	logger := e.logger
	logger.mu.RLock()
	cfg := logger.config
	logger.mu.RUnlock()

	if !shouldLog(cfg.Level, e.level) {
		putEntry(e)
		return
	}

	if e.addCaller {
		if _, file, line, ok := runtime.Caller(2); ok {
			e.fields = append(e.fields, Field{
				key: "caller",
				val: file + ":" + strconv.Itoa(line),
			})
		}
	}

	buf := getBuffer()

	switch cfg.Format {
	case "json":
		e.jsonFormat(buf, cfg)
	default:
		e.textFormat(buf, cfg)
	}

	// Write directly to the output writer.
	cfg.Writer.Write(buf.Bytes())
	putBuffer(buf)
	putEntry(e)
}

// stringToBytes converts a string to a byte slice without allocation.
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			int
		}{s, len(s)},
	))
}

// bytesToString converts a byte slice to a string without allocation.
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// writeJSONString writes a properly escaped JSON string into the buffer.
func writeJSONString(buf *bytes.Buffer, s string) {
	buf.WriteByte(quote)
	start := 0
	for i := 0; i < len(s); i++ {
		if b := s[i]; b < 0x20 || b == '\\' || b == '"' {
			buf.WriteString(s[start:i])
			switch b {
			case '"', '\\':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				buf.WriteString(`\u00`)
				buf.WriteByte(hexChars[b>>4])
				buf.WriteByte(hexChars[b&0x0f])
			}
			start = i + 1
		}
	}
	buf.WriteString(s[start:])
	buf.WriteByte(quote)
}

var hexChars = []byte("0123456789abcdef")

// itoa converts an int64 to a string and writes it to the buffer without allocation.
func itoa(buf *bytes.Buffer, i int64) {
	// For small integers, it's faster to use direct byte conversion.
	if i >= 0 && i < 10 {
		buf.WriteByte('0' + byte(i))
		return
	}

	// AppendInt writes directly to the buffer's internal slice without allocation.
	b := strconv.AppendInt(buf.AvailableBuffer(), i, 10)
	buf.Write(b)
}

// jsonFormat encodes the log entry as a compact JSON object.
func (e *Entry) jsonFormat(buf *bytes.Buffer, cfg Config) {
	buf.WriteByte(openBrace)

	isFirst := true

	// Helper to write a JSON key with proper formatting.
	writeKey := func(key string) {
		if !isFirst {
			buf.WriteByte(comma)
		}
		writeJSONString(buf, key)
		buf.WriteByte(colon)
		isFirst = false
	}

	// Write custom fields first.
	for _, f := range e.fields {
		writeKey(f.key)

		switch v := f.val.(type) {
		case string:
			writeJSONString(buf, v)
		case int:
			itoa(buf, int64(v))
		case int64:
			itoa(buf, v)
		case uint64:
			strconv.AppendUint(buf.AvailableBuffer(), v, 10)
		case float64:
			b := strconv.AppendFloat(buf.AvailableBuffer(), v, 'f', -1, 64)
			buf.Write(b)
		case bool:
			if v {
				buf.WriteString(trueStr)
			} else {
				buf.WriteString(falseStr)
			}
		case nil:
			buf.WriteString(nullStr)
		default:
			// Unsupported types are written as "null".
			buf.WriteString(nullStr)
		}
	}

	// Include timestamp if enabled.
	if e.addTime {
		writeKey(timeKey)
		ts := time.Now().Format(e.timeFmt)
		writeJSONString(buf, ts)
	}

	// Include log level if enabled.
	if e.addLevel {
		writeKey(levelKey)
		writeJSONString(buf, string(e.level))
	}

	// Include message if present.
	if e.msg != "" {
		writeKey(msgKey)
		writeJSONString(buf, e.msg)
	}

	buf.WriteByte(closeBrace)
	buf.WriteByte(newline)
}

// textFormat renders the log entry in plain text or key=value (slog-style) format.
func (e *Entry) textFormat(buf *bytes.Buffer, cfg Config) {
	sep := cfg.Separator
	if sep == "" {
		sep = " "
	}
	sepBytes := stringToBytes(sep)
	needSep := false
	isSlog := cfg.Format == "slog"

	// Append timestamp if enabled.
	if e.addTime {
		ts := time.Now().AppendFormat(buf.AvailableBuffer(), e.timeFmt)
		if isSlog {
			buf.WriteString("time=")
		}
		buf.Write(ts)
		needSep = true
	}

	// Append level if enabled.
	if e.addLevel {
		if needSep {
			buf.Write(sepBytes)
		}
		if isSlog {
			buf.WriteString("level=")
		}
		buf.WriteString(string(e.level))
		needSep = true
	}

	// Append custom fields.
	for _, f := range e.fields {
		if needSep {
			buf.Write(sepBytes)
		}
		if isSlog {
			buf.WriteString(f.key)
			buf.WriteByte('=')
		}
		switch v := f.val.(type) {
		case string:
			buf.WriteString(v)
		case int:
			itoa(buf, int64(v))
		case int64:
			itoa(buf, v)
		case uint64:
			b := buf.AvailableBuffer()
			b = strconv.AppendUint(b, v, 10)
			buf.Write(b)
		case float64:
			b := buf.AvailableBuffer()
			b = strconv.AppendFloat(b, v, 'f', -1, 64)
			buf.Write(b)
		case bool:
			if v {
				buf.WriteString(trueStr)
			} else {
				buf.WriteString(falseStr)
			}
		default:
			// Unsupported types are rendered as "null".
			buf.WriteString(nullStr)
		}
		needSep = true
	}

	// Append the message at the end.
	if e.msg != "" {
		if needSep {
			buf.Write(sepBytes)
		}
		if isSlog {
			buf.WriteString("msg=")
		}
		buf.WriteString(e.msg)
	}

	buf.WriteByte(newline)
}

// shouldLog determines whether the log entry should be output based on current log level.
func shouldLog(min, current Level) bool {
	return levelPriority(current) >= levelPriority(min)
}

// levelPriority maps log levels to numeric values for comparison.
func levelPriority(l Level) int {
	switch l {
	case DEBUG:
		return 1
	case INFO:
		return 2
	case WARN:
		return 3
	case ERROR:
		return 4
	default:
		return 0
	}
}

// TestLevelPriority exposes the numeric level of a Level (for tests).
func TestLevelPriority(l Level) int

// TestBytesToString returns a string view of bytes without allocation (for tests).
func TestBytesToString(b []byte) string
