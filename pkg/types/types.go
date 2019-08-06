// Package types provides types for the application
package types

// Params is a Struct which represents the cli-params
type Params struct {
	Version         bool
	Help            bool
	Verbose         bool
	IgnoreDefaults  bool
	SpacesAfterTabs bool
	DryRun          bool
	Excludes        string
	RawFiles        []string
	Disabled        DisabledChecks
}

type DisabledChecks struct {
	TrailingWhitspace bool
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
