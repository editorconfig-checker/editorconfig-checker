package logger

import (
	"fmt"
	"os"
)

// Colors which can be used
const (
	YELLOW = "\x1b[33;1m"
	GREEN  = "\x1b[32;1m"
	RED    = "\x1b[31;1m"
	RESET  = "\x1b[33;0m"
)

// Warning prints a warning message to Stdout in yellow
func Warning(message string) {
	print(message, YELLOW, os.Stdout)
}

// Error prints an error message to Stderr in red
func Error(message string) {
	print(message, RED, os.Stderr)
}

// Output prints a message on Stdout in 'normal' color
func Output(message string) {
	print(message, RESET, os.Stdout)
}

func print(message string, color string, stream *os.File) {
	fmt.Fprintf(stream, "%s%s%s\n", color, message, RESET)
}
