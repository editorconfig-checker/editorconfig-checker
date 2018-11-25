// Package utils provides ...
package utils

import (
	"os"
)

// IsRegularFile return wether a file is a regular file or not
func IsRegularFile(fi os.FileInfo) bool {
	return fi.Mode().IsRegular()
}
