// fileops is a Go package for various file operations, especially
// related to configuration files and command line operations.
package fileops

// Package wide global variable instructing functions whether to
// actually write to files or run commands.
var DryRun bool = false

// SetDryRun can be used to toggle package-wide dry-run-mode on or
// off. If state is true, no files or content will be
// persisted/written to disk and no commands will run. Instead,
// verbose output will be printed to stderr.
func SetDryRun(state bool) {
	DryRun = state
}
