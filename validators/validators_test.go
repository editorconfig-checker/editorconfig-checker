package validators

import (
	"testing"
)

func TestFinalNewline(t *testing.T) {
	if FinalNewline("x", false, "lf") != nil {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to false")
	}

	if FinalNewline("x\n", true, "lf") != nil {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to true and correct eol-char is used")
	}

	if FinalNewline("x\r", true, "cr") != nil {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to true and correct eol-char is used")
	}

	if FinalNewline("x\r\n", true, "crlf") != nil {
		t.Error("Expected FinalNewline to be true if insertFinalNewline is set to true and correct eol-char is used")
	}

	if FinalNewline("x\n", true, "cr") == nil {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	if FinalNewline("x\n", true, "crlf") == nil {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	if FinalNewline("x\r", true, "lf") == nil {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	if FinalNewline("x\r", true, "crlf") == nil {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}

	// TODO: This needs to be fixed
	// if FinalNewline("x\r\n", true, "lf")  == nil {
	// 	t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	// }

	if FinalNewline("x\r\n", true, "cr") == nil {
		t.Error("Expected FinalNewline to be false if insertFinalNewline is set to true and the wrong eol-char is used")
	}
}

func TestLineEnding(t *testing.T) {
	// Not existing
	if LineEnding("x", "y") != nil {
		t.Error("Expected to return true for not existing end of line char")
	}

	if LineEnding("x\n", "y") != nil {
		t.Error("Expected to return true for not existing end of line char")
	}

	if LineEnding("x\r", "y") != nil {
		t.Error("Expected to return true for not existing end of line char")
	}

	if LineEnding("x\r\n", "y") != nil {
		t.Error("Expected to return true for not existing end of line char")
	}

	// LF
	if LineEnding("x", "lf") != nil {
		t.Error("Expected to return true for line without linebreak")
	}

	if LineEnding("x\n", "lf") != nil {
		t.Error("Expected to return true for a valid file(lf)")
	}

	if LineEnding("x\r", "lf") == nil {
		t.Error("Expected to return false for an invalid file(lf) with \\r")
	}

	if LineEnding("x\r\n", "lf") == nil {
		t.Error("Expected to return false for an invalid file(lf) with \\r\\n")
	}

	if LineEnding("x\ry\nz\n", "lf") == nil {
		t.Error("Expected to return false for mixed file(lf)")
	}

	// CR
	if LineEnding("x", "cr") != nil {
		t.Error("Expected to return true for line without linebreak")
	}

	if LineEnding("x\r", "cr") != nil {
		t.Error("Expected to return true for a valid file(cr)")
	}

	if LineEnding("x\n", "cr") == nil {
		t.Error("Expected to return false for an invalid file(cr) with \\n")
	}

	if LineEnding("x\r\n", "cr") == nil {
		t.Error("Expected to return false for an invalid file(cr) with \\r\\n")
	}

	if LineEnding("x\ry\nz\n", "cr") == nil {
		t.Error("Expected to return false for mixed file(lf)")
	}

	// CRLF
	if LineEnding("x", "crlf") != nil {
		t.Error("Expected to return true for line without linebreak")
	}

	if LineEnding("x\r\n", "crlf") != nil {
		t.Error("Expected to return true for a valid file(crlf)")
	}

	if LineEnding("x\n", "crlf") == nil {
		t.Error("Expected to return false for an invalid file(crlf) with \\n")
	}

	if LineEnding("x\r", "crlf") == nil {
		t.Error("Expected to return false for an invalid file(crlf) with \\r")
	}

	if LineEnding("x\ry\nz\n", "crlf") == nil {
		t.Error("Expected to return false for mixed file(crlf)")
	}
}

func TestIndentation(t *testing.T) {
	if (Indentation("    x", "space", 4)) != nil {
		t.Error("Expected correctly indented line to return an nil")
	}

	if (Indentation("   x", "space", 4)) == nil {
		t.Error("Expected wrong indented line to return an error")
	}

	if (Indentation("	x", "tab", 0)) != nil {
		t.Error("Expected correctly indented line to return an nil")
	}

	if (Indentation("   x", "tab", 0)) == nil {
		t.Error("Expected wrong indented line to return an error")
	}

	if (Indentation("	x", "x", 0)) != nil {
		t.Error("Expected unknown indentation to return nil")
	}

	if (Indentation("   x", "x", 0)) != nil {
		t.Error("Expected unknown indentation to return nil")
	}
}

func TestSpace(t *testing.T) {
	if Space("", 4) != nil {
		t.Error("Expected empty line to return true regardless of parameter")
	}

	if Space("x", 0) != nil {
		t.Error("Expected call with indentSize 0 to always return nil")
	}

	if Space("x", 4) != nil {
		t.Error("Expected line which starts at the beginning to return true")
	}

	if Space("    x", 4) != nil {
		t.Error("Expected correctly indented line to return true")
	}

	if Space("     x", 4) == nil {
		t.Error("Expected falsy indented line to return false")
	}

	if Space("   x", 4) == nil {
		t.Error("Expected falsy indented line to return false")
	}

	if Space("     *", 4) != nil {
		t.Error("Expected correctly indented line to be true with empty block comment line")
	}

	if Space("     * some comment", 4) != nil {
		t.Error("Expected correctly indented line to be true with block comment")
	}
}

func TestTab(t *testing.T) {
	if Tab("") != nil {
		t.Error("Expected empty line to return true regardless of parameter")
	}

	if Tab("x") != nil {
		t.Error("Expected line which starts at the beginning to return true")
	}

	if Tab("	x") != nil {
		t.Error("Expected correctly indented line to return true")
	}

	if Tab(" x") == nil {
		t.Error("Expected falsy indented line to return false")
	}
}

func TestTrailingWhitespace(t *testing.T) {
	if TrailingWhitespace("", true) != nil {
		t.Error("Expected empty line to return true regardless of trimTrailingWhitespace parameter")
	}

	if TrailingWhitespace("", false) != nil {
		t.Error("Expected empty line to return true regardless of trimTrailingWhitespace parameter")
	}

	if TrailingWhitespace("x", true) != nil {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to true to return true")
	}

	if TrailingWhitespace("x", false) != nil {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to false to return true")
	}

	// Spaces
	if TrailingWhitespace("x ", true) == nil {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to true to return false")
	}

	if TrailingWhitespace("x ", false) != nil {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to false to return true")
	}

	if TrailingWhitespace("x .", true) != nil {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to true to return true")
	}

	if TrailingWhitespace("x .", false) != nil {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to false to return true")
	}

	// Tabs
	if TrailingWhitespace("x	", true) == nil {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to true to return false")
	}

	if TrailingWhitespace("x	", false) != nil {
		t.Error("Expected line with trailing space and trimTrailingWhitespace set to false to return true")
	}

	if TrailingWhitespace("x	.", true) != nil {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to true to return true")
	}

	if TrailingWhitespace("x	.", false) != nil {
		t.Error("Expected line with no trailing space and trimTrailingWhitespace set to false to return true")
	}
}
