package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/outputformat"
	"github.com/gkampitakis/go-snaps/snaps"
)

var rootConfigFilePath = []string{"../../.editorconfig-checker.json"}
var configWithIgnoredDefaults = []string{"../../testfiles/.editorconfig-checker.json"}

func TestNewConfig(t *testing.T) {
	actual, _ := NewConfig([]string{"abc"})
	var expected Config

	if actual == &expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestNoConfigFileFound(t *testing.T) {
	_, err := NewConfig([]string{"abc"})
	var expected = "no file found at abc"
	if err == nil || err.Error() != expected {
		t.Errorf("expected an error (%s), got %v", expected, err)
	}
}

func TestNoConfigs(t *testing.T) {
	_, err := NewConfig([]string{})
	var expected = "no config paths provided"
	if err == nil || err.Error() != expected {
		t.Errorf("expected an error (%s), got %v", expected, err)
	}
}

func TestNoConfigFileFoundInMultiplePaths(t *testing.T) {
	_, err := NewConfig([]string{"abc", "def"})
	var expected = "no file found at abc"
	if err == nil || err.Error() != expected {
		t.Errorf("expected an error (%s), got %v", expected, err)
	}
}

func TestConfigFileFirstFoundInMultiplePaths(t *testing.T) {
	c, _ := NewConfig([]string{"abc", rootConfigFilePath[0], configWithIgnoredDefaults[0]})
	if c.Path != rootConfigFilePath[0] {
		t.Errorf("expected %s, got %s", rootConfigFilePath[0], c.Path)
	}
}

func TestGetExcludesAsRegularExpression(t *testing.T) {
	c, _ := NewConfig(configWithIgnoredDefaults)
	err := c.Parse()
	if err != nil {
		t.Errorf("Should parse without an error, got: %v", err)
	}

	actual := c.GetExcludesAsRegularExpression()
	expected := "testfiles"

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	c.Exclude = append(c.Exclude, "stuff")

	actual = c.GetExcludesAsRegularExpression()
	expected = "testfiles|stuff"

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	c, _ = NewConfig(rootConfigFilePath)
	err = c.Parse()
	if err != nil {
		t.Errorf("Should parse without an error, got: %v", err)
	}

	actual = c.GetExcludesAsRegularExpression()
	expected = `testfiles|testdata|^\.yarn/|^yarn\.lock$|^package-lock\.json$|^composer\.lock$|^Cargo\.lock$|^Gemfile\.lock$|^\.pnp\.cjs$|^\.pnp\.js$|^\.pnp\.loader\.mjs$|\.snap$|\.otf$|\.woff$|\.woff2$|\.eot$|\.ttf$|\.gif$|\.png$|\.jpg$|\.jpeg$|\.webp$|\.avif$|\.pnm$|\.pbm$|\.pgm$|\.ppm$|\.mp4$|\.wmv$|\.svg$|\.ico$|\.bak$|\.bin$|\.pdf$|\.zip$|\.gz$|\.tar$|\.7z$|\.bz2$|\.log$|\.patch$|\.css\.map$|\.js\.map$|min\.css$|min\.js$`

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestMerge(t *testing.T) {
	modifiedConfig, err := NewConfig([]string{"../../.editorconfig-checker.json"})
	if err != nil {
		t.Errorf("Expected to create a config without errors, got %v", err)
	}

	err = modifiedConfig.Parse()
	if err != nil {
		t.Errorf("Expected to parse a config without errors, got %v", err)
	}

	emptyConfig := Config{}
	modifiedConfig.Merge(emptyConfig)

	parsedConfig, _ := NewConfig([]string{"../../.editorconfig-checker.json"})
	_ = parsedConfig.Parse()

	if !reflect.DeepEqual(modifiedConfig, parsedConfig) {
		t.Errorf("Expected a parsed config and a parsed config merged with an empty config to be the same config, got %v and %v", modifiedConfig, parsedConfig)
	}

	mergeConfig := Config{
		ShowVersion:         true,
		Version:             "v3.0.3",
		Help:                true,
		DryRun:              true,
		Path:                "some-other",
		Verbose:             true,
		Format:              "default",
		Debug:               true,
		NoColor:             true,
		IgnoreDefaults:      true,
		SpacesAfterTabs:     true,
		Exclude:             []string{"some-other"},
		PassedFiles:         []string{"src"},
		AllowedContentTypes: []string{"xml/"},
		Disable: DisabledChecks{
			TrimTrailingWhitespace: true,
			EndOfLine:              true,
			InsertFinalNewline:     true,
			Indentation:            true,
			IndentSize:             true,
			MaxLineLength:          true,
		},
	}

	modifiedConfig.Merge(mergeConfig)

	mergeConfig.AllowedContentTypes = []string{"text/", "application/octet-stream", "application/ecmascript", "application/json", "application/x-ndjson", "application/xml", "+json", "+xml", "xml/"}
	mergeConfig.Exclude = []string{"testfiles", "testdata", "some-other"}

	expected := mergeConfig
	// the following set the properties that cannot be specified directly in mergeConfig above, but would cause the test to fail if left unchanged
	expected.Logger.VerboseEnabled = true
	expected.Logger.DebugEnabled = true
	expected.Logger.NoColor = true
	expected.EditorconfigConfig = modifiedConfig.EditorconfigConfig

	if !reflect.DeepEqual(modifiedConfig, &expected) {
		t.Errorf("%#v", &expected)
		t.Errorf("%#v", modifiedConfig)
		t.Errorf("Expected, got %#v and %#v", modifiedConfig, &expected)
	}

	config := Config{Path: "./.editorconfig-checker.json"}
	err = config.Parse()

	if err == nil {
		t.Errorf("Expected an error to happen when parsing an unexisting file, got %v", err)
	}

	config = Config{Path: "./../../testfiles/.malformed.editorconfig-checker.json"}
	err = config.Parse()

	if err == nil {
		t.Errorf("Expected an error to happen when parsing an unexisting file, got %v", err)
	}
}

func TestParse(t *testing.T) {
	c, _ := NewConfig([]string{"../../testfiles/.editorconfig-checker.json"})
	_ = c.Parse()

	if c.Verbose != true ||
		c.Debug != true ||
		c.IgnoreDefaults != true ||
		!reflect.DeepEqual(c.Exclude, []string{"testfiles"}) ||
		!reflect.DeepEqual(c.AllowedContentTypes, []string{"text/", "application/octet-stream", "application/ecmascript", "application/json", "application/x-ndjson", "application/xml", "+json", "+xml", "hey"}) ||
		c.SpacesAfterTabs != true ||
		c.Disable.EndOfLine != false ||
		c.Disable.TrimTrailingWhitespace != false ||
		c.Disable.InsertFinalNewline != false ||
		c.Disable.Indentation != false {
		t.Error(c.AllowedContentTypes)
		t.Errorf("Expected config to have values from test file, got %v", c)
	}
}

func TestSave(t *testing.T) {
	dir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(dir)
	configFile := filepath.Join(dir, "config")
	c, _ := NewConfig([]string{configFile})
	if c.Save("VERSION") != nil {
		t.Error("Should create the config")
	}

	if c.Save("VERSION") == nil {
		t.Error("Should produce an error")
	}
}

func TestString(t *testing.T) {
	c, _ := NewConfig([]string{"../../.editorconfig-checker.json"})
	_ = c.Parse()
	c.Format = outputformat.Default

	snaps.MatchJSON(t, c.String())
}
