package fileops

import (
	"errors"
	"os/user"
)

// HomeDir returns the home directory of the specified user or an
// error if the user does not exist.
func HomeDir(username string) (string, error) {
	usr, err := user.Lookup(username)
	if err != nil {
		return "", orExit(errors.New("user not found"))
	}
	return usr.HomeDir, nil
}
