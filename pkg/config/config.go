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
	"min\\.js$"}

type Config struct {
	// CLI
	Version    bool
	Help       bool
	DryRun     bool
	ConfigPath string

	// CONFIG FILE
	Verbose           bool
	Debug             bool
	Ignore_Defaults   bool
	Spaces_After_tabs bool
	Exclude           []string
	PassedFiles       []string
	Disable           DisabledChecks

	// MISC
	Logger logger.Logger
}

// DisabledChecks is a Struct which represents disabled checks
type DisabledChecks struct {
	Trim_Trailing_Whitespace bool
	End_Of_Line              bool
	Insert_Final_Newline     bool
	Indentation              bool
}

func NewConfig(configPath string) (*Config, error) {
	var config Config

	if !utils.IsRegularFile(configPath) {
		return &config, fmt.Errorf("No file found at %s", configPath)
	}

	config.ConfigPath = configPath
	return &config, nil
}

func (c *Config) Parse() error {
	if c.ConfigPath != "" {
		configString, err := ioutil.ReadFile(c.ConfigPath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(configString, c)
		if err != nil {
			return err
		}

		if !c.Ignore_Defaults {
			c.Exclude = append(c.Exclude, defaultExcludes...)
		}

		if c.Debug {
			// TODO Print Config
			logger.Output("")
		}
	}

	return nil
}

func (c *Config) Merge(config Config) error {
	return nil
}

func (c Config) GetExcludesAsRegularExpression() string {
	return strings.Join(c.Exclude, "|")
}
