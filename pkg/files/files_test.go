package files

import (
	"os"
	"reflect"
	"regexp"
	"runtime"
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

var getContentTypeFilesTests = []getContentTypeFilesTest{
	{"8859_1_da.html", "^text/"},
	{"8859_1_de.html", "^text/"},
	{"8859_1_en.html", "^text/"},
	{"8859_1_es.html", "^text/"},
	{"8859_1_fr.html", "^text/"},
	{"8859_1_pt.html", "^text/"},
	{"ascii.txt", "^text/"},
	{"big5.html", "^text/"},
	{"candide-gb18030.txt", "^text/"},
	{"candide-utf-16le.txt", "^application/octet-stream$"}, // no BOM
	{"candide-utf-32be.txt", "^application/octet-stream$"}, // no BOM
	{"candide-utf-8.txt", "^text/"},                        // no BOM
	{"candide-windows-1252.txt", "^text/"},
	{"cp865.txt", "^text/"},
	{"euc_jp.html", "^text/"},
	{"euc_kr.html", "^text/"},
	{"gb18030.html", "^text/"},
	{"html.html", "^text/"},
	{"html.iso88591.html", "^text/"},
	{"html.svg.html", "^text/"},
	{"html.usascii.html", "^text/"},
	{"html.utf8bomdetect.html", "^text/"}, // has BOM
	{"html.utf8bom.html", "^text/"},       // has BOM
	{"html.utf8bomws.html", "^text/"},     // has BOM
	{"html.utf8.html", "^text/"},          // no BOM
	{"html.withbr.html", "^text/"},
	{"iso88591.txt", "^text/"},
	{"koi8_r.txt", "^text/"},
	{"latin1.txt", "^text/"},
	{"rashomon-euc-jp.txt", "^text/"},
	{"rashomon-iso-2022-jp.txt", "^text/"}, // byte 89 is an Esc (ASCII 27)
	{"rashomon-shift-jis.txt", "^text/"},
	{"rashomon-utf-8.txt", "^text/"}, // no BOM
	{"shift_jis.html", "^text/"},
	{"sunzi-bingfa-gb-levels-1-and-2-hz-gb2312.txt", "^text/"},
	{"sunzi-bingfa-gb-levels-1-and-2-utf-8.txt", "^text/"}, // no BOM
	{"sunzi-bingfa-simplified-gbk.txt", "^text/"},
	{"sunzi-bingfa-simplified-utf-8.txt", "^text/"}, // no BOM
	{"sunzi-bingfa-traditional-big5.txt", "^text/"},
	{"sunzi-bingfa-traditional-utf-8.txt", "^text/"}, // no BOM
	{"unsu-joh-eun-nal-euc-kr.txt", "^text/"},
	{"unsu-joh-eun-nal-utf-8.txt", "^text/"},       // no BOM
	{"utf16bebom.txt", "^text/"},                   // has BOM
	{"utf16lebom.txt", "^text/"},                   // has BOM
	{"utf16.txt", "^text/"},                        // has BOM
	{"utf32bebom.txt", "^text/"},                   // has BOM
	{"utf32lebom.txt", "^text/"},                   // has BOM
	{"utf8_bom.html", "^text/"},                    // has BOM
	{"utf8.html", "^text/"},                        // no BOM
	{"utf8-sdl.txt", "^text/"},                     // no BOM
	{"utf8.txt", "^text/"},                         // no BOM
	{"utf8.txt-encoding-test-files.txt", "^text/"}, // no BOM
}

func TestGetContentTypeFiles(t *testing.T) {
	for _, tt := range getContentTypeFilesTests {
		filePath := "../encoding/testdata/" + tt.filename
		contentType, err := GetContentType(filePath)
		if err != nil {
			t.Errorf("GetContentType(%q): expected %v, got %v", tt.filename, "nil", err.Error())
		}
		match, _ := regexp.MatchString(tt.regex, contentType)
		if !match {
			t.Errorf("GetContentType(%q): expected %v, got %v", tt.filename, tt.regex, contentType)
		}
	}
}
