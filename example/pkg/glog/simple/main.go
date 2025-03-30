package main

import (
	"github.com/jeffotoni/quick/pkg/gcolor"
	"github.com/jeffotoni/quick/pkg/glog"
	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	// Configure glog with slog-style output (field=value), custom timestamp and default level
	glog.Set(glog.Config{
		Format:     "text",                // Use slog formatting (key=value, colored)
		TimeFormat: "2006-01-02 15:04:05", // Human-friendly timestamp format
		Level:      glog.DEBUG,            // Minimum log level to display
		CustomFields: map[string]string{ // Fields injected into every log entry
			"service": "example-api",
			"env":     "production",
		},
		IncludeCaller: true,
		// Optional: use a custom pattern if desired
		// Pattern: "[${time}] ${level} ${msg}",
		// Pattern: "${time} ${level} ${msg} user=${user} retry=${retry} order_id=${order_id} total=${total}",
	})

	// Example 1: Trace ID only
	traceID := rand.TraceID()
	glog.Info(traceID) // string message only

	// Example 2: Trace ID as field
	glog.Info("Start request", glog.Fields{
		"TRACE": traceID,
	})

	// Example 3: Debug message with custom field
	glog.Debug("This is a debug message", glog.Fields{
		"user": gcolor.
			New().
			Fg("yellow").
			Sprint("jeffotoni"),
	})

	// Example 4: Formatted info message (uses sprintf-style)
	glog.Infof("User %s logged in successfully", "arthur")

	// Example 5: Warning with simple string
	glog.Warn("Low disk space warning")

	// Example 6: Error with custom retry field
	glog.Error("Database connection failed", glog.Fields{
		"retry": true,
	})

	// Example 7: Simulated order processing log
	glog.Info("Processing order", glog.Fields{
		"order_id": "ORD1234",
		"customer": "Alice",
		"total":    153.76,
	})
}
