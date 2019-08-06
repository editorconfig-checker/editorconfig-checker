package validation

import (
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
)

func TestProcessValidation(t *testing.T) {
	params := types.Params{Verbose: true}

	processValidationResult := ProcessValidation([]string{"./main.go"}, params)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Should not have errors when validating main.go, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/disabled-file.ext"}, params)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Disabled file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/empty-file.txt"}, params)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Empty file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/wrong-file.txt"}, params)
	if (len(processValidationResult) > 1) || (len(processValidationResult[0].Errors) != 1) {
		t.Error("Wrong file should have errors, got", processValidationResult)
	}
}

func TestValidateFile(t *testing.T) {
	params := types.Params{Verbose: true}

	result := ValidateFile("./main.go", params)
	if len(result) != 0 {
		t.Error("Should not have errors when validating main.go, got", result)
	}

	result = ValidateFile("./../../testfiles/wrong-file.txt", params)
	if len(result) != 1 {
		t.Error("Should have errors when validating file with one error, got", result)
	}

	params = types.Params{SpacesAfterTabs: true}
	result = ValidateFile("./../../testfiles/spaces-after-tabs.txt", params)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	params = types.Params{SpacesAfterTabs: false}
	result = ValidateFile("./../../testfiles/zero-indent.txt", params)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-line.txt", params)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	params = types.Params{SpacesAfterTabs: false}
	result = ValidateFile("./../../testfiles/spaces-after-tabs.txt", params)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	params = types.Params{Verbose: true}
	result = ValidateFile("./../../testfiles/trailing-whitespace.txt", params)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	params = types.Params{Verbose: true}
	params.Disabled.TrailingWhitspace = true
	result = ValidateFile("./../../testfiles/trailing-whitespace.txt", params)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	params = types.Params{Verbose: true}
	result = ValidateFile("./../../testfiles/final-newline-missing.txt", params)
	if len(result) != 1 {
		t.Error("Should have no error, got", result)
	}

	params = types.Params{Verbose: true}
	result = ValidateFile("./../../testfiles/wrong-line-ending.txt", params)
	if len(result) == 0 {
		t.Error("Should have one error, got", result)
	}

	params = types.Params{Verbose: true}
	params.Disabled.LineEnding = true
	result = ValidateFile("./../../testfiles/wrong-line-ending.txt", params)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}
}
