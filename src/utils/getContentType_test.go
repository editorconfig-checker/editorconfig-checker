package utils

import (
	"strings"
	"testing"
)

func TestGetContentType(t *testing.T) {
	contentType := GetContentType("./getContentType.go")
	if !strings.Contains(contentType, "text/plain") {
		t.Error("Expected getContentType.go to be of type text/plain")
	}

	contentType = GetContentType("./getContentType_test.go")
	if !strings.Contains(contentType, "application/octet-stream") {
		t.Error("Expected getContentType_test.go to be of type application/octet-stream")
	}
}
