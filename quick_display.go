// Package quick provides a lightweight and high-performance web framework for building web applications in Go.
// This file is responsible for displaying the Quick framework banner and startup configuration details.

package quick

import (
	"fmt"
	"net"
)

// ANSI terminal color codes used for banner styling
const (
	Reset  = "\033[0m"
	Blue   = "\033[34m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// QuickVersion represents the current version of the Quick framework.
const QuickVersion = "v0.0.1"

// Display prints a styled startup banner to the console showing essential information about the Quick server.
//
// It includes the following details:
//   - The current version of Quick
//   - Host and port where the server is running
//   - The total number of registered routes
//
// The banner is only printed if the `NoBanner` option is set to false in the configuration.
//
// Parameters:
//   - scheme: The protocol used by the server (e.g., "http" or "https").
//   - addr: The address the server is bound to (in the format "host:port").
func (q *Quick) Display(scheme, addr string) {
	if !q.config.NoBanner {

		// Counts the number of registered routes
		routeCount := len(q.GetRoute())

		// Extract port from addr
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			fmt.Println("Error separating host and port:", err)
			return
		}

		if len(host) == 0 {
			host = "127.0.0.1"
		}

		// Display the styled banner
		fmt.Println()
		fmt.Printf("%s%s   ██████╗ ██╗   ██╗██╗ ██████╗██╗  ██╗%s\n", Bold, Blue, Reset)
		fmt.Printf("%s  ██╔═══██╗██║   ██║██║██╔═══   ██║ ██╔╝%s\n", Blue, Reset)
		fmt.Printf("%s  ██║   ██║██║   ██║██║██║      █████╔╝ %s\n", Blue, Reset)
		fmt.Printf("%s  ██║▄▄ ██║██║   ██║██║██║      ██╔═██╗ %s\n", Blue, Reset)
		fmt.Printf("%s  ╚██████╔╝╚██████╔╝██║╚██████╔ ██║  ██╗%s\n", Blue, Reset)
		fmt.Printf("%s   ╚══▀▀═╝  ╚═════╝ ╚═╝ ╚═════╝ ╚═╝  ╚═╝%s\n", Blue, Reset)
		fmt.Println()
		fmt.Printf("%s%s Quick %s %s🚀 Fast & Minimal Web Framework%s\n", Bold, Cyan, QuickVersion, Yellow, Reset)
		fmt.Println("─────────────────── ───────────────────────────────")
		fmt.Printf("%s 🌎 Host : %s%s://%s%s\n", Yellow, Green, scheme, host, Reset)
		fmt.Printf("%s 📌 Port : %s%s%s\n", Yellow, Green, port, Reset)
		fmt.Printf("%s 🔀 Routes: %s%d%s\n", Yellow, Green, routeCount, Reset)
		fmt.Println("─────────────────── ───────────────────────────────")
		fmt.Println()
	}
}
