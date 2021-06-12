package files

import (
	"os"
	"reflect"
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
		t.Error("Expected getContentTypetest.go to be of type application/octet-stream")
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
	configuration := config.Config{AllowedContentTypes: []string{"text/", "application/octet-stream"}}
	isAllowedContentTypeTests := []struct {
		contentType string
		config      config.Config
		expected    bool
	}{
		{"bla", configuration, false},
		{"text/", configuration, true},
		{"text/xml abc", configuration, true},
	}

	for _, tt := range isAllowedContentTypeTests {
		actual := IsAllowedContentType(tt.contentType, tt.config)
		if actual != tt.expected {
			t.Errorf("IsAllowedContentType(%s, %+v): expected: %v, got: %v", tt.contentType, tt.config, tt.expected, actual)
		}
	}
}

func TestPathExists(t *testing.T) {
	pathExistsTests := []struct {
		path     string
		expected bool
	}{
		{".", true},
		{"notexisting", false},
	}

	for _, tt := range pathExistsTests {
		actual := PathExists(tt.path)
		if actual != tt.expected {
			t.Errorf("PathExists(%s): expected: %v, got: %v", tt.path, tt.expected, actual)
		}
	}
}

func TestGetRelativePath(t *testing.T) {
	// Should return paths that are already relative unchanged
	relativeFilePath, _ := GetRelativePath("bin/ec")
	if relativeFilePath != "bin/ec" {
		t.Errorf("GetRelativePath(%s): expected: %v, got: %v", "bin/ec", "bin/ec", relativeFilePath)
	}

	// Should convert absolute paths to be relative to current directory
	cwd, _ := os.Getwd()
	filePath := "/bin/ec"
	relativeFilePath, _ = GetRelativePath(cwd + filePath)

	if relativeFilePath != "bin/ec" {
		t.Errorf("GetRelativePath(%s): expected: %v, got: %v", cwd+filePath, "bin/ec", relativeFilePath)
	}

	DIR := "/tmp/stuff"
	os.Remove(DIR)
	err := os.Mkdir(DIR, 0755)
	if err != nil {
		panic(err)
	}

	err = os.Chdir(DIR)
	if err != nil {
		panic(err)
	}

	// Check with the current directory ("/tmp/stuff") in the middle of the given file path
	relativeFilePath, _ = GetRelativePath("/foo" + DIR + filePath)
	if relativeFilePath != "../../foo"+DIR+filePath {
		t.Errorf("GetRelativePath(%s): expected: %v, got: %v", "/foo"+DIR+filePath, "../../foo"+DIR+filePath, relativeFilePath)
	}

	err = os.Remove(DIR)
	if err != nil {
		panic(err)
	}

	_, err = GetRelativePath(cwd + filePath)

	if err == nil {
		t.Error("Expected an error for a not existing directory")
	}
}

func TestAddToFiles(t *testing.T) {
	configuration := config.Config{}
	excludedFileConfiguration := config.Config{Exclude: []string{"files.go"}}
	type Expected struct {
		filePaths []string
		err       error
	}
	addToFilesTests := []struct {
		filePaths []string
		filePath  string
		config    config.Config
		expected  Expected
	}{
		{[]string{},
			"./files.go",
			excludedFileConfiguration,
			Expected{[]string{}, nil}},
		{[]string{"./files.go"},
			"./files.go",
			configuration,
			Expected{[]string{"./files.go"}, nil}},
		{[]string{""},
			"./not-existing.go",
			configuration,
			Expected{[]string{""}, nil}},
		{[]string{"./files.go"},
			"./not-existing.go",
			configuration,
			Expected{[]string{"./files.go"}, nil}},
	}

	for _, tt := range addToFilesTests {
		actual, err := AddToFiles(tt.filePaths, tt.filePath, tt.config)

		if !reflect.DeepEqual(actual, tt.expected.filePaths) && err != tt.expected.err {
			t.Error(actual)
			t.Error(tt.expected)
			t.Errorf("AddToFiles(%s, %s, %+v): expected: %v, got: %v, err: %v", tt.filePaths, tt.filePath, tt.config, tt.expected, actual, err)
		}
	}
}

func TestGetFiles(t *testing.T) {
	configuration := config.Config{}
	files, err := GetFiles(configuration)

	if len(files) > 0 && err == nil {
		t.Errorf("Error expected")
	}
}
