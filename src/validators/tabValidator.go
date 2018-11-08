// Package validators provides ...
package validators

import (
	"regexp"
)

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
