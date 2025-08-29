package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

var NAME_TRACE_ID = "X-TRACE-ID"

func main() {

	q := quick.New()

	// Apply logger with JSON format
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "DEBUG",
		// Pattern: "[${X-TRACE-ID} ${service} [${time}] ${level} ${method} ${path} ${status} - ${latency}\n", // text and slog
	}))

	q.Post("/upload", func(c *quick.Ctx) error {
		traceID := c.GetTraceID(NAME_TRACE_ID)

		// set header and context in request
		c.SetTraceID(NAME_TRACE_ID, traceID)

		// logger
		c.Logger().
			Str("service-1", "user-service-1").
			Str("func-error-1", "error func2")

		c.Logger().
			Str("service-2", "user-service-2").
			Str("func-error-2", "error func2").
			Int("service-int", 3939).
			Bool("service-bool", true)

		return c.Status(200).JSON(quick.M{
			"msg":     "JSON logging example",
			"traceID": traceID,
			"context": fmt.Sprintf("%v", c.GetAllContextData()),

			"name":  c.FormValue("name"),
			"email": c.FormValue("email"),
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}

// $ curl -i -XPOST localhost:8080/v1/logger/json -d '{"name": "@jeffotoni"}'
// output
// {
//     "X-TRACE-ID": "BfWcvpWQgJxSoFPU",
//     "func-error-1": "error func2",
//     "func-error-2": "error func2",
//     "host": "localhost:8080",
//     "ip": "::1",
//     "latency": "93.584µs",
//     "level": "DEBUG",
//     "method": "POST",
//     "path": "/v1/logger/json",
//     "port": "8080",
//     "request_body": "{\"name\":\"jefferson\"}",
//     "request_headers": {
//         "Accept": [
//             "*/*"
//         ],
//         "Content-Length": [
//             "20"
//         ],
//         "Content-Type": [
//             "application/json"
//         ],
//         "User-Agent": [
//             "curl/8.7.1"
//         ]
//     },
//     "request_method": "POST",
//     "request_path": "/v1/logger/json",
//     "request_query": "",
//     "request_referer": "",
//     "request_size": 20,
//     "request_user_agent": "curl/8.7.1",
//     "response_body": "{\"context\":\"map[X-TRACE-ID:BfWcvpWQgJxSoFPU func-error-1:error func2 func-error-2:error func2 service-1:user-service-1 service-2:user-service-2 service-bool:true service-int:3939]\",\"msg\":\"JSON logging example\",\"traceID\":\"BfWcvpWQgJxSoFPU\"}",
//     "response_headers": {
//         "Content-Type": [
//             "application/json"
//         ],
//         "X-Trace-Id": [
//             "BfWcvpWQgJxSoFPU"
//         ]
//     },
//     "response_size": 239,
//     "response_status": 200,
//     "service-1": "user-service-1",
//     "service-2": "user-service-2",
//     "service-bool": true,
//     "service-int": 3939,
//     "status": 200,
//     "time": "2025-08-29T10:56:36-03:00"
// }
