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

func TestMainFunc(t *testing.T) {
	lastSeenCode := -1

	os.Args = []string{"--debug", "--verbose", "--exclude", "\\.git", "--exclude", "\\.exe$"}
	go main()
	lastSeenCode = <-mainHasRun
	if lastSeenCode != 0 {
		t.Errorf("main exited with return code %d, but we expected 0", lastSeenCode)
	}

	/*
		the following does not work yet, since flags can only be initialized once
		but keeping the flag parsing in an init func is not an option either, since it would not be executed the second time around
	*/ /*
		os.Args = []string{"--debug", "--verbose", "--exclude", "\\.git", "--exclude", "\\.exe$", "--config", "/nonexistant"}
		go main()
		lastSeenCode = <-mainHasRun
		t.Logf("Exit Code 1: %d", lastSeenCode)
	*/
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
