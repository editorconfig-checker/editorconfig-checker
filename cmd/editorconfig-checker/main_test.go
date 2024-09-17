package main

import (
	"os"
	"runtime"
	"testing"
)

var mainHasRun chan int

func captureReturncode(code int) {
	mainHasRun <- code
	runtime.Goexit()
}

func setArguments(args ...string) {
	// os.Args needs to have the binary
	os.Args = append([]string{"editorconfig-checker"}, args...)
}

func TestMainOurCodebase(t *testing.T) {
	cdRelativeToRepo(t, "")
	setArguments("--debug", "--verbose", "--exclude", "\\.git", "--exclude", "\\.exe$")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
	}
}

func TestMainMissingExplicitConfig(t *testing.T) {
	cdRelativeToRepo(t, "")
	setArguments("--debug", "--verbose", "--exclude", "\\.git", "--exclude", "\\.exe$", "--config", "/nonexistant")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 2 {
		t.Errorf("main exited with return code %d, but we expected 2", lastSeenCode)
	}
}
func TestMainWithFilesGiven(t *testing.T) {
	cdRelativeToRepo(t, "")
	setArguments("--debug", "--verbose", "README.md")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
	}
}

func TestMainInitializingANewConfig(t *testing.T) {
	dir := t.TempDir()
	setArguments("--debug", "--verbose", "--init", "--config", dir+"/testwriteconfig.json")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
	}

	// do the same test a second time, trying to "overwrite" the existing file
	setArguments("--debug", "--verbose", "--init", "--config", dir+"/testwriteconfig.json")
	go main()
	lastSeenCode = <-mainHasRun
	// but now we expect it to fail since it does not want to overwrite the existing file
	if lastSeenCode != 1 {
		t.Errorf("main exited with return code %d, but we expected 1", lastSeenCode)
	}
}

func TestMainLoadingAncientConfig(t *testing.T) {
	setArguments("--debug", "--verbose", "--config", "testdata/ancient-config.json")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 1 {
		t.Errorf("main exited with return code %d, but we expected 1", lastSeenCode)
	}
}

func TestMainShowVersion(t *testing.T) {
	setArguments("--version")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
	}
}

func TestMainShowHelp(t *testing.T) {
	setArguments("--help")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
	}
}

func TestMainDryRun(t *testing.T) {
	cdRelativeToRepo(t, "")
	setArguments("--dry-run")
	go main()
	lastSeenCode := <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
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

func TestMain(m *testing.M) {
	exitStub = captureReturncode
	mainHasRun = make(chan int)

	os.Exit(m.Run())
}
