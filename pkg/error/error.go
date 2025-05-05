// Package error contains functions and structs related to errors
package error

import (
	"encoding/json"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/outputformat"
	// x-release-please-end
)

// ValidationError represents one validation error
type ValidationError struct {
	LineNumber                    int
	Message                       error
	AdditionalIdenticalErrorCount int
}

// ValidationErrors represents which errors occurred in a file
type ValidationErrors struct {
	FilePath string
	Errors   []ValidationError
}

func (error1 *ValidationError) Equal(error2 ValidationError) bool {
	if error1.Message.Error() != error2.Message.Error() {
		return false
	}
	if error1.LineNumber != error2.LineNumber {
		return false
	}
	if error1.AdditionalIdenticalErrorCount != error2.AdditionalIdenticalErrorCount {
		return false
	}
	return true

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
			if nextError.Message.Error() == thisError.Message.Error() && nextError.LineNumber == thisError.LineNumber+thisError.AdditionalIdenticalErrorCount+1 {
				thisError.AdditionalIdenticalErrorCount++ // keep track of how many consecutive lines we've seen
				i++                                       // make sure the outer loop jumps over the consecutive errors we just found
			} else {
				break // if they are different errors we can stop comparing messages
			}
		}

		consolidatedErrors = append(consolidatedErrors, thisError)
	}

	return append(lineLessErrors, consolidatedErrors...)
}

func PrintErrorCount(errorCount int, config config.Config) {
	if errorCount == 0 {
		config.Logger.Verbose("\n%d errors found", errorCount)
		return
	}
	config.Logger.Error("\n%d errors found", errorCount)
}

func PrintErrorsAsHumanReadable(errors []ValidationErrors, config config.Config) {
	errorCount := 0
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			config.Logger.Error("%v", err.Error())
			continue
		}

		fileErrors.Errors = ConsolidateErrors(fileErrors.Errors, config)

		config.Logger.Warning("%s:", relativeFilePath)
		for _, singleError := range fileErrors.Errors {
			errorCount++

			if singleError.LineNumber == -1 {
				config.Logger.Error("\t%s", singleError.Message)
				continue
			}

			if singleError.AdditionalIdenticalErrorCount == 0 {
				config.Logger.Error("\t%d: %s", singleError.LineNumber, singleError.Message)
				continue
			}

			config.Logger.Error("\t%d-%d: %s", singleError.LineNumber, singleError.LineNumber+singleError.AdditionalIdenticalErrorCount, singleError.Message)
		}
	}
	PrintErrorCount(errorCount, config)
}

func PrintErrorsAsGHA(errors []ValidationErrors, config config.Config) {
	errorCount := 0
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			config.Logger.Error("%v", err.Error())
			continue
		}

		fileErrors.Errors = ConsolidateErrors(fileErrors.Errors, config)

		// github-actions: A format dedicated for usage in Github Actions
		for _, singleError := range fileErrors.Errors {
			errorCount++

			if singleError.LineNumber == -1 {
				config.Logger.Error("::error file=%s::%s", relativeFilePath, singleError.Message)
				continue
			}

			if singleError.AdditionalIdenticalErrorCount == 0 {
				config.Logger.Error("::error file=%s,line=%d::%s", relativeFilePath, singleError.LineNumber, singleError.Message)
				continue
			}

			config.Logger.Error("::error file=%s,line=%d,endLine=%d::%s", relativeFilePath, singleError.LineNumber, singleError.LineNumber+singleError.AdditionalIdenticalErrorCount, singleError.Message)
		}
	}
	PrintErrorCount(errorCount, config)
}

// gcc: A format mimicking the error format from GCC.
func PrintErrorsAsGCC(errors []ValidationErrors, config config.Config) {
	errorCount := 0
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			config.Logger.Error("%v", err.Error())
			continue
		}

		for _, singleError := range fileErrors.Errors {
			errorCount++

			lineNo := 0
			if singleError.LineNumber > 0 {
				lineNo = singleError.LineNumber
			}
			config.Logger.Error("%s:%d:%d: %s: %s", relativeFilePath, lineNo, 0, "error", singleError.Message)
		}
	}
	PrintErrorCount(errorCount, config)
}

// codeclimate: A format that is compatible with the codeclimate format for GitLab CI.
// https://docs.gitlab.com/ee/ci/testing/code_quality.html#implement-a-custom-tool
func PrintErrorsAsCodeclimate(errors []ValidationErrors, config config.Config) {
	var codeclimateIssues []CodeclimateIssue

	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			config.Logger.Error("%v", err.Error())
			continue
		}

		fileErrors.Errors = ConsolidateErrors(fileErrors.Errors, config)

		for _, singleError := range fileErrors.Errors {
			codeclimateIssues = append(codeclimateIssues, newCodeclimateIssue(singleError, relativeFilePath))
		}
	}

	if len(codeclimateIssues) > 0 {
		// marshall codeclimate issues to json
		codeclimateIssuesJSON, err := json.Marshal(codeclimateIssues)
		if err != nil {
			config.Logger.Error("Error creating codeclimate json: %s", err.Error())
		} else {
			config.Logger.Output("%s", string(codeclimateIssuesJSON))
		}
	}
}

// PrintErrors prints the errors to the console
func PrintErrors(errors []ValidationErrors, config config.Config) {
	switch config.Format {
	case outputformat.Codeclimate:
		// codeclimate: A format that is compatible with the codeclimate format for GitLab CI.
		// https://docs.gitlab.com/ee/ci/testing/code_quality.html#implement-a-custom-tool
		PrintErrorsAsCodeclimate(errors, config)
	case outputformat.GCC:
		// gcc: A format mimicking the error format from GCC.
		PrintErrorsAsGCC(errors, config)
	case outputformat.GithubActions:
		// github-actions: A format dedicated for usage in Github Actions
		PrintErrorsAsGHA(errors, config)
	default:
		// default: A human readable text format.
		PrintErrorsAsHumanReadable(errors, config)
	}
}
