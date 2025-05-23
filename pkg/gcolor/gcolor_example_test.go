package gcolor

import (
	"fmt"
	"log"
	"time"
)

// This function is named ExampleStyle_sprintf()
// it with the Examples type.
func ExampleStyle_sprintf() {
	traceID := "abc123"
	duration := 150 * time.Millisecond

	// Example 1: Color only the dynamic values inside log.Printf
	log.Printf(
		"[Trace-ID: %s] <- End of request duration:[(%v)]\n",
		New().Fg("cyan").Sprint(traceID),
		New().Fg("yellow").Sprintf("%v", duration),
	)

	// Example 2: Format and color the whole message with Sprintf
	msg := New().Fg("green").Sprintf("[Trace-ID: %s] <- End of request duration:[(%d)]", traceID, duration)
	log.Println(msg)

	// Example 3: Using predefined styles
	traceStyle := New().Fg("cyan").Bold()
	timeStyle := New().Fg("yellow")
	log.Printf(
		"[Trace-ID: %s] <- End of request duration:[(%v)]\n",
		traceStyle.Sprint(traceID),
		timeStyle.Sprintf("%v", duration),
	)

	// Example 4: Set a custom prefix for the logger
	log.SetPrefix(New().Fg("purple").Sprint("[APP] "))
	log.Println("Running with colored prefix")
}

// This function is named ExampleStyle_sprint()
// it with the Examples type.
func ExampleStyle_sprint() {
	New().
		Fg("red").
		Bg("black").
		Bold().
		Underline().
		Println("Hello, colorful Quick!")
}

// This function is named ExampleStyle_basic()
// it with the Examples type.
func ExampleStyle_basic() {
	New().Fg("green").Println("Success message")
	New().Fg("white").Bg("red").Println("Error with red background")
	New().Fg("yellow").Bold().Println("Warning in bold")
	New().Fg("cyan").Underline().Println("Link or reference")
	New().Fg("blue").Bg("white").Bold().Underline().Println("Styled and readable message")
}

// This function is named ExampleStyle_formatted()
// it with the Examples type.
func ExampleStyle_formatted() {
	user := "jeffotoni"
	fmt.Println(New().Fg("green").Sprintf("Welcome, %s!", user))
}

// This function is named ExampleStyle_logging()
// it with the Examples type.
func ExampleStyle_logging() {
	traceID := "abc123"
	duration := 215 * time.Millisecond

	log.Printf(
		"[Trace-ID: %s] <- Completed in %s\n",
		New().Fg("cyan").Sprint(traceID),
		New().Fg("yellow").Sprintf("%v", duration),
	)
}

// This function is named ExampleStyle_reusable()
// it with the Examples type.
func ExampleStyle_reusable() {
	warnStyle := New().Fg("yellow").Bold()
	warnStyle.Println("Disk space running low...")

	infoStyle := New().Fg("blue")
	infoStyle.Println("Server started successfully")
}

// This function is named ExampleStyle_loggerPrefix()
// it with the Examples type.
func ExampleStyle_loggerPrefix() {
	log.SetPrefix(New().Fg("purple").Sprint("[ "))
	log.Println("Logger initialized")
}

// This function is named ExampleStyle_buildAndStore()
// it with the Examples type.
func ExampleStyle_buildAndStore() {
	style := New().Fg("red").Bold().Underline()
	fmt.Println(style.Sprint("Reusable styled message"))
}
