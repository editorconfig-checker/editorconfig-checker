// Package main provides ...
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/editorconfig-checker/editorconfig-checker.go/src/slice"
	"github.com/editorconfig-checker/editorconfig-checker.go/src/types"
	"github.com/editorconfig-checker/editorconfig-checker.go/src/utils"
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

// Global variable to store the cli parameter
// only the init function should write to this variable
var params types.Params

// Init function, runs on start automagically
func init() {
	// define flags
	flag.BoolVar(&params.Version, "version", false, "print the version number")
	flag.BoolVar(&params.Version, "v", false, "print the version number")

	flag.BoolVar(&params.Help, "help", false, "print the help")
	flag.BoolVar(&params.Help, "h", false, "print the help")

	flag.BoolVar(&params.Verbose, "verbose", false, "print debugging information")

	// parse flags
	flag.Parse()

	// get remaining args as rawFiles
	var rawFiles = flag.Args()
	if len(rawFiles) == 0 {
		// set rawFiles to . (current working directory) if no parameters are passed
		rawFiles = []string{"."}
	}

	params.RawFiles = rawFiles
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
		strings.Contains(filepath.ToSlash(file), "vendor")
}

// Adds a file to a slice if it isn't already in there
// and returns the new slice
func addToFiles(files []string, file string) []string {
	contentType := utils.GetContentType(file)
	if !slice.Contains(files, file) &&
		(contentType == "application/octet-stream" || strings.Contains(contentType, "text/plain")) &&
		!isInDefaultExcludes(file) &&
		!isIgnoredByGitignore(file) {
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
	for _, rawFile := range params.RawFiles {
		absolutePath, err := filepath.Abs(rawFile)

		if err != nil {
			panic(err)
		}

		pathExist := utils.PathExists(absolutePath)

		if !pathExist {
			panic("The directory or file `" + absolutePath + "` does not exist or is not accessible.")
		}

		// if the path is an directory walk through it and add all files to files slice
		if utils.IsDirectory(absolutePath) {
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

func readLineNumbersOfFile(file string) []string {
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
func validateFile(file string) []types.ValidationError {
	var errors []types.ValidationError
	lines := readLineNumbersOfFile(file)
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
		errors = append(errors, types.ValidationError{LineNumber: -1, Message: "FINAL NEWLINE VALIDATOR FAILED"})
	}

	if !validators.LineEnding(fileContent, editorconfig.Raw["end_of_line"]) {
		errors = append(errors, types.ValidationError{LineNumber: -1, Message: "LINE ENDING VALIDATOR FAILED"})
	}

	for lineNumber, line := range lines {
		if !validators.TrailingWhitespace(line, editorconfig.Raw["trim_trailing_whitespace"] == "true") {
			errors = append(errors, types.ValidationError{LineNumber: lineNumber, Message: "TRAILING WHITESPACE VALIDATOR FAILED"})
		}

		indentSize, err := strconv.Atoi(editorconfig.Raw["indent_size"])
		if err != nil {
			panic(err)
		}

		if !validators.Space(line, editorconfig.Raw["indent_style"], indentSize) {
			errors = append(errors, types.ValidationError{LineNumber: lineNumber, Message: "SPACES VALIDATOR FAILED"})
		}

		if !validators.Tab(line, editorconfig.Raw["indent_style"]) {
			errors = append(errors, types.ValidationError{LineNumber: lineNumber, Message: "SPACES VALIDATOR FAILED"})
		}
	}

	return errors
}

// Validates all files and returns an array of validation errors
func validateFiles(files []string) []types.ValidationErrors {
	var validationErrors []types.ValidationErrors

	for _, file := range files {
		validationErrors = append(validationErrors, types.ValidationErrors{FilePath: file, Errors: validateFile(file)})
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
	for _, file := range errors {
		if len(file.Errors) > 0 {
			fmt.Println(file.FilePath)
			for _, errorr := range file.Errors {
				fmt.Printf("\t")

				if errorr.LineNumber != -1 {
					fmt.Printf("%d: ", errorr.LineNumber)
				}

				fmt.Printf("%s\n", errorr.Message)
			}
		}
	}
}

// Main function, dude
func main() {
	// Check for returnworthy params
	switch {
	case params.Version:
		fmt.Println(version)
		return
	case params.Help:
		fmt.Println("USAGE:")
		flag.PrintDefaults()
		return
	}

	// contains all files which should be checked
	files := getFiles()
	errors := validateFiles(files)
	errorCount := getErrorCount(errors)

	if params.Verbose {
		fmt.Printf("%d files found!\n", len(files))
	}

	if errorCount != 0 {
		printErrors(errors)
		fmt.Printf("\n%d errors found\n", errorCount)
		os.Exit(1)
	}

	os.Exit(0)
}
