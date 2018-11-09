// Package validators provides ...
package validators

import (
	"github.com/editorconfig-checker/editorconfig-checker.go/src/utils"
	"strings"
)

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
