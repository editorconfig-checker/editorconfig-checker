package slice

import (
	"testing"
)

func TestContains(t *testing.T) {
	testSlice := []string{"a", "b", "c"}

	if !Contains(testSlice, "a") {
		t.Error("Expected \"a\" to be part of the slice")
	}

	if Contains(testSlice, "z") {
		t.Error("Expected \"z\" to be not part of the slice")
	}
}
