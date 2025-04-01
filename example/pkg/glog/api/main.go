package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/pkg/glog"
	"github.com/jeffotoni/quick/pkg/rand"
)

const KeyName string = "TraceID"

// curl -i -XPOST -H "Content-Type:application/json" localhost:8080/v1/user -d '{"name": "jeff", "year": 2025}'
func main() {

	logger := glog.Set(glog.Config{
		Format:    "slog", // ou "slog", "json"
		Level:     glog.DEBUG,
		Separator: " | ",
	})

	q := quick.New()

	q.Post("/v1/user", func(c *quick.Ctx) error {
		// creating a trace
		traceID := c.Get(KeyName)
		if traceID == "" {
			traceID = rand.TraceID()
		}
		c.Set(KeyName, traceID)

		ctx, cancel := glog.NewCtx().
			Name(KeyName).
			Key(traceID).
			Timeout(10 * time.Second).
			Build()
		defer cancel()

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
			Caller().
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

		return c.Status(quick.StatusOK).Send(b)
	})

	q.Listen("0.0.0.0:8080")
}

func SaveSomeWhere(ctx context.Context, logger *glog.Logger, data any) (b []byte, err error) {

	traceID := glog.GetCtx(ctx)
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
		Str(KeyName, glog.GetCtx(ctx)).
		Str("func", "SendSQS").
		Str("status", "success").
		Send()

	return
}
