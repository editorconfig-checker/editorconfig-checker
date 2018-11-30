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
	"strconv"
	"strings"

	"gopkg.in/editorconfig/editorconfig-core-go.v1"

	"github.com/editorconfig-checker/editorconfig-checker.go/types"
	"github.com/editorconfig-checker/editorconfig-checker.go/utils"
	"github.com/editorconfig-checker/editorconfig-checker.go/validators"
)

// version
const version string = "0.0.1"

var defaultExcludes = [...]string{
	"/.git/",
	"/vendor/",
	"/yarn.lock",
	"/composer.lock",
	"node_modules/"}

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
// TODO: BETTER: elimante the need for this (better git filtering)
func isInDefaultExcludes(file string) bool {
	var result bool = false
	for _, exclude := range defaultExcludes {
		result = result || strings.Contains(filepath.ToSlash(file), exclude)
	}

	return result
}

// Adds a file to a slice if it isn't already in there
// and returns the new slice
func addFile(file string, addFileChan chan bool) {
	contentType := utils.GetContentType(file)

	addFileChan <- !isInDefaultExcludes(file) &&
		(contentType == "application/octet-stream" || strings.Contains(contentType, "text/plain")) &&
		!isIgnoredByGitignore(file)
}

// Returns all files which should be checked
// TODO: Manual excludes
func getFiles(getFileChan chan string) {
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
			err := filepath.Walk(absolutePath, func(path string, fi os.FileInfo, err error) error {
				if utils.IsRegularFile(fi) {
					addFileChan := make(chan bool)
					go addFile(path, addFileChan)
					if <-addFileChan {
						getFileChan <- path
					}
				}

				return nil
			})

			if err != nil {
				panic(err)
			}

			continue
		}

		// just add the absolutePath to files
		addFileChan := make(chan bool)
		go addFile(absolutePath, addFileChan)
		if <-addFileChan {
			getFileChan <- absolutePath
		}
	}

	close(getFileChan)
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
func validateFile(file string, errorChan chan types.ValidationError) {
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

	if currentError := validators.FinalNewline(fileContent, editorconfig.Raw["insert_final_newline"] == "true", editorconfig.Raw["end_of_line"]); currentError != nil {
		errorChan <- types.ValidationError{LineNumber: -1, Message: currentError}
	}

	if currentError := validators.LineEnding(fileContent, editorconfig.Raw["end_of_line"]); currentError != nil {
		errorChan <- types.ValidationError{LineNumber: -1, Message: currentError}
	}

	for lineNumber, line := range lines {
		if currentError := validators.TrailingWhitespace(line, editorconfig.Raw["trim_trailing_whitespace"] == "true"); currentError != nil {
			errorChan <- types.ValidationError{LineNumber: lineNumber + 1, Message: currentError}
		}

		var indentSize int
		indentSize, err = strconv.Atoi(editorconfig.Raw["indent_size"])

		// Set indentSize to zero if there is no indentSize set
		if err != nil {
			indentSize = 0
		}

		if currentError := validators.Indentation(line, editorconfig.Raw["indent_style"], indentSize); currentError != nil {
			errorChan <- types.ValidationError{LineNumber: lineNumber + 1, Message: currentError}
		}
	}

	close(errorChan)
}

func getErrorCount(errors []types.ValidationErrors, errorCountChan chan int) {
	for _, v := range errors {
		errorCountChan <- len(v.Errors)
	}

	close(errorCountChan)
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

	var errors []types.ValidationErrors
	var fileCount int = 0

	getFileChan := make(chan string)
	go getFiles(getFileChan)

	for {
		file, running := <-getFileChan

		if running == false {
			break
		}

		fileCount++

		errorChan := make(chan types.ValidationError)
		go validateFile(file, errorChan)
		var fileErrors []types.ValidationError

		for {
			singleError, running := <-errorChan

			if running == false {
				break
			}

			fileErrors = append(fileErrors, singleError)
		}

		errors = append(errors, types.ValidationErrors{FilePath: file, Errors: fileErrors})
	}

	errorCountChan := make(chan int)
	go getErrorCount(errors, errorCountChan)
	var errorCount int = 0

	for {
		errorCountPerFile, running := <-errorCountChan

		if running == false {
			break
		}

		errorCount += errorCountPerFile
	}

	if params.Verbose {
		fmt.Printf("%d files found!\n", fileCount)
	}

	if errorCount != 0 {
		printErrors(errors)
		fmt.Printf("\n%d errors found\n", errorCount)
		os.Exit(1)
	}

	os.Exit(0)
}
