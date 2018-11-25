package utils

import (
	"path/filepath"
	"testing"
)

func TestIsDirectory(t *testing.T) {
	if !IsDirectory(".") {
		t.Error("Expected . to be a directory")
	}

	absolutePath, _ := filepath.Abs("isDirectory_test.go")
	if IsDirectory(absolutePath) {
		t.Error("Expected this file to not be a directory")
	}
}
