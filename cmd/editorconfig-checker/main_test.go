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

func TestMain(m *testing.M) {
	exitProxy = captureReturnCode
	mainHasRun = make(chan int)

	os.Exit(m.Run())
}
