package fileops

import "os"

// MkdirAll creates path directory including all parent directories
// (similar to mkdir -p). If a sub directory does not exist, it will
// be created with mode 0755 by default or the value of the first item
// in the optional perm slice. Returns error on failure.
func MkdirAll(path string, perm ...os.FileMode) error {
	var permission os.FileMode = 0755
	if len(perm) > 0 {
		permission = perm[0]
	}
	return orExit(os.MkdirAll(path, permission))
}
