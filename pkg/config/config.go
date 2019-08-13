package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/logger"
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
	Version bool
	Help    bool
	DryRun  bool

	// UNDECIDED
	RawFiles   []string
	ConfigPath string
	Debug      bool
	Logger     logger.Logger

	// CONFIG FILE
	Verbose           bool
	Ignore_Defaults   bool
	Spaces_After_tabs bool
	Exclude           []string
	Disable           DisabledChecks
}

// DisabledChecks is a Struct which represents disabled checks
type DisabledChecks struct {
	Trim_Trailing_Whitespace bool
	End_Of_Line              bool
	Insert_Final_Newline     bool
	Indentation              bool
}

// TODO: EXTRACT THAT
func isRegularFile(filePath string) bool {
	absolutePath, _ := filepath.Abs(filePath)
	fi, err := os.Stat(absolutePath)

	return err == nil && fi.Mode().IsRegular()
}

func NewConfig(configPath string) (*Config, error) {
	var config Config

	if !isRegularFile(configPath) {
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
