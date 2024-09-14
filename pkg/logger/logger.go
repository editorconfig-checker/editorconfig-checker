// Package logger provides functions that are logging related
package logger

import (
	"fmt"
	"io"
	"os"
)

// Colors which can be used
const (
	YELLOW = "\x1b[33;1m"
	GREEN  = "\x1b[32;1m"
	RED    = "\x1b[31;1m"
	RESET  = "\x1b[33;0m"
)

type LogMessage struct {
	Level   string
	Message string
}

// Logger struct
type Logger struct {
	Verbosee bool
	Debugg   bool
	NoColor  bool
	writer   io.Writer
}

func GetLogger() Logger {
	logger := Logger{}
	logger.Init()
	return logger
}

// initialize the Logger to write to standard output
func (l *Logger) Init() {
	l.writer = os.Stdout
}

// ensure the Logger is initialized on first print
func (l *Logger) lazyInit() {
	if l.writer == nil {
		l.Init()
	}
}

// allow users to overwrite the writer used
func (l *Logger) SetWriter(w io.Writer) {
	l.writer = w
}

// Debug prints a message when Debugg is set to true on the Logger
func (l Logger) Debug(format string, a ...interface{}) {
	if l.Debugg {
		message := fmt.Sprintf(format, a...)
		l.println(message)
	}
}

// Verbose prints a message when Verbosee is set to true on the Logger
func (l Logger) Verbose(format string, a ...interface{}) {
	if l.Verbosee {
		message := fmt.Sprintf(format, a...)
		l.println(message)
	}
}

// Warning prints a warning message to Stdout in yellow
func (l Logger) Warning(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	if l.NoColor {
		l.println(message)
	} else {
		l.printlnColor(message, YELLOW)
	}
}

func (l Logger) PrintLogMessage(message LogMessage) {
	switch message.Level {
	case "error":
		l.Error(message.Message)
	case "warning":
		l.Warning(message.Message)
	case "debug":
		l.Debug(message.Message)
	case "verbose":
		l.Verbose(message.Message)
	default:
		l.Output(message.Message)
	}
}

func (l Logger) PrintLogMessages(messages []LogMessage) {
	for _, message := range messages {
		l.PrintLogMessage(message)
	}
}

// Output prints a message on Stdout in 'normal' color
func (l Logger) Output(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	l.println(message)
}

// Error prints an error message to Stdout in red
func (l Logger) Error(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	if l.NoColor {
		l.println(message)
	} else {
		l.printlnColor(message, RED)
	}
}

// println prints a message with a trailing newline
func (l Logger) println(message string) {
	l.lazyInit()
	fmt.Fprintf(l.writer, "%s\n", message)
}

// printColor prints a message in a given ANSI-color
func (l Logger) printColor(message string, color string) {
	l.lazyInit()
	fmt.Fprintf(l.writer, "%s%s%s", color, message, RESET)
}

// printlnColor prints a message in a given ANSI-color with a trailing newline
func (l Logger) printlnColor(message string, color string) {
	l.lazyInit()
	fmt.Fprintf(l.writer, "%s%s%s\n", color, message, RESET)
}
