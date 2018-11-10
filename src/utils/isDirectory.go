// Package utils provides ...
package utils

import (
	"os"
)

// IsDirectory returns wether a path is a directory or not
func IsDirectory(path string) bool {
	fi, _ := os.Stat(path)
	return fi.Mode().IsDir()
}
