// Package main provides ...
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/editorconfig/editorconfig-core-go.v1"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/utils"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/validators"
)

// version
const version string = "1.0.0"

// Global variable to store the cli parameter
// only the init function should write to this variable
var params types.Params

// Init function, runs on start automagically
func init() {
	// define flags
	flag.StringVar(&params.Excludes, "exclude", "", "a regex which files should be excluded from checking - needs to be a valid regular expression")
	flag.StringVar(&params.Excludes, "e", "", "a regex which files should be excluded from checking - needs to be a valid regular expression")

	flag.BoolVar(&params.IgnoreDefaults, "ignore", false, "ignore default excludes")
	flag.BoolVar(&params.IgnoreDefaults, "i", false, "ignore default excludes")

	flag.BoolVar(&params.Version, "version", false, "print the version number")

	flag.BoolVar(&params.Help, "help", false, "print the help")
	flag.BoolVar(&params.Help, "h", false, "print the help")

	flag.BoolVar(&params.Verbose, "verbose", false, "print debugging information")
	flag.BoolVar(&params.Verbose, "v", false, "print debugging information")

	// parse flags
	flag.Parse()

	// get remaining args as rawFiles
	var rawFiles = flag.Args()
	if len(rawFiles) == 0 {
		// set rawFiles to . (current working directory) if no parameters are passed
		rawFiles = []string{"."}
	}

	excludes := ""

	if !params.IgnoreDefaults {
		excludes = utils.DefaultExcludes
	}

	if utils.PathExists(".ecrc") {
		lines := readLineNumbersOfFile(".ecrc")
		if len(lines) > 0 {
			if excludes != "" {
				excludes = fmt.Sprintf("%s|%s", excludes, strings.Join(lines, "|"))
			} else {
				excludes = fmt.Sprintf("%s", strings.Join(lines, "|"))
			}

		}
	}

	if params.Excludes != "" {
		if excludes != "" {
			excludes = fmt.Sprintf("%s|%s", excludes, params.Excludes)
		} else {
			excludes = fmt.Sprintf("%s", params.Excludes)

		}
	}

	params.Excludes = excludes
	params.RawFiles = rawFiles
}

// Returns wether the file is inside an unwanted folder
func isExcluded(filePath string) bool {
	if params.Excludes == "" {
		return false
	}

	relativeFilePath, err := utils.GetRelativePath(filePath)
	if err != nil {
		panic(err)
	}

	result, err := regexp.MatchString(params.Excludes, relativeFilePath)
	if err != nil {
		panic(err)
	}

	return result
}

// Adds a file to a slice if it isn't already in there
// and returns the new slice
func addToFiles(files []string, filePath string, verbose bool) []string {
	contentType, err := utils.GetContentType(filePath)

	if err != nil {
		logger.Error(fmt.Sprintf("Could not get the ContentType of file: %s", filePath))
		logger.Error(err.Error())
	}

	if !isExcluded(filePath) && (contentType == "application/octet-stream" || strings.Contains(contentType, "text/plain")) {
		if verbose {
			logger.Output(fmt.Sprintf("Add %s to be checked", filePath))
		}
		return append(files, filePath)
	}

	if verbose {
		logger.Output(fmt.Sprintf("Don't add %s to be checked", filePath))
	}

	return files
}

// Returns all files which should be checked
func getFiles(verbose bool) []string {
	var files []string

	byteArray, err := exec.Command("git", "ls-tree", "-r", "--name-only", "HEAD").Output()
	if err != nil {
		// It is not a git repository.
		cwd, err := os.Getwd()
		if err != nil {
			panic("Could not get the current working directly")
		}

		_ = filepath.Walk(cwd, func(path string, fi os.FileInfo, err error) error {
			if fi.Mode().IsRegular() {
				files = addToFiles(files, path, verbose)
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
				files = addToFiles(files, filePath, verbose)
			}
		}
	}

	return files
}

func readLineNumbersOfFile(filePath string) []string {
	var lines []string
	fileHandle, _ := os.Open(filePath)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
}

// Validates a single file and returns the errors
func validateFile(filePath string, verbose bool) []types.ValidationError {
	var errors []types.ValidationError
	lines := readLineNumbersOfFile(filePath)

	// return if first line contains editorconfig-checker-disable-file
	if len(lines) == 0 || strings.Contains(lines[0], "editorconfig-checker-disable-file") {
		return errors
	}

	rawFileContent, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	fileContent := string(rawFileContent)

	editorconfig, err := editorconfig.GetDefinitionForFilename(filePath)
	if err != nil {
		panic(err)
	}

	if currentError := validators.FinalNewline(
		fileContent,
		editorconfig.Raw["insert_final_newline"] == "true",
		editorconfig.Raw["end_of_line"]); currentError != nil {
		if verbose {
			logger.Output(fmt.Sprintf("Final newline error found in %s", filePath))
		}
		errors = append(errors, types.ValidationError{LineNumber: -1, Message: currentError})
	}

	if currentError := validators.LineEnding(
		fileContent,
		editorconfig.Raw["end_of_line"]); currentError != nil {
		if verbose {
			logger.Output(fmt.Sprintf("Line ending error found in %s", filePath))
		}
		errors = append(errors, types.ValidationError{LineNumber: -1, Message: currentError})
	}

	for lineNumber, line := range lines {
		if strings.Contains(line, "editorconfig-checker-disable-line") {
			continue
		}

		if currentError := validators.TrailingWhitespace(
			line,
			editorconfig.Raw["trim_trailing_whitespace"] == "true"); currentError != nil {
			if verbose {
				logger.Output(fmt.Sprintf("Trailing whitespace error found in %s on line %d", filePath, lineNumber))
			}
			errors = append(errors, types.ValidationError{LineNumber: lineNumber + 1, Message: currentError})
		}

		var indentSize int
		indentSize, err = strconv.Atoi(editorconfig.Raw["indent_size"])

		// Set indentSize to zero if there is no indentSize set
		if err != nil {
			indentSize = 0
		}

		if currentError := validators.Indentation(
			line,
			editorconfig.Raw["indent_style"],
			indentSize); currentError != nil {
			if verbose {
				logger.Output(fmt.Sprintf("Indentation error found in %s on line %d", filePath, lineNumber))
			}
			errors = append(errors, types.ValidationError{LineNumber: lineNumber + 1, Message: currentError})
		}
	}

	return errors
}

// Validates all files and returns an array of validation errors
func processValidation(files []string, verbose bool) []types.ValidationErrors {
	var validationErrors []types.ValidationErrors

	for _, filePath := range files {
		if verbose {
			logger.Output(fmt.Sprintf("Validate %s", filePath))
		}
		validationErrors = append(validationErrors, types.ValidationErrors{FilePath: filePath, Errors: validateFile(filePath, verbose)})
	}

	return validationErrors
}

func getErrorCount(errors []types.ValidationErrors) int {
	var errorCount = 0

	for _, v := range errors {
		errorCount += len(v.Errors)
	}

	return errorCount
}

func printErrors(errors []types.ValidationErrors) {
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) > 0 {
			relativeFilePath, err := utils.GetRelativePath(fileErrors.FilePath)

			if err != nil {
				logger.Error(err.Error())
			}

			logger.Print(fmt.Sprintf("%s:", relativeFilePath), logger.YELLOW, os.Stderr)
			for _, singleError := range fileErrors.Errors {
				if singleError.LineNumber != -1 {
					logger.Error(fmt.Sprintf("\t%d: %s", singleError.LineNumber, singleError.Message))
				} else {
					logger.Error(fmt.Sprintf("\t%s", singleError.Message))
				}

			}
		}
	}
}

// Main function, dude
func main() {
	// Check for returnworthy params
	switch {
	case params.Version:
		logger.Output(version)
		return
	case params.Help:
		logger.Output("USAGE:")
		flag.PrintDefaults()
		return
	}

	if params.Verbose {
		logger.Output(fmt.Sprintf("Exclude Regexp: %s", params.Excludes))
	}

	// contains all files which should be checked
	files := getFiles(params.Verbose)
	errors := processValidation(files, params.Verbose)
	errorCount := getErrorCount(errors)

	if errorCount != 0 {
		printErrors(errors)
		logger.Error(fmt.Sprintf("\n%d errors found", errorCount))
	}

	if params.Verbose {
		logger.Output(fmt.Sprintf("%d files checked", len(files)))
	}

	if errorCount != 0 {
		os.Exit(1)
	}

	os.Exit(0)
}
