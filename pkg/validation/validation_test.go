package validation

import (
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
)

func TestProcessValidation(t *testing.T) {
	params := types.Params{}

	processValidationResult := ProcessValidation([]string{"./main.go"}, params)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Should not have errors when validating main.go, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/disabled-file.ext"}, params)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Disabled file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/empty-file.txt"}, params)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Empty file should have no errors, got", processValidationResult)
	}
}
