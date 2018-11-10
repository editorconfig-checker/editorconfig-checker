// Package utils provides ...
package utils

import (
	"net/http"
	"os"
)

// GetContentType returns the content type of a file
func GetContentType(path string) string {
	// TODO: Refactor this into somewhere else or return additionally an error
	fileStat, err := os.Stat(path)

	if err != nil {
		panic(err)
	}

	if fileStat.Size() == 0 {
		return ""
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		panic(err)
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer)
}
