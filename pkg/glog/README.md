# 🌟 glog — Fast, Flexible, Beautiful Logging for Go ![Quick Logo](./quick.png)

`glog` is a lightweight and powerful logging library for Go designed for simplicity, flexibility, and readability — built for humans and structured for machines.

Part of the [Quick Framework](https://github.com/jeffotoni/quick) ecosystem, `glog` supports colorized terminal output, `slog`-style key=value logs, JSON format, dynamic fields, custom patterns, and intelligent defaults.

---

## ✨ Features

- 🔥 Lightweight and zero-dependency
- 🎨 Supports `text`, `slog`, and `json` formats
- 🧩 Custom `Pattern` with placeholders (`${time}`, `${level}`, `${msg}`, etc)
- 🧠 Dynamic separator detection (` | `, `--`, `:`… based on your pattern)
- 🧵 Global `CustomFields` + per-log `Fields` (contextual)
- 🎯 Precise caller tracing with `${file}` (file:line) support
- 🎛️ Built-in log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`
- 🧪 Fully tested with 100% coverage of critical paths
- 🌈 Terminal colors for `level` field (in `text` and `slog`)
- ✅ Simple API: `Info()`, `Debugf()`, `Error()`, etc.

---

## 📦 Installation

```bash
$ go get github.com/jeffotoni/quick/pkg/glog
```

---

## 🚀 Usage

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

	glog.Info(rand.TraceID())

	glog.Info("Started request", glog.Fields{
		"TRACE": rand.TraceID(),
	})

	glog.Debug("This is a debug message", glog.Fields{"user": "jeff"})
	glog.Infof("User %s logged in successfully", "arthur")
	glog.Warn("Low disk space warning")
	glog.Error("Database connection failed", glog.Fields{"retry": true})

	glog.Info("Processing order", glog.Fields{
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

- All log levels (`Info`, `Debug`, `Error`, `Warn`)
- Pattern replacement and extra field rendering
- Separator auto-detection logic
- Custom field merging (global + per-call)
- Caller trace injection via `${file}`
- `Debugf`, `Errorf`, `Warnf`, `Infof` formatting
- JSON and slog output structure
- Writer redirection for test capture

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
