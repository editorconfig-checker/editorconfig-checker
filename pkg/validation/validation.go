// Package validation contains all validation functions
package validation

import (
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/encoding"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/error"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/validation/validators"

	// x-release-please-end

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
// Note: This function is not thread safe, so it should not be called concurrently
func ValidateFile(filePath string, config config.Config) []error.ValidationError {
	// idiomatic Go allows empty struct
	if config.EditorconfigConfig == nil {
		config.EditorconfigConfig = &editorconfig.Config{}
	}

	// EditorconfigConfig isn't thread safe, so we need to lock it
	def, warnings, err := config.EditorconfigConfig.LoadGraceful(filePath)
	if err != nil {
		config.Logger.Error("cannot load %s as .editorconfig: %s", filePath, err)
		return nil
	}
	if warnings != nil {
		config.Logger.Warning("%v", warnings.Error())
	}

	return ValidateFileWithDefinition(filePath, config, def)
}

// ValidateFileWithDefinition Validates a single file with a given editorconfig definition and returns the errors
func ValidateFileWithDefinition(filePath string, config config.Config, def *editorconfig.Definition) []error.ValidationError {
	const directivePrefix = "editorconfig-checker-"
	const directiveDisable = directivePrefix + "disable"
	const directiveDisableFile = directivePrefix + "disable-file"
	const directiveDisableLine = directivePrefix + "disable-line"
	const directiveDisableNextLine = directivePrefix + "disable-next-line"
	const directiveEnable = directivePrefix + "enable"

	var validationErrors []error.ValidationError
	var isDisabled bool = false

	rawFileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileContent := string(rawFileContent)
	mime, err := files.GetContentTypeBytes(rawFileContent)
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
				config.Logger.Error("%v", err.Error())
			}
			break
		}
	}
	lines := files.ReadLines(fileContent)

	// return if first line contains editorconfig-checker-disable-file
	if len(lines) == 0 || strings.Contains(lines[0], directiveDisableFile) {
		return validationErrors
	}

	fileInformation := files.FileInformation{Content: fileContent, FilePath: filePath, Editorconfig: def}
	validationError := ValidateFinalNewline(fileInformation, config)
	if validationError.Message != nil {
		validationErrors = append(validationErrors, validationError)
	}

	fileInformation = files.FileInformation{Content: fileContent, FilePath: filePath, Editorconfig: def}
	validationError = ValidateLineEnding(fileInformation, config)
	if validationError.Message != nil {
		validationErrors = append(validationErrors, validationError)
	}

	var disableNextLineFound bool // used to ignore the line when editorconfig-checker-disable-next-line was found on previous line
	for lineNumber, line := range lines {
		// search for editorconfig-checker-enable
		// but only if not disabled for performance reasons
		if isDisabled && strings.Contains(line, directiveEnable) {
			isDisabled = false
		}

		// check for the status of the previous line (it was the next line on previous loop iteration)
		if disableNextLineFound {
			// editorconfig-checker-disable-next-line was found on previous line

			// check for successive editorconfig-checker-disable-next-line
			disableNextLineFound = strings.Contains(line, directiveDisableNextLine)

			// there is no need to check for editorconfig-checker-disable-line here, since line will be skipped

			// skip current line
			continue
		}

		if isDisabled {
			// no need to check further if disabled, for performance reasons
			continue
		}

		if directiveIndex := strings.Index(line, directiveDisable); directiveIndex != -1 {
			// a directive STARTING with editorconfig-checker-disable was found
			// let's check the possible modifiers

			directiveText := line[directiveIndex:] // shorten the text for performance reasons

			// this variable is here for reability, code could have been simplified, but it would have been harder to read
			activateDisable := true

			// check for editorconfig-checker-disable-next-line, and set status for next line
			if strings.Contains(directiveText, directiveDisableNextLine) {
				disableNextLineFound = true
				// it's not a editorconfig-checker-disable, there is no reason to disable all the following lines
				activateDisable = false
			}

			if strings.Contains(directiveText, directiveDisableLine) {
				// found editorconfig-checker-disable-line, skip current line
				continue
			}

			if activateDisable {
				// found editorconfig-checker-disable, skip current line and all following
				isDisabled = true
				continue
			}
		}

		fileInformation = files.FileInformation{Line: line, FilePath: filePath, LineNumber: lineNumber, Editorconfig: def}
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
		config.Logger.Verbose("Final newline error found in %s", fileInformation.FilePath)
		return error.ValidationError{LineNumber: -1, Message: currentError}
	}

	return error.ValidationError{}
}

// ValidateLineEnding runs the line ending validator and processes the error into the proper type
func ValidateLineEnding(fileInformation files.FileInformation, config config.Config) error.ValidationError {
	if currentError := validators.LineEnding(
		fileInformation.Content,
		fileInformation.Editorconfig.Raw["end_of_line"]); !config.Disable.EndOfLine && currentError != nil {
		config.Logger.Verbose("Line ending error found in %s", fileInformation.FilePath)
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
		config.Logger.Verbose("Indentation error found in %s on line %d", fileInformation.FilePath, fileInformation.LineNumber)
		return error.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return error.ValidationError{}
}

// ValidateTrailingWhitespace runs the TrailingWhitespace validator and processes the error into the proper type
func ValidateTrailingWhitespace(fileInformation files.FileInformation, config config.Config) error.ValidationError {
	if currentError := validators.TrailingWhitespace(
		fileInformation.Line,
		fileInformation.Editorconfig.Raw["trim_trailing_whitespace"] == "true"); !config.Disable.TrimTrailingWhitespace && currentError != nil {
		config.Logger.Verbose("Trailing whitespace error found in %s on line %d", fileInformation.FilePath, fileInformation.LineNumber)
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
		config.Logger.Verbose("Max line length error found in %s on %d", fileInformation.FilePath, fileInformation.LineNumber)
		return error.ValidationError{LineNumber: fileInformation.LineNumber + 1, Message: currentError}
	}

	return error.ValidationError{}
}

// ProcessValidation Validates all files and returns an array of validation errors
func ProcessValidation(files []string, config config.Config) []error.ValidationErrors {
	// idiomatic Go allows empty struct
	if config.EditorconfigConfig == nil {
		config.EditorconfigConfig = &editorconfig.Config{}
	}

	var (
		validationErrors = make([]*error.ValidationErrors, len(files))
		wg               sync.WaitGroup
		lock             sync.Mutex
	)
	// Limit the number of concurrent goroutines to the number of CPUs
	limiter := make(chan struct{}, runtime.NumCPU())

	for i, filePath := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Wait for a slot to be available in the limiter, and release it via defer
			limiter <- struct{}{}
			defer func() { <-limiter }()

			config.Logger.Verbose("Validate %s", filePath)

			// EditorconfigConfig isn't thread safe, so we need to acquire a lock
			lock.Lock()
			def, warnings, err := config.EditorconfigConfig.LoadGraceful(filePath)
			lock.Unlock()
			if err != nil {
				config.Logger.Error("cannot load %s as .editorconfig: %s", filePath, err)
				return
			}
			if warnings != nil {
				config.Logger.Warning("%v", warnings.Error())
			}
			errors := ValidateFileWithDefinition(filePath, config, def)

			lock.Lock()
			validationErrors[i] = &error.ValidationErrors{FilePath: filePath, Errors: errors}
			lock.Unlock()
		}()
	}

	wg.Wait()

	// Remove all nil values
	result := make([]error.ValidationErrors, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		if validationError != nil {
			result = append(result, *validationError)
		}
	}

	return result
}
