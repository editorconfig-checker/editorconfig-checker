// Package main provides ...
package main

import (
	"flag"
	"fmt"
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
	params.sources = flag.Args()
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
