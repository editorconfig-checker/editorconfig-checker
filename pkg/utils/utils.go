// Package utils provides functions that are considered utility for being project independent
package utils

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// IsDirectory returns wether a path is a directory or not
func IsDirectory(path string) bool {
	fi, _ := os.Stat(path)
	return fi.Mode().IsDir()
}

// IsRegularFile return wether a file is a regular file or not
func IsRegularFile(fi os.FileInfo) bool {
	return fi.Mode().IsRegular()
}

// PathExists checks wether a path of a file or directory exists or not
func PathExists(filePath string) bool {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(absolutePath)

	if err == nil {
		return true
	}

	return false
}

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

func GetRelativePath(filePath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	relativePath := strings.Replace(filePath, cwd, ".", 1)
	return relativePath
}
