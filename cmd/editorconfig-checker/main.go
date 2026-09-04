// Package main provides ...
package main

import (
	"errors"
	"flag"
	"io/fs"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/gkampitakis/ciinfo"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/config"
	eccerror "github.com/editorconfig-checker/editorconfig-checker/v4/internal/error"
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/files"
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/fix"
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/outputformat"
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/utils"
	"github.com/editorconfig-checker/editorconfig-checker/v4/internal/validation"
	// x-release-please-end
)

// version is used for the help and to verify against the version stored in the config file
// version is dynamically set at compiletime
var version = "v4.0.1" // x-release-please-version

// defaultConfigFileNames determines the file names where the config is located
var defaultConfigFileNames = []string{".editorconfig-checker.json"}

// currentConfig is the config used in this run
var currentConfig *config.Config

// exitProxy is there to be replaced while running the tests
var exitProxy = os.Exit

//  loggerInjectionHook is there to be replaced while running the tests
var loggerInjectionHook = func() {}

const (
	exitCodeNormal             = iota
	exitCodeErrorOccurred      = iota
	exitCodeConfigFileNotFound = iota
)

// these must be globals, since they are referenced by init(), parseArguments
var (
	configFilePath  string
	cmdlineExclude  string
	cmdlineConfig   config.Config
	writeConfigFile bool
	cpuprofile      string
)

var colorExplicitlySet bool

var detectOutputFormat = func() (outputformat.OutputFormat, bool) {
	if ciinfo.IsVendor("GITHUB_ACTIONS") {
		return outputformat.GithubActions, true
	}
	return outputformat.Default, false
}

func enableNoColor(string) error {
	cmdlineConfig.NoColor = true
	colorExplicitlySet = true
	return nil
}

func disableNoColor(string) error {
	cmdlineConfig.NoColor = false
	colorExplicitlySet = true
	return nil
}

func init() {
	flag.BoolVar(&writeConfigFile, "init", false, "creates an initial configuration")
	flag.StringVar(&configFilePath, "config", "", "config")
	flag.StringVar(&cmdlineExclude, "exclude", "", "a regex which files should be excluded from checking - needs to be a valid regular expression. Combine patterns with | (pipe): -exclude \"vendor|testdata\"")
	flag.BoolVar(&cmdlineConfig.IgnoreDefaults, "ignore-defaults", false, "ignore default excludes")
	flag.BoolVar(&cmdlineConfig.DryRun, "dry-run", false, "show which files would be checked")
	flag.BoolVar(&cmdlineConfig.Fix, "fix", false, "fix supported EditorConfig violations before checking")
	flag.BoolVar(&cmdlineConfig.ShowVersion, "version", false, "print the version number")
	flag.BoolVar(&cmdlineConfig.Help, "help", false, "print the help")
	flag.BoolVar(&cmdlineConfig.Help, "h", false, "print the help")
	flag.TextVar(&cmdlineConfig.Format, "format", outputformat.Default, "specify the output format: "+outputformat.GetArgumentChoiceText())
	flag.TextVar(&cmdlineConfig.Format, "f", outputformat.Default, "specify the output format: "+outputformat.GetArgumentChoiceText())
	flag.BoolVar(&cmdlineConfig.Verbose, "verbose", false, "print debugging information")
	flag.BoolVar(&cmdlineConfig.Verbose, "v", false, "print debugging information")
	flag.BoolVar(&cmdlineConfig.Debug, "debug", false, "print debugging information")
	flag.BoolFunc("no-color", "disables printing color", enableNoColor)
	flag.BoolFunc("color", "enables printing color", disableNoColor)
	flag.BoolVar(&cmdlineConfig.Disable.TrimTrailingWhitespace, "disable-trim-trailing-whitespace", false, "disables the trailing whitespace check")
	flag.BoolVar(&cmdlineConfig.Disable.EndOfLine, "disable-end-of-line", false, "disables the end-of-line check")
	flag.BoolVar(&cmdlineConfig.Disable.InsertFinalNewline, "disable-insert-final-newline", false, "disables the final newline check")
	flag.BoolVar(&cmdlineConfig.Disable.Indentation, "disable-indentation", false, "disables the indentation check")
	flag.BoolVar(&cmdlineConfig.Disable.IndentSize, "disable-indent-size", false, "disables only the indent-size check")
	flag.BoolVar(&cmdlineConfig.Disable.MaxLineLength, "disable-max-line-length", false, "disables only the max-line-length check")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.BoolVar(&cmdlineConfig.Disable.Charset, "disable-charset", false, "disables only the charset check")
}

// parse the arguments from os.Args
func parseArguments() {
	// reset the global variables used to receive the arguments, so parseArguments can be called multiple times without reusing arguments from the previous run
	configFilePath = ""
	cmdlineExclude = ""
	cmdlineConfig = config.Config{}
	writeConfigFile = false
	colorExplicitlySet = false

	// check the NO_COLOR environment variable before parsing the arguments, so the arguments can override
	if nocolor := os.Getenv("NO_COLOR"); nocolor != "" {
		nocolorParsedAsBool, err := strconv.ParseBool(nocolor)
		if err != nil {
			// value did not parse as a boolean,
			// so the user intended to enable NoColor by setting an arbitrary value
			nocolorParsedAsBool = true
		}
		if nocolorParsedAsBool {
			_ = enableNoColor("")
		}
	}

	flag.Parse()

	configPaths := []string{}
	if configFilePath == "" {
		configPaths = append(configPaths, defaultConfigFileNames[:]...)
	} else {
		configPaths = append(configPaths, configFilePath)
	}

	currentConfig = config.NewConfig(configPaths)
	loggerInjectionHook()

	currentConfig.Logger.NoColor = cmdlineConfig.NoColor

	if writeConfigFile {
		err := currentConfig.Save(version)
		if err != nil {
			currentConfig.Logger.Error("%s", err)
			exitProxy(exitCodeErrorOccurred)
		}

		exitProxy(exitCodeNormal)
	}

	err := currentConfig.Parse()
	// this error should be suppressed if the configFilePath was not set by the user
	// since the default config paths could trigger this
	if err != nil && (configFilePath != "" || !errors.Is(err, fs.ErrNotExist)) {
		currentConfig.Logger.Error("%s", err)
		exitProxy(exitCodeConfigFileNotFound)
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

	if !cmdlineConfig.Format.IsValid() && !currentConfig.Format.IsValid() {
		if detected, ok := detectOutputFormat(); ok {
			cmdlineConfig.Format = detected
		}
	}

	// GitHub Actions annotations do not parse ANSI color codes; they render
	// as literal escape sequences and break the annotation format. Default to
	// no-color when the effective output format is `github-actions`, but let
	// an explicit --color / --no-color flag take precedence. See #537.
	effectiveFormat := currentConfig.Format
	if cmdlineConfig.Format.IsValid() {
		effectiveFormat = cmdlineConfig.Format
	}
	if effectiveFormat == outputformat.GithubActions && !colorExplicitlySet {
		cmdlineConfig.NoColor = true
	}

	currentConfig.Merge(cmdlineConfig)
}

// Main function, dude
func main() {
	parseArguments()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			currentConfig.Logger.Error("Creating CPU profile file %s: %v", cpuprofile, err.Error())
			exitProxy(exitCodeErrorOccurred)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			currentConfig.Logger.Error("Starting CPU profile: %v", err.Error())
			exitProxy(exitCodeErrorOccurred)
		}
	}

	config := *currentConfig
	// force the exclude regexp to be compiled and cached
	if _, err := config.CachedExcludesAsRegexp(); err != nil {
		config.Logger.Error("Compiling exclude regexp: %v", err.Error())
		exitProxy(exitCodeErrorOccurred)
	}

	config.Logger.Debug("Config: %s", config)
	config.Logger.Verbose("Exclude Regexp: %s", config.GetExcludesAsRegularExpression())

	if utils.FileExists(config.Path) && config.Version != "" && config.Version != version {
		config.Logger.Error("Version from config file is not the same as the version of the binary")
		config.Logger.Error("Binary: %s, Config %s", version, config.Version)
		exitProxy(exitCodeErrorOccurred)
	}

	// Check for returnworthy arguments
	shouldExit := ReturnableFlags(config)
	if shouldExit {
		exitProxy(exitCodeNormal)
	}

	// contains all files which should be checked
	filePaths, err := files.GetFiles(config)
	if err != nil {
		config.Logger.Error("%s", err)
		exitProxy(exitCodeErrorOccurred)
	}

	if config.DryRun {
		for _, file := range filePaths {
			config.Logger.Output("%s", file)
		}

		exitProxy(exitCodeNormal)
	}

	if config.Fix {
		for _, filePath := range filePaths {
			lock := config.EditorconfigConfig
			def, warnings, err := lock.LoadGraceful(filePath)
			if err != nil {
				config.Logger.Error("cannot load %s as .editorconfig: %s", filePath, err)
				exitProxy(exitCodeErrorOccurred)
			}
			if warnings != nil {
				config.Logger.Warning("%v", warnings.Error())
			}
			changed, err := fix.File(filePath, config, def)
			if err != nil {
				config.Logger.Error("cannot fix %s: %s", filePath, err)
				exitProxy(exitCodeErrorOccurred)
			}
			if changed {
				config.Logger.Verbose("Fixed %s", filePath)
			}
		}
	}

	errors := validation.ProcessValidation(filePaths, config)

	eccerror.PrintErrors(errors, config)

	config.Logger.Verbose("%d files checked", len(filePaths))

	if cpuprofile != "" {
		pprof.StopCPUProfile()
	}

	if eccerror.GetErrorCount(errors) != 0 {
		exitProxy(exitCodeErrorOccurred)
	}

	exitProxy(exitCodeNormal)
}

// ReturnableFlags returns whether a flag passed should exit the program
func ReturnableFlags(config config.Config) bool {
	switch {
	case config.ShowVersion:
		config.Logger.Output("%s", version)
	case config.Help:
		config.Logger.Output("USAGE:")
		config.Logger.Output("  editorconfig-checker [OPTIONS] [FILE...]")
		config.Logger.Output("")
		config.Logger.Output("With no FILE arguments, all files tracked by git are checked. When one or")
		config.Logger.Output("more FILE arguments are given, only those files are checked (the configured")
		config.Logger.Output("exclude patterns still apply).")
		config.Logger.Output("")
		config.Logger.Output("OPTIONS:")
		flag.CommandLine.SetOutput(config.Logger.GetWriter())
		flag.PrintDefaults()
		flag.CommandLine.SetOutput(nil)
	}

	return config.ShowVersion || config.Help
}
