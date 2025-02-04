// Package for having structured access to our output formats
package outputformat

import (
	"fmt"
	"slices"
	"strings"
)

type OutputFormat string

const (
	Default       = OutputFormat("default")
	Codeclimate   = OutputFormat("codeclimate")
	GCC           = OutputFormat("gcc")
	GithubActions = OutputFormat("github-actions")
)

var ValidOutputFormats = []OutputFormat{
	Default,
	Codeclimate,
	GCC,
	GithubActions,
}

func GetArgumentChoiceText() string {
	var outputStrings []string
	for _, f := range ValidOutputFormats {
		outputStrings = append(outputStrings, string(f))
	}
	return strings.Join(outputStrings, ", ")
}

func (format OutputFormat) MarshalText() ([]byte, error) {
	if !format.IsValid() {
		return nil, fmt.Errorf("%q is not a valid output format", format)
	}
	return []byte(format), nil
}

func (format *OutputFormat) UnmarshalText(data []byte) error {
	*format = OutputFormat(string(data))
	if !format.IsValid() {
		return fmt.Errorf("%q is not a valid output format", data)
	}
	return nil
}

func (format OutputFormat) IsValid() bool {
	return slices.Contains(ValidOutputFormats, format)
}

func (format OutputFormat) String() string {
	return string(format)
}
