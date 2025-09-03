package main

import (
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/glog"
	"github.com/jeffotoni/quick/rand"
)

var KeyName = "X-Trace-ID"
var b = []byte(`{"msg":"ok"}`)

func main() {

	logger := glog.New(glog.Config{
		Format: "json",
		Level:  glog.ERROR,
	})

	q := quick.New()

	q.Post("/v1/logger/json", func(c *quick.Ctx) error {
		// creating a trace
		traceID := c.Get(KeyName)
		if traceID == "" {
			traceID = rand.TraceID()
		}

		userID := "user3039"
		spanID := "span39393"

		_, cancel := glog.CreateCtx().
			Set("X-Trace-ID", traceID).
			Set("X-User-ID", userID).
			Set("X-Span-ID", spanID).
			Timeout(10 * time.Second).
			Build()
		defer cancel()

		c.Set("X-Trace-ID", traceID)
		c.Set("X-User-ID", userID)
		c.Set("X-Span-ID", spanID)

		// c.Set("Content-type", "application/json")
		// var d any
		// err := c.BodyParser(&d)
		// if err != nil {
		// 	logger.Error().
		// 		Time().
		// 		Level().
		// 		Str(KeyName, traceID).
		// 		Str("error", err.Error()).
		// 		Send()
		// 	return c.Status(500).JSON(quick.M{"msg": err.Error()})
		// }

		logger.Info().
			Time().
			Level().
			Str(KeyName, traceID).
			Str("func", "BodyParser").
			Str("status", "success").
			// Caller().
			Send()

		// all := glog.GetCtxAll(ctx)
		// fmt.Printf("X-Trace-ID:%s X-User-ID:%s X-Span-ID:%s\n", all["X-Trace-ID"], all["X-User-ID"], all["X-Span-ID"])

		return c.Status(quick.StatusOK).Send(b)
	})

	q.Listen("0.0.0.0:8080")
}
