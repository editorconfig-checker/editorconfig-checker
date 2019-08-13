// Package types provides types for the application
package types

import (
	"github.com/editorconfig/editorconfig-core-go/v2"
)

// FileInformation is a Struct which represents some FileInformation
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
