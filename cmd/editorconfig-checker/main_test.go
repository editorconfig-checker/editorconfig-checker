package main

import (
	"os"
	"os/exec"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	// run the binary b.N times
	for n := 0; n < b.N; n++ {
		dir, _ := os.Getwd()
		cmd := exec.Command("make", "run")
		// the test is executed where the `*_test.go` file is located
		cmd.Dir = dir + "/../../"
		cmd.Run()
	}
}
