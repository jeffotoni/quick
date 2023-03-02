package log

import (
	"os"
	"strings"
)

func Log(str ...string) {
	str = append(str, string("\n"))
	out := []byte(strings.Join(str, ""))
	_, err := os.Stdout.Write(out)
	if err != nil {
		panic(err)
	}
}
