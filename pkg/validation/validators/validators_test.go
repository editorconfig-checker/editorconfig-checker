package validators

import (
	"errors"
	"reflect"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config" // x-release-please-major
)

func TestFinalNewline(t *testing.T) {
	finalNewlineTests := []struct {
		line               string
		insertFinalNewline string
		lineEnding         string
		expected           error
	}{
		{"x\n", "true", "lf", nil},
		{"x\r", "true", "cr", nil},
		{"x\r\n", "true", "crlf", nil},

		{"x", "true", "lf", errors.New("Wrong line endings or no final newline")},
		{"x", "true", "cr", errors.New("Wrong line endings or no final newline")},
		{"x", "true", "crlf", errors.New("Wrong line endings or no final newline")},

		{"x\n", "true", "cr", errors.New("Wrong line endings or no final newline")},
		{"x\n", "true", "crlf", errors.New("Wrong line endings or no final newline")},

		{"x\r", "true", "lf", errors.New("Wrong line endings or no final newline")},
		{"x\r", "true", "crlf", errors.New("Wrong line endings or no final newline")},

		{"x\r\n", "true", "lf", errors.New("Wrong line endings or no final newline")},
		{"x\r\n", "true", "cr", errors.New("Wrong line endings or no final newline")},

		// insert_final_newline false
		{"x", "false", "lf", nil},
		{"x\n", "false", "lf", errors.New("No final newline expected")},
		{"x\r", "false", "lf", errors.New("No final newline expected")},
		{"x\r\n", "false", "lf", errors.New("No final newline expected")},

		{"x", "false", "cr", nil},
		{"x\n", "false", "cr", errors.New("No final newline expected")},
		{"x\r", "false", "cr", errors.New("No final newline expected")},
		{"x\r\n", "false", "cr", errors.New("No final newline expected")},

		{"x", "false", "crlf", nil},
		{"x\n", "false", "crlf", errors.New("No final newline expected")},
		{"x\r", "false", "crlf", errors.New("No final newline expected")},
		{"x\r\n", "false", "crlf", errors.New("No final newline expected")},

		// insert_final_newline not set
		{"x", "", "lf", nil},
		{"x", "", "cr", nil},
		{"x", "", "crlf", nil},
		{"x\n", "", "lf", nil},
		{"x\n", "", "cr", nil},
		{"x\n", "", "crlf", nil},
		{"x\r", "", "lf", nil},
		{"x\r", "", "cr", nil},
		{"x\r", "", "crlf", nil},
		{"x\r\n", "", "lf", nil},
		{"x\r\n", "", "cr", nil},
		{"x\r\n", "", "crlf", nil},

		// end_of_line not set
		{"x", "true", "", errors.New("Final newline expected")},
		{"x", "false", "", nil},
		{"x\n", "true", "", nil},
		{"x\n", "false", "", errors.New("No final newline expected")},
		{"x\r", "true", "", nil},
		{"x\r", "false", "", errors.New("No final newline expected")},
		{"x\r\n", "true", "", nil},
		{"x\r\n", "false", "", errors.New("No final newline expected")},
	}

	for _, tt := range finalNewlineTests {
		actual := FinalNewline(tt.line, tt.insertFinalNewline, tt.lineEnding)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("FinalNewline(%s, %s, %s): expected: %v, got: %v", tt.line, tt.insertFinalNewline, tt.lineEnding, tt.expected, actual)
		}
	}
}

func TestLineEnding(t *testing.T) {
	linedEndingTests := []struct {
		line       string
		lineEnding string
		expected   error
	}{
		{"x", "lf", nil},
		{"x\n", "lf", nil},
		{"x\r", "lf", errors.New("Not all lines have the correct end of line character")},
		{"x\r\n", "lf", errors.New("Not all lines have the correct end of line character")},
		{"x\ry\nz\n", "lf", errors.New("Not all lines have the correct end of line character")},

		{"x", "cr", nil},
		{"x\r", "cr", nil},
		{"x\n", "cr", errors.New("Not all lines have the correct end of line character")},
		{"x\r\n", "cr", errors.New("Not all lines have the correct end of line character")},
		{"x\ry\nz\n", "cr", errors.New("Not all lines have the correct end of line character")},

		{"x", "crlf", nil},
		{"x\r\n", "crlf", nil},
		{"x\r", "crlf", errors.New("Not all lines have the correct end of line character")},
		{"x\n", "crlf", errors.New("Not all lines have the correct end of line character")},
		{"x\ry\nz\n", "crlf", errors.New("Not all lines have the correct end of line character")},
	}

	for _, tt := range linedEndingTests {
		actual := LineEnding(tt.line, tt.lineEnding)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("LineEnding(%s, %s): expected: %v, got: %v", tt.line, tt.lineEnding, tt.expected, actual)
		}
	}
}

func TestIndentation(t *testing.T) {
	configuration := config.Config{SpacesAfterTabs: false}

	indentationTests := []struct {
		line        string
		indentStyle string
		indenSize   int
		expected    error
	}{
		{"    x", "space", 4, nil},
		{"   x", "space", 4, errors.New("Wrong amount of left-padding spaces(want multiple of 4)")},
		{"	x", "tab", 0, nil},
		{"   x", "tab", 0, errors.New("Wrong indentation type(spaces instead of tabs)")},
		{"	x", "x", 0, nil},
		{"   x", "x", 0, nil},
	}

	for _, tt := range indentationTests {
		actual := Indentation(tt.line, tt.indentStyle, tt.indenSize, configuration)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("Indentation(%s, %s, %d, %+v): expected: %v, got: %v", tt.line, tt.indentStyle, tt.indenSize, configuration, tt.expected, actual)
		}
	}
}

func TestSpace(t *testing.T) {
	enabledIndentSizeConfig := config.Config{}
	disabled := config.DisabledChecks{IndentSize: true}
	disabledIndentSizeConfig := config.Config{
		Disable: disabled,
	}

	spaceTests := []struct {
		line       string
		indentSize int
		config     config.Config
		expected   error
	}{
		{"", 4, enabledIndentSizeConfig, nil},
		{"x", 0, enabledIndentSizeConfig, nil},
		{"x", 4, enabledIndentSizeConfig, nil},
		{"    x", 4, enabledIndentSizeConfig, nil},
		// 5 spaces
		{"     x", 4, enabledIndentSizeConfig, errors.New("Wrong amount of left-padding spaces(want multiple of 4)")},
		// 3 spaces
		{"   x", 4, enabledIndentSizeConfig, errors.New("Wrong amount of left-padding spaces(want multiple of 4)")},
		// correct indented block comment, empty and non empty
		{"     *", 4, enabledIndentSizeConfig, nil},
		{"     * some comment", 4, enabledIndentSizeConfig, nil},
		{"    ", 4, enabledIndentSizeConfig, nil},
		// disabled indent size
		{"", 4, disabledIndentSizeConfig, nil},
		{" x ", 4, disabledIndentSizeConfig, nil},
		{"  x ", 4, disabledIndentSizeConfig, nil},
		{"    x ", 4, disabledIndentSizeConfig, nil},
		{"     x ", 4, disabledIndentSizeConfig, nil},
		{"	x ", 4, disabledIndentSizeConfig, errors.New("Wrong indent style found (tabs instead of spaces)")},
		{"    		x a", 4, disabledIndentSizeConfig, errors.New("Wrong indent style found (tabs instead of spaces)")},
	}

	for _, tt := range spaceTests {
		actual := Space(tt.line, tt.indentSize, tt.config)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("Space(%s, %d): expected: %v, got: %v", tt.line, tt.indentSize, tt.expected, actual)
		}
	}
}

func TestTab(t *testing.T) {
	spacesAllowed := config.Config{SpacesAfterTabs: true}
	spacesForbidden := config.Config{SpacesAfterTabs: false}
	tabTests := []struct {
		line     string
		config   config.Config
		expected error
	}{
		{" x", spacesAllowed, nil},
		{"	   bla", spacesAllowed, nil},
		{"	 bla", spacesAllowed, nil},
		{"		  xx", spacesAllowed, nil},

		{"", spacesForbidden, nil},
		{"x", spacesForbidden, nil},
		{"	x", spacesForbidden, nil},
		{"	", spacesForbidden, nil},
		{"		x", spacesForbidden, nil},
		{"  	a", spacesForbidden, errors.New("Wrong indentation type(spaces instead of tabs)")},
		{" *", spacesForbidden, nil},
		{"	 *", spacesForbidden, nil},
		{"	 * some comment", spacesForbidden, nil},
		{" */", spacesForbidden, nil},
		{"	 */", spacesForbidden, nil},
		{" *", spacesForbidden, nil},
	}

	for _, tt := range tabTests {
		actual := Tab(tt.line, tt.config)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("Tab(%s, %+v): expected: %v, got: %v", tt.line, tt.config, tt.expected, actual)
		}
	}
}

func TestTrailingWhitespace(t *testing.T) {
	trailingWhitespaceTests := []struct {
		line                   string
		trimTrailingWhitespace bool
		expected               error
	}{
		{"", true, nil},
		{"", false, nil},
		{"x", true, nil},
		{"x", false, nil},

		// Spaces
		{"x ", true, errors.New("Trailing whitespace")},
		{"x ", false, nil},
		{"x .", true, nil},
		{"x .", false, nil},

		// Tabs
		{"x	", true, errors.New("Trailing whitespace")},
		{"x	", false, nil},
		{"x	.", true, nil},
		{"x	.", false, nil},
	}

	for _, tt := range trailingWhitespaceTests {
		actual := TrailingWhitespace(tt.line, tt.trimTrailingWhitespace)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("TrailingWhitespace(%s, %v): expected: %v, got: %v", tt.line, tt.trimTrailingWhitespace, tt.expected, actual)
		}
	}
}

func TestMaxLineLength(t *testing.T) {
	maxLineLengthTest := []struct {
		line          string
		maxLineLength int
		charSet       string
		expected      error
	}{
		{"検索は次の", 5, "utf-8", nil},
		{"検索は次の", 2, "utf-8", errors.New("Line too long (5 instead of 2)")},
		{"\xEF\xBB\xBF検索は次の", 5, "utf-8-bom", nil},
		{"検索は次の", 5, "latin1", errors.New("Line too long (15 instead of 5)")},
		{"", 80, "latin1", nil},
		{"abc", 2, "latin1", errors.New("Line too long (3 instead of 2)")},
		{"   ", 2, "latin1", errors.New("Line too long (3 instead of 2)")},
		{"xx", 2, "latin1", nil},
	}

	for _, tt := range maxLineLengthTest {
		actual := MaxLineLength(tt.line, tt.maxLineLength, tt.charSet)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("Max(%s, %v): expected: %v, got: %v", tt.line, tt.maxLineLength, tt.expected, actual)
		}
	}
}
