// Package fix provides safe fixes for a subset of EditorConfig rules.
package fix

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	editorconfig "github.com/editorconfig/editorconfig-core-go/v2"

	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/config"
	"github.com/editorconfig-checker/editorconfig-checker/v3/pkg/encoding"
)

// FixFile applies the safe, supported fixes for filePath. It returns true when
// the file was changed. Rules which are not explicitly supported remain check-only.
// Files which are not valid UTF-8 (including binary-like files) are deliberately
// not modified: rewriting those files without an encoding-aware encoder could
// corrupt their contents.
func FixFile(filePath string, cfg config.Config, def *editorconfig.Definition) (bool, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		verbose(&cfg, "Skipping %s: symbolic links are not modified", filePath)
		return false, nil
	}
	if !info.Mode().IsRegular() {
		verbose(&cfg, "Skipping %s: not a regular file", filePath)
		return false, nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	// Empty files are intentionally left alone. The validation pipeline treats
	// them as having no lines, regardless of insert_final_newline.
	if len(content) == 0 {
		return false, nil
	}
	if bytes.Contains(content, []byte("\x00")) {
		verbose(&cfg, "Skipping %s: binary-like content", filePath)
		return false, nil
	}
	if hasDisableFile(content) {
		return false, nil
	}
	encodingName, _, _ := encoding.Detect(content)
	normalizedEncoding := strings.ToLower(strings.NewReplacer("-", "", "_", "").Replace(encodingName))
	if normalizedEncoding != "" && normalizedEncoding != "unknown" && normalizedEncoding != "ascii" && normalizedEncoding != "utf8" && normalizedEncoding != "utf8sig" && normalizedEncoding != "utf8bom" {
		verbose(&cfg, "Skipping %s: unsupported encoding %s", filePath, encodingName)
		return false, nil
	}
	if !utf8.Valid(content) {
		verbose(&cfg, "Skipping %s: content is not valid UTF-8", filePath)
		return false, nil
	}

	original := content
	bom := []byte(nil)
	if bytes.HasPrefix(content, []byte{0xef, 0xbb, 0xbf}) {
		bom, content = append([]byte(nil), content[:3]...), content[3:]
	}

	endOfLine := strings.ToLower(def.Raw["end_of_line"])
	if endOfLine == "unset" {
		endOfLine = ""
	}
	desiredEOL := ""
	switch endOfLine {
	case "lf":
		desiredEOL = "\n"
	case "cr":
		desiredEOL = "\r"
	case "crlf":
		desiredEOL = "\r\n"
	case "":
	default:
		// Unknown values are reported by validation; never guess a rewrite.
		return false, nil
	}

	trim := def.Raw["trim_trailing_whitespace"] == "true" && !cfg.Disable.TrimTrailingWhitespace
	insert := strings.ToLower(def.Raw["insert_final_newline"])
	if insert == "unset" {
		insert = ""
	}
	if insert != "" && insert != "true" && insert != "false" {
		return false, nil
	}
	if cfg.Disable.EndOfLine {
		desiredEOL = ""
	}
	if cfg.Disable.InsertFinalNewline {
		insert = ""
	}
	if !trim && desiredEOL == "" && insert == "" {
		return false, nil
	}

	parts := splitLines(content)
	var out bytes.Buffer
	for _, part := range parts {
		line := part.text
		if trim && !part.disabled {
			line = bytes.TrimRight(line, " \t")
		}
		terminator := part.terminator
		if desiredEOL != "" && len(terminator) != 0 {
			terminator = []byte(desiredEOL)
		}
		out.Write(line)
		out.Write(terminator)
	}
	fixed := out.Bytes()

	if insert == "true" {
		eol := desiredEOL
		if eol == "" {
			eol = detectFinalEOL(fixed)
			if eol == "" {
				eol = "\n"
			}
		}
		if !bytes.HasSuffix(fixed, []byte(eol)) {
			fixed = append(fixed, eol...)
		}
	} else if insert == "false" {
		for len(fixed) > 0 {
			if bytes.HasSuffix(fixed, []byte("\r\n")) {
				fixed = fixed[:len(fixed)-2]
			} else if fixed[len(fixed)-1] == '\r' || fixed[len(fixed)-1] == '\n' {
				fixed = fixed[:len(fixed)-1]
			} else {
				break
			}
		}
	}

	fixed = append(bom, fixed...)
	if bytes.Equal(original, fixed) {
		return false, nil
	}
	if err := atomicWrite(filePath, fixed, fileMode(info), info); err != nil {
		return false, err
	}
	return true, nil
}

func verbose(cfg *config.Config, format string, args ...any) {
	if cfg.Logger != nil {
		cfg.Logger.Verbose(format, args...)
	}
}

// atomicWrite writes in the target directory and then renames the complete
// file into place. This is crash-resistant and atomic on platforms/filesystems
// where os.Rename provides those guarantees (Go documents that non-Unix
// replacement is not atomic). It retains supported permission bits. The temp
// file is always cleaned up on failure. Symlinks are rejected by FixFile before
// this function is called. Replacing a path does not update other hard links.
func atomicWrite(path string, content []byte, mode os.FileMode, original os.FileInfo) (err error) {
	dir, base := splitPath(path)
	tmp, err := os.CreateTemp(dir, "."+base+".editorconfig-checker-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", path, err)
	}
	tmpName := tmp.Name()
	defer func() {
		if removeErr := os.Remove(tmpName); err == nil && removeErr != nil && !os.IsNotExist(removeErr) {
			err = fmt.Errorf("remove temporary file %s: %w", tmpName, removeErr)
		}
	}()
	if _, err = tmp.Write(content); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temporary file for %s: %w", path, err)
	}
	if err = tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temporary file for %s: %w", path, err)
	}
	// Set permissions after writing because Unix clears setuid/setgid bits
	// when a file's contents are modified.
	if err = tmp.Chmod(mode); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("set permissions on temporary file for %s: %w", path, err)
	}
	if err = tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temporary file metadata for %s: %w", path, err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("close temporary file for %s: %w", path, err)
	}
	current, statErr := os.Lstat(path)
	if statErr != nil {
		return fmt.Errorf("check %s before replacement: %w", path, statErr)
	}
	if !current.Mode().IsRegular() || !os.SameFile(original, current) || current.Size() != original.Size() || !current.ModTime().Equal(original.ModTime()) {
		return fmt.Errorf("refusing to replace %s: file changed while it was being fixed", path)
	}
	if err = os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}

func fileMode(info os.FileInfo) os.FileMode {
	return info.Mode() & (os.ModePerm | os.ModeSetuid | os.ModeSetgid | os.ModeSticky)
}

func splitPath(path string) (string, string) {
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	return dir, filepath.Base(path)
}

type linePart struct {
	text, terminator []byte
	disabled         bool
}

func splitLines(content []byte) []linePart {
	var result []linePart
	for len(content) > 0 {
		idx := bytes.IndexAny(content, "\r\n")
		if idx < 0 {
			result = append(result, linePart{text: append([]byte(nil), content...)})
			break
		}
		termLen := 1
		if content[idx] == '\r' && idx+1 < len(content) && content[idx+1] == '\n' {
			termLen = 2
		}
		result = append(result, linePart{text: append([]byte(nil), content[:idx]...), terminator: append([]byte(nil), content[idx:idx+termLen]...)})
		content = content[idx+termLen:]
	}
	// Match the line-scoped directives understood by validation. Global rules
	// (line endings and final newlines) are intentionally not line-scoped, but
	// trailing whitespace is.
	const (
		disable              = "editorconfig-checker-disable"
		disableLine          = "editorconfig-checker-disable-line"
		disableNextDirective = "editorconfig-checker-disable-next-line"
		enable               = "editorconfig-checker-enable"
	)
	disabled, disableNext := false, false
	for i := range result {
		line := string(result[i].text)
		if disabled && strings.Contains(line, enable) {
			disabled = false
		}
		result[i].disabled = disabled || disableNext
		if disableNext {
			disableNext = false
		}
		idx := strings.Index(line, disable)
		if idx < 0 {
			continue
		}
		directive := line[idx:]
		if strings.Contains(directive, disableNextDirective) {
			disableNext = true
		}
		if strings.Contains(directive, disableLine) {
			result[i].disabled = true
		}
		if !strings.Contains(directive, disableNextDirective) && !strings.Contains(directive, disableLine) {
			disabled = true
			// A block-disable directive is itself skipped by validation.
			result[i].disabled = true
		}
	}
	return result
}

func detectFinalEOL(content []byte) string {
	if bytes.HasSuffix(content, []byte("\r\n")) {
		return "\r\n"
	}
	if bytes.HasSuffix(content, []byte("\r")) {
		return "\r"
	}
	if bytes.HasSuffix(content, []byte("\n")) {
		return "\n"
	}
	return ""
}

func hasDisableFile(content []byte) bool {
	first := content
	if idx := bytes.IndexAny(first, "\r\n"); idx >= 0 {
		first = first[:idx]
	}
	return bytes.Contains(first, []byte("editorconfig-checker-disable-file"))
}
