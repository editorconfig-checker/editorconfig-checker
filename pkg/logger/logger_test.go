package logger

import (
	"testing"
)

// Wannbe tests

func TestLoggerDebug(t *testing.T) {
	logger := GetLogger()

	logger.Debug("hello")

	logger.Debugg = true
	logger.Debug("hello")
}

func TestLoggerVerbose(t *testing.T) {
	logger := GetLogger()

	logger.Verbose("hello")

	logger.Verbosee = true
	logger.Verbose("hello")
}

func TestLoggerWarning(t *testing.T) {
	logger := GetLogger()
	logger.Warning("bla%s", "hey")

	logger.NoColor = true
	logger.Warning("bla%s", "hey")
}

func TestLoggerOutput(t *testing.T) {
	logger := GetLogger()
	logger.Output("bla%s", "hey")
}

func TestLoggerError(t *testing.T) {
	logger := GetLogger()
	logger.Error("bla%s", "hey")

	logger.NoColor = true
	logger.Error("bla%s", "hey")
}

func TestPrintColor(t *testing.T) {
	logger := GetLogger()
	logger.printColor("Hello", RED)
}

func TestPrintLogMessage(t *testing.T) {
	logger := GetLogger()

	messages := []LogMessage{
		{Level: "debug", Message: "debug message"},
		{Level: "verbose", Message: "verbose message"},
		{Level: "warning", Message: "warning message"},
		{Level: "error", Message: "error message"},
		{Level: "output", Message: "normal message"},
	}

	logger.PrintLogMessages(messages)
}
