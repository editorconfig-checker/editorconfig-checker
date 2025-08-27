// Package files contains functions and structs related to files
package files

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/editorconfig/editorconfig-core-go/v2"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/utils"
	// x-release-please-end
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

	re, err := config.CachedExcludesAsRegexp()
	if err != nil {
		return true, err
	}
	return re.MatchString(relativeFilePath), nil
}

// AddToFiles adds a file to a slice if it isn't already in there
// and meets the requirements and returns the new slice
func AddToFiles(filePaths []string, filePath string, config config.Config) []string {
	config.Logger.Debug("AddToFiles: investigating file %s", filePath)

	isExcluded, err := IsExcluded(filePath, config)
	if err == nil && isExcluded {
		config.Logger.Verbose("Not adding %s to be checked, it is excluded", filePath)
		return filePaths
	}

	contentType, err := GetContentType(filePath)
	if err != nil {
		config.Logger.Error("Could not get the ContentType of file: %s", filePath)
		config.Logger.Error("%v", err.Error())
	}
	config.Logger.Debug("AddToFiles: detected ContentType %s on file %s", contentType, filePath)

	if err == nil && !IsAllowedContentType(contentType, config) {
		config.Logger.Verbose("Not adding %s to be checked, it does not have an allowed ContentType", filePath)
		return filePaths
	}

	config.Logger.Verbose("Adding %s to be checked", filePath)
	return append(filePaths, filePath)
}

// GetFilesFromDirectory returns all files from a directory and its subdirectories which should be checked
func GetFilesFromDirectory(rootDir string, config config.Config) ([]string, error) {
	filePaths := make([]string, 0)
	err := fs.WalkDir(os.DirFS(rootDir), ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fi, err := de.Info()
		if err != nil {
			return err
		}

		fullPath := filepath.Join(rootDir, path)
		if fi.Mode().IsRegular() {
			filePaths = AddToFiles(filePaths, fullPath, config)
		} else if fi.IsDir() {
			if isExcluded, err := IsExcluded(fullPath, config); err == nil && isExcluded {
				config.Logger.Verbose("Not adding %s and subentries to be checked, it is excluded", fullPath)
				return fs.SkipDir
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory %s: %w", rootDir, err)
	}

	return filePaths, nil
}

// GetFiles returns all files which should be checked
func GetFiles(config config.Config) ([]string, error) {
	filePaths := make([]string, 0)

	// Handle explicit passed files
	if len(config.PassedFiles) != 0 {
		for _, passedFile := range config.PassedFiles {
			if !utils.IsDirectory(passedFile) {
				filePaths = AddToFiles(filePaths, passedFile, config)
			} else {
				files, err := GetFilesFromDirectory(passedFile, config)
				if err != nil {
					return filePaths, err
				}
				filePaths = append(filePaths, files...)
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

		return GetFilesFromDirectory(cwd, config)
	}

	filesSlice := strings.SplitSeq(string(byteArray[:]), "\n")

	for filePath := range filesSlice {
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
func GetContentType(path string) (string, error) {
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

	fileContent, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer fileContent.Close()

	return GetContentTypeBytes(fileContent)
}

// GetContentTypeBytes returns the content type of a byte slice
func GetContentTypeBytes(fileContent io.Reader) (string, error) {
	mimeType, err := mimetype.DetectReader(fileContent)
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
	/*
		why not use mimetype.EqualsAny:
		it would only match types exactly, but we allow our users to give an entire type/ category
	*/
	for _, allowedContentType := range config.AllowedContentTypes {
		if strings.Contains(contentType, allowedContentType) {
			return true
		}
	}

	return false
}
