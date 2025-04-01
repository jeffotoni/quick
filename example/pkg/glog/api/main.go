package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/pkg/glog"
	"github.com/jeffotoni/quick/pkg/rand"
)

const KeyName string = "X-Trace-ID"

// curl -i -XPOST -H "Content-Type:application/json" localhost:8080/v1/user -d '{"name": "jeff", "year": 2025}'
func main() {

	logger := glog.New(glog.Config{
		Format: "json",
		Level:  glog.DEBUG,
	})

	q := quick.New()

	q.Post("/v1/user", func(c *quick.Ctx) error {
		// creating a trace
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = rand.TraceID()
		}

		userID := rand.AlgoDefault(9000, 9000)
		spanID := "span39393"

		ctx, cancel := glog.CreateCtx().
			Set("X-Trace-ID", traceID).
			Set("X-User-ID", userID).
			Set("X-Span-ID", spanID).
			Timeout(10 * time.Second).
			Build()
		defer cancel()

		c.Set("X-Trace-ID", traceID)
		c.Set("X-User-ID", userID)
		c.Set("X-Span-ID", spanID)

		c.Set("Content-type", "application/json")
		var d any
		err := c.BodyParser(&d)
		if err != nil {
			logger.Error().
				Time().
				Level().
				Str(KeyName, traceID).
				Str("error", err.Error()).
				Send()
			return c.Status(500).JSON(quick.M{"msg": err.Error()})
		}

		logger.Debug().
			Time().
			Level().
			Str(KeyName, traceID).
			Str("func", "BodyParser").
			Str("status", "success").
			// Caller().
			Send()

		// call metodh
		b, err := SaveSomeWhere(ctx, logger, d)
		if err != nil {
			logger.Error().
				Time().
				Level().
				Str(KeyName, traceID).
				Str("Error", err.Error()).
				Send()

			return c.Status(500).JSON(quick.M{"msg": err.Error()})
		}

		logger.Debug().
			Time().
			Level().
			Str(KeyName, traceID).
			Str("func", "SaveSomeWhere").
			Int("code", quick.StatusOK).
			Msg("api-post-fluent").
			Send()

		all := glog.GetCtxAll(ctx)
		fmt.Printf("X-Trace-ID:%s X-User-ID:%s X-Span-ID:%s\n", all["X-Trace-ID"], all["X-User-ID"], all["X-Span-ID"])

		return c.Status(quick.StatusOK).Send(b)
	})

	q.Listen("0.0.0.0:8080")
}

func SaveSomeWhere(ctx context.Context, logger *glog.Logger, data any) (b []byte, err error) {
	traceID := glog.GetCtx(ctx, KeyName)
	b, err = json.Marshal(data)
	if err != nil {
		logger.Error().
			Time().
			Level().
			Str(KeyName, traceID).
			Str("Error", err.Error()).
			Send()

		return
	}

	err = SendQueue(ctx, logger, b)
	if err != nil {
		logger.Error().
			Time().
			Level().
			Str(KeyName, traceID).
			Str("Error", err.Error()).
			Send()

		return nil, err
	}

	logger.Debug().
		Time().
		Level().
		Str(KeyName, traceID).
		Str("func", "Marshal").
		Str("status", "success").
		Send()

	return
}

func SendQueue(ctx context.Context, logger *glog.Logger, data []byte) (err error) {
	// send SQS
	time.Sleep(time.Millisecond * 100)

	logger.Debug().
		Time().
		Level().
		Str(KeyName, glog.GetCtx(ctx, KeyName)).
		Str("func", "SendSQS").
		Str("status", "success").
		Send()

	return
}
