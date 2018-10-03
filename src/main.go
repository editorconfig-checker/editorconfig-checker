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
	// source directories and/or files which should be validated
	sources []string
}

// Global variable to store the cli parameter
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

	// get remaining args as source directories
	var rawSources = flag.Args()
	if len(rawSources) == 0 {
		// set sources to . (current working directory) if no parameters are passed
		rawSources = []string{"."}
	}
	// loop over sources to make them absolute
	// and check if they exist
	for _, value := range rawSources {
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

		params.sources = append(params.sources, absolutePath)
	}
}

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

func makePathAbsolute(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	absolutePath, err := filepath.Abs(path)
	return absolutePath, err
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

	// do real stuff
	fmt.Println(params.sources)

	fmt.Println("Run Forrest, run!!!")
}
