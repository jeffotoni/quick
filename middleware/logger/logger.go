// Package logger provides a middleware for structured logging in Quick.
//
// This middleware captures request details such as HTTP method, path, status, response time,
// and additional metadata. It supports multiple logging formats, including:
// - "text": Standard text-based logs with customizable patterns.
// - "json": Structured JSON logs, ideal for log aggregation systems.
// - "slog": Uses Go's structured logging library (slog) with enhanced output styling.
//
// Features:
// - Supports different log formats (text, json, slog).
// - Customizable logging patterns with placeholders.
// - Captures request latency, status, user agent, and more.
// - Supports adding custom fields to logs.
package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jeffotoni/quick"
)

// ANSI color codes used for log output styling
const (
	ColorReset   = "\033[0m"  // Reset color to default
	ColorTime    = "\033[36m" // Cyan: Timestamp
	ColorLevel   = "\033[32m" // Green: Log level
	ColorMethod  = "\033[34m" // Blue: HTTP method
	ColorPath    = "\033[35m" // Magenta: Request path
	ColorStatus  = "\033[33m" // Yellow: HTTP status code
	ColorLatency = "\033[31m" // Red: Request latency
)

// Config defines the configuration for the logging middleware.
//
// Fields:
//   - Format: Log output format. Supported values: "text", "slog", "json".
//   - Pattern: The log format pattern for "text" and "slog" formats.
//   - Level: The log level threshold. Supported values: "DEBUG", "INFO", "WARN", "ERROR".
//   - CustomFields: Additional fields that will be included in log output.
type Config struct {
	Format       string            // Log format ("text", "slog", "json")
	Pattern      string            // Logging pattern
	Level        string            // Log level threshold
	CustomFields map[string]string // Additional custom fields for logging
}

// ConfigDefault provides the default logging configuration.
//
// Default values:
//   - Format: "text"
//   - Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency}\n"
var ConfigDefault = Config{
	Format:  "text",
	Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency}\n",
}

// ColorHandler is a slog.Handler that adds ANSI color to log output.
type ColorHandler struct {
	slog.Handler
	w io.Writer // Output writer for logging
}

// Handle processes the log record and applies ANSI colors based on log level.
func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = "\033[37m" // Gray for Debug
	case slog.LevelInfo:
		levelColor = "\033[32m" // Green for Info
	case slog.LevelWarn:
		levelColor = "\033[33m" // Yellow for Warning
	case slog.LevelError:
		levelColor = "\033[31m" // Red for Error
	default:
		levelColor = "\033[0m" // Default reset
	}
	reset := "\033[0m"

	timeStr := r.Time.Format(time.RFC3339)
	line := fmt.Sprintf("time=%s level=%s%s%s msg=%s\n",
		timeStr, levelColor, r.Level.String(), reset, r.Message)

	_, err := h.w.Write([]byte(line))
	return err
}

// loggerRespWriter is a ResponseWriter wrapper that captures the response status and size.
//
// Fields:
//   - status: HTTP response status code
//   - size: Total bytes written to response
type loggerRespWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Write captures the response size while writing to the underlying ResponseWriter.
func (w *loggerRespWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// WriteHeader captures the HTTP response status.
func (w *loggerRespWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// New initializes the Logger middleware for request logging.
//
// It supports different log formats including text, JSON, and slog.
//
// Parameters:
//   - config (optional): Custom logger configuration.
//
// Returns:
//   - Middleware function that wraps an HTTP handler.
//
// Example Usage:
//
//	q := quick.New()
//	q.Use(logger.New(logger.Config{
//	    Format: "text",
//	    Level: "INFO",
//	    Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency}",
//	}))
func New(config ...Config) func(http.Handler) http.Handler {
	cfg := ConfigDefault // set default value logger
	if len(config) > 0 {
		cfg = config[0]
	}

	var logger *slog.Logger // initialize logger
	var handlerOpts = &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	// set default value
	if len(cfg.Pattern) == 0 {
		cfg.Pattern = ConfigDefault.Pattern
	}

	// set default value
	if len(cfg.Level) == 0 {
		cfg.Level = "INFO"
	}

	// Select the appropriate logging format
	switch cfg.Format {
	case "slog":
		logger = slog.New(&ColorHandler{
			Handler: slog.NewTextHandler(os.Stdout, handlerOpts),
			w:       os.Stdout,
		})

	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))

	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOpts)) // Default to text format
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == quick.MethodOptions {
				return
			}
			start := time.Now()

			// Extract client IP and port from RemoteAddr
			ip, port, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
				port = "?"
			}

			// Capture request body size
			var bodySize int64
			if req.Body != nil {
				body, err := io.ReadAll(req.Body)
				if err == nil {
					bodySize = int64(len(body))
					req.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}

			// Wrap the response writer to capture status and size
			lrw := &loggerRespWriter{ResponseWriter: w}
			next.ServeHTTP(lrw, req)

			elapsed := time.Since(start)

			// Log data structure
			logData := map[string]interface{}{
				"level":         strings.ToUpper(cfg.Level),
				"time":          time.Now().Format(time.RFC3339),
				"ip":            ip,
				"port":          port,
				"method":        req.Method,
				"path":          req.URL.Path,
				"status":        lrw.status,
				"latency":       elapsed.String(),
				"body_size":     bodySize,
				"response_size": lrw.size,
				"user_agent":    req.UserAgent(),
				"referer":       req.Referer(),
				"query":         req.URL.RawQuery,
			}

			// Apply ANSI colors to log output in text mode
			colorLogData := map[string]string{
				"time":    ColorTime + logData["time"].(string) + ColorReset,
				"level":   ColorLevel + logData["level"].(string) + ColorReset,
				"method":  ColorMethod + logData["method"].(string) + ColorReset,
				"path":    ColorPath + logData["path"].(string) + ColorReset,
				"status":  ColorStatus + fmt.Sprintf("%v", logData["status"]) + ColorReset,
				"latency": ColorLatency + logData["latency"].(string) + ColorReset,
			}

			// Preserve uncolored fields
			colorLogData["ip"] = fmt.Sprintf("%v", logData["ip"])
			colorLogData["port"] = fmt.Sprintf("%v", logData["port"])
			colorLogData["body_size"] = fmt.Sprintf("%v", logData["body_size"])
			colorLogData["response_size"] = fmt.Sprintf("%v", logData["response_size"])
			colorLogData["user_agent"] = fmt.Sprintf("%v", logData["user_agent"])
			colorLogData["referer"] = fmt.Sprintf("%v", logData["referer"])
			colorLogData["query"] = fmt.Sprintf("%v", logData["query"])

			// Include custom fields
			for k, v := range cfg.CustomFields {
				logData[k] = v
				colorLogData[k] = v
			}

			switch cfg.Format {
			case "json":
				jsonData, _ := json.Marshal(logData)
				fmt.Printf("%s\n", string(jsonData)) // Log JSON format

			case "slog":
				pattern := cfg.Pattern
				for k, v := range colorLogData {
					pattern = strings.ReplaceAll(pattern, fmt.Sprintf("${%s}", k), fmt.Sprintf("%v", v))
				}

				switch strings.ToUpper(cfg.Level) {
				case "DEBUG":
					logger.Debug(pattern)
				case "WARN":
					logger.Warn(pattern)
				case "ERROR":
					logger.Error(pattern)
				default:
					logger.Info(pattern)
				}

			default:
				pattern := cfg.Pattern
				for k, v := range colorLogData {
					pattern = strings.ReplaceAll(pattern, fmt.Sprintf("${%s}", k), fmt.Sprintf("%v", v))
				}
				fmt.Printf("%s", pattern) // Log text format
			}
		})
	}
}
