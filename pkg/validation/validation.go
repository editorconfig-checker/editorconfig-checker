package validation

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/editorconfig/editorconfig-core-go/v2"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/validation/validators"
)

// ValidateFile Validates a single file and returns the errors
func ValidateFile(filePath string, config config.Config) []types.ValidationError {
	var validationErrors []types.ValidationError
	lines := files.ReadLines(filePath)

	// return if first line contains editorconfig-checker-disable-file
	if len(lines) == 0 || strings.Contains(lines[0], "editorconfig-checker-disable-file") {
		return validationErrors
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

	fileInformation := types.FileInformation{Content: fileContent, FilePath: filePath, Editorconfig: editorconfig}
	validationError := ValidateFinalNewline(fileInformation, config)
	if validationError.Message != nil {
		validationErrors = append(validationErrors, validationError)
	}

	fileInformation = types.FileInformation{Content: fileContent, FilePath: filePath, Editorconfig: editorconfig}
	validationError = ValidateLineEnding(fileInformation, config)
	if validationError.Message != nil {
		validationErrors = append(validationErrors, validationError)
	}

	for lineNumber, line := range lines {
		if strings.Contains(line, "editorconfig-checker-disable-line") {
			continue
		}

		fileInformation = types.FileInformation{Line: line, FilePath: filePath, LineNumber: lineNumber, Editorconfig: editorconfig}
		validationError = ValidateTrailingWhitespace(fileInformation, config)
		if validationError.Message != nil {
			validationErrors = append(validationErrors, validationError)
		}

		validationError = ValidateIndentation(fileInformation, config)
		if validationError.Message != nil {
			validationErrors = append(validationErrors, validationError)
		}

	}

	return validationErrors
}

// ValidateFinalNewline runs the final newline validator and processes the error into the proper type
func ValidateFinalNewline(fileInformation types.FileInformation, config config.Config) types.ValidationError {
	if currentError := validators.FinalNewline(
		fileInformation.Content,
		fileInformation.Editorconfig.Raw["insert_final_newline"],
		fileInformation.Editorconfig.Raw["end_of_line"]); !config.Disable.Insert_Final_Newline && currentError != nil {
		if config.Verbose {
			logger.Output(fmt.Sprintf("Final newline error found in %s", fileInformation.FilePath))
		}
		return types.ValidationError{LineNumber: -1, Message: currentError}
	}

	return types.ValidationError{}
}

// ValidateLineEnding runs the line ending validator and processes the error into the proper type
func ValidateLineEnding(fileInformation types.FileInformation, config config.Config) types.ValidationError {
	if currentError := validators.LineEnding(
		fileInformation.Content,
		fileInformation.Editorconfig.Raw["end_of_line"]); !config.Disable.End_Of_Line && currentError != nil {
		if config.Verbose {
			logger.Output(fmt.Sprintf("Line ending error found in %s", fileInformation.FilePath))
		}

		return types.ValidationError{LineNumber: -1, Message: currentError}
	}

	return types.ValidationError{}
}

// ValidateIndentation runs the Indentation validator and processes the error into the proper type
func ValidateIndentation(fileInformation types.FileInformation, config config.Config) types.ValidationError {
	var indentSize int
	indentSize, err := strconv.Atoi(fileInformation.Editorconfig.Raw["indent_size"])

	// Set indentSize to zero if there is no indentSize set
	if err != nil {
		indentSize = 0
	}

	if currentError := validators.Indentation(
		fileInformation.Line,
		fileInformation.Editorconfig.Raw["indent_style"],
		indentSize, config); !config.Disable.Indentation && currentError != nil {
		if config.Verbose {
			logger.Output(fmt.Sprintf("Indentation error found in %s on line %d", fileInformation.FilePath, fileInformation.LineNumber))
		}

		return types.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return types.ValidationError{}
}

// ValidateTrailingWhitespace runs the TrailingWhitespace validator and processes the error into the proper type
func ValidateTrailingWhitespace(fileInformation types.FileInformation, config config.Config) types.ValidationError {
	if currentError := validators.TrailingWhitespace(
		fileInformation.Line,
		fileInformation.Editorconfig.Raw["trim_trailing_whitespace"] == "true"); !config.Disable.Trim_Trailing_Whitespace && currentError != nil {
		if config.Verbose {
			logger.Output(fmt.Sprintf("Trailing whitespace error found in %s on line %d", fileInformation.FilePath, fileInformation.LineNumber))
		}
		return types.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return types.ValidationError{}
}

// ProcessValidation Validates all files and returns an array of validation errors
func ProcessValidation(files []string, config config.Config) []types.ValidationErrors {
	var validationErrors []types.ValidationErrors

	for _, filePath := range files {
		if config.Verbose {
			logger.Output(fmt.Sprintf("Validate %s", filePath))
		}
		validationErrors = append(validationErrors, types.ValidationErrors{FilePath: filePath, Errors: ValidateFile(filePath, config)})
	}

	return validationErrors
}
