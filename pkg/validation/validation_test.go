package validation

import (
	"os"
	"path/filepath"
	"testing"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	// x-release-please-end
)

func TestProcessValidation(t *testing.T) {
	configuration := config.NewConfig(nil)
	configuration.Verbose = true

	processValidationResult := ProcessValidation([]string{"./../../cmd/editorconfig-checker/main.go"}, *configuration)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Should not have errors when validating main.go, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/disabled-file.ext"}, *configuration)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Disabled file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/empty-file.txt"}, *configuration)
	if len(processValidationResult) > 1 || len(processValidationResult[0].Errors) != 0 {
		t.Error("Empty file should have no errors, got", processValidationResult)
	}

	processValidationResult = ProcessValidation([]string{"./../../testfiles/wrong-file.txt"}, *configuration)
	if (len(processValidationResult) > 1) || (len(processValidationResult[0].Errors) != 1) {
		t.Error("Wrong file should have errors, got", processValidationResult)
	}
}

func TestValidateFile(t *testing.T) {
	configuration := config.NewConfig(nil)
	configuration.Verbose = true

	result := ValidateFile("./../../cmd/editorconfig-checker/main.go", *configuration)
	if len(result) != 0 {
		t.Error("Should not have errors when validating main.go, got", result)
	}

	result = ValidateFile("./../../testfiles/wrong-file.txt", *configuration)
	if len(result) != 1 {
		t.Error("Should have errors when validating file with one error, got", result)
	}

	configuration.Disable.Indentation = true
	result = ValidateFile("./../../testfiles/wrong-file.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.SpacesAfterTabs = true
	result = ValidateFile("./../../testfiles/spaces-after-tabs.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.SpacesAfterTabs = false
	result = ValidateFile("./../../testfiles/zero-indent.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-line.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-next-line.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-block.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}

	result = ValidateFile("./../../testfiles/disabled-block-with-error.txt", *configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.SpacesAfterTabs = false
	result = ValidateFile("./../../testfiles/spaces-after-tabs.txt", *configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	result = ValidateFile("./../../testfiles/trailing-whitespace.txt", *configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	configuration.Disable.TrimTrailingWhitespace = true
	result = ValidateFile("./../../testfiles/trailing-whitespace.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	result = ValidateFile("./../../testfiles/final-newline-missing.txt", *configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	configuration.Disable.InsertFinalNewline = true
	result = ValidateFile("./../../testfiles/final-newline-missing.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	result = ValidateFile("./../../testfiles/wrong-line-ending.txt", *configuration)
	if len(result) == 0 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	result = ValidateFile("./../../testfiles/wrong-next-line.txt", *configuration)
	nbExpectedError := 2
	if len(result) != nbExpectedError {
		t.Errorf("Should have %d error, got %v", nbExpectedError, result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	configuration.Disable.EndOfLine = true
	configuration.Disable.InsertFinalNewline = true
	result = ValidateFile("./../../testfiles/wrong-line-ending.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	result = ValidateFile("./../../testfiles/line-to-long.txt", *configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	configuration.Disable.MaxLineLength = true
	result = ValidateFile("./../../testfiles/line-to-long.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	result = ValidateFile("./../../testfiles/spaces-with-tab.c", *configuration)
	if len(result) != 1 {
		t.Error("Should have one error, got", result)
	}

	configuration = config.NewConfig(nil)
	configuration.Verbose = true
	configuration.Disable.Indentation = true
	result = ValidateFile("./../../testfiles/favorites.gpx.txt", *configuration)
	if len(result) != 0 {
		t.Error("Should have no errors when validating valid file, got", result)
	}
}

// TestNonWildcardSectionOverridesGlob verifies that a specific filename
// section (e.g. [file.toml]) correctly overrides the [*] glob section.
//
// Regression for issue #100: editorconfig-checker was applying the [*]
// indent_style=tab rule to files matched by [file.toml] indent_style=space,
// incorrectly flagging space-indented TOML files as errors.
func TestNonWildcardSectionOverridesGlob(t *testing.T) {
	dir := t.TempDir()

	// Write .editorconfig with a global tab rule overridden for *.toml
	ecContent := `root = true

[*]
indent_style = tab
indent_size = 4

[*.toml]
indent_style = space
indent_size = 2
`
	if err := os.WriteFile(filepath.Join(dir, ".editorconfig"), []byte(ecContent), 0644); err != nil {
		t.Fatal("failed to write .editorconfig:", err)
	}

	// Write a TOML file that uses 2-space indentation — valid under [*.toml]
	tomlContent := `[section]\n  key = "value"\n`
	tomlPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(tomlPath, []byte(tomlContent), 0644); err != nil {
		t.Fatal("failed to write config.toml:", err)
	}

	cfg := config.NewConfig(nil)
	cfg.Verbose = true

	errs := ValidateFile(tomlPath, *cfg)
	if len(errs) != 0 {
		t.Errorf(
			"Expected no errors for space-indented TOML file under [*.toml] section, got: %v\n"+
				"This indicates [*.toml] is not overriding the [*] indent_style=tab rule.",
			errs,
		)
	}

	// Also verify that a non-TOML file still requires tabs
	goContent := "package main\n\nfunc main() {\n\t// tab-indented\n}\n"
	goPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatal("failed to write main.go:", err)
	}

	errs = ValidateFile(goPath, *cfg)
	if len(errs) != 0 {
		t.Errorf("Expected no errors for tab-indented Go file under [*] section, got: %v", errs)
	}
}
