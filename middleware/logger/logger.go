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

// ANSI color codes
const (
	ColorReset   = "\033[0m"
	ColorTime    = "\033[36m" // Cyan
	ColorLevel   = "\033[32m" // Green
	ColorMethod  = "\033[34m" // Blue
	ColorPath    = "\033[35m" // Magenta
	ColorStatus  = "\033[33m" // Yellow
	ColorLatency = "\033[31m" // Red
)

// Config allows customization of the logger middleware
type Config struct {
	Format       string            // "text", "slog", "json"
	Pattern      string            // Log format pattern
	Level        string            // "DEBUG", "INFO", "WARN", "ERROR"
	CustomFields map[string]string // Custom fields for logging
}

// Default configuration
var ConfigDefault = Config{
	Format:  "text",
	Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency}\n",
}

type ColorHandler struct {
	slog.Handler
	w io.Writer
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = "\033[37m" // Gray
	case slog.LevelInfo:
		levelColor = "\033[32m" // Green
	case slog.LevelWarn:
		levelColor = "\033[33m" // Yellow
	case slog.LevelError:
		levelColor = "\033[31m" // Red
	default:
		levelColor = "\033[0m"
	}
	reset := "\033[0m"

	timeStr := r.Time.Format(time.RFC3339)
	line := fmt.Sprintf("time=%s level=%s%s%s msg=%s\n",
		timeStr, levelColor, r.Level.String(), reset, r.Message)

	_, err := h.w.Write([]byte(line))
	return err
}

// loggerRespWriter captures response status and size
type loggerRespWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *loggerRespWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *loggerRespWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// New creates a logging middleware with a configurable output format
func New(config ...Config) func(http.Handler) http.Handler {
	cfg := ConfigDefault
	if len(config) > 0 {
		cfg = config[0]
	}

	var logger *slog.Logger
	var handlerOpts = &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	switch cfg.Format {
	case "slog":
		logger = slog.New(&ColorHandler{
			Handler: slog.NewTextHandler(os.Stdout, handlerOpts),
			w:       os.Stdout,
		})

	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOpts)) // Fallback to text
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == quick.MethodOptions {
				return
			}
			start := time.Now()

			ip, port, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
				port = "?"
			}

			var bodySize int64
			if req.Body != nil {
				body, err := io.ReadAll(req.Body)
				if err == nil {
					bodySize = int64(len(body))
					req.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}

			lrw := &loggerRespWriter{ResponseWriter: w}
			next.ServeHTTP(lrw, req)

			elapsed := time.Since(start)

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

			// Apply colors only in text mode
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

			for k, v := range cfg.CustomFields {
				logData[k] = v
				colorLogData[k] = v
			}

			switch cfg.Format {
			case "json":
				jsonData, _ := json.Marshal(logData)
				//logger.Info(string(jsonData))
				fmt.Println(string(jsonData))

			case "slog":
				// Apply pattern replacements
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
				fmt.Println(pattern)
			}
		})
	}
}
