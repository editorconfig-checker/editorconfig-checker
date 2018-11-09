package validators

import (
	"testing"
)

func TestLineEndingWithLf(t *testing.T) {
	if !LineEnding("x", "lf") {
		t.Error("Expected to return true for line without linebreak")
	}

	if !LineEnding("x\n", "lf") {
		t.Error("Expected to return true for a valid file(lf)")
	}

	if LineEnding("x\r", "lf") {
		t.Error("Expected to return false for an invalid file(lf) with \\r")
	}

	if LineEnding("x\r\n", "lf") {
		t.Error("Expected to return false for an invalid file(lf) with \\r\\n")
	}

	if LineEnding("x\ry\nz\n", "lf") {
		t.Error("Expected to return false for mixed file(lf)")
	}
}

func TestLineEndingWithCr(t *testing.T) {
	if !LineEnding("x", "cr") {
		t.Error("Expected to return true for line without linebreak")
	}

	if !LineEnding("x\r", "cr") {
		t.Error("Expected to return true for a valid file(cr)")
	}

	if LineEnding("x\n", "cr") {
		t.Error("Expected to return false for an invalid file(cr) with \\n")
	}

	if LineEnding("x\r\n", "cr") {
		t.Error("Expected to return false for an invalid file(cr) with \\r\\n")
	}

	if LineEnding("x\ry\nz\n", "cr") {
		t.Error("Expected to return false for mixed file(lf)")
	}
}

func TestLineEndingWithCrLf(t *testing.T) {
	if !LineEnding("x", "crlf") {
		t.Error("Expected to return true for line without linebreak")
	}

	if !LineEnding("x\r\n", "crlf") {
		t.Error("Expected to return true for a valid file(crlf)")
	}

	if LineEnding("x\n", "crlf") {
		t.Error("Expected to return false for an invalid file(crlf) with \\n")
	}

	if LineEnding("x\r", "crlf") {
		t.Error("Expected to return false for an invalid file(crlf) with \\r")
	}

	if LineEnding("x\ry\nz\n", "crlf") {
		t.Error("Expected to return false for mixed file(crlf)")
	}
}
