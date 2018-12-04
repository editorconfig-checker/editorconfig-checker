package utils

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	testSlice := []string{"a", "b", "c"}

	if !StringSliceContains(testSlice, "a") {
		t.Error("Expected \"a\" to be part of the slice")
	}

	if StringSliceContains(testSlice, "z") {
		t.Error("Expected \"z\" to be not part of the slice")
	}
}

func TestPathExists(t *testing.T) {
	if !PathExists(".") {
		t.Error("Expected . to be an existing path")
	}

	if PathExists("notexisting") {
		t.Error("Expected \"notexisting\" to not exist")
	}
}

func TestIsDirectory(t *testing.T) {
	if !IsDirectory(".") {
		t.Error("Expected . to be a directory")
	}

	absolutePath, _ := filepath.Abs("utils_test.go")
	if IsDirectory(absolutePath) {
		t.Error("Expected this file to not be a directory")
	}
}

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

func TestGetContentType(t *testing.T) {
	contentType := GetContentType("./utils.go")
	if !strings.Contains(contentType, "text/plain") {
		t.Error("Expected getContentType.go to be of type text/plain")
	}

	contentType = GetContentType("./../docs/logo.png")
	if !strings.Contains(contentType, "image/png") {
		t.Error("Expected getContentType_test.go to be of type application/octet-stream")
	}
}
