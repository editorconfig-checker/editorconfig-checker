// Package utils provides functions that are considered utility for being project independent
package utils

import (
	"strings"
)

// DefaultExcludes is the regular expression for ignored files
var DefaultExcludes = strings.Join(defaultExcludes, "|")

// defaultExcludes are an array to produce the correct string from
var defaultExcludes = []string{
	"yarn\\.lock$",
	"package-lock\\.json",
	"composer\\.lock$",
	"\\.snap$",
	"\\.otf$",
	"\\.woff$",
	"\\.woff2$",
	"\\.eot$",
	"\\.ttf$",
	"\\.gif$",
	"\\.png$",
	"\\.jpg$",
	"\\.jpeg$",
	"\\.mp4$",
	"\\.wmv$",
	"\\.svg$",
	"\\.ico$",
	"\\.bak$",
	"\\.bin$",
	"\\.pdf$",
	"\\.zip$",
	"\\.gz$",
	"\\.tar$",
	"\\.7z$",
	"\\.bz2$",
	"\\.log$",
	"\\.css\\.map$",
	"\\.js\\.map$",
	"min\\.css$",
	"min\\.js$"}

// GetEolChar returns the end of line character used in regexp
// depending on the end_of_line parameter
func GetEolChar(endOfLine string) string {
	switch endOfLine {
	case "lf":
		return "\n"

	case "cr":
		return "\r"
	case "crlf":
		return "\r\n"
	}

	return "\n"
}
