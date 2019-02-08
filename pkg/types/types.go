// Package types provides types for the application
package types

// Params is a Struct which represents the cli-params
type Params struct {
	Version        bool
	Help           bool
	Verbose        bool
	IgnoreDefaults bool
	// directories and/or files which should be validated
	RawFiles []string
	Excludes string
}

// ValidationError represents one validation error
type ValidationError struct {
	LineNumber int
	Message    error
}

// ValidationErrors represents which errors occurred in a file
type ValidationErrors struct {
	FilePath string
	Errors   []ValidationError
}
