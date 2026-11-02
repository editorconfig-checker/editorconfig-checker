package files

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/config"
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
	relativeFilePath, _ := GetRelativePath("bin/editorconfig-checker")
	if relativeFilePath != "bin/editorconfig-checker" {
		t.Errorf("GetRelativePath(%s): expected: %v, got: %v", "bin/editorconfig-checker", "bin/editorconfig-checker", relativeFilePath)
	}

	// Should convert absolute paths to be relative to current directory
	cwd, _ := os.Getwd()
	filePath := "/bin/editorconfig-checker"
	relativeFilePath, _ = GetRelativePath(cwd + filePath)

	if relativeFilePath != "bin/editorconfig-checker" {
		t.Errorf("GetRelativePath(%s): expected: %v, got: %v", cwd+filePath, "bin/editorconfig-checker", relativeFilePath)
	}

	if runtime.GOOS == "windows" {
		crossDrivePath := "Z:\\there\\drive\\file.txt"
		relCross, err := GetRelativePath(crossDrivePath)
		if err != nil {
			t.Errorf("GetRelativePath cross-drive unexpected error: %v", err)
		}
		if relCross != "Z:/there/drive/file.txt" {
			t.Errorf("GetRelativePath cross-drive expected Z:/there/drive/file.txt, got %s", relCross)
		}
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
		{
			[]string{},
			"./does/not/exist.txt",
			*configuration,
			[]string{},
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

func TestGetFilesFromGitWithSpecialCharacters(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	// git exports GIT_DIR and the like to hooks, and the tests may run from one
	// (see .pre-commit-config.yaml). Unset them, so that the git commands below
	// act on the test repository and not on the outer one.
	vars, err := exec.Command("git", "rev-parse", "--local-env-vars").Output()
	if err != nil {
		t.Fatalf("setup git rev-parse: %v", err)
	}
	for _, name := range strings.Fields(string(vars)) {
		t.Setenv(name, "") // so that it is restored after the test
		os.Unsetenv(name)
	}

	names := []string{
		"plain.txt",
		"with space.txt",
		// git ls-files prints these quoted and escaped unless -z is given
		"café.md",
		"日本語.txt",
	}
	if runtime.GOOS != "windows" {
		// same, but these characters are not allowed in file names on Windows
		names = append(names, `quo"te.txt`, `back\slash.txt`, "tab\there.txt", "new\nline.txt")
	}

	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("setup git %s: %v\n%s", args[0], err, out)
		}
	}

	dir := t.TempDir()
	t.Chdir(dir)
	git("init")
	for i, name := range names {
		if err := os.WriteFile(name, []byte("hello\n"), 0o644); err != nil {
			t.Fatalf("setup writefile: %v", err)
		}
		// track every other file, so that both --cached and --others are covered
		if i%2 == 0 {
			git("add", "--", name)
		}
	}

	files, err := GetFiles(*config.NewConfig(nil))
	if err != nil {
		t.Fatalf("GetFiles(git): unexpected error %s", err.Error())
	}

	sort.Strings(names)
	sort.Strings(files)
	if !reflect.DeepEqual(files, names) {
		t.Errorf("GetFiles(git): expected %q, got %q", names, files)
	}
}

func TestGetFilesGlobPattern(t *testing.T) {
	globConfig := config.NewConfig(nil)
	// Root-level markdown files (README.md, CHANGELOG.md, MAINTAINERS.md,
	// CONTRIBUTING.md) exist relative to pkg/files when we walk up two levels.
	globConfig.PassedFiles = []string{"./../../*.md"}

	files, err := GetFiles(*globConfig)
	if err != nil {
		t.Fatalf("GetFiles(glob): expected nil, got %s", err.Error())
	}
	if len(files) == 0 {
		t.Fatalf("GetFiles(glob): expected at least one match for ./../../*.md, got 0")
	}
	for _, f := range files {
		if !strings.HasSuffix(f, ".md") {
			t.Errorf("GetFiles(glob): expected .md matches only, got %q", f)
		}
	}
}

func TestGetFilesMissingPathWithoutMetaErrors(t *testing.T) {
	cfg := config.NewConfig(nil)
	cfg.PassedFiles = []string{"./does/not/exist.txt"}

	// A path the user named explicitly, with no glob metacharacters, is a typo
	// rather than a pattern that matched nothing, so it is an error.
	_, err := GetFiles(*cfg)
	if err == nil {
		t.Fatal("GetFiles(missing): expected an error for a path that does not exist, got nil")
	}
	if !strings.Contains(err.Error(), "./does/not/exist.txt") {
		t.Errorf("GetFiles(missing): error should name the path, got %q", err.Error())
	}
}

func TestGetFilesGlobNoMatchesKeepsPath(t *testing.T) {
	cfg := config.NewConfig(nil)
	cfg.PassedFiles = []string{"./nonexistent-prefix-*.xyz"}

	if _, err := GetFiles(*cfg); err != nil {
		t.Fatalf("GetFiles(no-match glob): unexpected error %s", err.Error())
	}
}

func TestHasGlobMetaDetectsAllMetachars(t *testing.T) {
	cases := []string{"*.go", "file?.txt", "log[0-9].txt", "a/b/c*"}
	for _, p := range cases {
		if !hasGlobMeta(p) {
			t.Errorf("hasGlobMeta(%q): expected true, got false", p)
		}
	}
	if hasGlobMeta("plain/path.txt") {
		t.Errorf("hasGlobMeta(plain/path.txt): expected false, got true")
	}
}

func TestHasGlobMetaSkipsEscapedMetachars(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("filepath.Match treats backslash as a path separator on Windows")
	}
	cases := []struct {
		path string
		want bool
	}{
		{`\*.go`, false},
		{`\?.txt`, false},
		{`\[abc\]`, false},
		{`dir/\*`, false},
		{`\`, false},
		{`foo\*bar*`, true},
	}
	for _, tc := range cases {
		if got := hasGlobMeta(tc.path); got != tc.want {
			t.Errorf("hasGlobMeta(%q): got %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestResolvePassedFileInvalidPatternErrors(t *testing.T) {
	_, err := resolvePassedFile("nonexistent-[unclosed")
	if err == nil {
		t.Fatalf("resolvePassedFile(invalid pattern): expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid pattern") {
		t.Errorf("resolvePassedFile(invalid pattern): expected error to mention \"invalid pattern\", got %q", err.Error())
	}
}

func TestGetFilesGlobPropagatesInvalidPattern(t *testing.T) {
	cfg := config.NewConfig(nil)
	cfg.PassedFiles = []string{"nonexistent-[unclosed"}

	if _, err := GetFiles(*cfg); err == nil {
		t.Fatalf("GetFiles(invalid pattern): expected error, got nil")
	}
}

func TestGetFilesGlobMatchesDirectory(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "dirA")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("setup mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("setup writefile: %v", err)
	}

	cfg := config.NewConfig(nil)
	cfg.PassedFiles = []string{filepath.Join(tmp, "dir*")}

	files, err := GetFiles(*cfg)
	if err != nil {
		t.Fatalf("GetFiles(dir glob): unexpected error %s", err.Error())
	}
	found := false
	for _, f := range files {
		if strings.HasSuffix(f, "a.txt") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetFiles(dir glob): expected dir*/a.txt to be walked, got %v", files)
	}
}

func TestGetMatchingExclude(t *testing.T) {
	cfg := config.NewConfig(nil)

	testCases := []struct {
		path            string
		expectedPattern string
		expectedMatch   bool
	}{
		{"changes.patch", `\.patch$`, true},
		{"dir/build.log", `\.log$`, true},
		{"node_modules/index.js", `(^|/)node_modules/`, true},
		{"src/index.js", "", false},
		{"main.go", "", false},
	}

	for _, tc := range testCases {
		pattern, isExcluded, err := GetMatchingExclude(tc.path, *cfg)
		if err != nil {
			t.Errorf("GetMatchingExclude(%q): unexpected error: %v", tc.path, err)
		}
		if isExcluded != tc.expectedMatch {
			t.Errorf("GetMatchingExclude(%q): expected match %v, got %v", tc.path, tc.expectedMatch, isExcluded)
		}
		if pattern != tc.expectedPattern {
			t.Errorf("GetMatchingExclude(%q): expected pattern %q, got %q", tc.path, tc.expectedPattern, pattern)
		}
	}

	customCfg := config.NewConfig(nil)
	customCfg.Exclude = []string{`custom/.*`}
	pattern, isExcluded, err := GetMatchingExclude("custom/file.txt", *customCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isExcluded || pattern != `custom/.*` {
		t.Errorf("expected match for custom pattern, got pattern %q, match %v", pattern, isExcluded)
	}

	ignoreDefaultsCfg := config.NewConfig(nil)
	ignoreDefaultsCfg.IgnoreDefaults = true
	pattern, isExcluded, err = GetMatchingExclude("changes.patch", *ignoreDefaultsCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isExcluded {
		t.Errorf("expected no match when IgnoreDefaults is true, got pattern %q", pattern)
	}

	invalidCfg := config.NewConfig(nil)
	invalidCfg.Exclude = []string{"["}
	_, _, err = GetMatchingExclude("main.go", *invalidCfg)
	if err == nil {
		t.Errorf("expected error for invalid regex pattern, got nil")
	}

	splitCfg := config.NewConfig(nil)
	splitCfg.IgnoreDefaults = true
	splitCfg.Exclude = []string{"(", ")"}
	pattern, isExcluded, err = GetMatchingExclude("main.go", *splitCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isExcluded || pattern != "" {
		t.Errorf("expected isExcluded=true and pattern=\"\", got match=%v, pattern=%q", isExcluded, pattern)
	}
}

func TestGetFilesExplicitExcludedFile(t *testing.T) {
	dir := t.TempDir()
	patchFile := filepath.Join(dir, "changes.patch")
	if err := os.WriteFile(patchFile, []byte("diff content\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	validFile := filepath.Join(dir, "valid.txt")
	if err := os.WriteFile(validFile, []byte("valid content\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("single explicit excluded file emits warning", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.NoColor = true
		cfg.PassedFiles = []string{patchFile}
		buf := new(bytes.Buffer)
		cfg.Logger.SetWriter(buf)

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files, got %v", files)
		}
		output := buf.String()
		expectedWarning := "skipped: " + patchFile + " (excluded by pattern \\.patch$)"
		if !strings.Contains(output, expectedWarning) {
			t.Errorf("expected warning %q in output, got: %q", expectedWarning, output)
		}
	})

	t.Run("multiple explicit files only includes non-excluded", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.NoColor = true
		cfg.PassedFiles = []string{patchFile, validFile}
		buf := new(bytes.Buffer)
		cfg.Logger.SetWriter(buf)

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 1 || files[0] != validFile {
			t.Errorf("expected [%s], got %v", validFile, files)
		}
		output := buf.String()
		expectedWarning := "skipped: " + patchFile + " (excluded by pattern \\.patch$)"
		if !strings.Contains(output, expectedWarning) {
			t.Errorf("expected warning %q in output, got: %q", expectedWarning, output)
		}
	})

	t.Run("duplicate explicit excluded file warns only once", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.NoColor = true
		cfg.PassedFiles = []string{patchFile, patchFile}
		buf := new(bytes.Buffer)
		cfg.Logger.SetWriter(buf)

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files, got %v", files)
		}
		output := buf.String()
		expectedWarning := "skipped: " + patchFile + " (excluded by pattern \\.patch$)"
		if strings.Count(output, expectedWarning) != 1 {
			t.Errorf("expected warning exactly once, got count %d in %q", strings.Count(output, expectedWarning), output)
		}
	})

	t.Run("custom exclude pattern warning", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.NoColor = true
		cfg.Exclude = []string{`valid\.txt$`}
		cfg.PassedFiles = []string{validFile}
		buf := new(bytes.Buffer)
		cfg.Logger.SetWriter(buf)

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files, got %v", files)
		}
		output := buf.String()
		expectedWarning := "skipped: " + validFile + " (excluded by pattern valid\\.txt$)"
		if !strings.Contains(output, expectedWarning) {
			t.Errorf("expected warning %q in output, got: %q", expectedWarning, output)
		}
	})

	t.Run("empty matching pattern warning", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.NoColor = true
		cfg.IgnoreDefaults = true
		cfg.Exclude = []string{"(", ")"}
		cfg.PassedFiles = []string{validFile}
		buf := new(bytes.Buffer)
		cfg.Logger.SetWriter(buf)

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files, got %v", files)
		}
		output := buf.String()
		expectedWarning := "skipped: " + validFile + " (excluded)"
		if !strings.Contains(output, expectedWarning) {
			t.Errorf("expected warning %q in output, got: %q", expectedWarning, output)
		}
	})

	t.Run("invalid exclude pattern returns error", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.Exclude = []string{"["}
		cfg.PassedFiles = []string{validFile}
		_, err := GetFiles(*cfg)
		if err == nil {
			t.Errorf("expected error when exclude regex is invalid")
		}
	})

	t.Run("nil logger does not panic when skipping excluded file", func(t *testing.T) {
		cfg := config.NewConfig(nil)
		cfg.Logger = nil
		cfg.PassedFiles = []string{patchFile}

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files, got %v", files)
		}
	})

	t.Run("relative explicit excluded file emits warning", func(t *testing.T) {
		tempDir := t.TempDir()
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		relDir, err := filepath.Rel(cwd, tempDir)
		if err != nil {
			t.Skip("cannot form relative path across drives")
		}

		relPatch := filepath.Join(relDir, "relative.patch")
		if err := os.WriteFile(filepath.Join(tempDir, "relative.patch"), []byte("diff\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg := config.NewConfig(nil)
		cfg.NoColor = true
		cfg.PassedFiles = []string{relPatch}
		buf := new(bytes.Buffer)
		cfg.Logger.SetWriter(buf)

		files, err := GetFiles(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files, got %v", files)
		}
		output := buf.String()
		expectedWarning := "skipped: " + relPatch + " (excluded by pattern \\.patch$)"
		if !strings.Contains(output, expectedWarning) {
			t.Errorf("expected warning %q in output, got: %q", expectedWarning, output)
		}
	})
}

func setup() {
	const testResultsJSON = "../encoding/test-results.json"

	f, err := os.Open(testResultsJSON)
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

func TestReadLines(t *testing.T) {
	// longer than the 64 KiB token limit of a default bufio.Scanner
	longLine := strings.Repeat("a", 128*1024)

	readLinesTests := []struct {
		content  string
		expected []string
	}{
		{"", nil},
		{"a\nb\n", []string{"a", "b"}},
		{"a\r\nb\r\n", []string{"a", "b"}},
		{"a\nb", []string{"a", "b"}},
		{"a\n" + longLine + "\nb \n", []string{"a", longLine, "b "}},
	}

	for _, tt := range readLinesTests {
		actual := ReadLines(tt.content)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("ReadLines(%.20q): expected %.20q, got %.20q", tt.content, tt.expected, actual)
		}
	}
}
