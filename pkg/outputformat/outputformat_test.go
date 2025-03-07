package outputformat

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestIsValid(t *testing.T) {
	for _, f := range ValidOutputFormats {
		if !f.IsValid() {
			t.Errorf("builtin value %s was not found valid. IsValid() is broken", string(f))
		}
	}

	if OutputFormat("nonexistant value").IsValid() {
		t.Error("failed to recognize an invalid value")
	}
}

func TestGetArgumentChoiceText(t *testing.T) {
	snaps.MatchSnapshot(t, GetArgumentChoiceText())
}

func TestMarshalling(t *testing.T) {
	for _, f := range ValidOutputFormats {
		// converting the builtin output formats to text must work
		m, marshalerror := f.MarshalText()
		if marshalerror != nil {
			t.Error(marshalerror)
		}

		var u OutputFormat
		unmarshalerror := u.UnmarshalText(m)
		if unmarshalerror != nil {
			t.Error(unmarshalerror)
		}

		if f != u {
			t.Errorf("marshalling and then unmarshalling of format %s failed to survive the roundtrip", f)
		}
	}
}

func TestMarshallingBrokenInputFails(t *testing.T) {
	broken := OutputFormat("invalid")
	_, err := broken.MarshalText()
	if err == nil {
		t.Error("marshalling did not recognize an invalid value and marshalled it anyway")
	}
}
func TestUnmarshallingBrokenInputFails(t *testing.T) {
	var broken OutputFormat
	err := broken.UnmarshalText([]byte("invalid"))
	if err == nil {
		t.Error("unmarshalling did not recognize an invalid value and unmarshalled it anyway")
	}
}

func TestUnmarshallingEmptyFormat(t *testing.T) {
	var working OutputFormat
	err := working.UnmarshalText([]byte(""))
	if err != nil {
		t.Errorf("unmarshalling an empty string as the default output format failed: %v", err)
	}
	if working != Default {
		t.Error("unmarshalling an empty string did not return the default output format")
	}
}
