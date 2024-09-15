package logger

import (
	"bytes"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

// interface for the closures passed to snapHelper
type loggerTest func()

// a helper function, that runs a closure with a logger, but redirects the Logger to write into a buffer, and then does snapshot testing on the buffer content
func snapHelper(t *testing.T, logger *Logger, test loggerTest) {
	t.Helper()
	buffer := bytes.Buffer{}
	logger.SetWriter(&buffer)
	test()
	snaps.MatchSnapshot(t, buffer.String())
}

func TestLoggerDebug(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Debug("hello")
	})

	snapHelper(t, &logger, func() {
		logger.Debugg = true
		logger.Debug("hello")
	})
}

func TestLoggerVerbose(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Verbose("hello")
	})

	snapHelper(t, &logger, func() {
		logger.Verbosee = true
		logger.Verbose("hello")
	})
}

func TestLoggerWarning(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Warning("bla%s", "hey")
	})

	snapHelper(t, &logger, func() {
		logger.NoColor = true
		logger.Warning("bla%s", "hey")
	})
}

func TestLoggerOutput(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Output("bla%s", "hey")
	})
}

func TestLoggerError(t *testing.T) {
	logger := GetLogger()
	snapHelper(t, &logger, func() {
		logger.Error("bla%s", "hey")
	})

	snapHelper(t, &logger, func() {
		logger.NoColor = true
		logger.Error("bla%s", "hey")
	})
}

func TestPrintColor(t *testing.T) {
	logger := GetLogger()
	snapHelper(t, &logger, func() {
		logger.printColor("Hello", RED)
	})
}
