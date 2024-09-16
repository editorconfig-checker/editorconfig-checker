package logger

import (
	"bytes"
	"os"
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
		logger.Debug("this text should not be printed by a logger with default config")
	})

	snapHelper(t, &logger, func() {
		logger.DebugEnabled = true
		logger.Debug("this text should be printed when debug was enabled")
	})
}

func TestLoggerVerbose(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Verbose("this text should not be printed by a logger with default config")
	})

	snapHelper(t, &logger, func() {
		logger.VerboseEnabled = true
		logger.Verbose("hello")
	})
}

func TestLoggerWarning(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Warning("this text should be printed by a logger with default config %s", "(and in color)")
	})

	snapHelper(t, &logger, func() {
		logger.NoColor = true
		logger.Warning("this text should be printed by a logger with default config %s", "(but not be colorized)")
	})
}

func TestLoggerOutput(t *testing.T) {
	logger := GetLogger()

	snapHelper(t, &logger, func() {
		logger.Output("plain output should be printed always %s", "(also supporting format strings)")
	})
}

func TestLoggerError(t *testing.T) {
	logger := GetLogger()
	snapHelper(t, &logger, func() {
		logger.Error("this text should be printed by a logger with default config %s", "(and in color)")
	})

	snapHelper(t, &logger, func() {
		logger.NoColor = true
		logger.Error("this text should be printed by a logger with default config %s", "(but not be colorized)")
	})
}

func TestConfigure(t *testing.T) {
	modifiableLogger := Logger{
		VerboseEnabled: false,
		DebugEnabled:   false,
		NoColor:        false,
		writer:         nil,
	}

	if modifiableLogger.VerboseEnabled {
		t.Errorf("Assumption broken: VerboseEnabled was true already")
	}
	if modifiableLogger.DebugEnabled {
		t.Errorf("Assumption broken: DebugEnabled was true already")
	}
	if modifiableLogger.NoColor {
		t.Errorf("Assumption broken: NoColor was true already")
	}
	if modifiableLogger.writer != nil {
		t.Errorf("Assumption broken: writer was set already")
	}

	configLogger := Logger{
		VerboseEnabled: true,
		DebugEnabled:   true,
		NoColor:        true,
		writer:         os.Stderr,
	}

	modifiableLogger.Configure(configLogger)

	if !modifiableLogger.VerboseEnabled {
		t.Errorf("Configuring a logger with another logger did not lead to VerboseEnabled becoming true")
	}
	if !modifiableLogger.DebugEnabled {
		t.Errorf("Configuring a logger with another logger did not lead to DebugEnabled becoming true")
	}
	if !modifiableLogger.NoColor {
		t.Errorf("Configuring a logger with another logger did not lead to NoColor becoming true")
	}
	if modifiableLogger.writer != os.Stderr {
		t.Errorf("Configuring a logger with another logger did not lead to writer becoming set to os.Stderr")
	}
}
