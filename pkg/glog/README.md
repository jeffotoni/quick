# 🌟 glog — Fast, Flexible, Beautiful Logging for Go ![Quick Logo](./quick.png)

`glog` is a lightweight and powerful logging library for Go designed for simplicity, flexibility, and readability — built for humans and structured for machines.

Part of the [Quick Framework](https://github.com/jeffotoni/quick) ecosystem, `glog` supports colorized terminal output, `slog`-style key=value logs, JSON format, dynamic fields, custom patterns, and intelligent defaults.

---

## ✨ Features

- 🔥 Lightweight and zero-dependency
- 🎨 Supports `text`, `slog`, and `json` formats
- 🧩 Custom `Pattern` with placeholders (`${time}`, `${level}`, `${msg}`, etc)
- 🧠 Dynamic separator detection (` | `, `--`, `:`… based on your pattern)
- 📋 Fluent log builder API with dynamic fields: `.Str()`, `.Int()`, `.Bool()`, `.Any()`, `.Msg()`
- 🎯 Built-in caller tracing: add `file:line` with `.Caller()`
- 🧵 Global `CustomFields` + per-log `Fields` (contextual)
- 🎯 Precise caller tracing with `${file}` (file:line) support
- 🎛️ Built-in log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`
- 🧪 Fully tested with 100% coverage of critical paths
- 🌈 Terminal colors for `level` field (in `text` and `slog`)
- 🧠 Built-in fluent context support: create and extract TraceID, X-User-ID, etc
- ✅ Simple API: `Info()`, `Debugf()`, `Error()`, etc.

---

## 📦 Installation

```bash
$ go get github.com/jeffotoni/quick/pkg/glog
```
---

## 🧠 Context Support (TraceID, etc.)

Sometimes you want to propagate values like `TraceID`, `X-Request-ID`, or `X-User-ID` across your services or middlewares. `glog` provides built-in helpers to work with `context.Context` safely and fluently.

```go
ctx, cancel := glog.CreateCtx().
	Name("X-Trace-ID").
	Key("abc-123").
	Timeout(10 * time.Second).
	Build()
defer cancel()

trace := glog.GetCtx(ctx,"X-Trace-ID") // returns "abc-123"
user := glog.GetCtx(ctx, "X-User-ID") // returns "" if not set
```

You can customize:
	• Name() → sets the context key (default: "TraceID")
	• Key() → sets the value to store in the context
	• Timeout() → context lifetime (default: 30s)


## 💡 If you don’t pass anything, it uses defaults:
```go
ctx, cancel := glog.CreateCtx().Key("abc-123").Build()
```

✅ Safe fallback behavior:
	• Returns "" if context is nil
	• Uses "TraceID" key if not specified
	• Timeout defaults to 30s if not provided
	• Internally avoids key collisions with a private key type

```go
glog.GetCtx(ctx, "custom")  // looks for key "custom"
```

## Example

```bash

   ██████╗ ██╗   ██╗██╗ ██████╗██╗  ██╗
  ██╔═══██╗██║   ██║██║██╔═══   ██║ ██╔╝
  ██║   ██║██║   ██║██║██║      █████╔╝
  ██║▄▄ ██║██║   ██║██║██║      ██╔═██╗
  ╚██████╔╝╚██████╔╝██║╚██████╔ ██║  ██╗
   ╚══▀▀═╝  ╚═════╝ ╚═╝ ╚═════╝ ╚═╝  ╚═╝

 Quick v0.0.1 🚀 Fast & Minimal Web Framework
─────────────────── ───────────────────────────────
 🌎 Host : http://0.0.0.0
 📌 Port : 8080
 🔀 Routes: 1
─────────────────── ───────────────────────────────

TraceID maF1ZfqvvId44Qka | func BodyParser | level DEBUG | msg api-fluent-example-post | status success | time 2025-03-30T16:07:15-03:00
TraceID maF1ZfqvvId44Qka | func SendSQS | level DEBUG | msg SendQueue | status success | time 2025-03-30T16:07:15-03:00
TraceID maF1ZfqvvId44Qka | func Marshal | level DEBUG | msg method SaveSomeWhere | status success | time 2025-03-30T16:07:15-03:00
TraceID maF1ZfqvvId44Qka | code 200 | func SaveSomeWhere | level DEBUG | msg api-fluent-example-post | time 2025-03-30T16:07:15-03:00

-- or

TraceID=yKlWprDAfPlCCO7A func=BodyParser level=DEBUG msg=api-fluent-example-post status=success time=2025-03-30T17:06:34-03:00
TraceID=yKlWprDAfPlCCO7A func=SendSQS level=DEBUG msg=SendQueue status=success time=2025-03-30T17:06:34-03:00
TraceID=yKlWprDAfPlCCO7A func=Marshal level=DEBUG msg=method SaveSomeWhere status=success time=2025-03-30T17:06:34-03:00
TraceID=yKlWprDAfPlCCO7A code=200 func=SaveSomeWhere level=DEBUG msg=api-fluent-example-post time=2025-03-30T17:06:34-03:00

```

## 🚀 Usage with fluent

```go
package main

import (
	"github.com/jeffotoni/quick/pkg/glog"
)

func main() {
	glog.Set(glog.Config{
		Format:        "text",
		Level:         glog.DEBUG,
		IncludeCaller: true,
		Pattern:       "[${time}] ${level} ${msg} ",
		TimeFormat:    "2006-01-02 15:04:05",
		CustomFields: map[string]string{
			"service": "example-api",
		},
	})

	glog.Debug("api-fluent-example").
		Int("TraceID", 123475).
		Str("func", "BodyParser").
		Str("status", "success").
		Send()

	glog.Info("api-fluent-example").
		Int("TraceID", 123475).
		Bool("error", false).
		Send()

	errTest := errors.New("something went wrong")
	ts := time.Now()
	dur := 1500 * time.Millisecond

	glog.Warn("Fluent log test").
		Str("user", "jeff").
		Int("retries", 3).
		Bool("authenticated", true).
		Float64("load", 87.4).
		Duration("elapsed", dur).
		Time("timestamp", ts).
		Err("error", errTest).
		Any("data", map[string]int{"a": 1}).
		Func("trace_id", func() any {
			return "abc123"
		}).
		Send()
}
```

🖨️ Sample Output (text):

🟢 INFO  
🔵 DEBUG  
🟡 WARN  
🔴 ERROR

```
2025-03-30 16:31:00 DEBUG api-fluent-example TraceID 123475 env production file proc.go:283 func BodyParser service example-api status success
2025-03-30 16:31:00 INFO api-fluent-example TraceID 123475 env production error false file proc.go:283 service example-api
2025-03-30 16:31:00 WARN Fluent log test authenticated true data map[a:1] elapsed 1.5s env production error something went wrong file proc.go:283 load 87.4 retries 3 service example-api timestamp 2025-03-30T16:31:00-03:00 trace_id abc123 user jeff
```

---

## 🚀 Usage normal but InfoT,DebugT,WarnT and ErrorT

```go
package main

import (
	"github.com/jeffotoni/quick/pkg/glog"
	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	glog.Set(glog.Config{
		Format:        "text",
		Level:         glog.DEBUG,
		IncludeCaller: true,
		Pattern:       "[${time}] ${level} ${msg} ",
		TimeFormat:    "2006-01-02 15:04:05",
		CustomFields: map[string]string{
			"service": "example-api",
		},
	})

	glog.InfoT(rand.TraceID())

	glog.InfoT("Started request", glog.Fields{
		"TRACE": rand.TraceID(),
	})

	glog.DebugT("This is a debug message", glog.Fields{"user": "jeff"})
	glog.Infof("User %s logged in successfully", "arthur")
	glog.WarnT("Low disk space warning")
	glog.ErrorT("Database connection failed", glog.Fields{"retry": true})

	glog.InfoT("Processing order", glog.Fields{
		"order_id": "ORD1234",
		"customer": "Alice",
		"total":    153.76,
	})
}
```

🖨️ Sample Output (text):

🟢 INFO  
🔵 DEBUG  
🟡 WARN  
🔴 ERROR

```
[2025-03-29 17:10:21] INFO 7KF5hlUUNic0K7Sr main.go:15 service example-api
[2025-03-29 17:10:21] INFO Started request main.go:17 TRACE zMxy1...
[2025-03-29 17:10:21] DEBUG This is a debug message main.go:19 user jeff service example-api
[2025-03-29 17:10:21] INFO User arthur logged in successfully main.go:20 service example-api
[2025-03-29 17:10:21] WARN Low disk space warning main.go:21 service example-api
[2025-03-29 17:10:21] ERROR Database connection failed main.go:22 retry true service example-api
[2025-03-29 17:10:21] INFO Processing order main.go:24 order_id ORD1234 customer Alice total 153.76 service example-api
```

---

## 🧪 Test Coverage

We implemented unit tests for:

- All log levels (`Info`, `Debug`, `Error`, `Warn`) with both fluent and legacy `*T` syntax
- Pattern replacement and extra field rendering
- Separator auto-detection logic (via pattern) and fallback to `" "` if not defined
- Custom field merging: global `CustomFields` + dynamic fluent fields
- Ordered field rendering in fluent logs (preserves insertion order)
- Caller trace injection via `${file}` when `IncludeCaller` is enabled
- `Debugf`, `Errorf`, `Warnf`, `Infof` formatted message handling
- JSON and slog output structure (key/value format with coloring for `slog`)
- Writer redirection for test capture and output validation
- Contextual fallback logic for `Separator` when `Pattern` is empty
- Edge case tests: missing keys, nil context, deadline timeout, custom key names
- Compatibility support for `Fields` maps via generic wrapper in legacy `*T` methods

Run tests with:

```bash
$ go test -v -cover
```

---

## 📚 Examples

In addition to tests, we included rich `Example_*()` functions following Go’s documentation pattern.

Explore them with:

```bash
$ go doc github.com/jeffotoni/quick/pkg/glog
```

In pkg.go.dev [quick pkg/glog](https://pkg.go.dev/github.com/jeffotoni/quick/pkg/glog)


---

## 🛣️ Roadmap

- 🧩 Named templates and reusable config profiles  
- 🧵 Fine-grained style customization (e.g. themes per level)  
- 🧠 More optimized color rendering for `slog` and `text`  
- 🚦 Buffered + async writer (opt-in)
- 🧪 Benchmark utilities

---

## 💬 Contribute

If you like this project, give it a ⭐ star and feel free to open issues or PRs!

Made with 💚 by [@jeffotoni](https://github.com/jeffotoni)  
Part of the [Quick Framework](https://github.com/jeffotoni/quick) ecosystem

---
