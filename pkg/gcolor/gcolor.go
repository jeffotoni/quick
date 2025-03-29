// Package gcolor offers a fluent, flexible, and expressive API
// for styling terminal output using ANSI escape codes.
//
// It allows you to combine foreground and background colors,
// text styles like bold and underline, and render the result
// to the terminal with full control over formatting.
//
// Example:
//
//	gcolor.New().
//	    Fg("red").
//	    Bg("black").
//	    Bold().
//	    Underline().
//	    Println("Hello, colorful Quick!")
package gcolor

import (
	"fmt"
	"strings"
)

// ANSI escape codes for styles and reset
const (
	ansiReset     = "\033[0m"
	ansiBold      = "\033[1m"
	ansiUnderline = "\033[4m"
)

// ANSI foreground color codes
var ansiFgColors = map[string]string{
	"black":  "\033[30m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",
}

// ANSI background color codes
var ansiBgColors = map[string]string{
	"black":  "\033[40m",
	"red":    "\033[41m",
	"green":  "\033[42m",
	"yellow": "\033[43m",
	"blue":   "\033[44m",
	"purple": "\033[45m",
	"cyan":   "\033[46m",
	"white":  "\033[47m",
}

// Style represents a combination of foreground, background and text styles.
type Style struct {
	fg        string
	bg        string
	bold      bool
	underline bool
}

// New returns a new empty Style instance.
func New() *Style {
	return &Style{}
}

// Fg sets the foreground color.
func (s *Style) Fg(color string) *Style {
	s.fg = ansiFgColors[color]
	return s
}

// Bg sets the background color.
func (s *Style) Bg(color string) *Style {
	s.bg = ansiBgColors[color]
	return s
}

// Bold enables bold text style.
func (s *Style) Bold() *Style {
	s.bold = true
	return s
}

// Underline enables underline text style.
func (s *Style) Underline() *Style {
	s.underline = true
	return s
}

// Sprint returns the styled string.
func (s *Style) Sprint(text string) string {
	var sb strings.Builder

	if s.fg != "" {
		sb.WriteString(s.fg)
	}
	if s.bg != "" {
		sb.WriteString(s.bg)
	}
	if s.bold {
		sb.WriteString(ansiBold)
	}
	if s.underline {
		sb.WriteString(ansiUnderline)
	}

	sb.WriteString(text)
	sb.WriteString(ansiReset)
	return sb.String()
}

// Print prints the styled text.
func (s *Style) Print(text string) {
	fmt.Print(s.Sprint(text))
}

// Println prints the styled text followed by a new line.
func (s *Style) Println(text string) {
	fmt.Println(s.Sprint(text))
}

// Sprintf returns a formatted and styled string.
func (s *Style) Sprintf(format string, args ...interface{}) string {
	formatted := fmt.Sprintf(format, args...)
	return s.Sprint(formatted)
}
