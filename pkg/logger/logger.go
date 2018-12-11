package logger

import (
	"fmt"
	"os"
)

const (
	YELLOW = "\x1b[33;1m"
	GREEN  = "\x1b[32;1m"
	RED    = "\x1b[31;1m"
	RESET  = "\x1b[33;0m"
)

func Warning(message string) {
	print(message, YELLOW, os.Stdout)
}

func Error(message string) {
	print(message, RED, os.Stderr)
}

func Output(message string) {
	print(message, RESET, os.Stdout)
}

func print(message string, color string, stream *os.File) {
	fmt.Fprintf(stream, "%s%s%s\n", color, message, RESET)
}
