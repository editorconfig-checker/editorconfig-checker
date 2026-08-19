package fix

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
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
			cfg := *config.NewConfig(nil)
			changed, err := FixFile(path, cfg, definition(tt.rules))
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
			changed, err = FixFile(path, cfg, definition(tt.rules))
			if err != nil {
				t.Fatal(err)
			}
			if changed {
				t.Fatal("fix is not idempotent")
			}
			mode, _ := os.Stat(path)
			if mode.Mode().Perm() != 0755 {
				t.Fatalf("permissions changed: %o", mode.Mode().Perm())
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
	changed, err := FixFile(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("disabled line should not be changed")
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
			changed, err := FixFile(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
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
		{"binary", []byte{'a', 0, ' ', '\n'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "file.txt")
			if err := os.WriteFile(path, tt.data, 0644); err != nil {
				t.Fatal(err)
			}
			cfg := *config.NewConfig(nil)
			changed, err := FixFile(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
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
	changed, err := FixFile(link, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
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
	changed, err := FixFile(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true", "end_of_line": "lf", "insert_final_newline": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("disabled checks changed file")
	}
	cfg.Disable = config.DisabledChecks{}
	changed, err = FixFile(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true", "end_of_line": "lf", "insert_final_newline": "true"}))
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
	changed, err := FixFile(path, cfg, definition(map[string]string{"end_of_line": "crlf", "insert_final_newline": "true"}))
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
	changed, err := FixFile(path, cfg, definition(map[string]string{"trim_trailing_whitespace": "true"}))
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
