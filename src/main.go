// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// version
const version string = "0.0.1"

// Params is a Struct which represents the cli-params
type Params struct {
	version bool
	help    bool
	// directories and/or files which should be validated
	rawFiles []string
}

// Global variable to store the cli parameter
// only the init function should write to this variable
var params Params

// Init function, runs on start automagically
func init() {
	// define flags
	flag.BoolVar(&params.version, "version", false, "print the number")
	flag.BoolVar(&params.version, "v", false, "print the version number")

	flag.BoolVar(&params.help, "help", false, "print the help")
	flag.BoolVar(&params.help, "h", false, "print the help")

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

// Converts a path to an absolute path, if already absolute it returns the original one
func makePathAbsolute(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	absolutePath, err := filepath.Abs(path)
	return absolutePath, err
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

// Adds a file to a slice if it isn't already in there
// and returns the new slice
func addToFiles(files []string, file string) []string {
	if !contains(files, file) {
		return append(files, file)
	}

	return files
}

// Returns all files which should be checked
// @TODO: filtering/gitignore
func getFiles() []string {
	var files []string

	// loop over rawFiles to make them absolute
	// and check if they exist
	for _, value := range params.rawFiles {
		absolutePath, err := makePathAbsolute(value)
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
			err := filepath.Walk(absolutePath, func(path string, f os.FileInfo, err error) error {
				if !isDirectory(path) {
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

func main() {
	// Check for returnworthy params
	switch {
	case params.version:
		fmt.Println(version)
		return
	case params.help:
		fmt.Println("Should print help!!!")
		flag.PrintDefaults()
		return
	}

	// contains all files which should be checked
	files := getFiles()

	// do real stuff aka validation
	fmt.Println(files)

	fmt.Println("Run Forrest, run!!!")
}
