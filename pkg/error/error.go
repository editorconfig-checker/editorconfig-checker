// Package error contains functions and structs related to errors
package error

import (
	"encoding/json"
	"fmt"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"       // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"        // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/logger"       // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/outputformat" // x-release-please-major
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
				i = i + j + 1                             // make sure the outer loop jumps over the consecutive errors we just found
			} else {
				break // if they are different errors we can stop comparing messages
			}
		}

		consolidatedErrors = append(consolidatedErrors, thisError)
	}

	return append(lineLessErrors, consolidatedErrors...)
}

func FormatErrorsAsHumanReadable(errors []ValidationErrors, config config.Config) []logger.LogMessage {
	var logMessages []logger.LogMessage

	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: err.Error()})
			continue
		}

		fileErrors.Errors = ConsolidateErrors(fileErrors.Errors, config)

		logMessages = append(logMessages, logger.LogMessage{Level: "warning", Message: fmt.Sprintf("%s:", relativeFilePath)})
		for _, singleError := range fileErrors.Errors {
			if singleError.LineNumber == -1 {
				logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("\t%s", singleError.Message)})
				continue
			}

			if singleError.AdditionalIdenticalErrorCount == 0 {
				logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("\t%d: %s", singleError.LineNumber, singleError.Message)})
				continue
			}

			logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("\t%d-%d: %s", singleError.LineNumber, singleError.LineNumber+singleError.AdditionalIdenticalErrorCount, singleError.Message)})
		}
	}

	return logMessages
}

func FormatErrorsAsGHA(errors []ValidationErrors, config config.Config) []logger.LogMessage {
	var logMessages []logger.LogMessage

	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: err.Error()})
			continue
		}

		fileErrors.Errors = ConsolidateErrors(fileErrors.Errors, config)

		// github-actions: A format dedicated for usage in Github Actions
		for _, singleError := range fileErrors.Errors {
			if singleError.LineNumber == -1 {
				logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("::error file=%s::%s", relativeFilePath, singleError.Message)})
				continue
			}

			if singleError.AdditionalIdenticalErrorCount == 0 {
				logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("::error file=%s,line=%d::%s", relativeFilePath, singleError.LineNumber, singleError.Message)})
				continue
			}

			logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("::error file=%s,line=%d,endLine=%d::%s", relativeFilePath, singleError.LineNumber, singleError.LineNumber+singleError.AdditionalIdenticalErrorCount, singleError.Message)})
		}
	}

	return logMessages
}

// gcc: A format mimicking the error format from GCC.
func FormatErrorsAsGCC(errors []ValidationErrors, config config.Config) []logger.LogMessage {
	var logMessages []logger.LogMessage

	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: err.Error()})
			continue
		}

		for _, singleError := range fileErrors.Errors {
			lineNo := 0
			if singleError.LineNumber > 0 {
				lineNo = singleError.LineNumber
			}
			logMessages = append(logMessages, logger.LogMessage{Level: "error", Message: fmt.Sprintf("%s:%d:%d: %s: %s", relativeFilePath, lineNo, 0, "error", singleError.Message)})
		}
	}

	return logMessages
}

// codeclimate: A format that is compatible with the codeclimate format for GitLab CI.
// https://docs.gitlab.com/ee/ci/testing/code_quality.html#implement-a-custom-tool
func FormatErrorsAsCodeclimate(errors []ValidationErrors, config config.Config) []logger.LogMessage {
	var codeclimateIssues []CodeclimateIssue

	for _, fileErrors := range errors {
		if len(fileErrors.Errors) == 0 {
			continue
		}

		relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)
		if err != nil {
			config.Logger.Error(err.Error())
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
			return []logger.LogMessage{{Level: "", Message: string(codeclimateIssuesJSON)}}
		}
	}
	return []logger.LogMessage{}
}

// FormatErrors prints the errors to the console
func FormatErrors(errors []ValidationErrors, config config.Config) []logger.LogMessage {
	switch config.Format {
	case outputformat.Codeclimate:
		// codeclimate: A format that is compatible with the codeclimate format for GitLab CI.
		// https://docs.gitlab.com/ee/ci/testing/code_quality.html#implement-a-custom-tool
		return FormatErrorsAsCodeclimate(errors, config)
	case outputformat.GCC:
		// gcc: A format mimicking the error format from GCC.
		return FormatErrorsAsGCC(errors, config)
	case outputformat.GithubActions:
		// github-actions: A format dedicated for usage in Github Actions
		return FormatErrorsAsGHA(errors, config)
	default:
		// default: A human readable text format.
		return FormatErrorsAsHumanReadable(errors, config)
	}
}
