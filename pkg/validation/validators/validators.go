// Package validators provides functions to validate if the rules of the `.editorconfig` are respected
package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config" // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/utils"  // x-release-please-major
)

// Indentation validates a files indentation
func Indentation(line string, indentStyle string, indentSize int, config config.Config) error {
	if indentStyle == "space" {
		return Space(line, indentSize, config)
	} else if indentStyle == "tab" {
		return Tab(line, config)
	}

	// if no indentStyle is given it should be valid
	return nil
}

// Space validates if a line is indented correctly respecting the indentSize
func Space(line string, indentSize int, config config.Config) error {
	if len(line) > 0 {
		// match recurring spaces and everything except tab characters
		regexpPattern := `^( )*([^ \t]|$)`
		matched, _ := regexp.MatchString(regexpPattern, line)

		if !matched {
			return fmt.Errorf("Wrong indent style found (tabs instead of spaces)")
		}

		if !config.Disable.IndentSize && indentSize > 0 {
			// match recurring spaces indentSize times - this can be recurring or never
			// match either a space followed by a * and maybe a space (block-comments)
			// or match everything despite a space or tab-character
			regexpPattern := fmt.Sprintf("^( {%d})*( \\* ?|[^ \t]|$)", indentSize)
			matched, _ := regexp.MatchString(regexpPattern, line)

			if !matched {
				return fmt.Errorf("Wrong amount of left-padding spaces(want multiple of %d)", indentSize)
			}
		}
	}

	return nil
}

// Tab validates if a line is indented with only tabs
func Tab(line string, config config.Config) error {
	if len(line) > 0 {
		// match starting with one or more tabs followed by a non-whitespace char
		// OR
		// match starting with one or more tabs, followed by one space and followed by at least one non-whitespace character
		// OR
		// match starting with a space followed by at least one non-whitespace character

		regexpPattern := "^(\t)*( \\* ?|[^ \t]|$)"

		if config.SpacesAfterTabs {
			regexpPattern = "(^(\t)*\\S)|(^(\t)+( )*\\S)|(^ \\S)"
		}

		matched, _ := regexp.MatchString(regexpPattern, line)

		if !matched {
			return errors.New("Wrong indentation type(spaces instead of tabs)")
		}

	}

	return nil
}

// TrailingWhitespace validates if a line has trailing whitespace
func TrailingWhitespace(line string, trimTrailingWhitespace bool) error {
	if trimTrailingWhitespace {
		regexpPattern := "^.*[ \t]+$"
		matched, _ := regexp.MatchString(regexpPattern, line)

		if matched {
			return errors.New("Trailing whitespace")
		}
	}

	return nil
}

// FinalNewline validates if a file has a final and correct newline
func FinalNewline(fileContent string, insertFinalNewline string, endOfLine string) error {
	if endOfLine != "" && insertFinalNewline == "true" {
		expectedEolChar := utils.GetEolChar(endOfLine)
		if !strings.HasSuffix(fileContent, expectedEolChar) || (expectedEolChar == "\n" && strings.HasSuffix(fileContent, "\r\n")) {
			return errors.New("Wrong line endings or no final newline")
		}
	} else {
		regexpPattern := "(\n|\r|\r\n)$"
		hasFinalNewline, _ := regexp.MatchString(regexpPattern, fileContent)

		if insertFinalNewline == "false" && hasFinalNewline {
			return errors.New("No final newline expected")
		}

		if insertFinalNewline == "true" && !hasFinalNewline {
			return errors.New("Final newline expected")
		}
	}

	return nil
}

// LineEnding validates if a file uses the correct line endings
func LineEnding(fileContent string, endOfLine string) error {
	if endOfLine != "" {
		expectedEolChar := utils.GetEolChar(endOfLine)
		expectedEols := len(strings.Split(fileContent, expectedEolChar))
		lfEols := len(strings.Split(fileContent, "\n"))
		crEols := len(strings.Split(fileContent, "\r"))
		crlfEols := len(strings.Split(fileContent, "\r\n"))

		switch endOfLine {
		case "lf":
			if !(expectedEols == lfEols && crEols == 1 && crlfEols == 1) {
				return errors.New("Not all lines have the correct end of line character")
			}
		case "cr":
			if !(expectedEols == crEols && lfEols == 1 && crlfEols == 1) {
				return errors.New("Not all lines have the correct end of line character")
			}
		case "crlf":
			// A bit hacky because \r\n matches \r and \n
			if !(expectedEols == crlfEols && lfEols == expectedEols && crEols == expectedEols) {
				return errors.New("Not all lines have the correct end of line character")
			}
		}
	}

	return nil
}

func MaxLineLength(line string, maxLineLength int, charSet string) error {
	var length int
	if charSet == "utf-8" || charSet == "utf-8-bom" {
		if charSet == "utf-8-bom" && strings.HasPrefix(line, "\xEF\xBB\xBF") {
			line = line[3:] // strip BOM
		}
		length = utf8.RuneCountInString(line)
	} else {
		// TODO: Handle utf-16be and utf-16le properly. Unfortunately, Go doesn't provide a utf16.RuneCountinString() function
		// Just go with byte count
		length = len(line)
	}

	if length > maxLineLength {
		return fmt.Errorf("Line too long (%d instead of %d)", length, maxLineLength)
	}

	return nil
}
