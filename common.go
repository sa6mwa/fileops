// fileops is a Go package for various file operations, especially
// related to configuration files, command line operations, and
// configuration management.
package fileops

import (
	"fmt"
	"os"
)

// Package wide variable instructing functions whether to actually
// write to files or run commands.
var DryRun bool = false

// Package wide variable instructing functions to call os.Exit(1) on
// error instead of return err. The error message will be printed to
// os.Stdout before terminating.
var ExitOnError bool = false

// SetDryRun can be used to toggle package-wide dry-run-mode on or
// off. If state is true, no files or content will be
// persisted/written to disk and no commands will run. Instead,
// verbose output will be printed to stderr.
func SetDryRun(state bool) {
	DryRun = state
}

// SetExitOnError can be used to set whether functions should exit on
// error instead of return err. See ExitOnError variable.
func SetExitOnError(state bool) {
	ExitOnError = state
}

func orExit(err error) error {
	if ExitOnError && err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		os.Exit(1)
	}
	return err
}
