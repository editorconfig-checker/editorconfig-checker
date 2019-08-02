// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/error"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/utils"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/validation"
)

// version
const version string = "1.2.1"

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

	flag.BoolVar(&params.DryRun, "dry-run", false, "show which files would be checked")
	flag.BoolVar(&params.DryRun, "d", false, "show which files would be checked")

	flag.BoolVar(&params.Version, "version", false, "print the version number")

	flag.BoolVar(&params.Help, "help", false, "print the help")
	flag.BoolVar(&params.Help, "h", false, "print the help")

	flag.BoolVar(&params.Verbose, "verbose", false, "print debugging information")
	flag.BoolVar(&params.Verbose, "v", false, "print debugging information")

	flag.BoolVar(&params.SpacesAfterTabs, "spaces-after-tabs", false, "allow spaces to be used as alignment after tabs")

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

	if files.PathExists(".ecrc") == nil {
		lines := files.ReadLineNumbers(".ecrc")
		if len(lines) > 0 {
			if excludes != "" {
				excludes = fmt.Sprintf("%s|%s", excludes, strings.Join(lines, "|"))
			} else {
				excludes = strings.Join(lines, "|")
			}

		}
	}

	if params.Excludes != "" {
		if excludes != "" {
			excludes = fmt.Sprintf("%s|%s", excludes, params.Excludes)
		} else {
			excludes = params.Excludes

		}
	}

	params.Excludes = excludes
	params.RawFiles = rawFiles
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
	filePaths := files.GetFiles(params)

	if params.DryRun {
		for _, file := range filePaths {
			logger.Output(file)
		}

		os.Exit(0)
	}

	errors := validation.ProcessValidation(filePaths, params)
	errorCount := error.GetErrorCount(errors)

	if errorCount != 0 {
		error.PrintErrors(errors)
		logger.Error(fmt.Sprintf("\n%d errors found", errorCount))
	}

	if params.Verbose {
		logger.Output(fmt.Sprintf("%d files checked", len(filePaths)))
	}

	if errorCount != 0 {
		os.Exit(1)
	}

	os.Exit(0)
}
