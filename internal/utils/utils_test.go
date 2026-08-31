package utils

import (
	"testing"
)

func TestGetEolChar(t *testing.T) {
	getEolCharTests := []struct {
		input    string
		expected string
	}{
		{"lf", "\n"},
		{"cr", "\r"},
		{"crlf", "\r\n"},
		{"", "\n"},
	}

	for _, tt := range getEolCharTests {
		actual := GetEolChar(tt.input)
		if actual != tt.expected {
			t.Errorf("GetEolChar(%s): expected: %v, got: %v", tt.input, tt.expected, actual)
		}
	}
}

func TestIsRegularFile(t *testing.T) {
	isRegularFileTests := []struct {
		input    string
		expected bool
	}{
		{"./utils.go", true},
		{"./notExisting.go", false},
		{".", false},
	}

	for _, tt := range isRegularFileTests {
		actual := IsRegularFile(tt.input)
		if actual != tt.expected {
			t.Errorf("IsRegularFile(%s): expected: %v, got: %v", tt.input, tt.expected, actual)
		}
	}
}

func TestIsDirectory(t *testing.T) {
	isDirectoryTests := []struct {
		input    string
		expected bool
	}{
		{".", true},
		{"./notExisting", false},
		{"./utils.go", false},
	}

	for _, tt := range isDirectoryTests {
		actual := IsDirectory(tt.input)
		if actual != tt.expected {
			t.Errorf("IsDirectory(%s): expected: %v, got: %v", tt.input, tt.expected, actual)
		}
	}
}

func TestFileExists(t *testing.T) {
	fileExistsTests := []struct {
		input    string
		expected bool
	}{
		{"./utils.go", true},
		{"./notExisting.go", false},
	}

	for _, tt := range fileExistsTests {
		actual := FileExists(tt.input)
		if actual != tt.expected {
			t.Errorf("FileExists(%s): expected: %v, got: %v", tt.input, tt.expected, actual)
		}
	}
}
