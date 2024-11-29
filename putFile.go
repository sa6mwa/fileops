package fileops

import "os"

func PutFile(textfile string, content string, perm os.FileMode) error {
	if err := os.WriteFile(textfile, []byte(content), perm); err != nil {
		return err
	}
	return os.Chmod(textfile, perm)
}
