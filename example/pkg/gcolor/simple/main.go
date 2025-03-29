package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/pkg/gcolor"
)

func main() {
	// Simple foreground color
	gcolor.New().Fg("green").Println("Success message")

	// Foreground + Background
	gcolor.New().Fg("white").Bg("red").Println("Error with red background")

	// Bold text
	gcolor.New().Fg("yellow").Bold().Println("Warning in bold")

	// Underline text
	gcolor.New().Fg("cyan").Underline().Println("Link or reference")

	// Full style chain
	gcolor.New().
		Fg("blue").
		Bg("white").
		Bold().
		Underline().
		Println("Styled and readable message")

	// Using Sprintf for formatted message
	user := "jeffotoni"
	fmt.Println(gcolor.New().Fg("green").Sprintf("Welcome, %s!", user))

	// Dynamic log formatting with colorized values
	traceID := "abc123"
	duration := 215 * time.Millisecond
	log.Printf(
		"[Trace-ID: %s] <- Completed in %s\n",
		gcolor.New().Fg("cyan").Sprint(traceID),
		gcolor.New().Fg("yellow").Sprintf("%v", duration),
	)

	// Reusable styles
	warnStyle := gcolor.New().Fg("yellow").Bold()
	warnStyle.Println("Disk space running low...")

	infoStyle := gcolor.New().Fg("blue")
	infoStyle.Println("Server started successfully")

	// Set custom prefix in log
	log.SetPrefix(gcolor.New().Fg("purple").Sprint("[GCOLOR] "))
	log.Println("Logger initialized")

	// Build and store reusable style
	style := gcolor.New().Fg("red").Bold().Underline()
	fmt.Println(style.Sprint("Reusable styled message"))

	// Color specific parts using Sprint()
	log.Printf(
		"[Trace-ID: %s] <- End of request duration:[(%v)]\n",
		gcolor.New().Fg("cyan").Sprint(traceID),
		gcolor.New().Fg("yellow").Sprint(fmt.Sprint(duration)),
	)

	// Colorize full message
	msg := gcolor.New().Fg("green").Sprint(
		fmt.Sprintf("[Trace-ID: %s] <- End of request duration:[(%v)]", traceID, duration),
	)
	log.Println(msg)

	// Custom reusable styles
	traceStyle := gcolor.New().Fg("cyan").Bold()
	timeStyle := gcolor.New().Fg("yellow")
	log.Printf(
		"[Trace-ID: %s] <- End of request duration:[(%v)]\n",
		traceStyle.Sprint(traceID),
		timeStyle.Sprint(fmt.Sprint(duration)),
	)

	// Colorized log prefix
	log.SetPrefix(gcolor.New().Fg("purple").Sprint("[APP] "))
	log.Println("This is an application log entry")

	// Full message color with Sprintf
	log.Print(
		gcolor.New().Fg("yellow").Bold().Sprintf("[Trace-ID: %s] Done in %v", traceID, duration),
	)
}
