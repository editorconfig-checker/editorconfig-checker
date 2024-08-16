package error

import (
	"errors"
	"slices"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
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

func TestConsolidateErrors(t *testing.T) {
	input := []ValidationError{
		// two messages that become one block
		{LineNumber: 1, Message: errors.New("message kind one")},
		{LineNumber: 2, Message: errors.New("message kind one")},
		// one message with a good line between it and the last bad line, but repeating the message
		{LineNumber: 4, Message: errors.New("message kind one")},
		// one message that breaks the continousnes of line 4 to line 6
		{LineNumber: 5, Message: errors.New("message kind two")},
		{LineNumber: 6, Message: errors.New("message kind one")},
		// one message without a line number, that will become sorted to the top
		{LineNumber: -1, Message: errors.New("file-level error")},
	}

	expected := []ValidationError{
		{LineNumber: -1, Message: errors.New("file-level error")},
		{LineNumber: 1, ConsecuitiveCount: 1, Message: errors.New("message kind one")},
		{LineNumber: 4, Message: errors.New("message kind one")},
		{LineNumber: 5, Message: errors.New("message kind two")},
		{LineNumber: 6, Message: errors.New("message kind one")},
	}

	actual := ConsolidateErrors(input, config.Config{})

	if !slices.EqualFunc(expected, actual, func(e1 ValidationError, e2 ValidationError) bool { return e1.Equal(e2) }) {
		t.Log("consolidation expectation          :", expected)
		t.Log("consolidation actual returned value:", actual)
		t.Error("returned list of validation errors deviated from the expected set")
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

	config3 := config.Config{Format: "github-actions"}
	PrintErrors(input, config3)
}
