package main

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

var mainHasRun chan int

func captureReturnCode(code int) {
	mainHasRun <- code
	runtime.Goexit()
}

func setArguments(t *testing.T, args ...string) {
	t.Helper()

	initialArgs := os.Args
	t.Cleanup(func() {
		os.Args = initialArgs
	})

	// os.Args needs to have the binary as the first element
	os.Args = append([]string{"editorconfig-checker"}, args...)
}

func runWithArguments(t *testing.T, args ...string) (string, int) {
	t.Helper()
	setArguments(t, args...)
	outputBuffer := new(bytes.Buffer)
	loggerInjectionHook = func() {
		currentConfig.Logger.SetWriter(outputBuffer)
	}
	go main()
	exitCode := <-mainHasRun // must not be inlined into the return statement, since we need to wait for main() to exist before trying to read the buffer
	return outputBuffer.String(), exitCode
}

func TestMainOurCodebase(t *testing.T) {
	cdRelativeToRepo(t, "")
	output, lastSeenCode := runWithArguments(t, "--debug", "--verbose")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainMissingExplicitConfig(t *testing.T) {
	cdRelativeToRepo(t, "")
	output, lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--config", "/nonexistant")
	if lastSeenCode != exitCodeConfigFileNotFound {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeConfigFileNotFound)
		t.Logf("Output:\n%s", output)
	}
}
func TestMainWithFilesGiven(t *testing.T) {
	cdRelativeToRepo(t, "")
	output, lastSeenCode := runWithArguments(t, "--debug", "--verbose", "README.md")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainInitializingANewConfig(t *testing.T) {
	dir := t.TempDir()

	output, lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--init", "--config", dir+"/testwriteconfig.json")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}

	// do the same test a second time, trying to "overwrite" the existing file
	output, lastSeenCode = runWithArguments(t, "--debug", "--verbose", "--init", "--config", dir+"/testwriteconfig.json")
	// but now we expect it to fail since it does not want to overwrite the existing file
	if lastSeenCode != exitCodeErrorOccurred {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeErrorOccurred)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainLoadingAncientConfig(t *testing.T) {
	output, lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--config", "testdata/ancient-config.json")
	if lastSeenCode != exitCodeErrorOccurred {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeErrorOccurred)
		t.Logf("Output:\n%s", output)
	}
	if !strings.Contains(output, "SpacesAftertabs") {
		t.Errorf("main did not produce a warning that SpacesAftertabs is deprecated\nOutput:\n%s", output)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainWithEcrc(t *testing.T) {
	// feed a symlink named .ecrc pointing to our actual .editorconfig-checker.json
	output, lastSeenCode := runWithArguments(t, "--config", "testdata/.ecrc")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}
	if !strings.Contains(output, "`.ecrc` is deprecated") {
		t.Error("main did not produce a warning that .ecrc is deprecated despite being give a file named .ecrc.")
		t.Logf("Output:\n%s", output)
	}
}

func TestMainShowVersion(t *testing.T) {
	output, lastSeenCode := runWithArguments(t, "--version")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainShowHelp(t *testing.T) {
	output, lastSeenCode := runWithArguments(t, "--help")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainDryRun(t *testing.T) {
	cdRelativeToRepo(t, "")
	output, lastSeenCode := runWithArguments(t, "--dry-run")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
		t.Logf("Output:\n%s", output)
	}
}

func TestMainColorSupport(t *testing.T) {
	type env map[string]string
	type args []string

	tests := []struct {
		name string
		env  env
		args args
	}{
		{"no-envvar-no-arg", env{}, args{}},
		{"envvar-no-arg", env{"NO_COLOR": "1"}, args{}},
		{"no-envvar-color-off", env{}, args{"--no-color"}},
		{"no-envvar-color-on", env{}, args{"--color"}},
		// test that arguments win over env var
		{"envvar-color-off", env{"NO_COLOR": "1"}, args{"--no-color"}},
		{"envvar-color-on", env{"NO_COLOR": "1"}, args{"--color"}},
		// test that the last argument wins
		{"no-envvar-color-offon", env{}, args{"--no-color", "--color"}},
		{"no-envvar-color-onoffon", env{}, args{"--color", "--no-color", "--color"}},
		// test that boolean values in the env var work correctly
		{"envvar-true", env{"NO_COLOR": "true"}, args{}},
		{"envvar-false", env{"NO_COLOR": "false"}, args{}},
		{"envvar-zero", env{"NO_COLOR": "0"}, args{}},
		// test that pure string values are correctly handled as setting NoColor
		{"envvar-yes", env{"NO_COLOR": "yes"}, args{}},
		{"envvar-no", env{"NO_COLOR": "no"}, args{}},
		{"envvar-stringval", env{"NO_COLOR": "some test value that nobody would set"}, args{}},
	}

	// we use the error message of a missing config file to test the coloredness of the output
	defaultArgs := []string{
		`--exclude=""`, "--ignore-defaults",
		"testdata/trailing-whitespace.txt",
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for name, value := range test.env {
				t.Setenv(name, value)
			}
			args := append(test.args, defaultArgs...)
			output, _ := runWithArguments(t, args...)
			snaps.MatchSnapshot(t, output)
		})
	}
}

// a little Helper to set the current working dir relative to the repository root,
// and return to the previous working directory once the test completes
func cdRelativeToRepo(t *testing.T, path string) {
	t.Helper()
	newdir := "../../" + path
	t.Chdir(newdir)
}

func TestReturnCodeInterface(t *testing.T) {
	// These constants are possibly used by external processes, so we must not change them
	if exitCodeNormal != 0 {
		t.Errorf("Return code for a normal condition was %d, but we expected 0", exitCodeNormal)
	}
	if exitCodeErrorOccurred != 1 {
		t.Errorf("Return code for an error condition was %d, but we expected 1", exitCodeErrorOccurred)
	}
	if exitCodeConfigFileNotFound != 2 {
		t.Errorf("Return code for a nonexistant config file was %d, but we expected 2", exitCodeConfigFileNotFound)
	}
}

func TestMain(m *testing.M) {
	exitProxy = captureReturnCode
	mainHasRun = make(chan int)

	os.Exit(m.Run())
}
