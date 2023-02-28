package print

import (
	"os"
	"strings"
)

func Stdout(str ...string) {
	_, err := os.Stdout.Write([]byte(strings.Join(str, "")))
	if err != nil {
		panic(err)
	}
}
