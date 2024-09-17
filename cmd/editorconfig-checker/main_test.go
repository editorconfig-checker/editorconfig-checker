package main

import (
	"os"
	"runtime"
	"testing"
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

func runWithArguments(t *testing.T, args ...string) int {
	t.Helper()
	setArguments(t, args...)
	go main()
	return <-mainHasRun
}

func TestMainOurCodebase(t *testing.T) {
	cdRelativeToRepo(t, "")
	lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--exclude", `\.git`, "--exclude", `\.exe$`)
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
	}
}

func TestMainMissingExplicitConfig(t *testing.T) {
	cdRelativeToRepo(t, "")
	lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--exclude", `\.git`, "--exclude", `\.exe$`, "--config", "/nonexistant")
	if lastSeenCode != exitCodeConfigFileNotFound {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeConfigFileNotFound)
	}
}
func TestMainWithFilesGiven(t *testing.T) {
	cdRelativeToRepo(t, "")
	lastSeenCode := runWithArguments(t, "--debug", "--verbose", "README.md")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
	}
}

func TestMainInitializingANewConfig(t *testing.T) {
	dir := t.TempDir()

	lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--init", "--config", dir+"/testwriteconfig.json")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
	}

	// do the same test a second time, trying to "overwrite" the existing file
	lastSeenCode = runWithArguments(t, "--debug", "--verbose", "--init", "--config", dir+"/testwriteconfig.json")
	// but now we expect it to fail since it does not want to overwrite the existing file
	if lastSeenCode != exitCodeErrorOccurred {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeErrorOccurred)
	}
}

func TestMainLoadingAncientConfig(t *testing.T) {
	lastSeenCode := runWithArguments(t, "--debug", "--verbose", "--config", "testdata/ancient-config.json")
	if lastSeenCode != exitCodeErrorOccurred {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeErrorOccurred)
	}
}

func TestMainShowVersion(t *testing.T) {
	lastSeenCode := runWithArguments(t, "--version")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
	}
}

func TestMainShowHelp(t *testing.T) {
	lastSeenCode := runWithArguments(t, "--help")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
	}
}

func TestMainDryRun(t *testing.T) {
	cdRelativeToRepo(t, "")
	lastSeenCode := runWithArguments(t, "--dry-run")
	if lastSeenCode != exitCodeNormal {
		t.Errorf("main exited with return code %d, but we expected %d", lastSeenCode, exitCodeNormal)
	}
}

// a little Helper to set the current working dir relative to the repository root,
// and return to the previous working directory once the test completes
func cdRelativeToRepo(t *testing.T, path string) {
	newdir := "../../" + path

	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not obtain current working directory: %s", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(startingDir); err != nil {
			t.Fatalf("Could not restore old working directory %s: %s", startingDir, err)
		}
	})
	if err := os.Chdir(newdir); err != nil {
		t.Fatalf("Could not chdir to %s: %s", newdir, err)
	}
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
