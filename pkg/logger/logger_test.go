package logger

import (
	"os"
	"testing"
)

// Wannbe tests

func TestWarning(t *testing.T) {
	Warning("Hello")
}

func TestError(t *testing.T) {
	Error("Hello")
}

func TestOutput(t *testing.T) {
	Output("Hello")
}

func TestPrint(t *testing.T) {
	Print("Hello", RED, os.Stdout)
	Print("Hello", GREEN, os.Stdout)
	Print("Hello", YELLOW, os.Stdout)
}
