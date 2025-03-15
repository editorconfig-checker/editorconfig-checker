// Package for having structured access to our output formats
package outputformat

import (
	"cmp"
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
	var output_strings []string
	for _, f := range ValidOutputFormats {
		output_strings = append(output_strings, string(f))
	}
	return strings.Join(output_strings, ", ")
}

func (format OutputFormat) MarshalText() ([]byte, error) {
	if !format.IsValid() {
		return nil, fmt.Errorf("%q is not a valid output format", format)
	}
	return []byte(format), nil
}

func (format *OutputFormat) UnmarshalText(data []byte) error {
	*format = OutputFormat(cmp.Or(string(data), "default"))
	if !format.IsValid() {
		return fmt.Errorf("%q is not a valid output format", data)
	}
	return nil
}

func (format OutputFormat) IsValid() bool {
	return slices.Contains(ValidOutputFormats, format)
}

func (f OutputFormat) String() string {
	return string(f)
}
