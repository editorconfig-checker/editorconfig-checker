package utils

import (
	"testing"
)

func TestGetEolChar(t *testing.T) {
	if GetEolChar("lf") != "\n" {
		t.Error("Expected end of line character to be \\n for \"lf\"")
	}

	if GetEolChar("cr") != "\r" {
		t.Error("Expected end of line character to be \\r for \"cr\"")
	}

	if GetEolChar("crlf") != "\r\n" {
		t.Error("Expected end of line character to be \\r\\n for \"crlf\"")
	}

	if GetEolChar("") != "\n" {
		t.Error("Expected end of line character to be \\n as a fallback")
	}
}

func TestIsRegularFile(t *testing.T) {
	if !IsRegularFile("./utils.go") {
		t.Error("Expected utils.go to be a regular file")
	}

	if IsRegularFile("./notExisting.go") {
		t.Error("Expected not existing file not to be a regular file")
	}

	if IsRegularFile(".") {
		t.Error("Expected a directory not to be a regular file")
	}
}

func TestIsDirectory(t *testing.T) {
	if !IsDirectory(".") {
		t.Error("Expected the current directory to be a directory")
	}

	if IsDirectory("./notExisting") {
		t.Error("Expected not existing directory not to be a directory")
	}

	if IsDirectory("./utils.go") {
		t.Error("Expected a file not to be a directory")
	}
}

func TestFileExists(t *testing.T) {
	if !FileExists("./utils.go") {
		t.Error("./utils.go should exist")
	}

	if FileExists("./notExisting.go") {
		t.Error("./notExisting.go should not exist")
	}
}
