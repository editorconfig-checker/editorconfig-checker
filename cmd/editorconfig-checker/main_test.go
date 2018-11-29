package main

import (
	"testing"
)

func BenchmarkMain(b *testing.B) {
	// run the main function b.N times
	for n := 0; n < b.N; n++ {
		main()
	}
}
