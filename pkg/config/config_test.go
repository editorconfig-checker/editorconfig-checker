package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const rootConfigFilePath = "../../.ecrc"
const configWithIgnoredDefaults = "../../testfiles/.ecrc"

func TestNewConfig(t *testing.T) {
	actual, _ := NewConfig("abc")
	var expected Config

	if actual == &expected {
		t.Errorf("expected %v, got %v", expected, actual)
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
	expected = `testfiles|^\.yarn/|^yarn\.lock$|^package-lock\.json$|^composer\.lock$|^Cargo\.lock$|^\.pnp\.cjs$|^\.pnp\.js$|^\.pnp\.loader\.mjs$|\.snap$|\.otf$|\.woff$|\.woff2$|\.eot$|\.ttf$|\.gif$|\.png$|\.jpg$|\.jpeg$|\.webp$|\.avif$|\.mp4$|\.wmv$|\.svg$|\.ico$|\.bak$|\.bin$|\.pdf$|\.zip$|\.gz$|\.tar$|\.7z$|\.bz2$|\.log$|\.patch$|\.css\.map$|\.js\.map$|min\.css$|min\.js$`

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestMerge(t *testing.T) {
	c1, err := NewConfig("../../.ecrc")
	if err != nil {
		t.Errorf("Expected to create a config without errors, got %v", err)
	}

	err = c1.Parse()
	if err != nil {
		t.Errorf("Expected to parse a config without errors, got %v", err)
	}

	mergeConfig := Config{}
	c1.Merge(mergeConfig)

	c2, _ := NewConfig("../../.ecrc")
	_ = c2.Parse()

	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("Expected a parsed config and a parsed config merged with an empty config to be the same config, got %v and %v", c1, c2)
	}

	mergeConfig = Config{
		ShowVersion:         true,
		Version:             "2.4.0",
		Help:                true,
		DryRun:              true,
		Path:                "some-other",
		Verbose:             true,
		Debug:               true,
		NoColor:             true,
		IgnoreDefaults:      true,
		SpacesAftertabs:     true,
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

	c1.Merge(mergeConfig)

	mergeConfig.AllowedContentTypes = []string{"text/", "application/octet-stream", "xml/"}
	mergeConfig.Exclude = []string{"testfiles", "some-other"}

	expected := mergeConfig
	expected.Logger.Verbosee = true
	expected.Logger.Debugg = true
	expected.Logger.NoColor = true

	if !reflect.DeepEqual(c1, &expected) {
		t.Errorf("%#v", &expected)
		t.Errorf("%#v", c1)
		t.Errorf("Expected, got %#v and %#v", c1, &expected)
	}

	config := Config{Path: "./.ecrc"}
	err = config.Parse()

	if err == nil {
		t.Errorf("Expected an error to happen when parsing an unexisting file, got %v", err)
	}

	config = Config{Path: "./../../testfiles/.malformed.ecrc"}
	err = config.Parse()

	if err == nil {
		t.Errorf("Expected an error to happen when parsing an unexisting file, got %v", err)
	}
}

func TestParse(t *testing.T) {
	c, _ := NewConfig("../../testfiles/.ecrc")
	_ = c.Parse()

	if c.Verbose != true ||
		c.Debug != true ||
		c.IgnoreDefaults != true ||
		!reflect.DeepEqual(c.Exclude, []string{"testfiles"}) ||
		!reflect.DeepEqual(c.AllowedContentTypes, []string{"text/", "application/octet-stream", "hey"}) ||
		c.SpacesAftertabs != true ||
		c.Disable.EndOfLine != false ||
		c.Disable.TrimTrailingWhitespace != false ||
		c.Disable.InsertFinalNewline != false ||
		c.Disable.Indentation != false {
		t.Error(c.AllowedContentTypes)
		t.Errorf("Expected config to have values from test file, got %v", c)
	}
}

func TestSave(t *testing.T) {
	dir, _ := ioutil.TempDir("", "example")
	defer os.RemoveAll(dir)
	configFile := filepath.Join(dir, "config")
	c, _ := NewConfig(configFile)
	if c.Save("VERSION") != nil {
		t.Error("Should create the config")
	}

	if c.Save("VERSION") == nil {
		t.Error("Should produce an error")
	}
}

func TestGetAsString(t *testing.T) {
	c, _ := NewConfig("../../.ecrc")
	_ = c.Parse()

	actual := c.GetAsString()
	expected := "Config: {ShowVersion:false Help:false DryRun:false Path:../../.ecrc Version:2.4.0 Verbose:false Debug:false IgnoreDefaults:false SpacesAftertabs:false NoColor:false Exclude:[testfiles] AllowedContentTypes:[text/ application/octet-stream] PassedFiles:[] Disable:{EndOfLine:false Indentation:false InsertFinalNewline:false TrimTrailingWhitespace:false IndentSize:false MaxLineLength:false} Logger:{Verbosee:false Debugg:false NoColor:false}}"

	if actual != expected {
		t.Errorf("Expected: %v, got: %v ", expected, actual)
	}
}
