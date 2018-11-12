package validators

import (
	"testing"
)

func TestSpaceValidator(t *testing.T) {
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
