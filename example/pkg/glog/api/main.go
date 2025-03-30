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
	glog.Set(glog.Config{
		Format: "slog",
		// Separator: " | ", // only in Format: text
		Level: glog.DEBUG,
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
			//glog.Error(traceID, glog.Fields{"error": err})
			glog.Error("api-example-post").
				Str(KeyName, traceID).
				Str("error", err.Error()).
				Send()
			return c.Status(500).JSON(quick.M{"msg": err.Error()})
		}

		glog.Debug("api-fluent-example-post").
			Str(KeyName, traceID).
			Str("func", "BodyParser").
			Str("status", "success").
			Send()

		// call metodh
		b, err := SaveSomeWhere(ctx, d)
		if err != nil {
			glog.ErrorT(traceID, glog.Fields{"error": err})
			return c.Status(500).JSON(quick.M{"msg": err.Error()})
		}

		glog.Debug("api-fluent-example-post").
			Str(KeyName, traceID).
			Str("func", "SaveSomeWhere").
			Int("code", quick.StatusOK).
			Send()

		return c.Status(quick.StatusOK).Send(b)
	})

	q.Listen("0.0.0.0:8080")
}

func SaveSomeWhere(ctx context.Context, data any) (b []byte, err error) {
	traceID := glog.GetCtx(ctx)
	b, err = json.Marshal(data)
	if err != nil {
		glog.ErrorT(traceID, glog.Fields{"error": err})
		return
	}

	err = SendQueue(ctx, b)
	if err != nil {
		glog.ErrorT(traceID, glog.Fields{"error": err})
		return nil, err
	}

	glog.Debug("method SaveSomeWhere").
		Str(KeyName, traceID).
		Str("func", "Marshal").
		Str("status", "success").
		Send()

	return
}

func SendQueue(ctx context.Context, data []byte) (err error) {
	// send SQS

	time.Sleep(time.Millisecond * 100)

	glog.Debug("SendQueue").
		Str(KeyName, glog.GetCtx(ctx)).
		Str("func", "SendSQS").
		Str("status", "success").
		Send()

	return
}
