package fileops

import (
	"os"
	"os/user"
)

// Exists checks if a path (e.g file or directory) exists. Returns
// true if the file or directory exists, false otherwise.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// UserExists is a frontend for user.Lookup(username) returning true
// if user exists, false if not.
func UserExists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}
