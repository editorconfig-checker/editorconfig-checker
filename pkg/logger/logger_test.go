package logger

import (
	"testing"
)

// Wannbe tests

func TestLoggerDebug(t *testing.T) {
	logger := Logger{}

	logger.Debug("hello")

	logger = Logger{Debugg: true}
	logger.Debug("hello")
}

func TestLoggerVerbose(t *testing.T) {
	logger := Logger{}

	logger.Verbose("hello")

	logger = Logger{Verbosee: true}
	logger.Verbose("hello")
}

func TestLoggerWarning(t *testing.T) {
	logger := Logger{}
	logger.Warning("bla%s", "hey")

	logger.No_Color = true
	logger.Warning("bla%s", "hey")
}

func TestWarning(t *testing.T) {
	Warning("bla%s", "hey")
}

func TestLoggerOutput(t *testing.T) {
	logger := Logger{}
	logger.Output("bla%s", "hey")
}

func TestOutput(t *testing.T) {
	Output("bla%s", "hey")
}

func TestLoggerError(t *testing.T) {
	logger := Logger{}
	logger.Error("bla%s", "hey")

	logger.No_Color = true
	logger.Error("bla%s", "hey")
}

func TestError(t *testing.T) {
	Error("bla%s", "hey")
}

func TestPrintColor(t *testing.T) {
	PrintColor("Hello", RED)
}

func TestPrint(t *testing.T) {
	Print("Hello")
}
