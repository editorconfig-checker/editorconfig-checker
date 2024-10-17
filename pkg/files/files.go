// Package files contains functions and structs related to files
package files

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/editorconfig/editorconfig-core-go/v2"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config" // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/utils"  // x-release-please-major
)

const DefaultMimeType = "application/octet-stream"

// FileInformation is a Struct which represents some FileInformation
type FileInformation struct {
	Line         string
	Content      string
	FilePath     string
	LineNumber   int
	Editorconfig *editorconfig.Definition
}

// IsExcluded returns whether the file is excluded via arguments or config file
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
func AddToFiles(filePaths []string, filePath string, config config.Config) []string {
	contentType, err := GetContentType(filePath, config)

	config.Logger.Debug("AddToFiles: filePath: %s, contentType: %s", filePath, contentType)

	if err != nil {
		config.Logger.Error("Could not get the ContentType of file: %s", filePath)
		config.Logger.Error(err.Error())
	}

	isExcluded, err := IsExcluded(filePath, config)

	if err == nil && !isExcluded && IsAllowedContentType(contentType, config) {
		config.Logger.Verbose("Add %s to be checked", filePath)
		return append(filePaths, filePath)
	}

	config.Logger.Verbose("Don't add %s to be checked", filePath)
	return filePaths
}

// GetFiles returns all files which should be checked
func GetFiles(config config.Config) ([]string, error) {
	filePaths := make([]string, 0)

	// Handle explicit passed files
	if len(config.PassedFiles) != 0 {
		for _, passedFile := range config.PassedFiles {
			if utils.IsDirectory(passedFile) {
				_ = fs.WalkDir(os.DirFS(passedFile), ".", func(path string, de fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					fi, err := de.Info()
					if err != nil {
						return err
					}

					if fi.Mode().IsRegular() {
						filePaths = AddToFiles(filePaths, filepath.Join(passedFile, path), config)
					}

					return nil
				})
			} else {
				filePaths = AddToFiles(filePaths, passedFile, config)
			}
		}

		return filePaths, nil
	}

	byteArray, err := exec.Command("git", "ls-files", "--cached", "--others", "--exclude-standard").Output()
	if err != nil {
		// It is not a git repository.
		cwd, err := os.Getwd()
		if err != nil {
			return filePaths, err
		}

		_ = fs.WalkDir(os.DirFS(cwd), ".", func(path string, de fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			fi, err := de.Info()
			if err != nil {
				return err
			}

			if fi.Mode().IsRegular() {
				filePaths = AddToFiles(filePaths, path, config)
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
				filePaths = AddToFiles(filePaths, filePath, config)
			}
		}
	}

	return filePaths, nil
}

// ReadLines returns the lines from a file as a slice
func ReadLines(content string) []string {
	var lines []string
	stringReader := strings.NewReader(content)
	fileScanner := bufio.NewScanner(stringReader)
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
}

// GetContentType returns the content type of a file
func GetContentType(path string, config config.Config) (string, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if fileStat.IsDir() {
		return "", fmt.Errorf("%s is a directory", path)
	}

	if fileStat.Size() == 0 {
		return "", nil
	}

	rawFileContent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return GetContentTypeBytes(rawFileContent, config)
}

// GetContentTypeBytes returns the content type of a byte slice
func GetContentTypeBytes(rawFileContent []byte, config config.Config) (string, error) {
	bytesReader := bytes.NewReader(rawFileContent)

	mimeType, err := mimetype.DetectReader(bytesReader)
	if err != nil {
		return "", err
	}

	mime := mimeType.String()
	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	parts := strings.Split(mime, ";")
	if len(parts) > 0 {
		mime = strings.TrimSpace(parts[0])
	}
	if mime == "" {
		return DefaultMimeType, nil
	}
	return mime, nil
}

// PathExists checks whether a path of a file or directory exists or not
func PathExists(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	_, err := os.Stat(absolutePath)

	return err == nil
}

// GetRelativePath returns the relative path of a file from the current working directory
func GetRelativePath(filePath string) (string, error) {
	filePath = filepath.FromSlash(filePath)
	if !filepath.IsAbs(filePath) {
		// Path is already relative. No changes needed
		return filepath.ToSlash(filePath), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Could not get the current working directory")
	}

	cwd = filepath.FromSlash(cwd)
	rel, err := filepath.Rel(cwd, filePath)
	return filepath.ToSlash(rel), err
}

// IsAllowedContentType returns whether the contentType is
// an allowed content type to check or not
func IsAllowedContentType(contentType string, config config.Config) bool {
	result := false

	for _, allowedContentType := range config.AllowedContentTypes {
		result = result || strings.Contains(contentType, allowedContentType)
	}

	return result
}
