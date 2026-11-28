package files

import (
	"errors"
	"io/fs"
	"syscall"
)

// errorInvalidName is the Windows ERROR_INVALID_NAME, returned instead of
// a missing file for a name that no file can have, like a glob pattern with
// a '*' that matched nothing
const errorInvalidName syscall.Errno = 123

// isNotExist reports whether err means that there is no file at the path
func isNotExist(err error) bool {
	return errors.Is(err, fs.ErrNotExist) || errors.Is(err, errorInvalidName)
}
