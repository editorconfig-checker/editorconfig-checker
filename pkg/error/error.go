// Package error contains functions and structs related to errors
package error

import (
	"encoding/json"
	"fmt"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"
)

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

// GetErrorCount returns the amount of errors
func GetErrorCount(errors []ValidationErrors) int {
	var errorCount = 0

	for _, v := range errors {
		errorCount += len(v.Errors)
	}

	return errorCount
}

// PrintErrors prints the errors to the console
func PrintErrors(errors []ValidationErrors, config config.Config) {
	codeclimateIssues := make([]CodeclimateIssue, 0)
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) > 0 {
			relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)

			if err != nil {
				config.Logger.Error(err.Error())
				continue
			}

			if config.Format == "gcc" {
				// gcc: A format mimicking the error format from GCC.
				for _, singleError := range fileErrors.Errors {
					var lineNo = 0
					if singleError.LineNumber > 0 {
						lineNo = singleError.LineNumber
					}
					config.Logger.Error(fmt.Sprintf("%s:%d:%d: %s: %s", relativeFilePath, lineNo, 0, "error", singleError.Message))
				}
			} else if config.Format == "codeclimate" {
				// codeclimate: A format that is compatible with the codeclimate format for GitLab CI.
				// https://docs.gitlab.com/ee/ci/testing/code_quality.html#implement-a-custom-tool
				for _, singleError := range fileErrors.Errors {
					codeclimateIssues = append(codeclimateIssues, newCodeclimateIssue(singleError, relativeFilePath))
				}
			} else {
				// default: A human readable text format.
				config.Logger.Warning(fmt.Sprintf("%s:", relativeFilePath))
				for _, singleError := range fileErrors.Errors {
					if singleError.LineNumber != -1 {
						config.Logger.Error(fmt.Sprintf("\t%d: %s", singleError.LineNumber, singleError.Message))
					} else {
						config.Logger.Error(fmt.Sprintf("\t%s", singleError.Message))
					}
				}
			}
		}
	}

	if len(codeclimateIssues) > 0 {
		// marshall codeclimate issues to json
		codeclimateIssuesJSON, err := json.Marshal(codeclimateIssues)
		if err != nil {
			config.Logger.Error("Error creating codeclimate json: %s", err.Error())
		} else {
			config.Logger.Output(string(codeclimateIssuesJSON))
		}
	}
}
