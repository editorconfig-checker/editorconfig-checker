// package config contains functions and structs related to config
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
	"github.com/editorconfig-checker/editorconfig-checker/pkg/utils"
)

// DefaultExcludes is the regular expression for ignored files
var DefaultExcludes = strings.Join(defaultExcludes, "|")

// defaultExcludes are an array to produce the correct string from
var defaultExcludes = []string{
	"yarn\\.lock$",
	"package-lock\\.json",
	"composer\\.lock$",
	"\\.snap$",
	"\\.otf$",
	"\\.woff$",
	"\\.woff2$",
	"\\.eot$",
	"\\.ttf$",
	"\\.gif$",
	"\\.png$",
	"\\.jpg$",
	"\\.jpeg$",
	"\\.mp4$",
	"\\.wmv$",
	"\\.svg$",
	"\\.ico$",
	"\\.bak$",
	"\\.bin$",
	"\\.pdf$",
	"\\.zip$",
	"\\.gz$",
	"\\.tar$",
	"\\.7z$",
	"\\.bz2$",
	"\\.log$",
	"\\.css\\.map$",
	"\\.js\\.map$",
	"min\\.css$",
	"min\\.js$",
}

var defaultAllowedContentTypes = []string{
	"text/",
	"application/octet-stream",
}

// Config struct, contains everything a config can contain
type Config struct {
	// CLI
	Version bool
	Help    bool
	DryRun  bool
	Path    string

	// CONFIG FILE
	Verbose               bool
	Debug                 bool
	Ignore_Defaults       bool
	Spaces_After_tabs     bool
	No_Color              bool
	Exclude               []string
	Allowed_Content_Types []string
	PassedFiles           []string
	Disable               DisabledChecks

	// MISC
	Logger logger.Logger
}

// DisabledChecks is a Struct which represents disabled checks
type DisabledChecks struct {
	End_Of_Line              bool
	Indentation              bool
	Insert_Final_Newline     bool
	Trim_Trailing_Whitespace bool
}

// NewConfig initializes a new config
func NewConfig(configPath string) (*Config, error) {
	var config Config
	config.Path = configPath

	if !utils.IsRegularFile(configPath) {
		return &config, fmt.Errorf("No file found at %s", configPath)
	}

	config.Exclude = []string{}
	config.Allowed_Content_Types = []string{}
	config.PassedFiles = []string{}

	return &config, nil
}

// Parse parses a config at a given path
func (c *Config) Parse() error {
	if c.Path != "" {
		configString, err := ioutil.ReadFile(c.Path)
		if err != nil {
			return err
		}

		err = json.Unmarshal(configString, c)
		if err != nil {
			return err
		}

		c.Allowed_Content_Types = append(defaultAllowedContentTypes, c.Allowed_Content_Types...)
	}

	return nil
}

// Merge merges a provided config with a config
func (c *Config) Merge(config Config) {
	if config.DryRun {
		c.DryRun = config.DryRun
	}

	if config.Version {
		c.Version = config.Version
	}

	if config.Help {
		c.Help = config.Help
	}

	if config.Verbose {
		c.Verbose = config.Verbose
	}

	if config.Debug {
		c.Debug = config.Debug
	}

	if config.Ignore_Defaults {
		c.Ignore_Defaults = config.Ignore_Defaults
	}

	if config.Spaces_After_tabs {
		c.Spaces_After_tabs = config.Spaces_After_tabs
	}

	if config.Path != "" {
		c.Path = config.Path
	}

	if len(config.Exclude) != 0 {
		c.Exclude = append(c.Exclude, config.Exclude...)
	}

	if len(config.Allowed_Content_Types) != 0 {
		c.Allowed_Content_Types = append(c.Allowed_Content_Types, config.Allowed_Content_Types...)
	}

	if len(config.PassedFiles) != 0 {
		c.PassedFiles = config.PassedFiles
	}

	if config.Disable.End_Of_Line {
		c.Disable.End_Of_Line = config.Disable.End_Of_Line
	}

	if config.Disable.Trim_Trailing_Whitespace {
		c.Disable.Trim_Trailing_Whitespace = config.Disable.Trim_Trailing_Whitespace
	}

	if config.Disable.Insert_Final_Newline {
		c.Disable.Insert_Final_Newline = config.Disable.Insert_Final_Newline
	}

	if config.Disable.Indentation {
		c.Disable.Indentation = config.Disable.Indentation
	}

	c.Logger = config.Logger
}

// GetExcludesAsRegularExpression returns the excludes as a combined regular expression
func (c Config) GetExcludesAsRegularExpression() string {
	if c.Ignore_Defaults {
		return strings.Join(c.Exclude, "|")
	}

	return strings.Join(append(c.Exclude, DefaultExcludes), "|")
}

// Save saves the config to it's Path
func (c Config) Save() error {
	if utils.IsRegularFile(c.Path) {
		return fmt.Errorf("File `%v` already exists", c.Path)
	}

	type writtenConfig struct {
		Verbose               bool
		Debug                 bool
		Ignore_Defaults       bool
		Spaces_After_tabs     bool
		No_Color              bool
		Exclude               []string
		Allowed_Content_Types []string
		PassedFiles           []string
		Disable               DisabledChecks
	}

	configJson, _ := json.MarshalIndent(writtenConfig{}, "", "  ")
	configString := strings.Replace(string(configJson[:]), "null", "[]", -1)
	err := ioutil.WriteFile(c.Path, []byte(configString), 0644)

	return err
}

// GetAsString returns the config in a readable form
func (c Config) GetAsString() string {
	return fmt.Sprintf("Config: %+v", c)
}
