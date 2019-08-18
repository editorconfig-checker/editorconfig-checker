package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/config"
)

func TestReturnableFlags(t *testing.T) {
	var result bool
	var expected bool

	config := config.Config{Help: true}
	result = ReturnableFlags(config)
	expected = true

	if result != expected {
		t.Errorf("Should exit the program, got %v", result)
	}

	config.Version = true
	result = ReturnableFlags(config)
	expected = true

	if result != expected {
		t.Errorf("Should exit the program, got %v", result)
	}

	config.Version = false
	config.Help = false

	result = ReturnableFlags(config)
	expected = false

	if result != expected {
		t.Errorf("Should not exit the program, got %v", result)
	}
}

func TestReturnableFlagsExitValue(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "-version")
	err := cmd.Run()
	if err != nil {
		t.Errorf("process ran with err %v, want exit status 0", err)
	}

	cmd = exec.Command("go", "run", "./main.go", "-help")
	err = cmd.Run()
	if err != nil {
		t.Errorf("process ran with err %v, want exit status 0", err)
	}

	cmd = exec.Command("go", "run", "./main.go", "-init", "-config", "stuff.json")
	err = cmd.Run()
	os.Remove("stuff.json")
	if err != nil {
		t.Errorf("process ran with err %v, want exit status 0", err)
	}
}
