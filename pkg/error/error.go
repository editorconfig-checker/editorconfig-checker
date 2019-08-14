package error

import (
	"fmt"
	"os"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
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
func PrintErrors(errors []ValidationErrors) {
	for _, fileErrors := range errors {
		if len(fileErrors.Errors) > 0 {
			relativeFilePath, err := files.GetRelativePath(fileErrors.FilePath)

			if err != nil {
				logger.Error(err.Error())
			}

			logger.Print(fmt.Sprintf("%s:", relativeFilePath), logger.YELLOW, os.Stderr)
			for _, singleError := range fileErrors.Errors {
				if singleError.LineNumber != -1 {
					logger.Error(fmt.Sprintf("\t%d: %s", singleError.LineNumber, singleError.Message))
				} else {
					logger.Error(fmt.Sprintf("\t%s", singleError.Message))
				}

			}
		}
	}
}
