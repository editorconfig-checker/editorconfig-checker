// Package main provides ...
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/editorconfig-checker/editorconfig-checker.go/src/validators"
	"gopkg.in/editorconfig/editorconfig-core-go.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// version
const version string = "0.0.1"

// Params is a Struct which represents the cli-params
type Params struct {
	version bool
	help    bool
	verbose bool
	// directories and/or files which should be validated
	rawFiles []string
}

// ValidationError represents one validation error
type ValidationError struct {
	line        int
	description string
}

// ValidationErrors represents which errors occurred in a file
type ValidationErrors struct {
	file   string
	errors []ValidationError
}

// Global variable to store the cli parameter
// only the init function should write to this variable
var params Params

// Init function, runs on start automagically
func init() {
	// define flags
	flag.BoolVar(&params.version, "version", false, "print the version number")
	flag.BoolVar(&params.version, "v", false, "print the version number")

	flag.BoolVar(&params.help, "help", false, "print the help")
	flag.BoolVar(&params.help, "h", false, "print the help")

	flag.BoolVar(&params.verbose, "verbose", false, "print debugging information")

	// parse flags
	flag.Parse()

	// get remaining args as rawFiles
	var rawFiles = flag.Args()
	if len(rawFiles) == 0 {
		// set rawFiles to . (current working directory) if no parameters are passed
		rawFiles = []string{"."}
	}

	params.rawFiles = rawFiles
}

// Returns wether a path is a directory or not
func isDirectory(path string) bool {
	fi, _ := os.Stat(path)
	return fi.Mode().IsDir()
}

// Checks wether a path of a file or directory exists or not
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Returns wether a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, sliceElement := range slice {
		if element == sliceElement {
			return true
		}
	}
	return false
}

// Checks if the file is ignored by the gitignore
// TODO: Remove dependency to git?
// TODO: In this state the application has to be run out of the repository root or some sub-folder
func isIgnoredByGitignore(file string) bool {
	cmd := exec.Command("git", "check-ignore", file)
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

// Returns wether the file is inside an unwanted folder
// TODO: This is only here for performance for now
// TODO: At least make this configurable i.e. .ecrc/.editorconfig-checkerrc
func isInDefaultExcludes(file string) bool {
	return strings.Contains(filepath.ToSlash(file), ".git/") ||
		strings.Contains(filepath.ToSlash(file), "node_modules/") ||
		strings.Contains(filepath.ToSlash(file), "vendor") ||
		strings.HasSuffix(file, ".png")
}

// Adds a file to a slice if it isn't already in there
// and returns the new slice
func addToFiles(files []string, file string) []string {
	if !contains(files, file) && !isInDefaultExcludes(file) && !isIgnoredByGitignore(file) {
		return append(files, file)
	}

	return files
}

// Returns all files which should be checked
// TODO: Manual excludes
func getFiles() []string {
	var files []string

	// loop over rawFiles to make them absolute
	// and check if they exist
	for _, value := range params.rawFiles {
		absolutePath, err := filepath.Abs(value)
		if err != nil {
			panic(err)
		}
		pathExist, err := pathExists(absolutePath)

		if !pathExist {
			panic("The directory or file `" + absolutePath + "` does not exist or is not accessible.")
		}

		if err != nil {
			panic(err)
		}

		// if the path is an directory walk through it and add all files to files slice
		if isDirectory(absolutePath) {
			// TODO: Performance optimization - this is the bottleneck and loops over every folder/file
			// and then checks if should be added. This needs some refactoring.
			err := filepath.Walk(absolutePath, func(path string, f os.FileInfo, err error) error {
				if !f.IsDir() {
					files = addToFiles(files, path)
				}
				return nil
			})

			if err != nil {
				panic(err)
			}

			continue
		}

		// just add the absolutePath to files
		files = addToFiles(files, absolutePath)
	}

	return files
}

func readLinesOfFile(file string) []string {
	var lines []string
	fileHandle, _ := os.Open(file)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
}

// Validates a single file and returns the errors
func validateFile(file string) []ValidationError {
	var errors []ValidationError
	lines := readLinesOfFile(file)
	rawFileContent, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	fileContent := string(rawFileContent)

	editorconfig, err := editorconfig.GetDefinitionForFilename(file)
	if err != nil {
		panic(err)
	}

	if !validators.FinalNewline(fileContent, editorconfig.Raw["insert_final_newline"] == "true", editorconfig.Raw["end_of_line"]) {
		errors = append(errors, ValidationError{-1, "TRAILING WHITESPACE VALIDATOR FAILED"})
	}

	if !validators.LineEnding(fileContent, editorconfig.Raw["end_of_line"]) {
		errors = append(errors, ValidationError{-1, "LINE ENDING VALIDATOR FAILED"})
	}

	for lineNumber, line := range lines {
		if !validators.TrailingWhitespace(line, editorconfig.Raw["trim_trailing_whitespace"] == "true") {
			errors = append(errors, ValidationError{lineNumber, "TRAILING WHITESPACE VALIDATOR FAILED"})
		}

		indentSize, err := strconv.Atoi(editorconfig.Raw["indent_size"])
		if err != nil {
			panic(err)
		}

		if !validators.Space(line, editorconfig.Raw["indent_style"], indentSize) {
			errors = append(errors, ValidationError{lineNumber, "SPACES VALIDATOR FAILED"})
		}

		if !validators.Tab(line, editorconfig.Raw["indent_style"]) {
			errors = append(errors, ValidationError{lineNumber, "SPACES VALIDATOR FAILED"})
		}
	}

	return errors
}

// Validates all files and returns an array of validation errors
func validateFiles(files []string) []ValidationErrors {
	var validationErrors []ValidationErrors

	for _, file := range files {
		validationErrors = append(validationErrors, ValidationErrors{file, validateFile(file)})
	}

	return validationErrors
}

func getErrorCount(errors []ValidationErrors) int {
	var errorCount = 0

	for _, v := range errors {
		errorCount += len(v.errors)
	}

	return errorCount
}

func printErrors(errors []ValidationErrors) {
	for _, v := range errors {
		fmt.Println(v.file)
		for _, errorr := range v.errors {
			if errorr.line != -1 {
				fmt.Printf("\t%d: ", errorr.line)
			}

			fmt.Printf("%s\n", errorr.description)
		}
	}
}

// Main function, dude
func main() {
	// Check for returnworthy params
	switch {
	case params.version:
		fmt.Println(version)
		return
	case params.help:
		fmt.Println("USAGE:")
		flag.PrintDefaults()
		return
	}

	// contains all files which should be checked
	files := getFiles()
	errors := validateFiles(files)
	errorCount := getErrorCount(errors)

	fmt.Printf("%d files found!\n", len(files))

	if errorCount != 0 {
		printErrors(errors)
		fmt.Printf("%d errors found\n", errorCount)
		os.Exit(1)
	}

	os.Exit(0)

	fmt.Println("Run Forrest, run!!!")
}
