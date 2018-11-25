// Package utils provides ...
package utils

import (
	"os"
)

// PathExists checks wether a path of a file or directory exists or not
func PathExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	return false
}
