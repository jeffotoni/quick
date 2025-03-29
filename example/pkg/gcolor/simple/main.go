package main

// More Examples

// Simple foreground color
// New().Fg("green").Println("Success message")

//Foreground + Background
// New().Fg("white").Bg("red").Println("Error with red background")

// Bold text
// New().Fg("yellow").Bold().Println("Warning in bold")

//Underline
// New().Fg("cyan").Underline().Println("Link or reference")

// Full style chain
// New().
//     Fg("blue").
//     Bg("white").
//     Bold().
//     Underline().
//     Println("Styled and readable message")

// Using Sprintf for formatted string
// user := "jeffotoni"
// New().Fg("green").Sprintf("Welcome, %s!", user)

// Dynamic log formatting
// traceID := "abc123"
// duration := 215 * time.Millisecond
// log.Printf(
//     "[Trace-ID: %s] <- Completed in %s\n",
//     New().Fg("cyan").Sprint(traceID),
//     New().Fg("yellow").Sprintf("%v", duration),

// Reusable styles
// warnStyle := New().Fg("yellow").Bold()
// warnStyle.Println("Disk space running low...")
// infoStyle := New().Fg("blue")
// infoStyle.Println("Server started successfully")

// Set custom prefix in log
// log.SetPrefix(New().Fg("purple").Sprint("[ "))
// log.Println("Logger initialized")

// Build and store style for later use
//  style := New().Fg("red").Bold().Underline()
// fmt.Println(style.Sprint("Reusable styled message"))

//////// more example
// Option 1: Using Sprint() to color only specific parts

// log.Printf(
//   "[Trace-ID: %s] <- End of request duration:[(%v)]\n",
//   New().Fg("cyan").Sprint(traceID),
//   New().Fg("yellow").Sprint(duration),
// )

// Option 2: Assemble the complete line with Sprint()
// msg := New().Fg("green").Sprint(
//   fmt.Sprintf("[Trace-ID: %s] <- End of request duration:[(%v)]", traceID, duration),
// )
// log.Println(msg)

// Option 3: Create custom themes
// traceStyle := New().Fg("cyan").Bold()
// timeStyle := New().Fg("yellow")

// log.Printf(
//   "[Trace-ID: %s] <- End of request duration:[(%v)]\n",
//   traceStyle.Sprint(traceID),
//   timeStyle.Sprint(duration),
// )

// Option 4: If you want to colorize even the log prefix
// log.SetPrefix(New().Fg("purple").Sprint("[APP] "))

// Option 5: Use Sprintf
// log.Print(
//   New().Fg("yellow").Bold().Sprintf("[Trace-ID: %s] Done in %v", traceID, duration),
// )
