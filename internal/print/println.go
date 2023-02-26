package print

import (
	"os"
	"strings"
)

func Stdout(str ...string) {
	os.Stdout.Write([]byte(strings.Join(str, "")))
}
