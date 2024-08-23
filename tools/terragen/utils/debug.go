package utils

import (
	"fmt"
	"strings"
)

func Debug(indent int, format string, a ...any) {
	if Flags.Verbose {
		prefix := strings.Repeat("    ", indent)
		fmt.Printf(prefix+format+"\n", a...)
	}
}
