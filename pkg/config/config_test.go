package config

import (
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
	expected = `testfiles|yarn\.lock$|package-lock\.json|composer\.lock$|\.snap$|\.otf$|\.woff$|\.woff2$|\.eot$|\.ttf$|\.gif$|\.png$|\.jpg$|\.jpeg$|\.mp4$|\.wmv$|\.svg$|\.ico$|\.bak$|\.bin$|\.pdf$|\.zip$|\.gz$|\.tar$|\.7z$|\.bz2$|\.log$|\.css\.map$|\.js\.map$|min\.css$|min\.js$`

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
		Version:           true,
		Help:              true,
		DryRun:            true,
		ConfigPath:        "some-other-path",
		Verbose:           true,
		Debug:             true,
		Ignore_Defaults:   true,
		Spaces_After_tabs: true,
		Exclude:           []string{"some-other"},
		PassedFiles:       []string{"src"},
		Disable: DisabledChecks{
			Trim_Trailing_Whitespace: true,
			End_Of_Line:              true,
			Insert_Final_Newline:     true,
			Indentation:              true,
		},
	}

	c1.Merge(mergeConfig)

	expected := mergeConfig
	expected.Exclude = []string{"testfiles", "some-other"}

	if !reflect.DeepEqual(c1, &expected) {
		t.Errorf("Expected, got %v and %v", c1, &expected)
	}

	config := Config{ConfigPath: "./.ecrc"}
	err = config.Parse()

	if err == nil {
		t.Errorf("Expected an error to happen when parsing an unexisting file, got %v", err)
	}

	config = Config{ConfigPath: "./../../testfiles/.malformed.ecrc"}
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
		c.Ignore_Defaults != true ||
		reflect.DeepEqual(c.Exclude, []string{"../../testfiles"}) ||
		c.Spaces_After_tabs != true ||
		c.Disable.End_Of_Line != false ||
		c.Disable.Trim_Trailing_Whitespace != false ||
		c.Disable.Insert_Final_Newline != false ||
		c.Disable.Indentation != false {
		t.Errorf("Expected config to have values from test file, got %v", c)
	}
}
