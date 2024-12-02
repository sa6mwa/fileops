package fileops

import "os"

// Exists checks if a path (e.g file or directory) exists. Returns
// true if the file or directory exists, false otherwise.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
