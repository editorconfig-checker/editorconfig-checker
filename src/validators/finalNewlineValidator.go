// Package validators provides ...
package validators

import (
	"fmt"
	"github.com/editorconfig-checker/editorconfig-checker.go/src/utils"
	"regexp"
)

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
