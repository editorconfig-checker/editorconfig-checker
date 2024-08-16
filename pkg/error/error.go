// Package error contains functions and structs related to errors
package error

import (
	"fmt"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"
)

// ValidationError represents one validation error
type ValidationError struct {
	LineNumber       int
	Message          error
	ConsecutiveCount int
}

// ValidationErrors represents which errors occurred in a file
type ValidationErrors struct {
	FilePath string
	Errors   []ValidationError
}

func (error1 *ValidationError) Equal(error2 ValidationError) bool {
	return error1.Message.Error() == error2.Message.Error() &&
		error1.LineNumber == error2.LineNumber &&
		error1.ConsecutiveCount == error2.ConsecutiveCount

}

// GetErrorCount returns the amount of errors
func GetErrorCount(errors []ValidationErrors) int {
	var errorCount = 0

	for _, v := range errors {
		errorCount += len(v.Errors)
	}

	return errorCount
}

func ConsolidateErrors(errors []ValidationError, config config.Config) []ValidationError {
	var lineLessErrors []ValidationError
	var errorsWithLines []ValidationError

	// filter the errors, so we do not need to care about LineNumber == -1 in the loop below
	for _, singleError := range errors {
		if singleError.LineNumber == -1 {
			lineLessErrors = append(lineLessErrors, singleError)
		} else {
			errorsWithLines = append(errorsWithLines, singleError)
		}
	}

	config.Logger.Debug("sorted errors: %d with line number -1, %d with a line number", len(lineLessErrors), len(errorsWithLines))

	var consolidatedErrors []ValidationError

	// scan through the errors
	for i := 0; i < len(errorsWithLines); i++ {
		thisError := errorsWithLines[i]
		config.Logger.Debug("investigating error %d(%s)", i, thisError.Message)
		// scan through the errors after the current one
		for j, nextError := range errorsWithLines[i+1:] {
			config.Logger.Debug("comparing against error %d(%s)", j, nextError.Message)
			if nextError.Message.Error() == thisError.Message.Error() && nextError.LineNumber == thisError.LineNumber+thisError.ConsecutiveCount+1 {
				thisError.ConsecutiveCount++ // keep track of how many consecutive lines we've seen
				i = i + j + 1                // make sure the outer loop jumps over the consecutive errors we just found
			} else {
				break // if they are different errors we can stop comparing messages
			}
		}

		consolidatedErrors = append(consolidatedErrors, thisError)
	}

	return append(lineLessErrors, consolidatedErrors...)
}

// PrintErrors prints the errors to the console
func PrintErrors(errors []ValidationErrors, config config.Config) {
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) > 0 {
			relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)

			if err != nil {
				config.Logger.Error(err.Error())
				continue
			}

			// for these formats the errors need to be consolidated first.
			if config.Format == "default" || config.Format == "github-actions" {
				fileErrors.Errors = ConsolidateErrors(fileErrors.Errors, config)
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
			} else if config.Format == "github-actions" {
				// github-actions: A format dedicated for usage in Github Actions
				for _, singleError := range fileErrors.Errors {
					if singleError.LineNumber != -1 {
						if singleError.ConsecutiveCount != 0 {
							config.Logger.Error(fmt.Sprintf("::error file=%s,line=%d,endLine=%d::%s", relativeFilePath, singleError.LineNumber, singleError.LineNumber+singleError.ConsecutiveCount, singleError.Message))
						} else {
							config.Logger.Error(fmt.Sprintf("::error file=%s,line=%d::%s", relativeFilePath, singleError.LineNumber, singleError.Message))
						}
					} else {
						config.Logger.Error(fmt.Sprintf("::error file=%s::%s", relativeFilePath, singleError.Message))
					}
				}
			} else {
				// default: A human readable text format.
				config.Logger.Warning(fmt.Sprintf("%s:", relativeFilePath))
				for _, singleError := range fileErrors.Errors {
					if singleError.LineNumber != -1 {
						if singleError.ConsecutiveCount != 0 {
							config.Logger.Error(fmt.Sprintf("\t%d-%d: %s", singleError.LineNumber, singleError.LineNumber+singleError.ConsecutiveCount, singleError.Message))
						} else {
							config.Logger.Error(fmt.Sprintf("\t%d: %s", singleError.LineNumber, singleError.Message))
						}
					} else {
						config.Logger.Error(fmt.Sprintf("\t%s", singleError.Message))
					}
				}
			}
		}
	}
}
