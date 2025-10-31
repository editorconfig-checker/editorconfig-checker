package files

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	// x-release-please-end
)

func TestGetContentType(t *testing.T) {
	inputFile := "./files.go"
	expected := "text/plain"
	contentType, _ := GetContentType(inputFile)
	if !strings.Contains(contentType, expected) {
		t.Errorf("GetContentType(%q): expected %v, got %v", inputFile, expected, contentType)
	}

	inputFile = "./../../docs/logo.png"
	expected = "image/png"
	contentType, _ = GetContentType(inputFile)
	if !strings.Contains(contentType, expected) {
		t.Errorf("GetContentType(%q): expected %v, got %v", inputFile, expected, contentType)
	}

	inputFile = "."
	_, err := GetContentType(inputFile)
	if err == nil {
		t.Errorf("GetContentType(%q): expected %v, got %v", inputFile, "an error", "nil")
	}

	inputFile = "a non-existent file"
	_, err = GetContentType(inputFile)
	if err == nil {
		t.Errorf("GetContentType(%q): expected %v, got %v", inputFile, "an error", "nil")
	}

	inputFile = "testdata/empty.txt"
	contentType, err = GetContentType(inputFile)
	if err != nil {
		t.Errorf("GetContentType(%q): expected %v, got %v", inputFile, "nil", err.Error())
	}
	expected = ""
	if contentType != expected {
		t.Errorf("GetContentType(%q): expected %v, got %v", inputFile, expected, contentType)
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

	if runtime.GOOS == "windows" {
		t.Skip("Windows fails if current directory is deleted")
	}

	var DIR string
	expectedPath := "../.."
	if runtime.GOOS == "darwin" {
		DIR = "/private"
		expectedPath += "/.."
	}
	DIR += "/tmp/stuff"
	os.Remove(DIR)
	err := os.Mkdir(DIR, 0o755)
	if err != nil {
		panic(err)
	}

	t.Chdir(DIR)

	arg := "/foo" + DIR + filePath
	expectedPath += arg

	// Check with the current directory ("/tmp/stuff") in the middle of the given file path
	relativeFilePath, _ = GetRelativePath(arg)
	if relativeFilePath != expectedPath {
		t.Errorf("GetRelativePath(%s): expected: %v, got: %v", arg, expectedPath, relativeFilePath)
	}

	err = os.Remove(DIR)
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == "darwin" {
		t.Skip("Darwin can obtain the current working directory even if it is deleted")
	}
	_, err = GetRelativePath(cwd + filePath)

	if err == nil {
		t.Error("Expected an error for a not existing directory")
	}
}

func TestAddToFiles(t *testing.T) {
	configuration := config.NewConfig(nil)
	configuration.AllowedContentTypes = nil
	excludedFileConfiguration := config.NewConfig(nil)
	excludedFileConfiguration.Exclude = []string{"files"}
	addToFilesTests := []struct {
		filePaths []string
		filePath  string
		config    config.Config
		expected  []string
	}{
		{
			[]string{},
			"./files.go",
			*excludedFileConfiguration,
			[]string{},
		},
		{
			[]string{"./files.go"},
			"./files.go",
			*configuration,
			[]string{"./files.go"},
		},
	}

	for _, tt := range addToFilesTests {
		actual := AddToFiles(tt.filePaths, tt.filePath, tt.config)

		if !reflect.DeepEqual(actual, tt.expected) {
			t.Error(actual)
			t.Error(tt.expected)
			t.Errorf("AddToFiles(%s, %s, %+v): expected: %v, got: %v", tt.filePaths, tt.filePath, tt.config, tt.expected, actual)
		}
	}
}

func TestGetFiles(t *testing.T) {
	docsConfig := config.NewConfig(nil)
	docsConfig.PassedFiles = []string{"./../../docs/"}
	configurations := []*config.Config{
		config.NewConfig(nil),
		docsConfig,
	}

	for _, configuration := range configurations {
		_, err := GetFiles(*configuration)
		if err != nil {
			t.Errorf("GetFiles(): expected nil, got %s", err.Error())
		}

		configuration.PassedFiles = []string{"."}
		files, err := GetFiles(*configuration)

		if len(files) > 0 && err != nil {
			t.Errorf("GetFiles(.): expected nil, got %s", err.Error())
		}
	}
}

type getContentTypeFilesTest struct {
	filename string
	regex    string
}

func setup() {
	const testResultsJson = "../encoding/test-results.json"

	f, err := os.Open(testResultsJson)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(f).Decode(&tests)
	if err != nil {
		panic(err)
	}
	f.Close()
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Filename < tests[j].Filename
	})
}

type test struct {
	Filename   string
	Encoding   string
	Charset    string
	Errored    bool
	Binary     bool
	Confidence float64
	Comment    string
}

var tests = []test{}

var exceptions = map[string]string{
	"testdata/text/candide-utf-16le.txt":                                "application/octet-stream",
	"testdata/text/candide-utf-32be.txt":                                "application/octet-stream",
	"testdata/uchardet/ja/utf-16be.txt":                                 "application/octet-stream",
	"testdata/uchardet/ja/utf-16le.txt":                                 "application/octet-stream",
	"testdata/uchardet/ko/iso-2022-kr.txt":                              "application/octet-stream",
	"testdata/wpt/legacy-mb-japanese/iso-2022-jp/iso2022jp_errors.html": "application/octet-stream",
	"testdata/wpt/resources/utf-32-big-endian-nobom.html":               "application/octet-stream",
	"testdata/wpt/resources/utf-32-big-endian-nobom.xml":                "application/octet-stream",
	"testdata/wpt/resources/utf-32-little-endian-nobom.html":            "application/octet-stream",
	"testdata/wpt/resources/utf-32-little-endian-nobom.xml":             "application/octet-stream",
}

func TestGetContentTypeFiles(t *testing.T) {
	setup()
	for _, tt := range tests {
		regex := "^text/"
		exception, ok := exceptions[tt.Filename]
		if ok {
			regex = exception
		}
		filePath := "../encoding/" + tt.Filename
		contentType, err := GetContentType(filePath)
		if err != nil {
			t.Errorf("GetContentType (%q): got %v, want %v", filePath, err.Error(), "nil")
		}
		match, _ := regexp.MatchString(regex, contentType)
		if !match {
			t.Errorf("GetContentType(%q): got %v, want %v", filePath, contentType, regex)
		}
	}
}
