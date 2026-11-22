// Package utils provides functions that are considered utility for being project independent
package utils

import (
	"bytes"
	"os"
	"path/filepath"
)

// LineBreakLen reports how many bytes at i are a line break.
// A carriage return inside a file that also contains LF is line content, not a
// break, unless it starts a CRLF pair or is the final byte. CR-only files still
// break on every carriage return.
func LineBreakLen(content []byte, i int, fileHasLF bool) int {
	if i < 0 || i >= len(content) {
		return 0
	}
	switch content[i] {
	case '\n':
		return 1
	case '\r':
		if i+1 < len(content) && content[i+1] == '\n' {
			return 2
		}
		if fileHasLF && i+1 < len(content) {
			return 0
		}
		return 1
	default:
		return 0
	}
}

// FileHasLF reports whether content contains an LF byte.
func FileHasLF(content []byte) bool {
	return bytes.Contains(content, []byte{'\n'})
}

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

// IsRegularFile returns whether a file is a regular file or not
func IsRegularFile(filePath string) bool {
	absolutePath, firstErr := filepath.Abs(filePath)
	fi, secondErr := os.Stat(absolutePath)
	return firstErr == nil && secondErr == nil && fi.Mode().IsRegular()
}

// IsDirectory returns whether a file is a directory or not
func IsDirectory(filePath string) bool {
	absolutePath, firstErr := filepath.Abs(filePath)
	fi, secondErr := os.Stat(absolutePath)
	return firstErr == nil && secondErr == nil && fi.Mode().IsDir()
}

// FileExists returns whether a file exists or not
func FileExists(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	fi, secondErr := os.Stat(absolutePath)
	return fi != nil && secondErr == nil
}
