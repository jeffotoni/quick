// Package gcolor provides a fluent and expressive API for styling terminal output
// using ANSI escape codes. It allows developers to apply foreground and background
// colors, text decorations such as bold and underline, and output styled strings
// to the terminal.
//
// gcolor is ideal for CLIs, logging systems, and other tools where color-coded
// or styled output can improve readability.
//
// Features:
//
//   - Foreground and background color support (8 base colors)
//   - Bold and underline styling
//   - Chainable (fluent) syntax
//   - Print, Println, Sprint, and Sprintf helpers
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

// ANSI escape codes for basic text styles and reset.
const (
	ansiReset     = "\033[0m" // Resets all styles (color, bold, underline, etc.)
	ansiBold      = "\033[1m" // Enables bold text
	ansiUnderline = "\033[4m" // Enables underline text
)

// ansiFgColors maps color names to ANSI escape sequences for foreground (text) colors.
// These can be used to change the text color in terminal output.
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

// ansiBgColors maps color names to ANSI escape sequences for background colors.
// These are used to change the background color behind terminal text.
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

// Style represents a combination of ANSI styles for terminal text,
// including foreground color, background color, bold, and underline.
// Use New() to construct a Style and apply styling methods fluently.
type Style struct {
	fg        string
	bg        string
	bold      bool
	underline bool
}

// New returns a new Style instance with no styles applied.
// You can chain methods like Fg, Bg, Bold, and Underline to apply formatting.
func New() *Style {
	return &Style{}
}

// Fg sets the foreground (text) color using a named ANSI color.
// Available options: black, red, green, yellow, blue, purple, cyan, white.
func (s *Style) Fg(color string) *Style {
	s.fg = ansiFgColors[color]
	return s
}

// Bg sets the background color using a named ANSI color.
// Available options: black, red, green, yellow, blue, purple, cyan, white.
func (s *Style) Bg(color string) *Style {
	s.bg = ansiBgColors[color]
	return s
}

// Bold enables bold formatting on the styled text.
func (s *Style) Bold() *Style {
	s.bold = true
	return s
}

// Underline enables underline formatting on the styled text.
func (s *Style) Underline() *Style {
	s.underline = true
	return s
}

// Sprint returns the input string with all the applied styles (colors and decorations).
// The result includes ANSI escape sequences and ends with a reset code.
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

// Print writes the styled string to standard output without a newline.
func (s *Style) Print(text string) {
	fmt.Print(s.Sprint(text))
}

// Println writes the styled string to standard output followed by a newline.
func (s *Style) Println(text string) {
	fmt.Println(s.Sprint(text))
}

// Sprintf formats a string using fmt.Sprintf and applies the style to the result.
func (s *Style) Sprintf(format string, args ...interface{}) string {
	formatted := fmt.Sprintf(format, args...)
	return s.Sprint(formatted)
}
