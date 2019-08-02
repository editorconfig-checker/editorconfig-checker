package error

import (
	"errors"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/types"
)

func TestGetErrorCount(t *testing.T) {
	count := GetErrorCount([]types.ValidationErrors{})
	if count != 0 {
		t.Error("Expected empty slice to have no errors, got", count)
	}

	input := []types.ValidationErrors{
		{
			FilePath: "some/path",
			Errors: []types.ValidationError{
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

	input = []types.ValidationErrors{
		{
			FilePath: "some/path",
			Errors: []types.ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
		{
			FilePath: "some/other/path",
			Errors: []types.ValidationError{
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

	input = []types.ValidationErrors{
		{
			FilePath: "some/path",
			Errors: []types.ValidationError{
				{
					LineNumber: 1,
					Message:    errors.New("WRONG"),
				},
			},
		},
		{
			FilePath: "some/other/path",
			Errors: []types.ValidationError{
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

	input = []types.ValidationErrors{
		{
			FilePath: "some/path",
			Errors:   []types.ValidationError{},
		},
		{
			FilePath: "some/other/path",
			Errors: []types.ValidationError{
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
	input := []types.ValidationErrors{
		{
			FilePath: "some/path",
			Errors:   []types.ValidationError{},
		},
		{
			FilePath: "some/other/path",
			Errors: []types.ValidationError{
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

	// wannabe test
	PrintErrors(input)
}
