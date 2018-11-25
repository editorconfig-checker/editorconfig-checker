package utils

import (
	// "path/filepath"
	"testing"
)

func TestPathExists(t *testing.T) {
	if !PathExists(".") {
		t.Error("Expected . to be an existing path")
	}

	if PathExists("notexisting") {
		t.Error("Expected \"notexisting\" to not exist")
	}
}
