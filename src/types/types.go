// Package types provides ...
package types

// Params is a Struct which represents the cli-params
type Params struct {
	Version bool
	Help    bool
	Verbose bool
	// directories and/or files which should be validated
	RawFiles []string
}

// ValidationError represents one validation error
type ValidationError struct {
	LineNumber int
	Message    string
}

// ValidationErrors represents which errors occurred in a file
type ValidationErrors struct {
	FilePath string
	Errors   []ValidationError
}
