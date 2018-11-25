// Package validators provides ...
package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker.go/utils"
)

// Space validates if a file is indented correctly if indentStyle is set to "space"
func Space(line string, indentStyle string, indentSize int) bool {
	if indentStyle == "space" && len(line) > 0 && indentSize > 0 {
		// match recurring spaces indentSize times - this can be recurring or never
		// match either a space followed by a * and maybe a space (block-comments)
		// or match everything despite a space or tab-character
		regexpPattern := fmt.Sprintf("^( {%d})*( \\* ?|[^ \t])", indentSize)

		matched, err := regexp.MatchString(regexpPattern, line)

		if err != nil {
			panic(err)
		}

		if matched {
			return true
		}

		return false
	}

	return true
}

// Tab validates if a file is indented correctly if indentStyle is set to "space"
func Tab(line string, indentStyle string) bool {
	if indentStyle == "tab" && len(line) > 0 {
		regexpPattern := "^\t*[^ \t]"
		matched, err := regexp.MatchString(regexpPattern, line)

		if err != nil {
			panic(err)
		}

		if matched {
			return true
		}

		return false
	}

	return true
}

// TrailingWhitespace validates if a file has trailing whitespace if the trimTrailingWhitespace parameter is true
func TrailingWhitespace(line string, trimTrailingWhitespace bool) bool {
	if trimTrailingWhitespace {
		regexpPattern := "^.*[ \t]+$"
		matched, err := regexp.MatchString(regexpPattern, line)

		if err != nil {
			panic(err)
		}

		if matched {
			return false
		}

		return true
	}

	return true
}

// FinalNewline validates if a file has a final newline if finalNewline is set to true
func FinalNewline(fileContent string, insertFinalNewline bool, endOfLine string) bool {
	if insertFinalNewline {
		regexpPattern := fmt.Sprintf("%s$", utils.GetEolChar(endOfLine))
		matched, err := regexp.MatchString(regexpPattern, fileContent)

		if err != nil {
			panic(err)
		}

		if matched {
			return true
		}

		return false
	}
	return true
}

// LineEnding validates if a file uses the correct line endings
func LineEnding(fileContent string, endOfLine string) bool {
	if endOfLine != "" {
		expectedEolChar := utils.GetEolChar(endOfLine)
		expectedEols := len(strings.Split(fileContent, expectedEolChar))
		lfEols := len(strings.Split(fileContent, "\n"))
		crEols := len(strings.Split(fileContent, "\r"))
		crlfEols := len(strings.Split(fileContent, "\r\n"))

		switch endOfLine {
		case "lf":
			return expectedEols == lfEols && crEols == 1 && crlfEols == 1
		case "cr":
			return expectedEols == crEols && lfEols == 1 && crlfEols == 1
		case "crlf":
			// A bit hacky because \r\n matches \r and \n
			return expectedEols == crlfEols && lfEols == expectedEols && crEols == expectedEols
		}
	}

	return true
}
