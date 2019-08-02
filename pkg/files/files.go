package files

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
)

// Returns wether the file is inside an unwanted folder
func IsExcluded(filePath string, params types.Params) bool {
	if params.Excludes == "" {
		return false
	}

	relativeFilePath, err := GetRelativePath(filePath)
	if err != nil {
		panic(err)
	}

	result, err := regexp.MatchString(params.Excludes, relativeFilePath)
	if err != nil {
		panic(err)
	}

	return result
}

// Adds a file to a slice if it isn't already in there and meets the requirements
// and returns the new slice
func AddToFiles(filePaths []string, filePath string, params types.Params) []string {
	contentType, err := GetContentType(filePath)

	if err != nil {
		logger.Error(fmt.Sprintf("Could not get the ContentType of file: %s", filePath))
		logger.Error(err.Error())
	}

	if !IsExcluded(filePath, params) && IsAllowedContentType(contentType) {
		if params.Verbose {
			logger.Output(fmt.Sprintf("Add %s to be checked", filePath))
		}
		return append(filePaths, filePath)
	}

	if params.Verbose {
		logger.Output(fmt.Sprintf("Don't add %s to be checked", filePath))
	}

	return filePaths
}

// Returns all files which should be checked
func GetFiles(params types.Params) []string {
	var filePaths []string

	byteArray, err := exec.Command("git", "ls-tree", "-r", "--name-only", "HEAD").Output()
	if err != nil {
		// It is not a git repository.
		cwd, err := os.Getwd()
		if err != nil {
			panic("Could not get the current working directly")
		}

		_ = filepath.Walk(cwd, func(path string, fi os.FileInfo, err error) error {
			if fi.Mode().IsRegular() {
				filePaths = AddToFiles(filePaths, path, params)
			}

			return nil
		})
	}

	filesSlice := strings.Split(string(byteArray[:]), "\n")

	for _, filePath := range filesSlice {
		if len(filePath) > 0 {
			fi, err := os.Stat(filePath)

			// The err would be a broken symlink for example,
			// so we want to program to continue but the file should not be checked
			if err == nil && fi.Mode().IsRegular() {
				filePaths = AddToFiles(filePaths, filePath, params)
			}
		}
	}

	return filePaths
}

// Returns the lines from a file as a slice
func ReadLineNumbers(filePath string) []string {
	var lines []string
	fileHandle, _ := os.Open(filePath)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
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
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s", err))
	}

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer), nil
}

// PathExists checks wether a path of a file or directory exists or not
func PathExists(filePath string) error {
	absolutePath, _ := filepath.Abs(filePath)
	_, err := os.Stat(absolutePath)

	return err
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

func IsAllowedContentType(contentType string) bool {
	return contentType == "application/octet-stream" || strings.Contains(contentType, "text/")
}
