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
	Print(message, YELLOW, os.Stdout)
}

// Error prints an error message to Stderr in red
func Error(message string) {
	Print(message, RED, os.Stderr)
}

// Output prints a message on Stdout in 'normal' color
func Output(message string) {
	Print(message, RESET, os.Stdout)
}

// Print prints a message to a given stream in a defined color
func Print(message string, color string, stream *os.File) {
	fmt.Fprintf(stream, "%s%s%s\n", color, message, RESET)
}
