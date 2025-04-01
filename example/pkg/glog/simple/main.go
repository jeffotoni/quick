package main

import (
	"errors"
	"time"

	"github.com/jeffotoni/quick/pkg/glog"
)

func main() {
	logger := glog.New(glog.Config{
		Format: "json",
		Level:  glog.DEBUG,
	})

	logger.Debug().
		Int("TraceID", 123475).
		Str("func", "BodyParser").
		Str("status", "success").
		Msg("api-fluent-example").
		Send()

	logger.Info().
		Int("TraceID", 123475).
		Bool("error", false).
		Msg("api-fluent-example").
		Send()

	errTest := errors.New("something went wrong")
	ts := time.Now()
	dur := 1500 * time.Millisecond

	logger.Warn().
		Str("user", "jeff").
		Int("retries", 3).
		Bool("authenticated", true).
		Float64("load", 87.4).
		Duration("elapsed", dur).
		AddTime("timestamp", ts).
		Err("error", errTest).
		Any("data", map[string]int{"a": 1}).
		// Func("trace_id", func() any {
		// 	return "abc123"
		// }). // soon
		Msg("Fluent log test").
		Send()

}
