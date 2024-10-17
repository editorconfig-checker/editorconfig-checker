// Package config contains functions and structs related to config
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/editorconfig/editorconfig-core-go/v2"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/logger"       // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/outputformat" // x-release-please-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/utils"        // x-release-please-major
)

// DefaultExcludes is the regular expression for ignored files
var DefaultExcludes = strings.Join(defaultExcludes, "|")

// defaultExcludes are an array to produce the correct string from
var defaultExcludes = []string{
	"^\\.yarn/",
	"^yarn\\.lock$",
	"^package-lock\\.json$",
	"^composer\\.lock$",
	"^Cargo\\.lock$",
	"^Gemfile\\.lock$",
	"^\\.pnp\\.cjs$",
	"^\\.pnp\\.js$",
	"^\\.pnp\\.loader\\.mjs$",
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
	"\\.webp$",
	"\\.avif$",
	"\\.pnm$",
	"\\.pbm$",
	"\\.pgm$",
	"\\.ppm$",
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
	"\\.patch$",
	"\\.css\\.map$",
	"\\.js\\.map$",
	"min\\.css$",
	"min\\.js$",
}

// keep synced with pkg/validation/validation.go#L20
var defaultAllowedContentTypes = []string{
	"text/",
	"application/octet-stream",
	"application/ecmascript",
	"application/json",
	"application/x-ndjson",
	"application/xml",
	"+json",
	"+xml",
}

// Config struct, contains everything a config can contain
type Config struct {
	// CLI
	ShowVersion bool
	Help        bool
	DryRun      bool
	Path        string

	// CONFIG FILE
	Version             string
	Verbose             bool
	Format              outputformat.OutputFormat
	Debug               bool
	IgnoreDefaults      bool
	SpacesAftertabs     *bool
	SpacesAfterTabs     bool
	NoColor             bool
	Exclude             []string
	AllowedContentTypes []string
	PassedFiles         []string
	Disable             DisabledChecks

	// MISC
	Logger             logger.Logger
	EditorconfigConfig *editorconfig.Config
}

// DisabledChecks is a Struct which represents disabled checks
type DisabledChecks struct {
	EndOfLine              bool
	Indentation            bool
	InsertFinalNewline     bool
	TrimTrailingWhitespace bool
	IndentSize             bool
	MaxLineLength          bool
}

// NewConfig initializes a new config
func NewConfig(configPaths []string) (*Config, error) {
	if len(configPaths) == 0 {
		return nil, fmt.Errorf("No config paths provided")
	}

	var config Config

	config.AllowedContentTypes = defaultAllowedContentTypes
	config.Exclude = []string{}
	config.PassedFiles = []string{}

	config.EditorconfigConfig = &editorconfig.Config{
		Parser: editorconfig.NewCachedParser(),
	}

	var configPath string = ""
	for _, path := range configPaths {
		if utils.IsRegularFile(path) {
			configPath = path
			break
		}
	}
	if configPath == "" {
		var configPath string = configPaths[0]
		config.Path = configPath
		return &config, fmt.Errorf("No file found at %s", configPath)
	}
	config.Path = configPath

	return &config, nil
}

// Parse parses a config at a given path
func (c *Config) Parse() error {
	if c.Path != "" {
		configString, err := os.ReadFile(c.Path)
		if err != nil {
			return err
		}

		tmpConfg := Config{}
		err = json.Unmarshal(configString, &tmpConfg)
		if err != nil {
			return err
		}

		c.Merge(tmpConfg)
	}

	return nil
}

// Merge merges a provided config with a config
func (c *Config) Merge(config Config) {
	if config.DryRun {
		c.DryRun = config.DryRun
	}

	if config.ShowVersion {
		c.ShowVersion = config.ShowVersion
	}

	if len(config.Version) > 0 {
		c.Version = config.Version
	}

	if config.Help {
		c.Help = config.Help
	}

	if config.Verbose {
		c.Verbose = config.Verbose
	}

	if config.Format.IsValid() {
		c.Format = config.Format
	}

	if config.Debug {
		c.Debug = config.Debug
	}

	if config.NoColor {
		c.NoColor = config.NoColor
	}

	if config.IgnoreDefaults {
		c.IgnoreDefaults = config.IgnoreDefaults
	}

	if config.SpacesAftertabs != nil {
		c.Logger.Warning("The configuration key `SpacesAftertabs` is deprecated. Use `SpacesAfterTabs` instead.")

		c.SpacesAfterTabs = *config.SpacesAftertabs
	}

	if config.SpacesAfterTabs {
		c.SpacesAfterTabs = config.SpacesAfterTabs
	}

	if config.Path != "" {
		c.Path = config.Path
	}

	if len(config.Exclude) != 0 {
		c.Exclude = append(c.Exclude, config.Exclude...)
	}

	if len(config.AllowedContentTypes) != 0 {
		c.AllowedContentTypes = append(c.AllowedContentTypes, config.AllowedContentTypes...)
	}

	if len(config.PassedFiles) != 0 {
		c.PassedFiles = config.PassedFiles
	}

	c.mergeDisabled(config.Disable)
	c.Logger = logger.Logger{
		Verbosee: c.Verbose || config.Verbose,
		Debugg:   c.Debug || config.Debug,
		NoColor:  c.NoColor || config.NoColor,
	}
}

// mergeDisabled merges the disabled checks into the config
// This is here because cyclomatic complexity of gocyclo was about 15 :/
func (c *Config) mergeDisabled(disabled DisabledChecks) {
	if disabled.EndOfLine {
		c.Disable.EndOfLine = disabled.EndOfLine
	}

	if disabled.TrimTrailingWhitespace {
		c.Disable.TrimTrailingWhitespace = disabled.TrimTrailingWhitespace
	}

	if disabled.InsertFinalNewline {
		c.Disable.InsertFinalNewline = disabled.InsertFinalNewline
	}

	if disabled.Indentation {
		c.Disable.Indentation = disabled.Indentation
	}

	if disabled.IndentSize {
		c.Disable.IndentSize = disabled.IndentSize
	}

	if disabled.MaxLineLength {
		c.Disable.MaxLineLength = disabled.MaxLineLength
	}
}

// GetExcludesAsRegularExpression returns the excludes as a combined regular expression
func (c Config) GetExcludesAsRegularExpression() string {
	if c.IgnoreDefaults {
		return strings.Join(c.Exclude, "|")
	}

	return strings.Join(append(c.Exclude, DefaultExcludes), "|")
}

// Save saves the config to it's Path
func (c Config) Save(version string) error {
	if utils.IsRegularFile(c.Path) {
		return fmt.Errorf("File `%v` already exists", c.Path)
	}

	type writtenConfig struct {
		Version             string
		Verbose             bool
		Format              string
		Debug               bool
		IgnoreDefaults      bool
		SpacesAfterTabs     bool
		NoColor             bool
		Exclude             []string
		AllowedContentTypes []string
		PassedFiles         []string
		Disable             DisabledChecks
	}

	configJSON, _ := json.MarshalIndent(writtenConfig{Version: version}, "", "  ")
	configString := strings.Replace(string(configJSON[:]), "null", "[]", -1)
	err := os.WriteFile(c.Path, []byte(configString), 0644)

	return err
}

// GetAsString returns the config in a readable form
func (c Config) GetAsString() string {
	return fmt.Sprintf("Config: %+v", c)
}
