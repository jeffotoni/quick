# 🎨 gcolor – Fluent Terminal Styling ![Quick Logo](./quick.png)

> A powerful and expressive ANSI color library for building stylish CLI tools in Go – with a fluent API, test coverage, and real-world usage patterns.

---

## 🚀 Features

	✅ Fluent API chaining  
	✅ Foreground & background colors  
	✅ Bold, underline & style composition  
	✅ Print, Println, Sprint, Sprintf  
	✅ 100% test coverage  
	✅ Plug & play with `log.Printf`, `log.SetPrefix`, and more  
	✅ Production-ready and extensible  

---

## 📦 Installation

```bash
$ go get github.com/jeffotoni/quick/pkg/gcolor
```

---

## ✨ Usage

```go
import (
    "fmt"
    "log"
    "time"

    "github.com/jeffotoni/quick/pkg/gcolor"
)

func main() {
    // ✅ Simple foreground color
    gcolor.New().Fg("green").Println("Success message")

    // 🎯 Foreground + Background
    gcolor.New().Fg("white").Bg("red").Println("Error with red background")

    // 💪 Bold
    gcolor.New().Fg("yellow").Bold().Println("Warning in bold")

    // 🔗 Underline
    gcolor.New().Fg("cyan").Underline().Println("Link or reference")

    // 🔥 Full style chain
    gcolor.New().
        Fg("blue").
        Bg("white").
        Bold().
        Underline().
        Println("Styled and readable message")

    // 🧠 Using Sprintf for dynamic formatting
    user := "jeffotoni"
    fmt.Println(gcolor.New().Fg("green").Sprintf("Welcome, %s!", user))

    // 📡 Dynamic log formatting
    traceID := "abc123"
    duration := 215 * time.Millisecond
    log.Printf(
        "[Trace-ID: %s] <- Completed in %s\n",
        gcolor.New().Fg("cyan").Sprint(traceID),
        gcolor.New().Fg("yellow").Sprintf("%v", duration),
    )

    // 🔁 Reusable styles
    warnStyle := gcolor.New().Fg("yellow").Bold()
    warnStyle.Println("Disk space running low...")

    infoStyle := gcolor.New().Fg("blue")
    infoStyle.Println("Server started successfully")

    // 🧾 Custom prefix in log
    log.SetPrefix(gcolor.New().Fg("purple").Sprint("[GCOLOR] "))
    log.Println("Logger initialized")

    // 🧱 Build and reuse style
    style := gcolor.New().Fg("red").Bold().Underline()
    fmt.Println(style.Sprint("Reusable styled message"))
}
```

---

## 🧪 Test Coverage

We implemented unit tests for:

- `Sprint()` and `Sprintf()`
- `Print()` and `Println()` using `os.Pipe()` capture
- Fluent API usage
- Empty/no-style rendering

Run with:

```bash
$ go test -v -cover 
```

---

## 📚 Examples

In addition to tests, we included rich `ExampleStyle_*()` functions following Go's documentation pattern.

Explore them with:

```bash
$ go doc github.com/jeffotoni/quick/pkg/gcolor
```

In pkg.go.dev [quick pkg/gcolor](https://pkg.go.dev/github.com/jeffotoni/quick/pkg/gcolor)

---

## 🛣️ Roadmap

- 🧩 Named themes and reusable style profiles  
- 🧵 Style reset handling and nested styling  
- 📤 Custom writer support (e.g. `io.Writer`)  
- 🧠 Terminal capability detection (true color, basic, etc)  
- 🌈 RGB and 256-color support *(coming soon)*  

---

## 💬 Contribute

If you like this project, give it a ⭐ star and feel free to open issues or PRs!

Made with 💚 by [@jeffotoni](https://github.com/jeffotoni)  
Part of the [Quick Framework](https://github.com/jeffotoni/quick) ecosystem


🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

---
