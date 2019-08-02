package utils

import (
	"testing"
)

func TestDefaultExcludes(t *testing.T) {
	result := "yarn\\.lock$|package-lock\\.json|composer\\.lock$|\\.snap$|\\.otf$|\\.woff$|\\.woff2$|\\.eot$|\\.ttf$|\\.gif$|\\.png$|\\.jpg$|\\.jpeg$|\\.mp4$|\\.wmv$|\\.svg$|\\.ico$|\\.bak$|\\.bin$|\\.pdf$|\\.zip$|\\.gz$|\\.tar$|\\.7z$|\\.bz2$|\\.log$|\\.css\\.map$|\\.js\\.map$|min\\.css$|min\\.js$"

	if DefaultExcludes != result {
		t.Error("Expected default excludes to match", result)
	}
}

func TestGetEolChar(t *testing.T) {
	if GetEolChar("lf") != "\n" {
		t.Error("Expected end of line character to be \\n for \"lf\"")
	}

	if GetEolChar("cr") != "\r" {
		t.Error("Expected end of line character to be \\r for \"cr\"")
	}

	if GetEolChar("crlf") != "\r\n" {
		t.Error("Expected end of line character to be \\r\\n for \"crlf\"")
	}

	if GetEolChar("") != "\n" {
		t.Error("Expected end of line character to be \\n as a fallback")
	}
}
