// Package files contains functions and structs related to files
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

	"github.com/editorconfig/editorconfig-core-go/v2"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/utils"
)

// FileInformation is a Struct which represents some FileInformation
type FileInformation struct {
	Line         string
	Content      string
	FilePath     string
	LineNumber   int
	Editorconfig *editorconfig.Definition
}

// IsExcluded returns wether the file is excluded via arguments or config file
func IsExcluded(filePath string, config config.Config) (bool, error) {
	if len(config.Exclude) == 0 && config.IgnoreDefaults {
		return false, nil
	}

	relativeFilePath, err := GetRelativePath(filePath)
	if err != nil {
		return true, err
	}

	result, err := regexp.MatchString(config.GetExcludesAsRegularExpression(), relativeFilePath)
	if err != nil {
		return true, err
	}

	return result, nil
}

// AddToFiles adds a file to a slice if it isn't already in there
// and meets the requirements and returns the new slice
func AddToFiles(filePaths []string, filePath string, config config.Config) ([]string, error) {
	contentType, err := GetContentType(filePath)

	config.Logger.Debug("AddToFiles: filePath: %s, contentType: %s", filePath, contentType)

	if err != nil {
		config.Logger.Error("Could not get the ContentType of file: %s", filePath)
		config.Logger.Error(err.Error())
		return filePaths, err
	}

	isExcluded, err := IsExcluded(filePath, config)

	if err == nil && !isExcluded && IsAllowedContentType(contentType, config) {
		config.Logger.Verbose("Add %s to be checked", filePath)
		return append(filePaths, filePath), err
	}

	config.Logger.Verbose("Don't add %s to be checked", filePath)
	return filePaths, err
}

// GetFiles returns all files which should be checked
func GetFiles(config config.Config) ([]string, error) {
	var filePaths []string
	var hasNonExistingFile error

	// Handle explicit passed files
	if len(config.PassedFiles) != 0 {
		for _, passedFile := range config.PassedFiles {
			if utils.IsDirectory(passedFile) {
				_ = filepath.Walk(passedFile, func(path string, fi os.FileInfo, err error) error {
					if fi.Mode().IsRegular() {
						filePaths, err = AddToFiles(filePaths, path, config)
						if hasNonExistingFile == nil {
							hasNonExistingFile = err
						}
					}

					return nil
				})
			} else {
				var err error
				filePaths, err = AddToFiles(filePaths, passedFile, config)
				if hasNonExistingFile == nil {
					hasNonExistingFile = err
				}
			}
		}

		return filePaths, hasNonExistingFile
	}

	byteArray, err := exec.Command("git", "ls-files", "--cached", "--others", "--modified", "--exclude-standard").Output()
	if err != nil {
		// It is not a git repository.
		cwd, err := os.Getwd()
		if err != nil {
			return filePaths, err
		}

		_ = filepath.Walk(cwd, func(path string, fi os.FileInfo, err error) error {
			if fi.Mode().IsRegular() {
				filePaths, err = AddToFiles(filePaths, path, config)
				if hasNonExistingFile == nil {
					hasNonExistingFile = err
				}
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
				filePaths, hasNonExistingFile = AddToFiles(filePaths, filePath, config)
			}
		}
	}

	return filePaths, hasNonExistingFile
}

// ReadLines returns the lines from a file as a slice
func ReadLines(filePath string) []string {
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
		panic(err)
	}

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer), nil
}

// PathExists checks wether a path of a file or directory exists or not
func PathExists(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	_, err := os.Stat(absolutePath)

	return err == nil
}

// GetRelativePath returns the relative path of a file from the current working directory
func GetRelativePath(filePath string) (string, error) {
	if !filepath.IsAbs(filePath) {
		// Path is already relative. No changes needed
		return filePath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Could not get the current working directory")
	}

	return filepath.Rel(cwd, filePath)
}

// IsAllowedContentType returns wether the contentType is
// an allowed content type to check or not
func IsAllowedContentType(contentType string, config config.Config) bool {
	result := false

	for _, allowedContentType := range config.AllowedContentTypes {
		result = result || strings.Contains(contentType, allowedContentType)
	}

	return result
}
