// Package utils provides functions that are considered utility for being project independent
package utils

import (
	"os"
	"path/filepath"
)

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
func IsRegularFile(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	fi, err := os.Stat(absolutePath)

	return err == nil && fi.Mode().IsRegular()
}

func IsDirectory(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	fi, err := os.Stat(absolutePath)

	return err == nil && fi.Mode().IsDir()
}
