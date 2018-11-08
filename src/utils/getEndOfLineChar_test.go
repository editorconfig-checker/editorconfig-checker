package utils

import (
	"testing"
)

func TestGetEolChar(t *testing.T) {
	if GetEolChar("lf") != "\n" {
		t.Error("Expected end of line character to be \\n for \"lf\"")
	}

	if GetEolChar("cr") != "\r" {
		t.Error("Expected end of line character to be \\r for \"cr\"")
	}

	if GetEolChar("crlf") != "\r\n" {
		t.Error("Expected end of line character to be \\r\\n for \"crlf\"")
	}

	if GetEolChar("") != "\n" {
		t.Error("Expected end of line character to be \\n as a fallback")
	}
}
