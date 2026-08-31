package fix

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	// x-release-please-start-major
	"github.com/editorconfig-checker/editorconfig-checker/v3/internal/config"
	// x-release-please-end
	editorconfig "github.com/editorconfig/editorconfig-core-go/v2"
)

func definition(raw map[string]string) *editorconfig.Definition {
	return &editorconfig.Definition{Raw: raw}
}

func TestFixFile(t *testing.T) {
	tests := []struct {
		name, input, want string
		rules             map[string]string
	}{
		{"lf", "a\r\nb\r", "a\nb\n", map[string]string{"end_of_line": "lf", "insert_final_newline": "true"}},
		{"crlf", "a\nb\r", "a\r\nb\r\n", map[string]string{"end_of_line": "crlf", "insert_final_newline": "true"}},
		{"cr", "a\nb\r\n", "a\rb\r", map[string]string{"end_of_line": "cr", "insert_final_newline": "true"}},
		{"mixed", "a\r\nb\nc\r", "a\nb\nc", map[string]string{"end_of_line": "lf", "insert_final_newline": "false"}},
		{"trailing", "a \n b\t", "a\n b", map[string]string{"trim_trailing_whitespace": "true", "insert_final_newline": "false"}},
		{"empty-final-newline", "", "", map[string]string{"insert_final_newline": "true"}},
		{"no-final-newline", "a\n\n", "a", map[string]string{"insert_final_newline": "false"}},
		{"false-policy", "a\n", "a", map[string]string{"insert_final_newline": "false"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			if err := os.WriteFile(path, []byte(tt.input), 0755); err != nil {
				t.Fatal(err)
			}
			before, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			changed, err := File(path, cfg, definition(tt.rules))
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			if !changed && tt.name == "empty-final-newline" {
				return
			}
			if !changed {
				t.Fatal("expected file to change")
			}
			changed, err = File(path, cfg, definition(tt.rules))
			if err != nil {
				t.Fatal(err)
			}
			if changed {
				t.Fatal("fix is not idempotent")
			}
			after, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if after.Mode().Perm() != before.Mode().Perm() {
				t.Fatalf("permissions changed: %o, want %o", after.Mode().Perm(), before.Mode().Perm())
			}
		})
	}
}

func TestFixFileRespectsDisable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	input := "// editorconfig-checker-disable-next-line\na \n"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := *config.NewConfig(nil)
	changed, err := File(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("disabled line should not be changed")
	}
}

func TestFixFileSkipsEmptyAndNonRegularFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := *config.NewConfig(nil)
	if changed, err := File(dir, cfg, definition(map[string]string{"insert_final_newline": "true"})); err != nil || changed {
		t.Fatalf("directory: changed=%v, err=%v", changed, err)
	}

	empty := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(empty, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if changed, err := File(empty, cfg, definition(map[string]string{"insert_final_newline": "true"})); err != nil || changed {
		t.Fatalf("empty file: changed=%v, err=%v", changed, err)
	}

	if _, err := File(filepath.Join(dir, "missing.txt"), cfg, definition(nil)); err == nil {
		t.Fatal("missing file should return an error")
	}

	noOp := filepath.Join(dir, "no-op.txt")
	if err := os.WriteFile(noOp, []byte("already clean\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if changed, err := File(noOp, cfg, definition(nil)); err != nil || changed {
		t.Fatalf("no-op policy: changed=%v, err=%v", changed, err)
	}
}

func TestFixFileSkipsInvalidRules(t *testing.T) {
	tests := []struct {
		name  string
		rules map[string]string
	}{
		{"invalid-eol", map[string]string{"end_of_line": "dos"}},
		{"invalid-final-newline", map[string]string{"insert_final_newline": "sometimes"}},
		{"unset-rules", map[string]string{"end_of_line": "unset", "insert_final_newline": "unset"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			input := []byte("value\r\n")
			if err := os.WriteFile(path, input, 0644); err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			changed, err := File(path, cfg, definition(tt.rules))
			if err != nil {
				t.Fatal(err)
			}
			if changed {
				t.Fatal("invalid or unset rules should not change the file")
			}
			got, _ := os.ReadFile(path)
			if string(got) != string(input) {
				t.Fatalf("got %q, want %q", got, input)
			}
		})
	}
}

func TestFixFilePreservesExistingNewlinesWhenEOLDisabled(t *testing.T) {
	for _, input := range []string{"value\n", "value\r\n"} {
		t.Run(fmt.Sprintf("%x", input), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			if err := os.WriteFile(path, []byte(input), 0644); err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			cfg.Disable.EndOfLine = true
			changed, err := File(path, cfg, definition(map[string]string{
				"end_of_line":          "crlf",
				"insert_final_newline": "true",
			}))
			if err != nil {
				t.Fatal(err)
			}
			if changed {
				t.Fatal("already-present final newline should be preserved")
			}
			got, _ := os.ReadFile(path)
			if string(got) != input {
				t.Fatalf("got %q, want %q", got, input)
			}
		})
	}
}

func TestFixFileSupportsConfiguredPolicies(t *testing.T) {
	tests := []struct {
		name, input, want, eol string
	}{
		{"lf", "a\r\nb", "a\nb\n", "\n"},
		{"crlf", "a\nb", "a\r\nb\r\n", "\r\n"},
		{"cr", "a\n b\r\n", "a\r b\r", "\r"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			if err := os.WriteFile(path, []byte(tt.input), 0644); err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			changed, err := File(path, cfg, definition(map[string]string{
				"end_of_line":          tt.name,
				"insert_final_newline": "true",
			}))
			if err != nil || !changed {
				t.Fatalf("changed=%v, err=%v", changed, err)
			}
			got, _ := os.ReadFile(path)
			if string(got) != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFixHelpers(t *testing.T) {
	for _, test := range []struct {
		raw, want string
	}{
		{"", ""}, {"unset", ""}, {"lf", "\n"}, {"cr", "\r"}, {"crlf", "\r\n"},
	} {
		got, ok := configuredEOL(test.raw)
		if !ok || got != test.want {
			t.Errorf("configuredEOL(%q) = %q, %v", test.raw, got, ok)
		}
	}
	if _, ok := configuredEOL("invalid"); ok {
		t.Fatal("invalid EOL should be rejected")
	}
	for _, test := range []string{"", "unset", "true", "false"} {
		if _, ok := configuredFinalNewline(test); !ok {
			t.Errorf("configuredFinalNewline(%q) rejected", test)
		}
	}
	if _, ok := configuredFinalNewline("invalid"); ok {
		t.Fatal("invalid final-newline value should be rejected")
	}
	if !isSupportedEncoding("UTF-8-SIG") || !isSupportedEncoding("ascii") || !isSupportedEncoding("unknown") {
		t.Fatal("expected supported encodings")
	}
	if isSupportedEncoding("latin1") {
		t.Fatal("latin1 should remain check-only")
	}

	for _, test := range []struct {
		input, want, policy string
	}{
		{"a", "a\n", "true"},
		{"a\r\n", "a", "false"},
		{"a\r", "a\r", "true"},
	} {
		got := fixFinalNewline([]byte(test.input), test.policy, "")
		if string(got) != test.want {
			t.Errorf("fixFinalNewline(%q, %q) = %q, want %q", test.input, test.policy, got, test.want)
		}
	}
	for _, test := range []struct {
		input, want string
	}{
		{"value\r\n", "\r\n"},
		{"value\n", "\n"},
		{"value\r", "\r"},
		{"value", ""},
	} {
		if got := detectFinalEOL([]byte(test.input)); got != test.want {
			t.Errorf("detectFinalEOL(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestAtomicWriteRefusesLostUpdate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("before"), 0644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("changed concurrently"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := atomicWrite(path, []byte("replacement"), fileMode(info), info); err == nil {
		t.Fatal("expected replacement to be rejected after a concurrent change")
	}
	got, _ := os.ReadFile(path)
	if string(got) != "changed concurrently" {
		t.Fatalf("lost update changed file to %q", got)
	}
}

func TestAtomicWriteFailureBeforeReplacement(t *testing.T) {
	dir := t.TempDir()
	if err := atomicWrite(filepath.Join(dir, "missing.txt"), []byte("replacement"), 0644, nil); err == nil {
		t.Fatal("expected missing target to be rejected")
	}
	if err := atomicWrite(filepath.Join(dir, "missing-parent", "file.txt"), []byte("replacement"), 0644, nil); err == nil {
		t.Fatal("expected temporary-file creation to fail for a missing parent")
	}
}

func TestFixFileDirectiveScopes(t *testing.T) {
	tests := []struct {
		name, input, want string
	}{
		{"disable-file", "// editorconfig-checker-disable-file \na \n", "// editorconfig-checker-disable-file \na \n"},
		{"disable-block", "// editorconfig-checker-disable  \na \n// editorconfig-checker-enable  \nb \n", "// editorconfig-checker-disable  \na \n// editorconfig-checker-enable\nb\n"},
		{"disable-line", "a  // editorconfig-checker-disable-line  \nb  \n", "a  // editorconfig-checker-disable-line  \nb\n"},
		{"disable-next-line", "// editorconfig-checker-disable-next-line  \na  \nb  \n", "// editorconfig-checker-disable-next-line\na  \nb\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			if err := os.WriteFile(path, []byte(tt.input), 0644); err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			changed, err := File(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
			if err != nil {
				t.Fatal(err)
			}
			got, _ := os.ReadFile(path)
			if string(got) != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			if changed == (tt.name == "disable-file") {
				t.Fatalf("changed=%v for %s", changed, tt.name)
			}
		})
	}
}

func TestFixFileSkipsUnsafeFilesAndSymlinks(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"non-utf8", []byte{0xff, 'a', ' ', '\n'}},
		{"invalid-utf8", []byte{0x81, 'a', ' ', '\n'}},
		{"binary", []byte{'a', 0, ' ', '\n'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			if err := os.WriteFile(path, tt.data, 0644); err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			changed, err := File(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
			if err != nil {
				t.Fatal(err)
			}
			if changed {
				t.Fatal("unsafe file was changed")
			}
		})
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	link := filepath.Join(dir, "link")
	if err := os.WriteFile(target, []byte("a \n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	cfg := *config.NewConfig(nil)
	changed, err := File(link, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("symlink was changed")
	}
	got, _ := os.ReadFile(target)
	if string(got) != "a \n" {
		t.Fatalf("symlink target changed: %q", got)
	}
}

func TestFixFilePreservesBOMAndHonorsDisabledChecks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	input := append([]byte{0xef, 0xbb, 0xbf}, []byte("a \r\nb")...)
	if err := os.WriteFile(path, input, 0644); err != nil {
		t.Fatal(err)
	}
	cfg := *config.NewConfig(nil)
	cfg.Disable.TrimTrailingWhitespace = true
	cfg.Disable.EndOfLine = true
	cfg.Disable.InsertFinalNewline = true
	changed, err := File(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true", "end_of_line": "lf", "insert_final_newline": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("disabled checks changed file")
	}
	cfg.Disable = config.DisabledChecks{}
	changed, err = File(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true", "end_of_line": "lf", "insert_final_newline": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("enabled checks did not change file")
	}
	got, _ := os.ReadFile(path)
	want := append([]byte{0xef, 0xbb, 0xbf}, []byte("a\nb\n")...)
	if string(got) != string(want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFixFileDisableEndOfLineStillAddsNewline(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := *config.NewConfig(nil)
	cfg.Disable.EndOfLine = true
	changed, err := File(path, cfg, definition(map[string]string{"end_of_line": "crlf", "insert_final_newline": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("final newline was not fixed")
	}
	got, _ := os.ReadFile(path)
	if string(got) != "a\n" {
		t.Fatalf("got %q, want LF without configured EOL enforcement", got)
	}
}

func TestFixFilePreservesSpecialPermissionBits(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose Unix special permission bits")
	}
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("a "), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0755|os.ModeSetuid); err != nil {
		t.Fatal(err)
	}
	cfg := *config.NewConfig(nil)
	changed, err := File(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("file was not fixed")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0755 || info.Mode()&os.ModeSetuid == 0 {
		t.Fatalf("permission bits changed: got %v", info.Mode())
	}
}
