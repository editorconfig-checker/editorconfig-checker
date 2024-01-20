package error

import (
	"errors"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/v2/pkg/config"
)

func TestGetErrorCount(t *testing.T) {
	count := GetErrorCount([]ValidationErrors{})
	if count != 0 {
		t.Error("Expected empty slice to have no errors, got", count)
	}

	input := []ValidationErrors{
		{
			FilePath: "some/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
	}

	count = GetErrorCount(input)
	if count != 1 {
		t.Error("Expected one error slice to have exactly one erorr errors, got", count)
	}

	input = []ValidationErrors{
		{
			FilePath: "some/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
		{
			FilePath: "some/other/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
	}

	count = GetErrorCount(input)
	if count != 2 {
		t.Error("Expected two error slice to have exactly one erorr errors, got", count)
	}

	input = []ValidationErrors{
		{
			FilePath: "some/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
		{
			FilePath: "some/other/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
				{
					LineNumber: 2,
					Message:    errors.New("WRONG"),
				},
			},
		},
	}

	count = GetErrorCount(input)
	if count != 3 {
		t.Error("Expected three error slice to have exactly one erorr errors, got", count)
	}

	input = []ValidationErrors{
		{
			FilePath: "some/path",
			Errors:   []ValidationError{},
		},
		{
			FilePath: "some/other/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
				{
					LineNumber: 2,
					Message:    errors.New("WRONG"),
				},
			},
		},
	}

	count = GetErrorCount(input)
	if count != 2 {
		t.Error("Expected three error slice to have exactly one erorr errors, got", count)
	}
}

func TestPrintErrors(t *testing.T) {
	input := []ValidationErrors{
		{
			FilePath: "some/path",
			Errors:   []ValidationError{},
		},
		{
			FilePath: "/proc/cpuinfo",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
		{
			FilePath: "/proc/cpuinfoNOT",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
		{
			FilePath: "some/other/path",
			Errors: []ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
				{
					LineNumber: -1,
					Message:    errors.New("WRONG"),
				},
			},
		},
	}

	// wannabe test
	config1 := config.Config{}
	PrintErrors(input, config1)

	config2 := config.Config{Format: "gcc"}
	PrintErrors(input, config2)
}
