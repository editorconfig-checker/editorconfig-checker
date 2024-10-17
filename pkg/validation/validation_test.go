package validation

import (
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config" // x-release-please-major
)

func TestProcessValidation(t *testing.T) {
	configuration := config.Config{
		Verbose: true,
	}

	processValidationResult := ProcessValidation([]string{"./../../cmd/editorconfig-checker/main.go"}, configuration)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Should not have errors when validating main.go, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/disabled-file.ext"}, configuration)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Disabled file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/empty-file.txt"}, configuration)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Empty file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/wrong-file.txt"}, configuration)
	if (len(processValidationResult) > 1) || (len(processValidationResult[0].Errors) != 1) {
		t.Error("Wrong file should have errors, got", processValidationResult)
	}
}

func TestValidateFile(t *testing.T) {
	configuration := config.Config{Verbose: true}

	result := ValidateFile("./../../cmd/editorconfig-checker/main.go", configuration)
	if len(result) != 0 {
		t.Error("Should not have errors when validating main.go, got", result)
	}

	result = ValidateFile("./../../testfiles/wrong-file.txt", configuration)
	if len(result) != 1 {
		t.Error("Should have errors when validating file with one error, got", result)
	}

	configuration.Disable.Indentation = true
	result = ValidateFile("./../../testfiles/wrong-file.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no errors, got", result)
	}

	configuration = config.Config{SpacesAfterTabs: true}
	result = ValidateFile("./../../testfiles/spaces-after-tabs.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	configuration = config.Config{SpacesAfterTabs: false}
	result = ValidateFile("./../../testfiles/zero-indent.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-line.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-next-line.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-block.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-block-with-error.txt", configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.Config{SpacesAfterTabs: false}
	result = ValidateFile("./../../testfiles/spaces-after-tabs.txt", configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.Config{Verbose: true}
	result = ValidateFile("./../../testfiles/trailing-whitespace.txt", configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.Config{Verbose: true}
	configuration.Disable.TrimTrailingWhitespace = true
	result = ValidateFile("./../../testfiles/trailing-whitespace.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.Config{Verbose: true}
	result = ValidateFile("./../../testfiles/final-newline-missing.txt", configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.Config{Verbose: true}
	configuration.Disable.InsertFinalNewline = true
	result = ValidateFile("./../../testfiles/final-newline-missing.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.Config{Verbose: true}
	result = ValidateFile("./../../testfiles/wrong-line-ending.txt", configuration)
	if len(result) == 0 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.Config{Verbose: true}
	result = ValidateFile("./../../testfiles/wrong-next-line.txt", configuration)
	nbExpectedError := 2
	if len(result) != nbExpectedError {
		t.Errorf("Should have %d error, got %v", nbExpectedError, result)
	}

	configuration = config.Config{Verbose: true}
	configuration.Disable.EndOfLine = true
	configuration.Disable.InsertFinalNewline = true
	result = ValidateFile("./../../testfiles/wrong-line-ending.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.Config{Verbose: true}
	result = ValidateFile("./../../testfiles/line-to-long.txt", configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.Config{Verbose: true}
	configuration.Disable.MaxLineLength = true
	result = ValidateFile("./../../testfiles/line-to-long.txt", configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.Config{Verbose: true}
	result = ValidateFile("./../../testfiles/spaces-with-tab.c", configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}
}
