package validators

import (
	"regexp"
)

// Validates if a file has trailing whitespace if the trimTrailingWhitespace parameter is true
func TrailingWhitespace(line string, trimTrailingWhitespace bool) bool {
	if trimTrailingWhitespace {
		matched, err := regexp.MatchString("^.*[ \t]+$", line)
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
