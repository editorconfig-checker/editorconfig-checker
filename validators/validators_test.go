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

func TestLineEnding(t *testing.T) {
	// Not existing
	if !LineEnding("x", "y") {
		t.Error("Expected to return true for not existing end of line char")
	}

	if !LineEnding("x\n", "y") {
		t.Error("Expected to return true for not existing end of line char")
	}

	if !LineEnding("x\r", "y") {
		t.Error("Expected to return true for not existing end of line char")
	}

	if !LineEnding("x\r\n", "y") {
		t.Error("Expected to return true for not existing end of line char")
	}

	// LF
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

	// CR
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

	// CRLF
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

func TestSpace(t *testing.T) {
	if !Space("", "space", 4) {
		t.Error("Expected empty line to return true regardless of parameter")
	}

	if !Space("x", "space", 4) {
		t.Error("Expected line which starts at the beginning to return true")
	}

	if !Space("    x", "space", 4) {
		t.Error("Expected correctly indented line to return true")
	}

	if Space("     x", "space", 4) {
		t.Error("Expected falsy indented line to return false")
	}

	if Space("   x", "space", 4) {
		t.Error("Expected falsy indented line to return false")
	}

	if !Space("     *", "space", 4) {
		t.Error("Expected correctly indented line to be true with empty block comment line")
	}

	if !Space("     * some comment", "space", 4) {
		t.Error("Expected correctly indented line to be true with block comment")
	}

	if !Space("", "tab", 4) {
		t.Error("Expected if indentStyle is set to tab to return true")
	}

	if !Space("x", "tab", 4) {
		t.Error("Expected if indentStyle is set to tab to return true")
	}

	if !Space("    x", "tab", 4) {
		t.Error("Expected if indentStyle is set to tab to return true")
	}

	if !Space("   x", "tab", 4) {
		t.Error("Expected if indentStyle is set to tab to return true")
	}
}

func TestTab(t *testing.T) {
	if !Tab("", "tab") {
		t.Error("Expected empty line to return true regardless of parameter")
	}

	if !Tab("x", "tab") {
		t.Error("Expected line which starts at the beginning to return true")
	}

	if !Tab("	x", "tab") {
		t.Error("Expected correctly indented line to return true")
	}

	if Tab(" x", "tab") {
		t.Error("Expected falsy indented line to return false")
	}

	if !Tab("", "space") {
		t.Error("Expected if indentStyle is set to tab to return true")
	}

	if !Tab("x", "space") {
		t.Error("Expected if indentStyle is set to tab to return true")
	}

	if !Tab("	x", "space") {
		t.Error("Expected if indentStyle is set to tab to return true")
	}

	if !Tab("   x", "space") {
		t.Error("Expected if indentStyle is set to tab to return true")
	}
}

func TestTrailingWhitespace(t *testing.T) {
	if !TrailingWhitespace("", true) {
		t.Error("Expected empty line to return true regardless of trimTrailingWhitespace parameter")
	}

	if !TrailingWhitespace("", false) {
		t.Error("Expected empty line to return true regardless of trimTrailingWhitespace parameter")
	}

	if !TrailingWhitespace("x", true) {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to true to return true")
	}

	if !TrailingWhitespace("x", false) {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to false to return true")
	}

	// Spaces
	if TrailingWhitespace("x ", true) {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to true to return false")
	}

	if !TrailingWhitespace("x ", false) {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to false to return true")
	}

	if !TrailingWhitespace("x .", true) {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to true to return true")
	}

	if !TrailingWhitespace("x .", false) {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to false to return true")
	}

	// Tabs
	if TrailingWhitespace("x	", true) {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to true to return false")
	}

	if !TrailingWhitespace("x	", false) {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to false to return true")
	}

	if !TrailingWhitespace("x	.", true) {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to true to return true")
	}

	if !TrailingWhitespace("x	.", false) {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to false to return true")
	}
}
