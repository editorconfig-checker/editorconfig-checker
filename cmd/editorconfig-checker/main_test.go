package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestProcessValidation(t *testing.T) {
	// Should not have errors when validating main.go
	processValidationResult := processValidation([]string{"./main.go"}, false, false)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Expected something, got", processValidationResult)
	}

	processValidationResult = processValidation([]string{"./main.go"}, true, false)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Expected something, got", processValidationResult)
	}

	processValidationResult = processValidation([]string{"./../../testfiles/disabled-file.ext"}, true, false)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Expected something, got", processValidationResult)
	}

	// empty file should have no errors
	processValidationResult = processValidation([]string{"./../../testfiles/empty-file.txt"}, true, false)
	if len(processValidationResult) > 1 && len(processValidationResult[0].Errors) != 0 {
		t.Error("Expected something, got", processValidationResult)
	}
}

func BenchmarkMain(b *testing.B) {
	// run the binary b.N times
	for n := 0; n < b.N; n++ {
		dir, _ := os.Getwd()
		cmd := exec.Command("make", "run")
		// the test is executed where the `*_test.go` file is located
		cmd.Dir = dir + "/../../"
		err := cmd.Run()

		if err != nil {
			panic("Something went wrong")
		}
	}
}
