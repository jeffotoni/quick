package quick

import (
	"fmt"
)

// ANSI color definition
const (
	Reset  = "\033[0m"
	Blue   = "\033[34m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Quick version
const QuickVersion = "v0.0.1"

func (q *Quick) Display(scheme, port string) {
	if !q.config.NoBanner {

		// Counts the number of registered routes
		routeCount := len(q.GetRoute())

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
		fmt.Printf("%s 🌎 Host : %s%s://127.0.0.1:%s%s\n", Yellow, Green, scheme, port, Reset)
		fmt.Printf("%s 📌 Port : %s%s%s\n", Yellow, Green, port, Reset)
		fmt.Printf("%s 🔀 Routes: %s%d%s\n", Yellow, Green, routeCount, Reset)
		fmt.Println("─────────────────── ───────────────────────────────")
		fmt.Println()
	}
}
