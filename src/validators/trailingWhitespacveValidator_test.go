package validators

import (
	"testing"
)

func TestTrailingWhitespaceWithCommon(t *testing.T) {
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
}

func TestTrailingWhitespaceWithSpace(t *testing.T) {
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
}

func TestTrailingWhitespaceWithTab(t *testing.T) {
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
