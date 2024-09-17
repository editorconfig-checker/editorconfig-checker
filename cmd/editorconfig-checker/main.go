// Package main provides ...
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/error"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/files"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/outputformat"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/utils"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/validation"
)

// version is used vor the help
// version is dynamically set at compiletime
var version string

// defaultConfigFileNames determines the file names where the config is located
var defaultConfigFileNames = []string{".editorconfig-checker.json", ".ecrc"}

// currentConfig is the config used in this run
var currentConfig *config.Config

// exitProxy is there to be replaced while running the tests
var exitProxy = os.Exit

// these must be globals, since they are referenced by init(), parseArguments
var configFilePath string
var cmdlineExclude string
var cmdlineConfig config.Config
var writeConfigFile bool

func init() {
	flag.BoolVar(&writeConfigFile, "init", false, "creates an initial configuration")
	flag.StringVar(&configFilePath, "config", "", "config")
	flag.StringVar(&cmdlineExclude, "exclude", "", "a regex which files should be excluded from checking - needs to be a valid regular expression")
	flag.BoolVar(&cmdlineConfig.IgnoreDefaults, "ignore-defaults", false, "ignore default excludes")
	flag.BoolVar(&cmdlineConfig.DryRun, "dry-run", false, "show which files would be checked")
	flag.BoolVar(&cmdlineConfig.ShowVersion, "version", false, "print the version number")
	flag.BoolVar(&cmdlineConfig.Help, "help", false, "print the help")
	flag.BoolVar(&cmdlineConfig.Help, "h", false, "print the help")
	flag.TextVar(&cmdlineConfig.Format, "format", outputformat.Default, "specify the output format: "+outputformat.GetArgumentChoiceText())
	flag.TextVar(&cmdlineConfig.Format, "f", outputformat.Default, "specify the output format: "+outputformat.GetArgumentChoiceText())
	flag.BoolVar(&cmdlineConfig.Verbose, "verbose", false, "print debugging information")
	flag.BoolVar(&cmdlineConfig.Verbose, "v", false, "print debugging information")
	flag.BoolVar(&cmdlineConfig.Debug, "debug", false, "print debugging information")
	flag.BoolVar(&cmdlineConfig.NoColor, "no-color", false, "dont print colors")
	flag.BoolVar(&cmdlineConfig.Disable.TrimTrailingWhitespace, "disable-trim-trailing-whitespace", false, "disables the trailing whitespace check")
	flag.BoolVar(&cmdlineConfig.Disable.EndOfLine, "disable-end-of-line", false, "disables the trailing whitespace check")
	flag.BoolVar(&cmdlineConfig.Disable.InsertFinalNewline, "disable-insert-final-newline", false, "disables the final newline check")
	flag.BoolVar(&cmdlineConfig.Disable.Indentation, "disable-indentation", false, "disables the indentation check")
	flag.BoolVar(&cmdlineConfig.Disable.IndentSize, "disable-indent-size", false, "disables only the indent-size check")
	flag.BoolVar(&cmdlineConfig.Disable.MaxLineLength, "disable-max-line-length", false, "disables only the max-line-length check")
}

// parse the arguments from os.Args
func parseArguments() {
	// reset the global variables used to receive the arguments, so parseArguments can be called multiple times without reusing arguments from the previous run
	configFilePath = ""
	cmdlineExclude = ""
	cmdlineConfig = config.Config{}
	writeConfigFile = false

	flag.Parse()

	var configPaths = []string{}
	if configFilePath == "" {
		configPaths = append(configPaths, defaultConfigFileNames[:]...)
	} else {
		configPaths = append(configPaths, configFilePath)
	}

	currentConfig, _ = config.NewConfig(configPaths)

	if strings.HasSuffix(currentConfig.Path, ".ecrc") {
		currentConfig.Logger.Warning("The default configuration file name `.ecrc` is deprecated. Use `.editorconfig-checker.json` instead. You can simply rename it")
	}

	if writeConfigFile {
		err := currentConfig.Save(version)
		if err != nil {
			currentConfig.Logger.Error(err.Error())
			exitProxy(1)
		}

		exitProxy(0)
	}

	err := currentConfig.Parse()
	// this error should be surpressed if the configFilePath was not set by the user
	// since the default config paths could trigger this
	if err != nil && !(configFilePath == "" && errors.Is(err, fs.ErrNotExist)) {
		currentConfig.Logger.Error(err.Error())
		exitProxy(2)
	}

	if cmdlineExclude != "" {
		cmdlineConfig.Exclude = append(cmdlineConfig.Exclude, cmdlineExclude)
	}

	// Some wrapping tools pass an empty string as arguments so
	// our file searching algorithm will break because it thinks there are
	// empty files and will cause the program to crash
	for _, arg := range flag.Args() {
		if arg != "" {
			cmdlineConfig.PassedFiles = append(cmdlineConfig.PassedFiles, arg)
		}
	}

	currentConfig.Merge(cmdlineConfig)
}

// Main function, dude
func main() {
	parseArguments()

	config := *currentConfig
	config.Logger.Debug("Config: %s", config)
	config.Logger.Verbose("Exclude Regexp: %s", config.GetExcludesAsRegularExpression())

	if utils.FileExists(config.Path) && config.Version != "" && config.Version != version {
		config.Logger.Error("Version from config file is not the same as the version of the binary")
		config.Logger.Error(fmt.Sprintf("Binary: %s, Config %s", version, config.Version))
		exitProxy(1)
	}

	// Check for returnworthy arguments
	shouldExit := ReturnableFlags(config)
	if shouldExit {
		exitProxy(0)
	}

	// contains all files which should be checked
	filePaths, err := files.GetFiles(config)

	if err != nil {
		config.Logger.Error(err.Error())
		exitProxy(1)
	}

	if config.DryRun {
		for _, file := range filePaths {
			config.Logger.Output(file)
		}

		exitProxy(0)
	}

	errors := validation.ProcessValidation(filePaths, config)

	error.PrintErrors(errors, config)

	config.Logger.Verbose("%d files checked", len(filePaths))

	if error.GetErrorCount(errors) != 0 {
		exitProxy(1)
	}

	exitProxy(0)
}

// ReturnableFlags returns whether a flag passed should exit the program
func ReturnableFlags(config config.Config) bool {
	switch {
	case config.ShowVersion:
		config.Logger.Output(version)
	case config.Help:
		config.Logger.Output("USAGE:")
		flag.PrintDefaults()
	}

	return config.ShowVersion || config.Help
}
