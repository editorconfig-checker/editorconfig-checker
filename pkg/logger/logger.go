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

type Logger struct {
	Verbosee bool
	Debugg   bool
}

func (l Logger) Debug(message string) {
	if l.Debugg {
		Print(message, RESET, os.Stdout)
	}
}

func (l Logger) Verbose(message string) {
	if l.Verbosee {
		Print(message, RESET, os.Stdout)
	}
}

// Warning prints a warning message to Stdout in yellow
func (l Logger) Warning(message string) {
	Print(message, YELLOW, os.Stdout)
}

// Output prints a message on Stdout in 'normal' color
func (l Logger) Output(message string) {
	Print(message, RESET, os.Stdout)
}

// Output prints a message on Stdout in 'normal' color
func Output(message string) {
	Print(message, RESET, os.Stdout)
}

// Error prints an error message to Stderr in red
func (l Logger) Error(message string) {
	Print(message, RED, os.Stderr)
}

func Error(message string) {
	Print(message, RED, os.Stderr)
}

// Print prints a message to a given stream in a defined color
func Print(message string, color string, stream *os.File) {
	fmt.Fprintf(stream, "%s%s%s\n", color, message, RESET)
}
