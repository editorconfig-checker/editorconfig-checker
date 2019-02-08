// Package utils provides functions that are considered utility for being project independent
package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	"\\.jepg$",
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

// PathExists checks wether a path of a file or directory exists or not
func PathExists(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	_, err := os.Stat(absolutePath)

	if err == nil {
		return true
	}

	return false
}

// GetContentType returns the content type of a file
func GetContentType(path string) (string, error) {
	fileStat, err := os.Stat(path)

	if err != nil {
		return "", err
	}

	if fileStat.Size() == 0 {
		return "", nil
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer), nil
}

// GetRelativePath returns the relative path of a file from the current working directory
func GetRelativePath(filePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Could not get the current working directly")
	}

	relativePath := strings.Replace(filePath, cwd, ".", 1)
	return relativePath, nil
}
