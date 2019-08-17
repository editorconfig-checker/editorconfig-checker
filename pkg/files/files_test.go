package files

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/config"
)

func TestGetContentType(t *testing.T) {
	contentType, _ := GetContentType("./files.go")
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

func TestIsAllowedContentType(t *testing.T) {
	config := config.Config{Allowed_Content_Types: []string{"text/", "application/octet-stream"}}

	if IsAllowedContentType("bla", config) {
		t.Error("Bla shouldn't be an allowed contentType")
	}

	if !IsAllowedContentType("text/", config) {
		t.Error("text/ shouldn't be an allowed contentType")
	}

	if !IsAllowedContentType("text/xml abc", config) {
		t.Error("text/xml shouldn't be an allowed contentType")
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

func TestGetRelativePath(t *testing.T) {
	cwd, _ := os.Getwd()
	filePath := "/bin/ec"
	relativeFilePath, _ := GetRelativePath(cwd + filePath)

	if relativeFilePath != "."+filePath {
		t.Error("Expected the relative filePath to match")
	}

	DIR := "/tmp/stuff"
	os.Remove(DIR)
	err := os.Mkdir(DIR, 0755)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s", err))
	}

	err = os.Chdir(DIR)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s", err))
	}

	err = os.Remove(DIR)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s", err))
	}

	_, err = GetRelativePath(cwd + filePath)

	if err == nil {
		t.Error("Expected an error for a not existing directory")
	}
}

func TestIsExcluded(t *testing.T) {
	result := IsExcluded("./cmd/editorconfig-checker/main.go", config.Config{})

	if result {
		t.Error("Should return true if no excludes are given, got", result)
	}
}
