// Package validation contains all validation functions
package validation

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/encoding"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/error"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/validation/validators"
	"github.com/editorconfig/editorconfig-core-go/v2"
)

// keep synced with /pkg/config/config.go#L59
var textRegexes = []string{
	"^text/",
	"^application/ecmascript$",
	"^application/json$",
	"^application/x-ndjson$",
	"^application/xml$",
	"+json",
	"+xml$",
}

// ValidateFile Validates a single file and returns the errors
func ValidateFile(filePath string, config config.Config) []error.ValidationError {
	const directivePrefix = "editorconfig-checker-"
	const directiveDisable = directivePrefix + "disable"
	const directiveDisableFile = directivePrefix + "disable-file"
	const directiveDisableLine = directivePrefix + "disable-line"
	const directiveEnable = directivePrefix + "enable"

	var validationErrors []error.ValidationError
	var isDisabled bool = false

	rawFileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileContent := string(rawFileContent)
	mime, err := files.GetContentTypeBytes(rawFileContent, config)
	if err != nil {
		panic(err)
	}
	for _, regex := range textRegexes {
		match, _ := regexp.MatchString(regex, mime)
		if match {
			var charset string
			fileContent, charset, err = encoding.DecodeBytes(rawFileContent)
			if err != nil {
				if charset == "" {
					charset = "unknown"
				}
				config.Logger.Error("Could not decode the %s encoded file: %s", charset, filePath)
				config.Logger.Error(err.Error())
			}
			break
		}
	}
	lines := files.ReadLines(fileContent)

	// return if first line contains editorconfig-checker-disable-file
	if len(lines) == 0 || strings.Contains(lines[0], directiveDisableFile) {
		return validationErrors
	}

	// idiomatic Go allows empty struct
	if config.EditorconfigConfig == nil {
		config.EditorconfigConfig = &editorconfig.Config{}
	}

	definition, warnings, err := config.EditorconfigConfig.LoadGraceful(filePath)
	if err != nil {
		config.Logger.Error("cannot load %s as .editorconfig: %s", filePath, err)
		return validationErrors
	}
	if warnings != nil {
		config.Logger.Warning(warnings.Error())
	}

	fileInformation := files.FileInformation{Content: fileContent, FilePath: filePath, Editorconfig: definition}
	validationError := ValidateFinalNewline(fileInformation, config)
	if validationError.Message != nil {
		validationErrors = append(validationErrors, validationError)
	}

	fileInformation = files.FileInformation{Content: fileContent, FilePath: filePath, Editorconfig: definition}
	validationError = ValidateLineEnding(fileInformation, config)
	if validationError.Message != nil {
		validationErrors = append(validationErrors, validationError)
	}

	for lineNumber, line := range lines {
		if strings.Contains(line, directiveEnable) {
			isDisabled = false
		} else if directiveIndex := strings.Index(line, directiveDisable); directiveIndex != -1 {
			if strings.HasPrefix(line[directiveIndex:], directiveDisableLine) {
				// found editorconfig-checker-disable-line, so just skip current line
				continue
			} else {
				isDisabled = true
			}
		}

		if isDisabled {
			continue
		}

		fileInformation = files.FileInformation{Line: line, FilePath: filePath, LineNumber: lineNumber, Editorconfig: definition}
		validationError = ValidateTrailingWhitespace(fileInformation, config)
		if validationError.Message != nil {
			validationErrors = append(validationErrors, validationError)
		}

		validationError = ValidateIndentation(fileInformation, config)
		if validationError.Message != nil {
			validationErrors = append(validationErrors, validationError)
		}

		validationError = ValidateMaxLineLength(fileInformation, config)
		if validationError.Message != nil {
			validationErrors = append(validationErrors, validationError)
		}
	}

	return validationErrors
}

// ValidateFinalNewline runs the final newline validator and processes the error into the proper type
func ValidateFinalNewline(fileInformation files.FileInformation, config config.Config) error.ValidationError {
	if currentError := validators.FinalNewline(
		fileInformation.Content,
		fileInformation.Editorconfig.Raw["insert_final_newline"],
		fileInformation.Editorconfig.Raw["end_of_line"]); !config.Disable.InsertFinalNewline && currentError != nil {
		config.Logger.Verbose(fmt.Sprintf("Final newline error found in %s", fileInformation.FilePath))
		return error.ValidationError{LineNumber: -1, Message: currentError}
	}

	return error.ValidationError{}
}

// ValidateLineEnding runs the line ending validator and processes the error into the proper type
func ValidateLineEnding(fileInformation files.FileInformation, config config.Config) error.ValidationError {
	if currentError := validators.LineEnding(
		fileInformation.Content,
		fileInformation.Editorconfig.Raw["end_of_line"]); !config.Disable.EndOfLine && currentError != nil {
		config.Logger.Verbose(fmt.Sprintf("Line ending error found in %s", fileInformation.FilePath))
		return error.ValidationError{LineNumber: -1, Message: currentError}
	}

	return error.ValidationError{}
}

// ValidateIndentation runs the Indentation validator and processes the error into the proper type
func ValidateIndentation(fileInformation files.FileInformation, config config.Config) error.ValidationError {
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
		config.Logger.Verbose(fmt.Sprintf("Indentation error found in %s on line %d", fileInformation.FilePath, fileInformation.LineNumber))
		return error.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return error.ValidationError{}
}

// ValidateTrailingWhitespace runs the TrailingWhitespace validator and processes the error into the proper type
func ValidateTrailingWhitespace(fileInformation files.FileInformation, config config.Config) error.ValidationError {
	if currentError := validators.TrailingWhitespace(
		fileInformation.Line,
		fileInformation.Editorconfig.Raw["trim_trailing_whitespace"] == "true"); !config.Disable.TrimTrailingWhitespace && currentError != nil {
		config.Logger.Verbose(fmt.Sprintf("Trailing whitespace error found in %s on line %d", fileInformation.FilePath, fileInformation.LineNumber))
		return error.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return error.ValidationError{}
}

// ValidateMaxLineLength runs the max line length validator and processes the error into the proper type
func ValidateMaxLineLength(fileInformation files.FileInformation, config config.Config) error.ValidationError {
	maxLineLength, err := strconv.Atoi(fileInformation.Editorconfig.Raw["max_line_length"])
	if err != nil {
		return error.ValidationError{}
	}

	charSet := fileInformation.Editorconfig.Raw["charset"]

	if currentError := validators.MaxLineLength(fileInformation.Line, maxLineLength, charSet); !config.Disable.MaxLineLength && currentError != nil {
		config.Logger.Verbose(fmt.Sprintf("Max line length error found in %s on %d", fileInformation.FilePath, fileInformation.LineNumber))
		return error.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return error.ValidationError{}
}

// ProcessValidation Validates all files and returns an array of validation errors
func ProcessValidation(files []string, config config.Config) []error.ValidationErrors {
	var validationErrors []error.ValidationErrors

	for _, filePath := range files {
		config.Logger.Verbose(fmt.Sprintf("Validate %s", filePath))
		validationErrors = append(validationErrors, error.ValidationErrors{FilePath: filePath, Errors: ValidateFile(filePath, config)})
	}

	return validationErrors
}
