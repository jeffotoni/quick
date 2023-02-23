package quick

import (
	"os"
	"strings"
)

func Print(str ...string) {
	os.Stdout.Write([]byte(strings.Join(str, "")))
}
