// Package logger provides functions that are logging related
package logger

import (
	"fmt"
	"io"
	"os"
)

// Colors which can be used
const (
	escSeqYellow = "\x1b[33;1m"
	escSeqGreen  = "\x1b[32;1m"
	escSeqRed    = "\x1b[31;1m"
	escSeqReset  = "\x1b[33;0m"
)

// Logger struct
type Logger struct {
	VerboseEnabled bool
	DebugEnabled   bool
	NoColor        bool
	writer         io.Writer
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

func (l Logger) GetWriter() io.Writer {
	return l.writer
}

// allow users to overwrite the writer used
func (l *Logger) SetWriter(w io.Writer) {
	l.writer = w
}

// apply the settings from the Logger given to the instance
func (l *Logger) Configure(newLogger Logger) {
	l.VerboseEnabled = newLogger.VerboseEnabled
	l.DebugEnabled = newLogger.DebugEnabled
	l.NoColor = newLogger.NoColor
	if newLogger.writer != nil {
		l.SetWriter(newLogger.writer)
	}
}

// Debug prints a message when Debugg is set to true on the Logger
func (l Logger) Debug(format string, a ...interface{}) {
	if l.DebugEnabled {
		message := fmt.Sprintf(format, a...)
		l.println(message)
	}
}

// Verbose prints a message when Verbosee is set to true on the Logger
func (l Logger) Verbose(format string, a ...interface{}) {
	if l.VerboseEnabled {
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
		l.printlnColor(message, escSeqYellow)
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
		l.printlnColor(message, escSeqRed)
	}
}

// println prints a message with a trailing newline
func (l Logger) println(message string) {
	l.lazyInit()
	fmt.Fprintln(l.writer, message)
}

// printlnColor prints a message in a given ANSI-color with a trailing newline
func (l Logger) printlnColor(message string, color string) {
	l.lazyInit()
	fmt.Fprintf(l.writer, "%s%s%s\n", color, message, escSeqReset)
}
