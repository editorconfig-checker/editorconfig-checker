package error

import (
	"crypto/md5"
	"fmt"
)

// CodeclimageLines represents the lines of an issue in codeclimate format
type CodeclimateLines struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}

// CodeclimateLocation represents the location of an issue in codeclimate format
type CodeclimateLocation struct {
	Path  string           `json:"path"`
	Lines CodeclimateLines `json:"lines"`
}

// CodeclimateIssue represents an issue in codeclimate format
type CodeclimateIssue struct {
	Check       string              `json:"check_name"`
	Description string              `json:"description"`
	Fingerprint string              `json:"fingerprint"`
	Severity    string              `json:"severity"`
	Location    CodeclimateLocation `json:"location"`
}

const (
	checkName = "editorconfig-checker"
	severity  = "minor"
)

func newCodeclimateIssue(err ValidationError, path string) CodeclimateIssue {
	toHash := fmt.Sprintf("%s:%d:%d:%s", path, err.LineNumber, err.AdditionalIdenticalErrorCount, err.Message.Error())
	fingerprint := fmt.Sprintf("%x", md5.Sum([]byte(toHash)))
	return CodeclimateIssue{
		Check:       checkName,
		Description: err.Message.Error(),
		Fingerprint: fingerprint,
		Severity:    severity,
		Location: CodeclimateLocation{
			Path: path,
			Lines: CodeclimateLines{
				Begin: err.LineNumber,
				End:   err.AdditionalIdenticalErrorCount,
			},
		},
	}
}
