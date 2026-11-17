//go:build !windows

package files

import (
	"errors"
	"io/fs"
)

// isNotExist reports whether err means that there is no file at the path
func isNotExist(err error) bool {
	return errors.Is(err, fs.ErrNotExist)
}
