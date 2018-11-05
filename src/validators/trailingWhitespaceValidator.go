package validators

import (
	"regexp"
)

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
