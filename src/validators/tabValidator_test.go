package validators

import (
	"testing"
)

func TestTabValidator(t *testing.T) {
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
