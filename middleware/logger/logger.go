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
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jeffotoni/quick"
)

// Global storage for context data per request
var (
	requestContextData sync.Map // map[*http.Request]map[string]string
)

// var defaulJSON string = `{
// 				"level":  "",
// 				"time":    "",
// 				"ip":      "",
// 				"port":    "",
// 				"path":    "",
// 				"status":  "",
// 				"latency": "",
// 				"host":    "",
// 				"method": "",
// 				"headers": "",
// 				"body": "",
// 				"size": "",
// 				"user_agent": "",
// 				"referer": "",
// 				"query": "",

// 				"request_method":    "",
// 				"request_headers":    "",
// 				"request_body":     "",
// 				"request_size":       "",
// 				"request_user_agent": "",
// 				"request_referer":   "",
// 				"request_query":      "",
// 				"request_path":       "",

// 				"response_status": "",
// 				"response_size":    "",
// 				"response_headers": "",
// 				"response_body":    ""
// 			}` // json format`

// setRequestContextData stores context data for a specific request (called by SetContext)
func setRequestContextData(req *http.Request, data map[string]any) {
	// Instead of replacing, accumulate the data
	if existing, ok := requestContextData.Load(req); ok {
		if existingMap, ok := existing.(map[string]any); ok {
			// Merge with existing data
			merged := make(map[string]any)
			for k, v := range existingMap {
				merged[k] = v
			}
			for k, v := range data {
				merged[k] = v // New values override existing ones
			}
			requestContextData.Store(req, merged)
			return
		}
	}
	// No existing data, store as-is
	requestContextData.Store(req, data)
}

// SetRequestContextData - exported function for quick package to call
func SetRequestContextData(req *http.Request, data map[string]any) {
	setRequestContextData(req, data)
}

// init registers the callback with the quick package
func init() {
	// Set the callback in the quick package to capture context data
	quick.ContextDataCallback = setRequestContextData
}

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
	Pattern: "[[${time}] ${level} ${method} ${path} ${status} - ${latency}\n",
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
//   - body: Captured response body
//   - headers: Captured response headers
type loggerRespWriter struct {
	http.ResponseWriter
	status  int
	size    int
	body    []byte
	headers http.Header
	request *http.Request // Store request reference to access updated context
}

// Write captures the response size and body while writing to the underlying ResponseWriter.
func (w *loggerRespWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	// Capture response body
	w.body = append(w.body, b...)
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// WriteHeader captures the HTTP response status and headers.
func (w *loggerRespWriter) WriteHeader(status int) {
	w.status = status
	// Capture response headers before they are written
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	for k, v := range w.ResponseWriter.Header() {
		w.headers[k] = v
	}
	w.ResponseWriter.WriteHeader(status)
}

// SetRequest updates the stored request (used by context updates)
func (w *loggerRespWriter) SetRequest(req *http.Request) {
	w.request = req
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
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}

			port := getPort(req)

			// Capture request body size and detect multipart uploads
			var bodySize int64
			var bodyVal string
			var isMultipartForm bool
			var multipartInfo map[string]any

			contentType := strings.ToLower(req.Header.Get("Content-Type"))
			isMultipartForm = strings.HasPrefix(contentType, "multipart/form-data")

			if req.Body != nil {
				body, err := io.ReadAll(req.Body)
				if err == nil {
					bodySize = int64(len(body))

					if isMultipartForm {
						// For multipart uploads, don't store the actual body content
						bodyVal = "--multipart/form-data--"
						// Extract multipart file information
						multipartInfo = extractMultipartInfo(req, body)
					} else {
						bodyVal = string(body)
					}

					req.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}

			// Wrap the response writer to capture status, size, body, and headers
			lrw := &loggerRespWriter{
				ResponseWriter: w,
				headers:        make(http.Header),
				request:        req, // Store request to access updated context later
			}
			next.ServeHTTP(lrw, req)

			elapsed := time.Since(start)

			dynamicContextData := make(map[string]any)

			ctx := req.Context()
			if lrw.request != nil {
				ctx = lrw.request.Context() // Use updated context if available
			}

			if ctxData := ctx.Value("__quick_context_data__"); ctxData != nil {
				if contextMap, ok := ctxData.(map[string]any); ok {
					for key, value := range contextMap {
						if value != nil {
							dynamicContextData[key] = value
						}
					}
				}
			} else {
				if data, ok := requestContextData.Load(req); ok {
					if contextMap, ok := data.(map[string]any); ok {
						for key, value := range contextMap {
							if value != nil {
								dynamicContextData[key] = value
							}
						}
					}
					// Clean up the global map to prevent memory leaks
					requestContextData.Delete(req)
				}
			}

			// Prepare response body (limit size for logging)
			responseBody := string(lrw.body)
			if len(responseBody) > 1000 {
				responseBody = responseBody[:1000] + "..."
			}

			// Log data structure
			var logData = map[string]any{
				"level":   strings.ToUpper(cfg.Level),
				"time":    time.Now().Format(time.RFC3339),
				"ip":      ip,
				"port":    port,
				"path":    req.URL.Path,
				"status":  lrw.status,
				"latency": elapsed.String(),
				"host":    req.Host,
				"method":  req.Method,

				// request
				"request_method":     req.Method,
				"request_headers":    sanitizeHeaders(req.Header),
				"request_body":       bodyVal,
				"request_size":       bodySize,
				"request_user_agent": req.UserAgent(),
				"request_referer":    req.Referer(),
				"request_query":      req.URL.RawQuery,
				"request_path":       req.URL.Path,

				// response
				"response_status":  lrw.status,
				"response_size":    lrw.size,
				"response_headers": sanitizeHeaders(lrw.headers),
				"response_body":    responseBody,
			}

			// Add multipart file information for JSON format only
			if cfg.Format == "json" && isMultipartForm && multipartInfo != nil {
				for key, value := range multipartInfo {
					logData[key] = value
				}
			}

			// Add dynamic context data to logData
			for key, value := range dynamicContextData {
				logData[key] = value
			}

			// Create complete colorLogData based on logData (same fields for all formats)
			colorLogData := make(map[string]string)

			// Convert all logData fields to colorLogData with appropriate colors/formatting
			for key, value := range logData {
				valueStr := fmt.Sprintf("%v", value)

				// Apply colors to specific fields for better visualization
				switch key {
				case "time":
					colorLogData[key] = ColorTime + valueStr + ColorReset
				case "level":
					colorLogData[key] = ColorLevel + valueStr + ColorReset
				case "method":
					colorLogData[key] = ColorMethod + valueStr + ColorReset
				case "path":
					colorLogData[key] = ColorPath + valueStr + ColorReset
				case "status":
					colorLogData[key] = ColorStatus + valueStr + ColorReset
				case "latency":
					colorLogData[key] = ColorLatency + valueStr + ColorReset
				default:
					//fmt.Println("key:", key, "val:", valueStr)
					// All other fields (including dynamic context data) without colors
					colorLogData[key] = valueStr
				}
			}

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

				// If no pattern is defined, create a comprehensive default pattern
				if pattern == "" {
					// Create a pattern that includes all available fields
					var fields []string
					for k, v := range colorLogData {
						fields = append(fields, fmt.Sprintf("%s=%v", k, v))
					}
					fmt.Printf("%s\n", strings.Join(fields, " "))
				} else {
					// Use custom pattern and replace placeholders
					replacedFields := make(map[string]bool)

					for k, v := range colorLogData {
						placeholder := fmt.Sprintf("${%s}", k)
						if strings.Contains(pattern, placeholder) {
							pattern = strings.ReplaceAll(pattern, placeholder, fmt.Sprintf("%v", v))
							replacedFields[k] = true
						}
					}

					// Add dynamic fields that weren't in the original pattern
					var extraFields []string
					for k, v := range colorLogData {
						if !replacedFields[k] && !isStandardField(k) {
							extraFields = append(extraFields, fmt.Sprintf("%s=%v", k, v))
						}
					}

					// Print pattern + extra dynamic fields
					output := pattern
					if len(extraFields) > 0 {
						// Remove trailing newline if exists to add extra fields on same line
						output = strings.TrimSuffix(output, "\n")
						output += " | " + strings.Join(extraFields, " ") + "\n"
					}
					fmt.Printf("%s", output)
				}
			}
		})
	}
}

func sanitizeHeaders(headers http.Header) map[string][]string {
	sanitized := make(map[string][]string)
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"set-cookie":    true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			sanitized[key] = []string{"[********]"}
		} else {
			sanitized[key] = values
		}
	}

	return sanitized
}

func getPort(req *http.Request) string {
	if host := req.Host; host != "" {
		_, port, err := net.SplitHostPort(host)
		if err == nil && port != "" {
			return port
		}
	}

	if req.TLS != nil {
		return "443"
	}
	return "80"
}

// isStandardField checks if a field is part of the standard log fields
func isStandardField(fieldName string) bool {
	standardFields := map[string]bool{
		"time": true, "level": true, "method": true, "path": true,
		"status": true, "latency": true, "ip": true, "port": true,
		"host": true, "headers": true, "body": true, "size": true,
		"user_agent": true, "referer": true, "query": true,

		// Request fields
		"request_method": true, "request_headers": true, "request_body": true,
		"request_size": true, "request_user_agent": true, "request_referer": true,
		"request_query": true, "request_path": true,

		// Response fields
		"response_status": true, "response_size": true, "response_headers": true,
		"response_body": true,

		// Multipart fields
		"file_name": true, "file_size": true, "file_type": true,
	}
	return standardFields[fieldName]
}

// extractMultipartInfo extracts file information from multipart form data
func extractMultipartInfo(req *http.Request, bodyBytes []byte) map[string]any {
	info := make(map[string]any)

	// Parse Content-Type to get boundary
	contentType := req.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return info
	}

	boundary, ok := params["boundary"]
	if !ok {
		return info
	}

	// Create multipart reader
	reader := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)

	var fileNames []string
	var totalFileSize int64
	var fileTypes []string

	// Parse multipart form
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		filename := part.FileName()
		if filename != "" {
			// Remove path and get just the filename
			if idx := strings.LastIndex(filename, "/"); idx != -1 {
				filename = filename[idx+1:]
			}
			if idx := strings.LastIndex(filename, "\\"); idx != -1 {
				filename = filename[idx+1:]
			}
			
			// Handle "blob" case - use the form field name instead
			if filename == "blob" || filename == "" {
				if formName := part.FormName(); formName != "" {
					filename = formName
				}
			}
			
			fileNames = append(fileNames, filename)

			content, err := io.ReadAll(part)
			if err == nil {
				partSize := int64(len(content))
				totalFileSize += partSize

				// Detect content type
				contentType := http.DetectContentType(content)
				if contentType != "" {
					fileTypes = append(fileTypes, contentType)
				}
			}
		}
		part.Close()
	}

	if len(fileNames) > 0 {
		if len(fileNames) == 1 {
			info["file_name"] = fileNames[0]
		} else {
			info["file_name"] = strings.Join(fileNames, ", ")
		}

		info["file_size"] = totalFileSize

		if len(fileTypes) == 1 {
			info["file_type"] = fileTypes[0]
		} else if len(fileTypes) > 1 {
			info["file_type"] = strings.Join(fileTypes, ", ")
		} else {
			info["file_type"] = "unknown"
		}
	}

	return info
}
