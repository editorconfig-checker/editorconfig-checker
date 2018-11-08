package validators

import (
	"testing"
)

func TestFinalNewline(t *testing.T) {
	if !FinalNewline("x", false, "lf") {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to false")
	}

	if !FinalNewline("x\n", true, "lf") {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to true and correct eol-char is used")
	}

	if !FinalNewline("x\r", true, "cr") {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to true and correct eol-char is used")
	}

	if !FinalNewline("x\r\n", true, "crlf") {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to true and correct eol-char is used")
	}

	if FinalNewline("x\n", true, "cr") {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	if FinalNewline("x\n", true, "crlf") {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	if FinalNewline("x\r", true, "lf") {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	if FinalNewline("x\r", true, "crlf") {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	// TODO: This needs to be fixed
	// if FinalNewline("x\r\n", true, "lf") {
	// 	t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	// }

	if FinalNewline("x\r\n", true, "cr") {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}
}
