// Package validators provides ...
package validators

import (
	"fmt"
	"regexp"
)

// Space validates if a file is indented correctly if indentStyle is set to "space"
func Space(line string, indentStyle string, indentSize int) bool {
	if indentStyle == "space" && len(line) > 0 && indentSize > 0 {
		regexpPattern := fmt.Sprintf("^( {%d})*[^ \t]", indentSize)
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
