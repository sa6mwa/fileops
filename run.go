package fileops

import (
	"fmt"
	"os"
	"os/exec"

	"al.essio.dev/pkg/shellescape"
)

func Run(command string) error {
	shell := `/bin/sh`
	shellCommandOption := `-c`
	cmd := exec.Command(shell, shellCommandOption, command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running %q %q %q: %w", shell, shellCommandOption, command, err)
	}
	return nil
}

// Escape is an alias for shellescape.Quote(s) used to escape a
// variable for use in shell command. Returns s escaped.
func Escape(s string) string {
	return shellescape.Quote(s)
}
