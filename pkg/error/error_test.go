package error

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/outputformat"

	// x-release-please-end

	"github.com/gkampitakis/go-snaps/snaps"
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

func TestValidationErrorEqual(t *testing.T) {
	baseError := ValidationError{
		LineNumber: -1,
		Message:    errors.New("a message"),
	}
	wrongLineNumberError := ValidationError{
		LineNumber: 2,
		Message:    errors.New("a message"),
	}
	differentMessageError := ValidationError{
		LineNumber: -1,
		Message:    errors.New("different message"),
	}
	differentCountError := ValidationError{
		LineNumber:                    -1,
		Message:                       errors.New("a message"),
		AdditionalIdenticalErrorCount: 1,
	}
	if !baseError.Equal(baseError) {
		t.Error("failed to detect an error being equal to itself")
	}
	if baseError.Equal(wrongLineNumberError) {
		t.Error("failed to detect a difference in the LineNumber")
	}
	if baseError.Equal(differentMessageError) {
		t.Error("failed to detect a difference in the Message")
	}
	if baseError.Equal(differentCountError) {
		t.Error("failed to detect a difference in the ConsequtiveCount")
	}
}

func TestConsolidateErrors(t *testing.T) {
	input := []ValidationError{
		// two messages that become one block
		{LineNumber: 1, Message: errors.New("message kind one")},
		{LineNumber: 2, Message: errors.New("message kind one")},
		// one message with a good line between it and the last bad line, but repeating the message
		{LineNumber: 4, Message: errors.New("message kind one")},
		// one message that breaks the continuousness of line 4 to line 6
		{LineNumber: 5, Message: errors.New("message kind two")},
		{LineNumber: 6, Message: errors.New("message kind one")},
		// one message without a line number, that will become sorted to the top
		{LineNumber: -1, Message: errors.New("file-level error")},
	}

	expected := []ValidationError{
		{LineNumber: -1, Message: errors.New("file-level error")},
		{LineNumber: 1, AdditionalIdenticalErrorCount: 1, Message: errors.New("message kind one")},
		{LineNumber: 4, Message: errors.New("message kind one")},
		{LineNumber: 5, Message: errors.New("message kind two")},
		{LineNumber: 6, Message: errors.New("message kind one")},
	}

	actual := ConsolidateErrors(input, *config.NewConfig(nil))

	if !slices.EqualFunc(expected, actual, func(e1 ValidationError, e2 ValidationError) bool { return e1.Equal(e2) }) {
		t.Log("consolidation expectation          :", expected)
		t.Log("consolidation actual returned value:", actual)
		t.Error("returned list of validation errors deviated from the expected set")
	}
}

func TestConsolidateErrorCounts(t *testing.T) {
	input := []ValidationError{
		// A block of errors, followed by a small number of errors
		{LineNumber: 2, Message: errors.New("wrong indentation type")},
		{LineNumber: 3, Message: errors.New("wrong indentation type")},
		{LineNumber: 4, Message: errors.New("wrong indentation type")},
		{LineNumber: 5, Message: errors.New("wrong indentation type")},
		// one line gap
		{LineNumber: 7, Message: errors.New("wrong indentation type")},
		// one line gap
		{LineNumber: 9, Message: errors.New("wrong indentation type")},
		{LineNumber: 10, Message: errors.New("wrong indentation type")},
		{LineNumber: 11, Message: errors.New("wrong indentation type")},
	}

	expected := []ValidationError{
		{LineNumber: 2, AdditionalIdenticalErrorCount: 3, Message: errors.New("wrong indentation type")},
		{LineNumber: 7, Message: errors.New("wrong indentation type")},
		{LineNumber: 9, AdditionalIdenticalErrorCount: 2, Message: errors.New("wrong indentation type")},
	}

	actual := ConsolidateErrors(input, *config.NewConfig(nil))

	if !slices.EqualFunc(expected, actual, func(e1 ValidationError, e2 ValidationError) bool { return e1.Equal(e2) }) {
		t.Log("consolidation expectation          :", expected)
		t.Log("consolidation actual returned value:", actual)
		t.Error("returned list of validation errors deviated from the expected set")
	}
}

func TestConsolidatingInterleavedErrors(t *testing.T) {
	t.Skip("Consolidating non-consecutive errors is not supported by the current implementation")
	/*
		an assumption made about the possible future implement:
		it is implied that the implementation will sort the error messages by their error message
		If the implementation does not sort, this test will randomly fail when the implementation uses a map
	*/
	input := []ValidationError{
		{LineNumber: 1, Message: errors.New("message kind 2")},

		{LineNumber: 2, Message: errors.New("message kind 1")},
		{LineNumber: 2, Message: errors.New("message kind 2")},

		{LineNumber: 3, Message: errors.New("message kind 1")},
		{LineNumber: 3, Message: errors.New("message kind 2")},
		{LineNumber: 3, Message: errors.New("message kind 3")},

		{LineNumber: 4, Message: errors.New("message kind 4")},

		{LineNumber: 5, Message: errors.New("message kind 1")},

		{LineNumber: -1, Message: errors.New("file-level error")},
	}

	expected := []ValidationError{
		{LineNumber: -1, Message: errors.New("file-level error")},
		{LineNumber: 1, AdditionalIdenticalErrorCount: 2, Message: errors.New("message kind 2")},
		{LineNumber: 2, AdditionalIdenticalErrorCount: 1, Message: errors.New("message kind 1")},
		{LineNumber: 3, AdditionalIdenticalErrorCount: 0, Message: errors.New("message kind 3")},
		{LineNumber: 4, AdditionalIdenticalErrorCount: 1, Message: errors.New("message kind 4")},
		{LineNumber: 5, AdditionalIdenticalErrorCount: 0, Message: errors.New("message kind 1")},
	}

	actual := ConsolidateErrors(input, *config.NewConfig(nil))

	if !slices.EqualFunc(expected, actual, func(e1 ValidationError, e2 ValidationError) bool { return e1.Equal(e2) }) {
		t.Log("consolidation expectation          :", expected)
		t.Log("consolidation actual returned value:", actual)
		t.Error("returned list of validation errors deviated from the expected set")
	}
}

func TestFormatErrors(t *testing.T) {
	/*
		why change directory?
		The relative path conversion done by FormatErrors() changes how absolute paths
		are displayed, depending on in which directory the test is run in.

		We however still need to keep track of where we started, else we could to give snaps the path to the snapshots.
	*/
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not obtain current working directory: %s", err)
	}
	t.Chdir("/")

	/*
		why care about the path separators?

		The absolute paths in the testdata below come out differently, depending on the path separator
		used by the system that runs the tests.

		This leads to the problem, that one cannot verify a snapshot taken under Linux when using Windows
		and vice versa, unless the tests are differentiated by the path separator.
	*/
	var safePathSep string
	if os.PathSeparator == '/' {
		safePathSep = "slash"
	} else if os.PathSeparator == '\\' {
		safePathSep = "backslash"
	} else {
		t.Fatal("current path separator is unexpected - please fix test to handle this path separator")
	}
	s := snaps.WithConfig(
		snaps.Dir(filepath.Join(startingDir, "__snapshots__", "pathseparator-"+safePathSep)),
	)

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
		{
			FilePath: "some/file/with/consecutive/errors",
			Errors: []ValidationError{
				{LineNumber: 1, Message: errors.New("message kind one")},
				{LineNumber: 2, Message: errors.New("message kind one")},
				{LineNumber: 4, Message: errors.New("message kind one")},
				{LineNumber: 5, Message: errors.New("message kind two")},
				{LineNumber: 6, Message: errors.New("message kind one")},
				{LineNumber: -1, Message: errors.New("file-level error")},
			},
		},
	}

	for _, format := range outputformat.ValidOutputFormats {
		t.Run(string(format), func(t *testing.T) {
			buffer := bytes.Buffer{}
			config := config.NewConfig(nil)
			config.Format = format
			config.Logger.SetWriter(&buffer)
			PrintErrors(input, *config)
			s.MatchSnapshot(t, buffer.String())
		})
	}
}

func TestPrintErrorCount(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		errorcount int
	}{
		{"nonverbose-zero", false, 0},
		{"verbose-zero", true, 0},
		{"nonverbose-ten", false, 10},
		{"verbose-ten", true, 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			config := config.NewConfig(nil)
			config.Logger.VerboseEnabled = test.verbose
			config.Logger.SetWriter(&buffer)
			PrintErrorCount(test.errorcount, *config)
			snaps.MatchSnapshot(t, buffer.String())
		})
	}
}

func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}
