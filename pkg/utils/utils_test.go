package utils

import (
	"os"
	"strings"
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
	contentType, _ := GetContentType("./utils.go")
	if !strings.Contains(contentType, "text/plain") {
		t.Error("Expected getContentType.go to be of type text/plain")
	}

	contentType, _ = GetContentType("./../../docs/logo.png")
	if !strings.Contains(contentType, "image/png") {
		t.Error("Expected getContentType_test.go to be of type application/octet-stream")
	}

	_, err := GetContentType(".")
	if err == nil {
		t.Error("Expected to return an error for a directory")
	}

	_, err = GetContentType("/abc/!@#")
	if err == nil {
		t.Error("Expected to return an error for a not existing file")
	}

	emptyFile, _ := os.Create("empty.txt")
	defer emptyFile.Close()
	defer os.Remove("empty.txt")

	contentType, err = GetContentType("empty.txt")
	if contentType != "" || err != nil {
		t.Error("Expected to return an empty string for an empty file and no error")
	}
}

func TestGetRelativePath(t *testing.T) {
	cwd, _ := os.Getwd()
	filePath := "/bin/ec"
	relativeFilePath, _ := GetRelativePath(cwd + filePath)

	if relativeFilePath != "."+filePath {
		t.Error("Expected the relative filePath to match")
	}

	DIR := "/tmp/stuff"
	os.Remove(DIR)
	os.Mkdir(DIR, 0755)
	os.Chdir(DIR)
	os.Remove(DIR)

	_, err := GetRelativePath(cwd + filePath)

	if err == nil {
		t.Error("Expected an error for a not existing directory")
	}
}
