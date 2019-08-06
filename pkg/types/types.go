// Package types provides types for the application
package types

import (
	"github.com/editorconfig/editorconfig-core-go/v2"
)

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

// DisabledChecks is a Struct which represents disabled checks
type DisabledChecks struct {
	TrailingWhitspace bool
	LineEnding        bool
	FinalNewline      bool
	Indentation       bool
}

// FileInformation is a Struct wich represents some FileInformation
type FileInformation struct {
	Line         string
	Content      string
	FilePath     string
	LineNumber   int
	Editorconfig *editorconfig.Definition
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
