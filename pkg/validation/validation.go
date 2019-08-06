package validation

import (
	"fmt"
	"github.com/editorconfig/editorconfig-core-go/v2"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/validation/validators"
)

// Validates a single file and returns the errors
func ValidateFile(filePath string, params types.Params) []types.ValidationError {
	var errors []types.ValidationError
	lines := files.ReadLines(filePath)

	// return if first line contains editorconfig-checker-disable-file
	if len(lines) == 0 || strings.Contains(lines[0], "editorconfig-checker-disable-file") {
		return errors
	}

	rawFileContent, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	fileContent := string(rawFileContent)

	editorconfig, err := editorconfig.GetDefinitionForFilename(filePath)
	if err != nil {
		panic(err)
	}

	if currentError := validators.FinalNewline(
		fileContent,
		editorconfig.Raw["insert_final_newline"],
		editorconfig.Raw["end_of_line"]); currentError != nil {
		if params.Verbose {
			logger.Output(fmt.Sprintf("Final newline error found in %s", filePath))
		}
		errors = append(errors, types.ValidationError{LineNumber: -1, Message: currentError})
	}

	if currentError := validators.LineEnding(
		fileContent,
		editorconfig.Raw["end_of_line"]); currentError != nil {
		if params.Verbose {
			logger.Output(fmt.Sprintf("Line ending error found in %s", filePath))
		}
		errors = append(errors, types.ValidationError{LineNumber: -1, Message: currentError})
	}

	for lineNumber, line := range lines {
		if strings.Contains(line, "editorconfig-checker-disable-line") {
			continue
		}

		if currentError := validators.TrailingWhitespace(
			line,
			editorconfig.Raw["trim_trailing_whitespace"] == "true"); currentError != nil {
			if params.Verbose {
				logger.Output(fmt.Sprintf("Trailing whitespace error found in %s on line %d", filePath, lineNumber))
			}
			errors = append(errors, types.ValidationError{LineNumber: lineNumber + 1, Message: currentError})
		}

		var indentSize int
		indentSize, err = strconv.Atoi(editorconfig.Raw["indent_size"])

		// Set indentSize to zero if there is no indentSize set
		if err != nil {
			indentSize = 0
		}

		if currentError := validators.Indentation(
			line,
			editorconfig.Raw["indent_style"],
			indentSize, params); currentError != nil {
			if params.Verbose {
				logger.Output(fmt.Sprintf("Indentation error found in %s on line %d", filePath, lineNumber))
			}
			errors = append(errors, types.ValidationError{LineNumber: lineNumber + 1, Message: currentError})
		}
	}

	return errors
}

// Validates all files and returns an array of validation errors
func ProcessValidation(files []string, params types.Params) []types.ValidationErrors {
	var validationErrors []types.ValidationErrors

	for _, filePath := range files {
		if params.Verbose {
			logger.Output(fmt.Sprintf("Validate %s", filePath))
		}
		validationErrors = append(validationErrors, types.ValidationErrors{FilePath: filePath, Errors: ValidateFile(filePath, params)})
	}

	return validationErrors
}
